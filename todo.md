# arcline-docs — Knowledge Base & Migration Guides

SEO-friendly documentation, migration guides, and self-hosting tutorials.
High-conversion content that directly supports the Arcline pitch.

## Format
- Static HTML (generated from Markdown via a simple Go builder, or hand-written)
- Arcline design system (same CSS as website)
- Hosted at: docs.arclineit.com

## Content Plan

### Migration Guides (priority — high SEO, high conversion)
- [ ] Migrate from GoDaddy Shared Hosting to Arcline
- [ ] Migrate from Bluehost to Arcline
- [ ] Migrate from SiteGround to Arcline
- [ ] Migrate from Namecheap to Arcline
- [ ] Migrate from HostGator to Arcline
- [ ] Migrate WordPress from managed WP hosts (WP Engine, Kinsta) to Arcline VPS
- [ ] Transfer a domain to Arcline (registrar transfer walkthrough)

### Getting Started
- [ ] How to connect to your server via SSH
- [ ] How to upload files via SFTP (FileZilla, Cyberduck)
- [ ] How to create and restore a MySQL database backup
- [ ] How to set up email (MX records, cPanel email accounts)
- [ ] How to point your domain's nameservers to Arcline
- [ ] How to install an SSL certificate (Let's Encrypt via cPanel)

### WordPress
- [ ] Installing WordPress on shared hosting
- [ ] Installing WordPress on a VPS (LAMP stack)
- [ ] Configuring W3 Total Cache without a CDN
- [ ] WordPress security hardening (no third-party CDN required)
- [ ] Setting up WooCommerce on Arcline

### VPS Guides
- [ ] Initial VPS setup (Debian/Ubuntu): users, SSH keys, ufw firewall
- [ ] Install Nginx + PHP-FPM + MySQL on Debian
- [ ] Deploy a static site with Nginx
- [ ] Deploy a Node.js app with PM2 + Nginx
- [ ] Deploy a Go binary as a systemd service
- [ ] Set up automated backups with restic
- [ ] Set up fail2ban for SSH brute-force protection

### Privacy & Self-Hosting
- [ ] Why you shouldn't put Cloudflare in front of everything
- [ ] How to check if your host is actually self-hosted (arcline-check walkthrough)
- [ ] Self-hosting without a CDN: performance tips
- [ ] What SPF, DKIM, and DMARC actually do (and how to set them up)

### Reference
- [ ] Arcline nameservers and DNS records
- [ ] Supported PHP versions
- [ ] Resource limits by plan
- [ ] Acceptable Use Policy summary
- [ ] How to open a support ticket

## Site Builder
- [ ] Simple Go Markdown → HTML builder (goldmark)
- [ ] Template: Arcline design system header/footer/sidebar
- [ ] Search: client-side (pagefind or simple JSON index)
- [ ] Sitemap.xml generation
- [ ] RSS feed for new guides
- [ ] Build script: watch mode for local dev
