package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"arclineit.com/docs/internal/store"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
)

// ── Config ─────────────────────────────────────────────────────────────────────

type Config struct {
	Port         string
	BillingDB    string
	DocsDB       string
	AdminEmail   string
	BillingURL   string
	ContentDir   string
	StaticDir    string
	TemplatesDir string
}

func configFromEnv() Config {
	get := func(key, def string) string {
		if v := os.Getenv(key); v != "" {
			return v
		}
		return def
	}
	return Config{
		Port:         get("PORT", "8080"),
		BillingDB:    get("BILLING_DB", "local/billing.db"),
		DocsDB:       get("DOCS_DB", "local/docs.db"),
		AdminEmail:   get("ADMIN_EMAIL", ""),
		BillingURL:   get("BILLING_URL", "https://portal.arcline.it"),
		ContentDir:   get("CONTENT_DIR", "content"),
		StaticDir:    get("STATIC_DIR", "static"),
		TemplatesDir: get("TEMPLATES_DIR", "templates"),
	}
}

// ── Markdown ───────────────────────────────────────────────────────────────────

var md = goldmark.New(
	goldmark.WithExtensions(extension.Table, extension.Strikethrough, extension.TaskList),
	goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	goldmark.WithRendererOptions(goldmarkhtml.WithUnsafe()),
)

var frontmatterRe = regexp.MustCompile(`(?s)^---\n(.+?)\n---\n?`)

type frontmatter struct {
	Title       string
	Description string
	Section     string
	Order       int
}

func parseFrontmatter(src []byte) (frontmatter, []byte) {
	m := frontmatterRe.FindSubmatch(src)
	if m == nil {
		return frontmatter{}, src
	}
	var fm frontmatter
	for _, line := range strings.Split(string(m[1]), "\n") {
		k, v, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		v = strings.TrimSpace(v)
		switch strings.TrimSpace(k) {
		case "title":
			fm.Title = strings.Trim(v, `"`)
		case "description":
			fm.Description = strings.Trim(v, `"`)
		case "section":
			fm.Section = v
		case "order":
			fm.Order, _ = strconv.Atoi(v)
		}
	}
	return fm, src[len(m[0]):]
}

func renderMarkdown(src []byte) (template.HTML, error) {
	var buf bytes.Buffer
	if err := md.Convert(src, &buf); err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

// ── Content loading ────────────────────────────────────────────────────────────

type publicPage struct {
	Title       string
	Description string
	Section     string
	Slug        string
	URL         string
	Order       int
	RawContent  []byte
}

type publicSection struct {
	Slug  string
	Title string
	Pages []*publicPage
}

var sectionMeta = []struct{ Slug, Title string }{
	{"getting-started", "Getting Started"},
	{"migrate", "Migration Guides"},
	{"wordpress", "WordPress"},
	{"vps", "VPS Guides"},
	{"privacy", "Privacy & Self-Hosting"},
	{"reference", "Reference"},
}

func loadContent(contentDir string) ([]*publicSection, error) {
	bySlug := map[string]*publicSection{}

	err := filepath.WalkDir(contentDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".md") {
			return err
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		fm, body := parseFrontmatter(raw)

		rel, _ := filepath.Rel(contentDir, path)
		parts := strings.Split(filepath.ToSlash(rel), "/")
		if len(parts) != 2 {
			return nil
		}
		secSlug := parts[0]
		slug := strings.TrimSuffix(parts[1], ".md")

		if _, ok := bySlug[secSlug]; !ok {
			title := secSlug
			for _, m := range sectionMeta {
				if m.Slug == secSlug {
					title = m.Title
					break
				}
			}
			bySlug[secSlug] = &publicSection{Slug: secSlug, Title: title}
		}

		bySlug[secSlug].Pages = append(bySlug[secSlug].Pages, &publicPage{
			Title:       fm.Title,
			Description: fm.Description,
			Section:     secSlug,
			Slug:        slug,
			URL:         "/" + secSlug + "/" + slug + "/",
			Order:       fm.Order,
			RawContent:  body,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	var sections []*publicSection
	for _, m := range sectionMeta {
		if sec, ok := bySlug[m.Slug]; ok {
			sort.Slice(sec.Pages, func(i, j int) bool {
				return sec.Pages[i].Order < sec.Pages[j].Order
			})
			sections = append(sections, sec)
		}
	}
	return sections, nil
}

// ── Template data ──────────────────────────────────────────────────────────────

type navPage    struct{ Title, URL string; Active bool }
type navSection struct{ Title, URL string; Pages []navPage }
type breadcrumb struct{ Label, URL string }
type pageLink   struct{ Title, URL string }

// TD is the single template data struct used by every template.
type TD struct {
	// Page meta
	Title, Description string
	Content            template.HTML
	Canonical          string
	Year               int

	// Auth
	Customer *store.Customer
	IsAdmin  bool

	// Public doc sidebar/nav (non-nil on public doc pages only)
	Nav         []navSection
	Breadcrumbs []breadcrumb
	Prev, Next  *pageLink

	// Client docs listing
	Pages []*store.Page

	// Single client/admin page
	Page *store.Page

	// Admin edit form
	PlanOptions []store.PlanOption
	Customers   []store.Customer
	Error       string
	CSRFToken   string
}

// ── Handler ────────────────────────────────────────────────────────────────────

type handler struct {
	cfg      Config
	st       *store.Store
	tmpl     *template.Template
	sections []*publicSection
}

func newHandler(cfg Config, st *store.Store) (*handler, error) {
	sections, err := loadContent(cfg.ContentDir)
	if err != nil {
		return nil, fmt.Errorf("load content: %w", err)
	}

	files, err := filepath.Glob(cfg.TemplatesDir + "/*.html")
	if err != nil || len(files) == 0 {
		return nil, fmt.Errorf("no templates in %s", cfg.TemplatesDir)
	}
	funcs := template.FuncMap{
		"visLabel":  store.VisibilityLabel,
		"hasPrefix": strings.HasPrefix,
	}
	tmpl, err := template.New("").Funcs(funcs).ParseFiles(files...)
	if err != nil {
		return nil, fmt.Errorf("parse templates: %w", err)
	}

	return &handler{cfg: cfg, st: st, tmpl: tmpl, sections: sections}, nil
}

func (h *handler) render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tmpl.ExecuteTemplate(w, name, data); err != nil {
		slog.Error("template error", "name", name, "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// ── Auth ───────────────────────────────────────────────────────────────────────

func (h *handler) currentCustomer(r *http.Request) (*store.Customer, *store.Subscription) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return nil, nil
	}
	c, err := h.st.GetCustomerBySession(cookie.Value)
	if err != nil || c == nil {
		return nil, nil
	}
	sub, _ := h.st.GetSubscription(c.ID)
	return c, sub
}

func (h *handler) isAdmin(c *store.Customer) bool {
	return c != nil && h.cfg.AdminEmail != "" && c.Email == h.cfg.AdminEmail
}

func (h *handler) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c, _ := h.currentCustomer(r); c == nil {
			http.Redirect(w, r, h.cfg.BillingURL+"/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *handler) requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := h.currentCustomer(r)
		if !h.isAdmin(c) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ── Nav ────────────────────────────────────────────────────────────────────────

func (h *handler) buildNav(activeSec, activeSlug string) []navSection {
	nav := make([]navSection, 0, len(h.sections))
	for _, sec := range h.sections {
		pages := make([]navPage, 0, len(sec.Pages))
		for _, p := range sec.Pages {
			pages = append(pages, navPage{
				Title:  p.Title,
				URL:    p.URL,
				Active: sec.Slug == activeSec && p.Slug == activeSlug,
			})
		}
		nav = append(nav, navSection{Title: sec.Title, URL: "/" + sec.Slug + "/", Pages: pages})
	}
	return nav
}

// ── Public handlers ────────────────────────────────────────────────────────────

func (h *handler) home(w http.ResponseWriter, r *http.Request) {
	c, _ := h.currentCustomer(r)

	var sb strings.Builder
	sb.WriteString(`<div class="docs-article"><h1>Arcline Documentation</h1>` +
		`<p>Guides for getting started with Arcline hosting, migrating from other providers, and managing your account.</p>`)
	for _, sec := range h.sections {
		sb.WriteString(fmt.Sprintf(
			`<div class="docs-home__section"><h2 class="docs-home__section-title"><a href="/%s/">%s</a></h2><ul class="docs-home__list">`,
			sec.Slug, sec.Title,
		))
		for _, p := range sec.Pages {
			sb.WriteString(fmt.Sprintf(`<li><a href="%s">%s</a>`, p.URL, p.Title))
			if p.Description != "" {
				sb.WriteString(fmt.Sprintf(` <span class="docs-home__desc">— %s</span>`, p.Description))
			}
			sb.WriteString(`</li>`)
		}
		sb.WriteString(`</ul></div>`)
	}
	if c != nil {
		sb.WriteString(`<div class="docs-home__section"><h2 class="docs-home__section-title">` +
			`<a href="/client/">My Docs</a></h2>` +
			`<p class="docs-home__desc">Private guides and documentation specific to your account.</p></div>`)
	}
	sb.WriteString(`</div>`)

	h.render(w, "page.html", TD{
		Title:     "Arcline Documentation",
		Canonical: "https://docs.arcline.it/",
		Nav:       h.buildNav("", ""),
		Year:      time.Now().Year(),
		Content:   template.HTML(sb.String()),
		Customer:  c, IsAdmin: h.isAdmin(c),
	})
}

// pathParts splits a URL path into its non-empty segments.
// "/getting-started/ssh/" → ["getting-started", "ssh"]
func pathParts(p string) []string {
	var parts []string
	for _, s := range strings.Split(strings.Trim(p, "/"), "/") {
		if s != "" {
			parts = append(parts, s)
		}
	}
	return parts
}

func (h *handler) sectionIndex(w http.ResponseWriter, r *http.Request) {
	parts := pathParts(r.URL.Path)
	if len(parts) != 1 {
		http.NotFound(w, r)
		return
	}
	slug := parts[0]
	var sec *publicSection
	for _, s := range h.sections {
		if s.Slug == slug {
			sec = s
			break
		}
	}
	if sec == nil {
		http.NotFound(w, r)
		return
	}
	c, _ := h.currentCustomer(r)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(
		`<div class="docs-article"><h1>%s</h1><div class="section-index">`, sec.Title,
	))
	for _, p := range sec.Pages {
		sb.WriteString(fmt.Sprintf(
			`<a href="%s" class="section-index__card"><span class="section-index__title">%s</span>`+
				`<span class="section-index__desc">%s</span></a>`,
			p.URL, p.Title, p.Description,
		))
	}
	sb.WriteString(`</div></div>`)

	h.render(w, "page.html", TD{
		Title:       sec.Title + " — Arcline Docs",
		Canonical:   "https://docs.arcline.it/" + slug + "/",
		Nav:         h.buildNav(slug, ""),
		Breadcrumbs: []breadcrumb{{Label: sec.Title}},
		Year:        time.Now().Year(),
		Content:     template.HTML(sb.String()),
		Customer:    c, IsAdmin: h.isAdmin(c),
	})
}

func (h *handler) publicPage(w http.ResponseWriter, r *http.Request) {
	parts := pathParts(r.URL.Path)
	if len(parts) != 2 {
		http.NotFound(w, r)
		return
	}
	secSlug, pageSlug := parts[0], parts[1]

	var sec *publicSection
	for _, s := range h.sections {
		if s.Slug == secSlug {
			sec = s
			break
		}
	}
	if sec == nil {
		http.NotFound(w, r)
		return
	}

	var pg *publicPage
	var idx int
	for i, p := range sec.Pages {
		if p.Slug == pageSlug {
			pg, idx = p, i
			break
		}
	}
	if pg == nil {
		http.NotFound(w, r)
		return
	}

	content, err := renderMarkdown(pg.RawContent)
	if err != nil {
		http.Error(w, "render error", http.StatusInternalServerError)
		return
	}

	var prev, next *pageLink
	if idx > 0 {
		prev = &pageLink{sec.Pages[idx-1].Title, sec.Pages[idx-1].URL}
	}
	if idx < len(sec.Pages)-1 {
		next = &pageLink{sec.Pages[idx+1].Title, sec.Pages[idx+1].URL}
	}

	c, _ := h.currentCustomer(r)
	h.render(w, "page.html", TD{
		Title: pg.Title + " — Arcline Docs", Description: pg.Description,
		Canonical: "https://docs.arcline.it/" + secSlug + "/" + pageSlug + "/",
		Nav:       h.buildNav(secSlug, pageSlug),
		Breadcrumbs: []breadcrumb{
			{Label: sec.Title, URL: "/" + secSlug + "/"},
			{Label: pg.Title},
		},
		Prev: prev, Next: next,
		Year:     time.Now().Year(),
		Content:  template.HTML(`<div class="prose">`) + content + `</div>`,
		Customer: c, IsAdmin: h.isAdmin(c),
	})
}

// ── Client handlers ────────────────────────────────────────────────────────────

func (h *handler) clientIndex(w http.ResponseWriter, r *http.Request) {
	c, sub := h.currentCustomer(r)
	pages, err := h.st.ListVisiblePages(c, sub)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	h.render(w, "client_index.html", TD{
		Year: time.Now().Year(), Customer: c, IsAdmin: h.isAdmin(c), Pages: pages,
	})
}

func (h *handler) clientPage(w http.ResponseWriter, r *http.Request) {
	c, sub := h.currentCustomer(r)
	pg, err := h.st.GetPageBySlug(r.PathValue("slug"))
	if err != nil || pg == nil {
		http.NotFound(w, r)
		return
	}
	if !store.CanSee(c, sub, pg.Visibility) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	content, err := renderMarkdown([]byte(pg.Content))
	if err != nil {
		http.Error(w, "render error", http.StatusInternalServerError)
		return
	}
	h.render(w, "client_page.html", TD{
		Year: time.Now().Year(), Customer: c, IsAdmin: h.isAdmin(c),
		Page:    pg,
		Content: template.HTML(`<div class="prose">`) + content + `</div>`,
	})
}

// ── Admin handlers ─────────────────────────────────────────────────────────────

func (h *handler) adminPages(w http.ResponseWriter, r *http.Request) {
	c, _ := h.currentCustomer(r)
	pages, err := h.st.ListPages()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	h.render(w, "admin_pages.html", TD{
		Year: time.Now().Year(), Customer: c, IsAdmin: true, Pages: pages,
		CSRFToken: csrfGet(w, r),
	})
}

func (h *handler) adminNewPage(w http.ResponseWriter, r *http.Request) {
	c, _ := h.currentCustomer(r)
	customers, _ := h.st.ListCustomers()
	h.render(w, "admin_edit.html", TD{
		Year: time.Now().Year(), Customer: c, IsAdmin: true,
		PlanOptions: store.PlanOptions(), Customers: customers,
		CSRFToken: csrfGet(w, r),
	})
}

func (h *handler) adminCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil || !csrfCheck(r) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	pg := pageFromForm(r, 0)
	if pg.Title == "" || pg.Slug == "" {
		h.adminEditError(w, r, pg, "Title and slug are required.")
		return
	}
	if _, err := h.st.CreatePage(pg); err != nil {
		msg := "Could not save page."
		if strings.Contains(err.Error(), "UNIQUE") {
			msg = "A page with that slug already exists."
		}
		h.adminEditError(w, r, pg, msg)
		return
	}
	http.Redirect(w, r, "/admin/", http.StatusSeeOther)
}

func (h *handler) adminEditPage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	pg, err := h.st.GetPageByID(id)
	if err != nil || pg == nil {
		http.NotFound(w, r)
		return
	}
	c, _ := h.currentCustomer(r)
	customers, _ := h.st.ListCustomers()
	h.render(w, "admin_edit.html", TD{
		Year: time.Now().Year(), Customer: c, IsAdmin: true, Page: pg,
		PlanOptions: store.PlanOptions(), Customers: customers,
		CSRFToken: csrfGet(w, r),
	})
}

func (h *handler) adminUpdate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil || !csrfCheck(r) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	pg := pageFromForm(r, id)
	if pg.Title == "" || pg.Slug == "" {
		h.adminEditError(w, r, pg, "Title and slug are required.")
		return
	}
	if err := h.st.UpdatePage(pg); err != nil {
		msg := "Could not update page."
		if strings.Contains(err.Error(), "UNIQUE") {
			msg = "A page with that slug already exists."
		}
		h.adminEditError(w, r, pg, msg)
		return
	}
	http.Redirect(w, r, "/admin/", http.StatusSeeOther)
}

func (h *handler) adminDelete(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil || !csrfCheck(r) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	_ = h.st.DeletePage(id)
	http.Redirect(w, r, "/admin/", http.StatusSeeOther)
}

func (h *handler) adminEditError(w http.ResponseWriter, r *http.Request, pg *store.Page, msg string) {
	c, _ := h.currentCustomer(r)
	customers, _ := h.st.ListCustomers()
	h.render(w, "admin_edit.html", TD{
		Year: time.Now().Year(), Customer: c, IsAdmin: true, Page: pg,
		PlanOptions: store.PlanOptions(), Customers: customers,
		Error: msg, CSRFToken: csrfGet(w, r),
	})
}

// ── Logout ─────────────────────────────────────────────────────────────────────

func (h *handler) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name: "session", Value: "", Path: "/", MaxAge: -1, HttpOnly: true,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ── CSRF ───────────────────────────────────────────────────────────────────────

func csrfGet(w http.ResponseWriter, r *http.Request) string {
	if c, err := r.Cookie("csrf"); err == nil && len(c.Value) == 64 {
		return c.Value
	}
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	token := hex.EncodeToString(b)
	http.SetCookie(w, &http.Cookie{
		Name: "csrf", Value: token, Path: "/",
		SameSite: http.SameSiteStrictMode,
	})
	return token
}

func csrfCheck(r *http.Request) bool {
	c, err := r.Cookie("csrf")
	if err != nil {
		return false
	}
	return c.Value != "" && r.FormValue("csrf_token") == c.Value
}

// ── Form helpers ───────────────────────────────────────────────────────────────

func pageFromForm(r *http.Request, id int64) *store.Page {
	order, _ := strconv.Atoi(r.FormValue("display_order"))
	var visibility string
	switch r.FormValue("vis_type") {
	case "plan":
		visibility = "plan:" + r.FormValue("vis_plan")
	case "customer":
		visibility = "customer:" + r.FormValue("vis_customer")
	default:
		visibility = "public"
	}
	return &store.Page{
		ID:           id,
		Title:        strings.TrimSpace(r.FormValue("title")),
		Slug:         strings.TrimSpace(r.FormValue("slug")),
		Description:  strings.TrimSpace(r.FormValue("description")),
		Section:      strings.TrimSpace(r.FormValue("section")),
		Content:      r.FormValue("content"),
		Visibility:   visibility,
		DisplayOrder: order,
	}
}

// ── Main ───────────────────────────────────────────────────────────────────────

func main() {
	cfg := configFromEnv()

	st, err := store.New(cfg.BillingDB, cfg.DocsDB)
	if err != nil {
		slog.Error("store init failed", "err", err)
		os.Exit(1)
	}
	defer st.Close()

	h, err := newHandler(cfg, st)
	if err != nil {
		slog.Error("handler init failed", "err", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()

	// Auth
	mux.HandleFunc("GET /logout", h.logout)

	// Admin (all fixed-prefix, no wildcard subtree conflicts)
	mux.Handle("GET /admin/", h.requireAdmin(http.HandlerFunc(h.adminPages)))
	mux.Handle("GET /admin/pages/new", h.requireAdmin(http.HandlerFunc(h.adminNewPage)))
	mux.Handle("POST /admin/pages", h.requireAdmin(http.HandlerFunc(h.adminCreate)))
	mux.Handle("GET /admin/pages/{id}/edit", h.requireAdmin(http.HandlerFunc(h.adminEditPage)))
	mux.Handle("POST /admin/pages/{id}", h.requireAdmin(http.HandlerFunc(h.adminUpdate)))
	mux.Handle("POST /admin/pages/{id}/delete", h.requireAdmin(http.HandlerFunc(h.adminDelete)))

	// Client docs (fixed prefix /client/, no conflict with admin)
	mux.Handle("GET /client/", h.requireAuth(http.HandlerFunc(h.clientIndex)))
	mux.Handle("GET /client/{slug}/", h.requireAuth(http.HandlerFunc(h.clientPage)))

	// Root catch-all: handles home, public doc pages, and static files.
	// Wildcard subtree patterns like /{section}/{slug}/ conflict with every
	// other subtree in Go 1.22's mux, so we route public docs manually here.
	staticServer := http.FileServer(http.Dir(cfg.StaticDir))
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path
		// Static assets live in the static dir under /css/, /js/, /public/.
		if strings.HasPrefix(p, "/css/") || strings.HasPrefix(p, "/js/") || strings.HasPrefix(p, "/public/") {
			staticServer.ServeHTTP(w, r)
			return
		}
		parts := pathParts(p)
		switch len(parts) {
		case 0:
			h.home(w, r)
		case 1:
			h.sectionIndex(w, r)
		case 2:
			h.publicPage(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	addr := ":" + cfg.Port
	slog.Info("docs server starting", "addr", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		slog.Error("server error", "err", err)
		os.Exit(1)
	}
}
