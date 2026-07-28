---
title: "Point Your Domain to Arcline"
description: "How to update your nameservers at GoDaddy, Namecheap, Google Domains, and other registrars to point your domain to Arcline hosting."
section: getting-started
order: 5
---

# Point Your Domain to Arcline

When you register a domain, it points to your registrar's nameservers by default. To use your domain with Arcline hosting, you need to update the nameservers so that DNS queries for your domain are answered by Arcline's servers.

---

## What nameservers do

Nameservers are the authoritative source for your domain's DNS records. They answer questions like "where does `yourdomain.com` point?" and "what mail server handles `@yourdomain.com`?"

When you switch to Arcline nameservers:
- Arcline controls all DNS records for your domain
- Your website, email, and subdomains all resolve through Arcline
- Any DNS records you had at your registrar will need to be recreated in cPanel's Zone Editor

If you want to keep DNS at your registrar (or use Cloudflare), see [Using external DNS](#using-external-dns) below.

---

## Arcline nameservers

Use these two nameservers for all domains hosted on Arcline:

```
ns1.arcline.it
ns2.arcline.it
```

Your welcome email also lists these. If you have a reseller or custom nameserver setup, contact support for the correct values.

---

## Update nameservers at GoDaddy

1. Log in at [godaddy.com](https://godaddy.com) and go to **My Products → Domains**.
2. Click the domain you want to update.
3. Scroll to **Nameservers** and click **Change**.
4. Select **I'll use my own nameservers**.
5. Remove the existing nameservers and enter:
   - `ns1.arcline.it`
   - `ns2.arcline.it`
6. Click **Save**.

GoDaddy shows a warning that custom nameservers disable their DNS. That's expected — Arcline will handle DNS instead.

---

## Update nameservers at Namecheap

1. Log in at [namecheap.com](https://namecheap.com) and go to **Domain List**.
2. Click **Manage** next to your domain.
3. Under **Nameservers**, change the dropdown from **Namecheap BasicDNS** to **Custom DNS**.
4. Enter:
   - `ns1.arcline.it`
   - `ns2.arcline.it`
5. Click the green checkmark to save.

---

## Update nameservers at Google Domains / Squarespace Domains

1. Go to [domains.google.com](https://domains.google.com) (or your Squarespace Domains dashboard if it's been migrated).
2. Select your domain and go to **DNS**.
3. At the top, switch from **Google Domains DNS** to **Custom name servers**.
4. Delete the existing entries and add:
   - `ns1.arcline.it`
   - `ns2.arcline.it`
5. Click **Save**.

---

## Update nameservers at Cloudflare (as registrar)

If you transferred your domain to Cloudflare as the registrar (not just using Cloudflare for DNS):

1. Log in to the Cloudflare dashboard and select your account.
2. Go to **Domain Registration → Manage Domains**.
3. Click **Manage** next to your domain, then **Configuration → Name servers**.
4. Select **Use custom nameservers** and enter:
   - `ns1.arcline.it`
   - `ns2.arcline.it`
5. Save changes.

> Note: if you're using Cloudflare for DNS only (not as registrar), see [Using external DNS](#using-external-dns) — you'll keep Cloudflare nameservers and point the A record to Arcline instead.

---

## How long does it take?

DNS propagation typically takes **15 minutes to a few hours**. In rare cases it can take up to 48 hours, depending on your registrar and the TTL on the old records.

During propagation, some visitors may still see the old hosting while others see the new. This is normal and temporary.

You can check propagation status at [dnschecker.org](https://dnschecker.org) — enter your domain and select the A record type. Green checkmarks mean that location has the updated DNS.

---

## After switching nameservers

Once propagation is complete:

1. **Verify your site loads** — visit your domain in a browser. If it doesn't load immediately, wait a bit longer.
2. **Set up email** — if you were using the registrar's email, you'll need to recreate email accounts in cPanel and add MX records (see [Set Up Email](/getting-started/email-setup/)).
3. **Check SSL** — cPanel's AutoSSL should issue a certificate within a few minutes. If it doesn't, check **cPanel → Security → SSL/TLS Status** (see [SSL Certificate](/getting-started/ssl-certificate/)).

---

## Using external DNS

If you want to keep DNS management at your registrar or use Cloudflare, don't change your nameservers. Instead, update the A record at your current DNS provider to point to your Arcline server IP:

| Type | Name | Value | TTL |
|---|---|---|---|
| A | `@` | your Arcline server IP | 300 |
| A | `www` | your Arcline server IP | 300 |

Find your server IP in cPanel → **Server Information**, or in your Arcline welcome email.

You'll also need to add MX records and any other DNS records manually. The downside: if you move servers, you have to update your external DNS yourself rather than Arcline handling it automatically.

---

## Subdomain setup

After switching to Arcline nameservers, you can manage all DNS in cPanel:

- **cPanel → Domains → Zone Editor** — add, edit, or remove any DNS record
- **cPanel → Domains → Subdomains** — create subdomains that map to subdirectories of your hosting

Common records to add:
- `www` → your root domain (CNAME or A record, usually already set)
- `mail` → A record pointing to your server IP (needed for email client connections)
- Custom subdomains like `app`, `blog`, `shop`
