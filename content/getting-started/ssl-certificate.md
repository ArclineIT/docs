---
title: "SSL Certificate and HTTPS"
description: "How to enable SSL on your Arcline domain, force HTTPS, fix mixed content warnings, and troubleshoot certificate problems."
section: getting-started
order: 6
---

# SSL Certificate and HTTPS

Every Arcline hosting account includes free SSL certificates via AutoSSL (powered by Let's Encrypt). Certificates are issued automatically for your domain and renewed every 90 days without any action required on your part.

---

## Check your SSL status

In cPanel, go to **Security → SSL/TLS Status**. You'll see a list of all domains on your account and their certificate status:

- **Certificate valid** — green, working
- **Certificate expiring soon** — will auto-renew before it expires
- **Failed** — usually a DNS problem; see [Troubleshooting](#troubleshooting) below

If AutoSSL hasn't run yet after you pointed your domain to Arcline, click **Run AutoSSL** at the top of the page to trigger it immediately.

---

## Force HTTPS

Having an SSL certificate doesn't automatically redirect `http://` to `https://`. You need to add a redirect so all visitors get the secure version.

### WordPress

Install the **Really Simple SSL** plugin — it handles the redirect and fixes most mixed content issues in one click. It's free and widely trusted.

Or add this to the top of your `.htaccess` file (before the `# BEGIN WordPress` block):

```apache
RewriteEngine On
RewriteCond %{HTTPS} off
RewriteRule ^(.*)$ https://%{HTTP_HOST}%{REQUEST_URI} [R=301,L]
```

### Non-WordPress sites

Add the same redirect to `.htaccess` in your document root (`/home/username/public_html/.htaccess`):

```apache
RewriteEngine On
RewriteCond %{HTTPS} off
RewriteRule ^(.*)$ https://%{HTTP_HOST}%{REQUEST_URI} [R=301,L]
```

### Via cPanel

cPanel has a built-in redirect tool: go to **Domains → Redirects**, set type to **Permanent (301)**, enter `http://yourdomain.com` as source, `https://yourdomain.com` as destination. Repeat for `www`.

Alternatively, cPanel → **Security → SSL/TLS** → **Manage SSL Sites** has a "Force HTTPS Redirect" checkbox for domains with an active certificate.

---

## WordPress mixed content

After switching to HTTPS, WordPress sometimes still loads some resources over HTTP — images, scripts, or stylesheets embedded in post content. This causes browser warnings ("Not Secure" or padlock with warning).

**Quick fix:** Really Simple SSL catches most of these automatically.

**Manual fix for URLs stored in the database:**

Use the **Better Search Replace** plugin to change `http://yourdomain.com` to `https://yourdomain.com` across all tables. Run it in dry-run mode first to see what would change.

Or via WP-CLI over SSH:

```bash
wp search-replace 'http://yourdomain.com' 'https://yourdomain.com' --all-tables
```

After replacing, go to **Settings → General** and make sure both WordPress Address and Site Address are set to `https://`.

---

## Let's Encrypt certificate details

Arcline uses Let's Encrypt certificates:

- **Validity:** 90 days (renewed automatically ~30 days before expiry)
- **Covers:** your domain and `www.yourdomain.com` (and other configured domains/subdomains on the same account)
- **CA:** Let's Encrypt (trusted by all modern browsers)
- **Cost:** free, included with all Arcline hosting plans

There's no manual renewal step. cPanel's AutoSSL daemon runs nightly and renews certificates that are within 30 days of expiry.

---

## Installing a third-party certificate

If you have a purchased certificate (e.g., an Extended Validation or wildcard cert from a commercial CA), you can install it manually:

1. cPanel → **Security → SSL/TLS → Manage SSL Sites**
2. Select your domain from the dropdown
3. Paste the certificate (`.crt`) and private key (`.key`) into the fields
4. Paste the CA bundle / intermediate certificate if required by your CA
5. Click **Install Certificate**

The Let's Encrypt certificate will be replaced. AutoSSL won't overwrite a manually installed certificate that isn't yet expired.

---

## Troubleshooting

### Certificate not issued

AutoSSL requires that your domain resolve to your Arcline server's IP before it can issue a certificate. Common causes of failure:

1. **DNS hasn't propagated** — wait up to 24 hours after pointing nameservers, then run AutoSSL again
2. **Using external DNS with wrong A record** — verify that `yourdomain.com` and `www.yourdomain.com` A records point to your Arcline server IP
3. **CAA record blocking Let's Encrypt** — if you have a CAA record that doesn't include `letsencrypt.org`, Let's Encrypt can't issue. Add: `0 issue "letsencrypt.org"` or remove the CAA record

### ERR_SSL_PROTOCOL_ERROR or site not loading over HTTPS

Your certificate may not be installed yet. Check **SSL/TLS Status** in cPanel and run AutoSSL if needed. Temporarily access the site over `http://` while waiting.

### Certificate valid but browser shows "Not Secure"

Mixed content — some resources are still loading over HTTP. See [WordPress mixed content](#wordpress-mixed-content) above.

### Certificate expired

AutoSSL should prevent this, but if it fails repeatedly (usually due to DNS issues), the certificate can expire. Fix the underlying DNS problem, then run AutoSSL manually in cPanel. If the domain is unreachable from outside, AutoSSL cannot complete the Let's Encrypt challenge.

### HSTS issues after switching back to HTTP

If your site previously sent `Strict-Transport-Security` (HSTS) headers, browsers will refuse to load it over HTTP even if you remove SSL. This is the HSTS preload taking effect. The fix is to restore the certificate — don't remove SSL from a site that had HSTS enabled.
