# arcline-docs

Knowledge base and migration guides for arclineit.com. Hosted at `docs.arclineit.com`.

Covers migration from major shared hosts, getting-started guides for new customers, VPS setup walkthroughs, and self-hosting/privacy content that supports the Arcline pitch.

## Status

Planned. Not yet started.

## Format

Static HTML generated from Markdown (Go builder using goldmark). Uses the Arcline design system (same CSS as the main website). No JavaScript framework, no CMS.

## Content areas

| Area | Priority |
|---|---|
| Migration guides (GoDaddy, Bluehost, SiteGround, Namecheap, HostGator, WP Engine) | High |
| Getting started (SSH, SFTP, MySQL backup, email, DNS, SSL) | High |
| WordPress (install, WooCommerce, caching without CDN, security) | Medium |
| VPS setup (Nginx, PHP-FPM, Node.js, Go services, backups, fail2ban) | Medium |
| Privacy & self-hosting (no-CDN tips, arcline-check walkthrough, email auth) | Medium |
| Reference (nameservers, PHP versions, plan limits, AUP summary) | Low |

## Planned features

- Client-side search (pagefind or simple JSON index)
- Sitemap and RSS feed generation
- Watch mode for local development

See [todo.md](todo.md) for the full content plan and builder task list.
