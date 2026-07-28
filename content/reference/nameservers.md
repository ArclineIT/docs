---
title: "Arcline Nameservers and DNS Records"
description: "List of Arcline nameservers and instructions for pointing your domain to Arcline."
section: reference
order: 1
---

# Arcline Nameservers and DNS Records

This page lists the Arcline authoritative nameservers and the DNS records you'll need for common services.

---

## Authoritative nameservers

Point your domain to Arcline using these nameservers:

| Purpose | Nameserver |
|---------|------------|
| Primary | `ns1.arcline.it` |
| Secondary | `ns2.arcline.it` |

Both nameservers serve all zones authoritatively. They run NSD on Arcline's own infrastructure.

---

## How to point your domain

At your registrar's control panel, replace the existing nameservers with:

```
ns1.arcline.it
ns2.arcline.it
```

DNS propagation typically takes 15 minutes to a few hours, but can take up to 48 hours in rare cases. See [Point Your Domain to Arcline](/getting-started/nameservers/) for registrar-specific walkthroughs.

---

## Common DNS records

### Website (A record)

| Type | Name | Value |
|------|------|-------|
| A | `@` | `your-server-ip` |
| A | `www` | `your-server-ip` |

### Email (MX records)

| Type | Name | Priority | Value |
|------|------|----------|-------|
| MX | `@` | 0 | `mail.yourdomain.com` |

### Email authentication (TXT records)

| Type | Name | Value |
|------|------|-------|
| TXT | `@` | `v=spf1 mx include:arclineit.com ~all` |
| TXT | `arcline._domainkey` | `v=DKIM1; h=sha256; p=...` |
| TXT | `_dmarc` | `v=DMARC1; p=quarantine; rua=mailto:dmarc@yourdomain.com` |

### Subdomain examples

| Type | Name | Value |
|------|------|-------|
| A | `portal` | `your-server-ip` |
| A | `git` | `your-server-ip` |
| CNAME | `docs` | `yourdomain.com` |
| CNAME | `status` | `yourdomain.com` |

---

## Managing DNS records

In Arcline cPanel, go to **Domains → Zone Editor** to add, edit, or remove DNS records. Changes are instant — NSD reloads the zone automatically.

