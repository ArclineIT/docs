# arcline-docs — Knowledge Base

Security guides, privacy and self-hosting explainers, and private client pages.
Hosted at: docs.arcline.it

## Status

**Running.** The Rails app serves the documentation site with search, sitemap, RSS feed, private client pages, and sign-in through Nexus. The hosting content is gone and the security content is still small; see Next.

## Format
- Ruby on Rails 8.1 app, rendering Markdown on request (kramdown, GitHub-flavoured)
- Arcline design system (same CSS as website)
- Client-side search (JSON index with keyboard navigation)
- RSS feed (rss.xml) for new guides
- Sitemap generation (sitemap.xml)

## Content

Hosting was dropped in October 2026. The hosting guides (migration, getting
started, WordPress and VPS setup, PHP versions, resource limits, nameservers)
and the legal pages were removed; they are in the repository history.

### Security ✅
- [x] Set up fail2ban for SSH brute-force protection
- [x] Set up automated backups with restic
- [x] WordPress security hardening

### Privacy & Self-Hosting ✅
- [x] Why you shouldn't put Cloudflare in front of everything
- [x] How to check if your host is actually self-hosted (arcline-check walkthrough)
- [x] What SPF, DKIM, and DMARC actually do (and how to set them up)

### Reference ✅
- [x] How to open a support ticket

### Next
- [ ] New legal pages for security engagements (MSA, DPA, SLA, acceptable use), written with a lawyer
- [ ] Guides for the services Arcline sells: WireGuard remote access, Wazuh monitoring, ZFS backups and restore tests, compliance readiness

## Site ✅
- [x] Port from Go to Ruby on Rails
- [x] Markdown rendering (tables, strikethrough, task lists, heading IDs)
- [x] Arcline design system header/footer/sidebar, breadcrumbs, prev/next
- [x] Search: client-side JSON index (`/search.json`, ⌘K)
- [x] Sitemap (`/sitemap.xml`) and RSS feed (`/rss.xml`)
- [x] Client pages limited by plan or customer, with an admin editor
- [x] Sign-in through Nexus (OpenID Connect)

## Ideas
- [ ] Full-text search over guide bodies, not just the opening
- [ ] Reload guides without a deploy

## Development

```bash
bin/setup --skip-server
bin/rails server
bin/rails test
```
