package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
)

// ── types ─────────────────────────────────────────────────────────────────────

type Page struct {
	Title       string
	Description string
	Section     string
	Order       int
	Slug        string
	URL         string
	Content     template.HTML
	Excerpt     string
}

type Section struct {
	Slug  string
	Title string
	Pages []*Page
}

type SidebarPage struct {
	Title  string
	URL    string
	Active bool
}

type SidebarSection struct {
	Title string
	URL   string
	Pages []SidebarPage
}

type Breadcrumb struct {
	Label string
	URL   string
}

type PageLink struct {
	Title string
	URL   string
}

type PageData struct {
	Title       string
	Description string
	Content     template.HTML
	Root        string
	Canonical   string
	Nav         []SidebarSection
	Breadcrumbs []Breadcrumb
	Prev        *PageLink
	Next        *PageLink
	Year        int
	Customer    any
	IsAdmin     bool
}

type SearchDoc struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Section string `json:"section"`
	Excerpt string `json:"excerpt"`
}

type RSSItem struct {
	Title       string
	URL         string
	Description string
	Section     string
}

// ── config ────────────────────────────────────────────────────────────────────

const (
	siteBaseURL = "https://docs.arcline.it"
	distDir     = "dist"
	contentDir  = "content"
	staticDir   = "static"
)

var sectionMeta = []struct {
	Slug  string
	Title string
}{
	{"getting-started", "Getting Started"},
	{"migrate", "Migration Guides"},
	{"wordpress", "WordPress"},
	{"vps", "VPS Guides"},
	{"privacy", "Privacy & Self-Hosting"},
	{"reference", "Reference"},
}

// ── main ──────────────────────────────────────────────────────────────────────

func main() {
	watch := false
	for _, arg := range os.Args[1:] {
		if arg == "--watch" {
			watch = true
		}
	}

	if watch {
		runWatch()
	} else {
		build()
	}
}

func build() {
	start := time.Now()

	if err := os.RemoveAll(distDir); err != nil {
		die("rm dist: %v", err)
	}

	pages, err := collectPages()
	if err != nil {
		die("collect pages: %v", err)
	}

	sections := buildSections(pages)

	tmpl := template.New("").Funcs(template.FuncMap{
		"hasPrefix": strings.HasPrefix,
	})
	tmpl, err = template.ParseFiles("templates/page.html", "templates/layout.html")
	if err != nil {
		die("parse template: %v", err)
	}

	for _, sec := range sections {
		for i, page := range sec.Pages {
			var prev, next *PageLink
			if i > 0 {
				prev = &PageLink{Title: sec.Pages[i-1].Title, URL: sec.Pages[i-1].URL}
			}
			if i < len(sec.Pages)-1 {
				next = &PageLink{Title: sec.Pages[i+1].Title, URL: sec.Pages[i+1].URL}
			}
			if err := renderPage(tmpl, page, sections, sec, prev, next); err != nil {
				die("render %s: %v", page.URL, err)
			}
		}
		if err := renderSectionIndex(tmpl, sec, sections); err != nil {
			die("render section index %s: %v", sec.Slug, err)
		}
	}

	if err := renderHome(tmpl, sections); err != nil {
		die("render home: %v", err)
	}
	if err := render404(tmpl); err != nil {
		die("render 404: %v", err)
	}
	if err := copyDir(staticDir, distDir); err != nil {
		die("copy static: %v", err)
	}

	fontSrc := filepath.Join("..", "website", "static", "public", "fonts")
	if info, err := os.Stat(fontSrc); err == nil && info.IsDir() {
		fontDst := filepath.Join(distDir, "public", "fonts")
		if err := copyDir(fontSrc, fontDst); err != nil {
			fmt.Fprintf(os.Stderr, "warn: copy fonts: %v\n", err)
		}
	}

	if err := generateSearch(pages, sections); err != nil {
		die("generate search: %v", err)
	}
	if err := generateSitemap(pages, sections); err != nil {
		die("generate sitemap: %v", err)
	}
	if err := generateRSS(pages, sections); err != nil {
		die("generate rss: %v", err)
	}

	total := 0
	for _, sec := range sections {
		total += len(sec.Pages)
	}
	fmt.Printf("Built %d pages in %s → %s/\n", total, time.Since(start).Round(time.Millisecond), distDir)
}

func runWatch() {
	fmt.Println("Watching for changes (Ctrl+C to stop)...")
	build()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		die("fsnotify: %v", err)
	}
	defer watcher.Close()

	for _, dir := range []string{contentDir, "templates", staticDir} {
		if err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return watcher.Add(path)
			}
			return nil
		}); err != nil {
			die("watch add %s: %v", dir, err)
		}
	}

	debounce := time.NewTimer(0)
	if !debounce.Stop() {
		<-debounce.C
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&(fsnotify.Create|fsnotify.Write|fsnotify.Remove) != 0 {
				debounce.Reset(300 * time.Millisecond)
			}
		case err := <-watcher.Errors:
			slog.Warn("watch error", "err", err)
		case <-debounce.C:
			fmt.Println("\nChange detected, rebuilding...")
			build()
			fmt.Println("Watching for changes (Ctrl+C to stop)...")
		}
	}
}

// ── collect ───────────────────────────────────────────────────────────────────

var md = goldmark.New(
	goldmark.WithExtensions(
		extension.Table,
		extension.Strikethrough,
		extension.TaskList,
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
	goldmark.WithRendererOptions(
		goldmarkhtml.WithUnsafe(),
	),
)

func collectPages() ([]*Page, error) {
	var pages []*Page

	err := filepath.WalkDir(contentDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".md") {
			return err
		}

		raw, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		fm, body := parseFrontmatter(raw)

		var buf bytes.Buffer
		if err := md.Convert(body, &buf); err != nil {
			return fmt.Errorf("convert %s: %w", path, err)
		}

		rel, _ := filepath.Rel(contentDir, path)
		parts := strings.Split(filepath.ToSlash(rel), "/")

		var section, slug string
		if len(parts) == 1 {
			return nil
		}
		section = parts[0]
		slug = strings.TrimSuffix(parts[len(parts)-1], ".md")

		url := "/" + section + "/" + slug + "/"
		order, _ := strconv.Atoi(fm["order"])
		htmlContent := buf.String()

		pages = append(pages, &Page{
			Title:       fm["title"],
			Description: fm["description"],
			Section:     section,
			Order:       order,
			Slug:        slug,
			URL:         url,
			Content:     template.HTML(htmlContent),
			Excerpt:     htmlExcerpt(htmlContent, 220),
		})
		return nil
	})

	return pages, err
}

func parseFrontmatter(raw []byte) (map[string]string, []byte) {
	fm := make(map[string]string)
	s := string(raw)
	if !strings.HasPrefix(s, "---") {
		return fm, raw
	}
	rest := s[3:]
	if rest != "" && rest[0] == '\n' {
		rest = rest[1:]
	}
	end := strings.Index(rest, "---")
	if end < 0 {
		return fm, raw
	}
	block := rest[:end]
	body := rest[end+3:]
	if len(body) > 0 && body[0] == '\n' {
		body = body[1:]
	}

	for _, line := range strings.Split(block, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		idx := strings.Index(line, ":")
		if idx < 0 {
			continue
		}
		k := strings.TrimSpace(line[:idx])
		v := strings.TrimSpace(line[idx+1:])
		v = strings.Trim(v, `"'`)
		fm[k] = v
	}
	return fm, []byte(body)
}

var tagRe = regexp.MustCompile(`<[^>]+>`)

func htmlExcerpt(h string, maxLen int) string {
	text := tagRe.ReplaceAllString(h, " ")
	text = strings.Join(strings.Fields(text), " ")
	if len([]rune(text)) > maxLen {
		r := []rune(text)[:maxLen]
		text = string(r) + "…"
	}
	return text
}

// ── sections ──────────────────────────────────────────────────────────────────

func buildSections(pages []*Page) []*Section {
	bySlug := make(map[string]*Section)
	for _, sm := range sectionMeta {
		bySlug[sm.Slug] = &Section{Slug: sm.Slug, Title: sm.Title}
	}

	for _, p := range pages {
		sec := bySlug[p.Section]
		if sec == nil {
			sec = &Section{
				Slug:  p.Section,
				Title: strings.ReplaceAll(strings.ToTitle(p.Section[:1])+p.Section[1:], "-", " "),
			}
			bySlug[p.Section] = sec
		}
		sec.Pages = append(sec.Pages, p)
	}

	for _, sec := range bySlug {
		sort.Slice(sec.Pages, func(i, j int) bool {
			if sec.Pages[i].Order != sec.Pages[j].Order {
				return sec.Pages[i].Order < sec.Pages[j].Order
			}
			return sec.Pages[i].Title < sec.Pages[j].Title
		})
	}

	var out []*Section
	for _, sm := range sectionMeta {
		if sec := bySlug[sm.Slug]; sec != nil && len(sec.Pages) > 0 {
			out = append(out, sec)
		}
	}
	return out
}

func buildNav(sections []*Section, currentURL string) []SidebarSection {
	nav := make([]SidebarSection, 0, len(sections))
	for _, sec := range sections {
		pages := make([]SidebarPage, 0, len(sec.Pages))
		for _, p := range sec.Pages {
			pages = append(pages, SidebarPage{
				Title:  p.Title,
				URL:    p.URL,
				Active: p.URL == currentURL,
			})
		}
		nav = append(nav, SidebarSection{
			Title: sec.Title,
			URL:   "/" + sec.Slug + "/",
			Pages: pages,
		})
	}
	return nav
}

// ── render ────────────────────────────────────────────────────────────────────

func renderPage(tmpl *template.Template, page *Page, sections []*Section, sec *Section, prev, next *PageLink) error {
	outPath := filepath.Join(distDir, page.Section, page.Slug, "index.html")
	f := mustCreate(outPath)
	defer f.Close()

	return tmpl.Execute(f, PageData{
		Title:       page.Title + " — Arcline Docs",
		Description: page.Description,
		Content:     page.Content,
		Root:        "../../",
		Canonical:   siteBaseURL + page.URL,
		Nav:         buildNav(sections, page.URL),
		Breadcrumbs: []Breadcrumb{
			{Label: sec.Title, URL: "/" + sec.Slug + "/"},
			{Label: page.Title},
		},
		Prev: prev,
		Next: next,
		Year: time.Now().Year(),
	})
}

func renderSectionIndex(tmpl *template.Template, sec *Section, sections []*Section) error {
	outPath := filepath.Join(distDir, sec.Slug, "index.html")
	f := mustCreate(outPath)
	defer f.Close()

	var buf bytes.Buffer
	buf.WriteString(`<div class="section-index">`)
	for _, p := range sec.Pages {
		fmt.Fprintf(&buf, `<a href="%s" class="section-index__card">`, p.URL)
		fmt.Fprintf(&buf, `<span class="section-index__title">%s</span>`, template.HTMLEscapeString(p.Title))
		if p.Description != "" {
			fmt.Fprintf(&buf, `<span class="section-index__desc">%s</span>`, template.HTMLEscapeString(p.Description))
		}
		buf.WriteString(`</a>`)
	}
	buf.WriteString(`</div>`)

	return tmpl.Execute(f, PageData{
		Title:       sec.Title + " — Arcline Docs",
		Description: "Guides in the " + sec.Title + " section.",
		Content:     template.HTML(buf.String()),
		Root:        "../",
		Canonical:   siteBaseURL + "/" + sec.Slug + "/",
		Nav:         buildNav(sections, "/"+sec.Slug+"/"),
		Breadcrumbs: []Breadcrumb{{Label: sec.Title}},
		Year:        time.Now().Year(),
	})
}

func renderHome(tmpl *template.Template, sections []*Section) error {
	outPath := filepath.Join(distDir, "index.html")
	f := mustCreate(outPath)
	defer f.Close()

	var buf bytes.Buffer
	buf.WriteString(`<p class="docs-home__intro">Find guides on getting connected, migrating from other hosts, and managing your Arcline hosting account.</p>`)
	for _, sec := range sections {
		fmt.Fprintf(&buf, `<div class="docs-home__section"><h2 class="docs-home__section-title"><a href="/%s/">%s</a></h2><ul class="docs-home__list">`,
			sec.Slug, template.HTMLEscapeString(sec.Title))
		for _, p := range sec.Pages {
			fmt.Fprintf(&buf, `<li><a href="%s">%s</a>`, p.URL, template.HTMLEscapeString(p.Title))
			if p.Description != "" {
				fmt.Fprintf(&buf, ` — <span class="docs-home__desc">%s</span>`, template.HTMLEscapeString(p.Description))
			}
			buf.WriteString(`</li>`)
		}
		buf.WriteString(`</ul></div>`)
	}

	return tmpl.Execute(f, PageData{
		Title:       "Arcline Docs — Knowledge Base",
		Description: "Guides, tutorials, and reference docs for Arcline hosting customers.",
		Content:     template.HTML(buf.String()),
		Root:        "",
		Canonical:   siteBaseURL + "/",
		Nav:         buildNav(sections, "/"),
		Year:        time.Now().Year(),
	})
}

func render404(tmpl *template.Template) error {
	outPath := filepath.Join(distDir, "404.html")
	f := mustCreate(outPath)
	defer f.Close()

	content := template.HTML(`
<div class="docs-article docs-404">
  <h1>Page not found</h1>
  <p>The page you're looking for doesn't exist or has been moved.</p>
  <p><a href="/" class="btn btn--primary">&larr; Back to docs home</a></p>
</div>`)

	return tmpl.Execute(f, PageData{
		Title:       "Page not found — Arcline Docs",
		Description: "The page you're looking for doesn't exist or has been moved.",
		Content:     content,
		Root:        "",
		Canonical:   "",
		Nav:         nil,
		Year:        time.Now().Year(),
	})
}

// ── static assets ─────────────────────────────────────────────────────────────

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		return copyFile(path, target)
	})
}

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func mustCreate(path string) *os.File {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		die("mkdir %s: %v", filepath.Dir(path), err)
	}
	f, err := os.Create(path)
	if err != nil {
		die("create %s: %v", path, err)
	}
	return f
}

// ── search ────────────────────────────────────────────────────────────────────

func generateSearch(pages []*Page, sections []*Section) error {
	secTitle := make(map[string]string, len(sections))
	for _, sec := range sections {
		secTitle[sec.Slug] = sec.Title
	}

	docs := make([]SearchDoc, 0, len(pages))
	for _, p := range pages {
		docs = append(docs, SearchDoc{
			Title:   p.Title,
			URL:     p.URL,
			Section: secTitle[p.Section],
			Excerpt: p.Excerpt,
		})
	}

	data, err := json.Marshal(docs)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(distDir, "search.json"), data, 0644)
}

// ── sitemap ───────────────────────────────────────────────────────────────────

func generateSitemap(pages []*Page, sections []*Section) error {
	var b strings.Builder
	b.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	b.WriteString("<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n")

	add := func(url, priority string) {
		fmt.Fprintf(&b, "  <url>\n    <loc>%s%s</loc>\n    <changefreq>monthly</changefreq>\n    <priority>%s</priority>\n  </url>\n",
			siteBaseURL, url, priority)
	}

	add("/", "1.0")
	for _, sec := range sections {
		add("/"+sec.Slug+"/", "0.8")
		for _, p := range sec.Pages {
			add(p.URL, "0.7")
		}
	}

	b.WriteString("</urlset>\n")
	return os.WriteFile(filepath.Join(distDir, "sitemap.xml"), []byte(b.String()), 0644)
}

// ── RSS ───────────────────────────────────────────────────────────────────────

func generateRSS(pages []*Page, sections []*Section) error {
	secTitle := make(map[string]string, len(sections))
	for _, sec := range sections {
		secTitle[sec.Slug] = sec.Title
	}

	var items []RSSItem
	for _, sec := range sections {
		for _, p := range sec.Pages {
			items = append(items, RSSItem{
				Title:       p.Title,
				URL:         siteBaseURL + p.URL,
				Description: p.Description,
				Section:     secTitle[p.Section],
			})
		}
	}

	now := time.Now().Format(time.RFC1123Z)

	var b strings.Builder
	b.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	b.WriteString("<rss version=\"2.0\" xmlns:atom=\"http://www.w3.org/2005/Atom\">\n")
	b.WriteString("<channel>\n")
	fmt.Fprintf(&b, "  <title>Arcline Docs</title>\n")
	fmt.Fprintf(&b, "  <link>%s/</link>\n", siteBaseURL)
	fmt.Fprintf(&b, "  <description>Guides, tutorials, and reference docs for Arcline hosting customers.</description>\n")
	fmt.Fprintf(&b, "  <language>en-us</language>\n")
	fmt.Fprintf(&b, "  <lastBuildDate>%s</lastBuildDate>\n", now)
	fmt.Fprintf(&b, "  <atom:link href=\"%s/rss.xml\" rel=\"self\" type=\"application/rss+xml\"/>\n", siteBaseURL)

	for _, item := range items {
		fmt.Fprintf(&b, "  <item>\n")
		fmt.Fprintf(&b, "    <title>%s</title>\n", escapeXML(item.Title))
		fmt.Fprintf(&b, "    <link>%s</link>\n", escapeXML(item.URL))
		fmt.Fprintf(&b, "    <guid>%s</guid>\n", escapeXML(item.URL))
		fmt.Fprintf(&b, "    <description>%s</description>\n", escapeXML(item.Description))
		fmt.Fprintf(&b, "    <category>%s</category>\n", escapeXML(item.Section))
		fmt.Fprintf(&b, "  </item>\n")
	}

	b.WriteString("</channel>\n")
	b.WriteString("</rss>\n")
	return os.WriteFile(filepath.Join(distDir, "rss.xml"), []byte(b.String()), 0644)
}

func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

// ── helpers ───────────────────────────────────────────────────────────────────

func die(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(1)
}
