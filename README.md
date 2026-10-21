# Arcline Docs

Security guides, privacy and self-hosting explainers, and client pages for Arcline,
hosted at `docs.arcline.it`.

This is the Ruby on Rails rewrite of the original Go site builder and server
(preserved in the repository history). Pages are rendered on request, so there
is no separate static build.

## What it serves

- **Public guides**: the Markdown files under `content/`, one folder per
  section. 7 guides across security, privacy and self-hosting, and reference.
  The hosting guides and the hosting-era legal pages were removed when Arcline
  stopped offering hosting; they remain in the repository history.
- **My docs** (`/client`): private pages for signed-in customers. Each page is
  visible to everyone signed in, to customers on one plan, or to one customer.
- **Admin** (`/admin/pages`): write and edit the private pages.
- **Search** in the header (⌘K), backed by `/search.json`.
- `/sitemap.xml` and `/rss.xml`.

## Writing a guide

Add a Markdown file at `content/<section>/<slug>.md`. It is served at
`/<section>/<slug>/`.

```markdown
---
title: "Set Up Fail2ban for SSH Brute-Force Protection"
description: "Protect any internet-facing Linux server from SSH brute-force attacks with fail2ban."
section: security
order: 1
---

# Set Up Fail2ban for SSH Brute-Force Protection

Fail2ban monitors system logs for repeated failed login attempts...
```

`order` sorts guides within a section. The page shows the title from the
frontmatter, so a leading `# Title` line is dropped. Markdown is
GitHub-flavoured: tables, strikethrough, task lists, fenced code, and heading
IDs for linking. Raw HTML is passed through.

Sections and their display order are listed in `app/models/guide.rb`. A new
folder appears automatically, after the listed ones.

In production guides are read once at boot, so deploy to publish a change.

## Sign-in

Sign-in goes through Nexus (OpenID Connect, authorization code with PKCE).
Register this site in Nexus under Connected Apps with the redirect URI
`https://docs.arcline.it/auth/nexus/callback`, then set `NEXUS_CLIENT_ID` and
`NEXUS_CLIENT_SECRET`. Without them the public guides still work and the
sign-in link is hidden.

A user is an admin if they hold the Nexus `admin` role or their email matches
`ADMIN_EMAIL`.

Plan-limited pages need `BILLING_DB_PATH`: the billing app's SQLite database,
opened read-only. A customer's plan is found by the email they sign in with.
A page limited to one customer is matched on that same email.

## Requirements

- Ruby 4.0+
- SQLite 3 (bundled via the `sqlite3` gem)

## Setup

```bash
bin/setup --skip-server   # install gems and prepare the database
cp .env.example .env      # then fill in your values
bin/rails server          # http://localhost:3000, or PORT if set
```

## Environment variables

| Variable | Default | Description |
|---|---|---|
| `PORT` | `3000` | HTTP listen port (the original service used `8080`) |
| `DATABASE_PATH` | `storage/<env>.sqlite3` | SQLite database file (client pages) |
| `DOCS_BASE_URL` | `https://docs.arcline.it` | Public address, for canonical links, sitemap and feed |
| `NEXUS_ISSUER` | `https://nexus.arcline.it` | Nexus base URL |
| `NEXUS_CLIENT_ID` | — | OAuth client ID from Nexus |
| `NEXUS_CLIENT_SECRET` | — | OAuth client secret from Nexus |
| `ADMIN_EMAIL` | — | This user may manage client pages |
| `BILLING_DB_PATH` | — | Billing's SQLite database, read-only, for plan lookups |

## Tests

```bash
bin/rails test
bin/rubocop
bin/brakeman
```

The suite renders every guide in `content/`, so a guide that fails to parse
fails the build.

## Deployment

The app ships with a `Dockerfile` and a Kamal config (`config/deploy.yml`).
The SQLite database lives in `storage/`, which is mounted as a volume.

## License

MIT. See [LICENSE](LICENSE).
