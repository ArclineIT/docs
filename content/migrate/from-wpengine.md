---
title: "Migrate from WP Engine or Kinsta"
description: "How to move a WordPress site from managed hosting (WP Engine, Kinsta) to Arcline's VPS or shared hosting without downtime."
section: migrate
order: 6
---

# Migrate from WP Engine or Kinsta

WP Engine and Kinsta are managed WordPress hosts. They don't provide cPanel or direct server access — instead you get a custom dashboard with WordPress-specific tools. Migrating to Arcline involves exporting your site through WordPress itself and setting it up fresh.

The process is similar for both WP Engine and Kinsta. Differences are called out where they matter.

---

## Before you start

- WP Engine User Portal or Kinsta MyKinsta dashboard access
- WordPress admin access (for plugin-based backups)
- Arcline cPanel account (for shared hosting) or VPS access
- 45–90 minutes depending on your site size
- **Important:** Managed hosts sometimes disallow certain plugins. Check the [WP Engine disallowed plugins list](https://wpengine.com/support/disallowed-plugins/) or [Kinsta banned plugins list](https://kinsta.com/knowledgebase/banned-plugins/) — if your backup plugin is on the list, use the manual export method instead.

---

## Step 1 — Export your WordPress site

You have three options for exporting. Choose one.

### Option A: UpdraftPlus (easiest)

1. Install the **UpdraftPlus** plugin from your WordPress admin → Plugins → Add New
2. Go to **Settings → UpdraftPlus Backups**
3. Click **Backup Now** and select both **Database** and **Files**
4. When the backup completes, download all five files from the **Existing Backups** tab:
   - Database
   - Plugins
   - Themes
   - Uploads
   - Others

### Option B: All-in-One WP Migration

1. Install **All-in-One WP Migration** from Plugins → Add New
2. Go to **All-in-One WP Migration → Export**
3. Choose **Export To → File**
4. Download the `.wpm` file when it's ready

> The free version has a 512 MB upload limit for importing. If your site is larger, use UpdraftPlus or manual export instead.

### Option C: Manual export

**Export files from WP Engine:**

WP Engine provides SFTP access — you can connect with FileZilla or Cyberduck:

1. WP Engine User Portal → **Sites → your site → SFTP**
2. Note the host, port (2222), username, and password shown
3. Connect with FileZilla and download the entire `wp-content` folder

**Export files from Kinsta:**

Kinsta provides SFTP and SSH access:

1. MyKinsta → **Sites → your site → Info**
2. Under **SFTP/SSH**, note the host, port, username, and password
3. Connect with FileZilla and download the `wp-content` folder

**Export the database:**

Both hosts include phpMyAdmin access:

1. **WP Engine:** User Portal → Sites → your site → **phpMyAdmin** (under Utilities)
2. **Kinsta:** MyKinsta → Sites → your site → **Info** → Open phpMyAdmin

Select your database in the left sidebar → **Export → Quick → SQL → Go**.

---

## Step 2 — Set up WordPress on Arcline

### On shared hosting (cPanel)

1. Arcline cPanel → **Software → WordPress Manager by Softaculous**, or install manually
2. Create a fresh WordPress installation in the root `public_html` directory (or a subdirectory if you prefer to test there first)

### Manual WordPress install

1. Download WordPress from [wordpress.org](https://wordpress.org/download/)
2. Upload the files to `/home/username/public_html/` via SFTP
3. Create a database in cPanel → **Databases → MySQL Databases**
4. Create a database user and add it to the database with **All Privileges**
5. Visit `https://yourdomain.com` and complete the WordPress installation wizard
6. Use the database name, username, and password from step 3

### On a VPS

Follow the VPS setup guides for installing WordPress on a LAMP or LEMP stack.

---

## Step 3 — Import your site

### If you used UpdraftPlus

1. Install and activate UpdraftPlus on your fresh Arcline WordPress installation
2. Go to **Settings → UpdraftPlus Backups → Existing Backups → Upload backup files**
3. Upload all five files you downloaded in Step 1
4. After upload, click **Restore** and select all components
5. UpdraftPlus will restore your database and files
6. Log out and log back in — you'll use your **original** WordPress username and password now

### If you used All-in-One WP Migration

1. Install **All-in-One WP Migration** on the fresh Arcline WordPress installation
2. Go to **All-in-One WP Migration → Import**
3. Upload the `.wpm` file
4. Confirm the overwrite warning
5. Log out and log back in with your original credentials

### If you exported manually

1. Delete the `wp-content` folder on your new Arcline installation (via SFTP or file manager)
2. Upload the `wp-content` folder from your managed host backup
3. In phpMyAdmin on Arcline, select the new WordPress database → **Import** → upload the `.sql` file
4. Update `wp-config.php` to match the database credentials you created in Arcline cPanel:

```php
define( 'DB_NAME',     'arclineuser_wpdb' );
define( 'DB_USER',     'arclineuser_wpuser' );
define( 'DB_PASSWORD', 'your-password' );
define( 'DB_HOST',     'localhost' );
```

On WP Engine and Kinsta, `DB_HOST` is **not** `localhost` — it's a remote database server address. You must change this to `localhost` for Arcline.

---

## Step 4 — Update URLs in the database

Managed hosts often use a temporary or staging domain before your real domain points to them. If your WordPress database still references an old URL, update it.

After importing, log in to your WordPress admin and install **Better Search Replace**. Search for the old URL (e.g., `yourdomain.wpengine.com` or `staging-yourdomain.kinsta.cloud`) and replace it with `yourdomain.com` across all tables.

Alternatively, use WP-CLI over SSH:

```bash
wp search-replace 'yourdomain.wpengine.com' 'yourdomain.com' --all-tables
wp search-replace 'staging-yourdomain.kinsta.cloud' 'yourdomain.com' --all-tables
```

---

## Step 5 — Disable managed-host-specific plugins

WP Engine and Kinsta install plugins that connect to their infrastructure. These won't work on Arcline and should be removed:

**WP Engine:**
- **WP Engine Smart Plugin Manager** — manages updates through WP Engine's infrastructure
- **WP Engine Automated Migration** — only connects to WP Engine

**Kinsta:**
- **Kinsta MU Plugin** (must-use plugin) — handles Kinsta's caching and CDN integration. Delete it from `/wp-content/mu-plugins/`
- **Kinsta Cache Plugin** — only works on Kinsta's infrastructure

Remove these from WordPress admin → Plugins. For must-use plugins, delete them via SFTP from `wp-content/mu-plugins/`.

---

## Step 6 — Set up caching

Managed hosts include server-level caching. On Arcline, you'll need to set this up yourself:

1. Install **W3 Total Cache** or **WP Super Cache** (both free)
2. Configure page caching and browser caching — default settings are fine for most sites
3. For VPS users: also consider PHP-FPM opcache and Nginx fastcgi_cache

---

## Step 7 — Test before switching DNS

Before pointing your domain to Arcline, test using your `hosts` file.

Find your server IP in Arcline cPanel → **Server Information** (or use your VPS IP).

Add this line to your hosts file:

```
YOUR_ARCLINE_IP  yourdomain.com www.yourdomain.com
```

Now visit `https://yourdomain.com` in your browser. Verify:

- The site loads correctly
- All pages, posts, and media display
- WordPress admin works (use your original credentials)
- Plugins and themes function
- Any e-commerce checkout, forms, and membership features work

Remove the `hosts` entry after testing.

---

## Step 8 — Switch DNS/nameservers

**If your domain is managed at the managed host:**

- **WP Engine:** WP Engine doesn't do domain registration. Your domain is likely at a separate registrar (GoDaddy, Namecheap, etc.). Update nameservers there.
- **Kinsta:** Kinsta doesn't register domains. Manage DNS at your registrar.

**At your registrar**, change nameservers to:

```
ns1.arcline.it
ns2.arcline.it
```

See [Point Your Domain to Arcline](/getting-started/nameservers/) for registrar-specific instructions.

**If you use Cloudflare:**

Keep Cloudflare nameservers and update the A record to point to your Arcline IP. Remove any WP Engine or Kinsta-specific DNS records.

DNS propagation takes 15 minutes to a few hours. Keep your managed host active during the transition.

---

## Step 9 — Email

Managed WordPress hosts typically **don't** include email hosting. If you had email set up elsewhere (Google Workspace, Microsoft 365, your registrar), it won't be affected by this migration.

If you want to set up email on Arcline after migration, see [Set Up Email on Your Domain](/getting-started/email-setup/). You'll create accounts in Arcline cPanel and, if your DNS is managed externally, add the required MX records.

---

## WP Engine-specific notes

**Git push deployment** — WP Engine supports deploying via Git. This doesn't transfer; you'll use SFTP or standard WordPress workflows on Arcline.

**WP Engine CDN** — WP Engine includes a CDN. On Arcline, you can optionally set up Cloudflare (free tier) if you want CDN coverage.

**WP Engine's redirect rules** — if you set up redirects through WP Engine's rules engine, recreate them in your `.htaccess` file on Arcline.

**Staging → Production** — WP Engine's three-environment setup (dev/staging/production) is managed by their platform. On Arcline, set up a staging subdomain manually if you need a test environment.

---

## Kinsta-specific notes

**Kinsta CDN** — Kinsta includes Cloudflare-powered CDN. Similar to WP Engine, you can optionally set up Cloudflare independently on Arcline.

**Kinsta's Redis caching** — Kinsta uses Redis for object caching. On Arcline, install a Redis object cache plugin (like **Redis Object Cache**) if your plan or VPS includes Redis. For shared hosting without Redis, the standard page caching plugins are sufficient.

**Kinsta's New Relic monitoring** — Kinsta includes performance monitoring. On Arcline, use Query Monitor (free WordPress plugin) for basic performance debugging.

**PHP 8.x** — Kinsta pushes PHP 8.x aggressively. Arcline supports PHP 8.1 and 8.2 on all shared plans and VPS instances. Your site should run without PHP version issues.

**IonCube and other loaders** — Kinsta doesn't support IonCube or SourceGuardian. Neither does standard Arcline shared hosting. If your site requires these encoders, you'll need a VPS where you can install them.

---

## After migration checklist

- [ ] Site loads at `https://yourdomain.com`
- [ ] WordPress admin works with original credentials
- [ ] All plugins and themes function
- [ ] Managed-host plugins removed (WP Engine MU plugin, Kinsta cache plugin)
- [ ] Caching plugin installed and configured
- [ ] SSL certificate active (AutoSSL or Let's Encrypt)
- [ ] Email accounts set up (if using Arcline email)
- [ ] DNS propagated (check at dnschecker.org)
- [ ] Managed host account kept active for 48 hours as fallback
- [ ] Managed host billing cancelled only after confirming everything works

