---
title: "Configure W3 Total Cache Without a CDN"
description: "Speed up your WordPress site with W3 Total Cache on Arcline — page caching, browser caching, and opcode caching without third-party CDNs."
section: wordpress
order: 3
---

# Configure W3 Total Cache Without a CDN

W3 Total Cache (W3TC) is a free WordPress caching plugin that speeds up your site by storing pre-rendered pages, compressing assets, and leveraging browser caching. This guide configures W3TC for good performance **without** a third-party CDN — everything runs on your Arcline server.

---

## Why no CDN?

CDNs add a third party between your visitors and your server. You may not need one:

- Arcline servers are fast and colocated in well-connected data centers
- A properly cached WordPress site on Arcline loads in under a second for most visitors
- CDNs introduce an additional cost, privacy concern, and point of failure

If you do want CDN coverage later, W3TC supports Cloudflare, BunnyCDN, and generic pull CDNs — but you don't need one to get started.

---

## Install W3 Total Cache

1. WordPress admin → **Plugins → Add New**
2. Search for "W3 Total Cache"
3. Click **Install Now** → **Activate**

The plugin adds a **Performance** menu to the sidebar. All configuration lives there.

---

## Page cache (most important)

Page caching saves fully rendered HTML pages so WordPress doesn't process PHP and query the database for every request. A cached page is served in milliseconds.

1. **Performance → General Settings**
2. Under **Page Cache**, check **Enable**
3. Set **Page Cache Method** to **Disk: Enhanced**
4. Click **Save all settings**

**Disk: Enhanced** writes static `.html` files to the cache directory and serves them directly via `.htaccess` or Nginx rules — WordPress isn't even loaded for cached pages. This is the fastest option that doesn't require extra server software.

### Verify it's working

Visit your site in an incognito window (so you're not logged in as admin) and view the page source. Scroll to the bottom — you should see:

```html
<!-- Performance optimized by W3 Total Cache. Page Caching using disk: enhanced -->
```

---

## Browser cache

Browser caching tells visitors' browsers to store images, CSS, and JavaScript files locally so they don't re-download on every page view.

1. **Performance → General Settings**
2. Under **Browser Cache**, check **Enable**
3. Click **Save all settings**

Then configure each section:

**Performance → Browser Cache → CSS & JS:**
- Set **Expires header lifetime** to `31536000` seconds (1 year)
- Check **Set cache control header**
- Set **Cache Control policy** to `cache with max-age`

**Performance → Browser Cache → HTML & XML:**
- Set **Expires header lifetime** to `3600` seconds (1 hour)
- Check **Set cache control header**
- Set **Cache Control policy** to `cache with max-age`

**Performance → Browser Cache → Media & Other Files:**
- Set **Expires header lifetime** to `31536000` seconds (1 year)
- **Cache Control policy** to `cache with max-age`

Click **Save all settings** after each tab.

---

## Minify (optional — test carefully)

Minification reduces file sizes by stripping whitespace and comments from HTML, CSS, and JS. It can improve load times but also **can break your site** if not configured correctly.

Start with HTML minification only — it's the safest:

1. **Performance → General Settings → Minify → Enable**
2. Set **Minify mode** to **Manual** (not Auto)
3. Click **Save all settings**

Then in **Performance → Minify → HTML & XML:**
- Check **Enable** for HTML minify
- Leave JS and CSS minify disabled for now
- Click **Save all settings**

Test your site thoroughly. If anything looks wrong, disable minification for that type and try a different combination.

**If your theme CSS or JS files have incorrect paths after minifying**, you may need to add them to the "Never minify" list in the JS or CSS settings tab.

---

## Object cache

Object caching stores database query results in memory, reducing repeated database queries. On Arcline shared hosting, use **Disk** as the caching method:

1. **Performance → General Settings → Object Cache → Enable**
2. Set **Object Cache Method** to **Disk**
3. Click **Save all settings**

On a VPS with Redis installed, set the method to **Redis** instead and enter `127.0.0.1:6379` as the server. Redis object caching is significantly faster than disk-based caching.

---

## Database cache (use on VPS only)

Database caching stores query results. On shared hosting, it can slow things down if the disk is under load. On a VPS, it helps:

1. **Performance → General Settings → Database Cache → Enable**
2. Set method to **Disk** (or **Redis** on a VPS with Redis)
3. Click **Save all settings**

Skip this on shared hosting unless your site has heavy database usage (WooCommerce, membership sites, forums).

---

## Exclude pages from caching

Some pages should never be cached:

- **Cart, checkout, and account pages** (for WooCommerce)
- **Login and registration pages**
- **Admin pages**

**Performance → Page Cache → Advanced:**

In the **"Never cache the following pages"** field, add:

```
wp-login.php
wp-admin/*
cart/*
checkout/*
my-account/*
```

For WooCommerce specifically, W3TC should detect it and add these automatically. If not, add `/cart/`, `/checkout/`, `/my-account/` individually — one per line.

---

## Clear the cache

You'll need to clear the cache whenever you make significant site changes (new theme, updated plugins, content restructuring):

- **Performance → Dashboard → Empty all caches**
- Or use the admin bar: **Performance → Purge All Caches**

Set up automatic purging for new posts:

- **Performance → Page Cache → Purge Policy:**
- Check **Front page**, **Posts page**, and **Post page**

This keeps your cache fresh without manual intervention after publishing new content.

---

## Testing your cache setup

After configuring, test your site's performance:

1. Visit your site in an incognito window
2. Open browser DevTools → Network tab
3. Reload the page and check:
   - **HTML document:** should load in under 200 ms
   - **CSS/JS files:** should show "304 Not Modified" or "(disk cache)" on second load
   - **Images:** similar — cached after the first load

For more detailed testing, use [PageSpeed Insights](https://pagespeed.web.dev) or [GTmetrix](https://gtmetrix.com). Both give specific recommendations for improvement.

---

## Troubleshooting

**Site looks broken after enabling minify** — disable minification for the type that broke (JS, CSS, or HTML). Minify is the most likely setting to cause issues.

**Logged-in users see stale pages** — W3TC should skip caching for logged-in users by default. Check **Performance → Page Cache → Advanced → "Don't cache pages for logged in users"** is checked.

**Cache files filling up disk space** — W3TC has garbage collection that runs on WordPress cron. On a busy site, the cache directory can grow. Set a reasonable **Garbage collection interval** in **Performance → Page Cache → Advanced** (the default 3600 seconds / 1 hour is fine).

**CDN tab references** — ignore everything in the **CDN** settings section. That's for Cloudflare, BunnyCDN, or generic pull CDNs, which this guide intentionally avoids.

**"Disk: Enhanced" not available** — your server may not support the enhanced mode. Switch to **Disk: Basic** instead, which uses PHP to serve cached pages. It's slightly slower but works everywhere.

