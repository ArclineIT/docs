# arcline-docs — Knowledge Base & Migration Guides

SEO-friendly documentation, migration guides, and self-hosting tutorials.
Hosted at: docs.arcline.it

## Status

**Complete.** All content sections are written. The Rails app serves the documentation site with search, sitemap, RSS feed, private client pages, and sign-in through Nexus.

## Format
- Ruby on Rails 8.1 app, rendering Markdown on request (kramdown, GitHub-flavoured)
- Arcline design system (same CSS as website)
- Client-side search (JSON index with keyboard navigation)
- RSS feed (rss.xml) for new guides
- Sitemap generation (sitemap.xml)

## Content

### Migration Guides ✅
- [x] Migrate from GoDaddy Shared Hosting to Arcline
- [x] Migrate from Bluehost to Arcline
- [x] Migrate from SiteGround to Arcline
- [x] Migrate from Namecheap to Arcline
- [x] Migrate from HostGator to Arcline
- [x] Migrate from WP Engine to Arcline VPS
- [x] Transfer a domain to Arcline (registrar transfer walkthrough)

### Getting Started ✅
- [x] How to connect to your server via SSH
- [x] How to upload files via SFTP (FileZilla, Cyberduck)
- [x] How to create and restore a MySQL database backup
- [x] How to set up email (MX records, cPanel email accounts)
- [x] How to point your domain's nameservers to Arcline
- [x] How to install an SSL certificate (Let's Encrypt via cPanel)

### WordPress ✅
- [x] Installing WordPress on shared hosting
- [x] Installing WordPress on a VPS (LAMP stack)
- [x] Configuring W3 Total Cache without a CDN
- [x] WordPress security hardening
- [x] Setting up WooCommerce on Arcline

### VPS Guides ✅
- [x] Initial VPS setup (Debian/Ubuntu): users, SSH keys, ufw firewall
- [x] Install Nginx + PHP-FPM + MySQL on Debian
- [x] Deploy a static site with Nginx
- [x] Deploy a Node.js app with PM2 + Nginx
- [x] Deploy a Go binary as a systemd service
- [x] Set up automated backups with restic
- [x] Set up fail2ban for SSH brute-force protection

### Privacy & Self-Hosting ✅
- [x] Why you shouldn't put Cloudflare in front of everything
- [x] How to check if your host is actually self-hosted (arcline-check walkthrough)
- [x] Self-hosting without a CDN: performance tips
- [x] What SPF, DKIM, and DMARC actually do (and how to set them up)

### Reference ✅
- [x] Arcline nameservers and DNS records
- [x] Supported PHP versions
- [x] Resource limits by plan
- [x] Acceptable Use Policy summary
- [x] How to open a support ticket

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
