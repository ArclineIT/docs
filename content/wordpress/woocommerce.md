---
title: "Setting Up WooCommerce on Arcline"
description: "Install and configure WooCommerce on your Arcline VPS or shared hosting plan."
section: wordpress
order: 5
---

# Setting Up WooCommerce on Arcline

WooCommerce is the most popular e-commerce platform for WordPress. This guide covers installing it on Arcline shared hosting or a VPS and optimizing it for performance without a CDN.

---

## Prerequisites

- WordPress installed and running (see [Install WordPress](/wordpress/install-shared/) for shared hosting or [Install WordPress on a VPS](/wordpress/install-vps/) for VPS)
- A domain pointed to your Arcline server
- SSL certificate installed (Let's Encrypt via cPanel or Certbot)
- PHP memory limit of at least 256 MB (512 MB recommended for VPS)

---

## Step 1 — Install WooCommerce

### Via the WordPress admin

1. Log in to your WordPress admin dashboard at `https://yourdomain.com/wp-admin`
2. Go to **Plugins → Add New**
3. Search for "WooCommerce"
4. Click **Install Now** → **Activate**

WooCommerce will launch the setup wizard on activation.

### Via WP-CLI (faster, especially on a VPS)

```bash
# Install and activate WooCommerce
wp plugin install woocommerce --activate

# Install a recommended theme (e.g., Storefront)
# wp theme install storefront --activate
```

---

## Step 2 — Run the setup wizard

The WooCommerce setup wizard will walk you through:

1. **Store location**: Set your business address, currency (USD), and selling location
2. **Industry**: What types of products you sell (physical, digital, membership, etc.)
3. **Product types**: Simple products, variable products, external/affiliate products
4. **Business details**: Whether you're already selling elsewhere
5. **Payment methods**: WooPayments, PayPal, Stripe, or offline payments
6. **Shipping**: Set up shipping zones and rates
7. **Tax**: Configure basic tax settings
8. **Personalize**: Choose a theme and install free extensions

You can skip any step and configure it later.

---

## Step 3 — Choose payment methods

### WooPayments (built-in, recommended)

WooPayments is Stripe-based and included with WooCommerce. It accepts credit cards, Apple Pay, and Google Pay. No additional plugin is needed.

To set it up:
1. Go to **WooCommerce → Settings → Payments**
2. Click **Set up** next to WooPayments
3. Follow the prompts to connect your Stripe account

### PayPal

1. Go to **WooCommerce → Settings → Payments**
2. Toggle **PayPal** on
3. Click **Set up** and enter your PayPal email address

### Stripe (standalone plugin)

If you prefer the standalone Stripe plugin instead of WooPayments:

```bash
wp plugin install woocommerce-gateway-stripe --activate
```

Then configure it in **WooCommerce → Settings → Payments → Stripe**.

---

## Step 4 — Configure shipping

1. Go to **WooCommerce → Settings → Shipping**
2. Click **Add shipping zone**
3. Name the zone (e.g., "United States")
4. Select the zone regions (countries, states, or postcode ranges)
5. Click **Add shipping method**:
   - **Flat rate** — Single rate per order (e.g., $5.99)
   - **Free shipping** — Free shipping with a minimum order amount
   - **Local pickup** — Customer picks up at your location

---

## Step 5 — Performance optimization for WooCommerce

WooCommerce adds database queries and page weight. Optimize it carefully:

### Essential caching

1. Install a caching plugin. See [W3 Total Cache Configuration](/wordpress/w3-total-cache/) for detailed setup.
2. **Important**: In W3 Total Cache, do NOT enable page caching for the cart, checkout, and my-account pages. Add these to the **Never cache the following pages** list: `/cart/*`, `/checkout/*`, `/my-account/*`, `/wc-api/*`

### Disable unused features

Go to **WooCommerce → Settings → Advanced → Features** and disable:
- Coupons (if you don't use them)
- If digital only, disable shipping and tax

### Enable native cart fragments (if needed)

Cart fragments (the little cart icon that updates via AJAX) are expensive on shared hosting. If your theme doesn't need real-time cart updates, disable it:

```bash
wp option set woocommerce_cart_fragments_enabled no
```

Customers will see a non-JS fallback link to the cart page.

### Image optimization

- Product images should be no larger than 1200px on the longest side
- Use WebP format for product images (convert with `cwebp`)
- Install a plugin like **Smush** or **Imagify** for automatic compression
- Set WooCommerce image sizes appropriately: **WooCommerce → Settings → Products → Display**

---

## Step 6 — Essential WooCommerce plugins

| Plugin | Purpose |
|--------|---------|
| [WooCommerce](https://wordpress.org/plugins/woocommerce/) | Core e-commerce platform |
| [Akismet Anti-Spam](https://wordpress.org/plugins/akismet/) | Block spam product reviews |
| [W3 Total Cache](https://wordpress.org/plugins/w3-total-cache/) | Page cache, DB cache, object cache |
| [UpdraftPlus](https://wordpress.org/plugins/updraftplus/) | Scheduled backups of your store |
| [WooCommerce Stripe](https://wordpress.org/plugins/woocommerce-gateway-stripe/) | Stripe credit card payments |

---

## Step 7 — SSL and security

- **WooCommerce forces SSL on checkout** automatically. Ensure your SSL certificate is valid.
- **Force HTTPS for the entire site** in W3 Total Cache → General Settings → Page Cache → Enable HTTP(S) support
- **Disable file editing** in WordPress admin by adding to `wp-config.php`:
  ```php
  define('DISALLOW_FILE_EDIT', true);
  ```
- **Set up fail2ban**: See [Set Up Fail2ban](/vps/fail2ban/) for SSH brute-force protection
- **Regular backups**: Use UpdraftPlus or [Automated Backups with Restic](/vps/automated-backups/) to back up your database and uploads

---

## Troubleshooting

**Cart and checkout pages not working with caching**: Make sure cart, checkout, and my-account URLs are excluded from the page cache.

**500 error after installing WooCommerce**: Increase PHP memory limit to 512 MB. WooCommerce is memory-intensive.

**SSL not working on checkout**: Go to **WooCommerce → Settings → Advanced → Pages** and verify all pages are set. Then **Settings → General → WordPress Address URL and Site Address URL** must start with `https://`.

**PayPal IPN not working**: In PayPal, set the IPN URL to `https://yourdomain.com/wc-api/ipn-handler/`.

---

## What's next

- [WordPress security hardening](/wordpress/security/)
- [W3 Total Cache configuration](/wordpress/w3-total-cache/)
- [Self-hosting without a CDN: performance tips](/privacy/self-hosting-performance/)

