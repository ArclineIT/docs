# arcline-docs

Knowledge base, migration guides, and self-hosting tutorials for Arcline. Hosted at `docs.arclineit.com`.

## Status

**Complete.** All content sections are written (33 guides total across 6 sections). The Go static site builder generates a fully functional documentation site with search, sitemap, RSS feed, and 404 pages.

## Content

| Area | Status | Guides |
|------|--------|--------|
| Migration guides (GoDaddy, Bluehost, SiteGround, Namecheap, HostGator, WP Engine, domain transfer) | ✅ | 7 |
| Getting started (SSH, SFTP, MySQL backup, email, DNS, SSL) | ✅ | 6 |
| WordPress (install shared/VPS, WooCommerce, caching, security) | ✅ | 5 |
| VPS setup (initial setup, Nginx+PHP+MySQL, static/Node.js/Go deployment, backups, fail2ban) | ✅ | 7 |
| Privacy & self-hosting (no-CDN, arcline-check, performance, email auth) | ✅ | 4 |
| Reference (nameservers, PHP versions, plan limits, AUP, support) | ✅ | 5 |
| Legal (DPA, MSA, SLA) | ✅ | 3 |

## Site Builder Features

- **Markdown → HTML** — goldmark renderer with tables, strikethroughs, task lists
- **Arcline design system** — same CSS as the main website
- **Client-side search** — JSON index with keyboard navigation (⌘K)
- **Sitemap** — `sitemap.xml` for SEO
- **RSS feed** — `rss.xml` for content syndication
- **404 page** — custom error page with search
- **Navigation** — sidebar with active page highlighting, breadcrumbs, prev/next pagers
- **Watch mode** — `--watch` flag rebuilds on content/template changes

## Development

```bash
# Build the dynamic server
make build

# Run the dev server
make run

# Build the static site (outputs to dist/)
make build-static

# Watch mode — rebuilds on content/template changes
make build-static-watch
```

## Tech Stack

- **Language**: Go
- **Markdown**: goldmark (tables, strikethrough, task lists, auto heading IDs)
- **Templates**: Go html/template
- **Storage**: SQLite (for client/admin pages)
- **Watch mode**: fsnotify
- **Static output**: `dist/` directory (self-hostable via any web server)

## License

MIT — see [LICENSE](LICENSE).
