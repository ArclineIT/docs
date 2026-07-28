---
title: "WordPress Security Hardening"
description: "Lock down your WordPress site on Arcline — file permissions, login protection, XML-RPC hardening, and no third-party CDN required."
section: wordpress
order: 4
---

# WordPress Security Hardening

Most WordPress compromises happen through outdated plugins or weak passwords — not through server vulnerabilities. These steps harden a standard WordPress install on Arcline against the most common attacks.

---

## Keep everything updated

The single most effective security measure is updating WordPress core, plugins, and themes promptly. Every Arcline cPanel account includes:

- **WordPress Toolkit** (cPanel → Software) — shows update status for all your WordPress sites at a glance. Click **Update** to apply security patches across all sites in one go.
- **Softaculous** — can auto-update WordPress core. Go to Softaculous → **WordPress Manager → Settings** and enable auto-updates.

In WordPress admin, enable auto-updates for plugins and themes you trust. Go to **Plugins → Installed Plugins** and click **Enable auto-updates** next to each plugin.

---

## File permissions

WordPress files should be readable by the web server but not writable by anyone other than your cPanel user. Incorrect permissions are the most common way an attacker who gains access through a plugin vulnerability escalates to full site takeover.

Via cPanel **File Manager** or SFTP:

- **Directories:** `755` (rwxr-xr-x)
- **Files:** `644` (rw-r--r--)
- **wp-config.php:** `640` or `600` — the most sensitive file in your install
- **wp-content/uploads/:** `755` — must be writable for media uploads

To fix permissions via SSH:

```bash
find /home/username/public_html -type d -exec chmod 755 {} \;
find /home/username/public_html -type f -exec chmod 644 {} \;
chmod 640 /home/username/public_html/wp-config.php
```

Run these as your cPanel user — not as root.

---

## Block XML-RPC

XML-RPC is a legacy API that's almost never needed by modern WordPress sites. It's heavily abused for brute-force attacks and DDoS amplification. Most sites can disable it entirely.

**Via .htaccess** (shared hosting):

```apache
<Files xmlrpc.php>
    Order Deny,Allow
    Deny from all
</Files>
```

**Via Nginx** (VPS) — add to your site config:

```nginx
location = /xmlrpc.php { deny all; }
```

**Plugins that need XML-RPC** (don't disable it if you use these):
- Jetpack (some features)
- The WordPress mobile app
- Trackbacks and pingbacks (disabled anyway on most sites)

If you use the WordPress mobile app, you need XML-RPC. For everyone else, disabling it has no downside.

---

## Protect wp-config.php

`wp-config.php` contains your database credentials. Anyone who reads this file owns your database.

**.htaccess** protection (add at the top of `.htaccess`):

```apache
<Files wp-config.php>
    Order Deny,Allow
    Deny from all
</Files>
```

For defense in depth, move `wp-config.php` one directory **above** `public_html` — WordPress looks there automatically. If it's currently at `/home/username/public_html/wp-config.php`, move it to `/home/username/wp-config.php`. WordPress will find it.

---

## Disable file editing from the admin panel

By default, any WordPress administrator can edit theme and plugin files directly from the admin panel. If an attacker compromises an admin account, this lets them inject arbitrary PHP code and take over the entire server.

Add to `wp-config.php`:

```php
define( 'DISALLOW_FILE_EDIT', true );
```

This removes the **Appearance → Theme File Editor** and **Plugins → Plugin File Editor** menu items for everyone. You'll make file changes via SFTP or cPanel File Manager instead.

---

## Disable plugin and theme installation from the admin panel (advanced)

On a production site where you manage installations through SFTP, you can completely disable the ability to install plugins and themes from the admin panel:

```php
define( 'DISALLOW_FILE_MODS', true );
```

This blocks plugin/theme installs, updates, and deletions from the WordPress admin. Updates must be done via WP-CLI, WordPress Toolkit, or manually via SFTP. This is aggressive but very effective — it's a trade-off between convenience and security.

---

## Limit login attempts

WordPress has no built-in rate limiting on the login page, so attackers can try thousands of passwords without restriction.

Install **Wordfence Security** (free) or **Limit Login Attempts Reloaded**. Both block IPs after a configurable number of failed attempts.

**Wordfence settings:**
- **Wordfence → Firewall → Brute Force Protection**
- Set **Lock out after how many login failures** to `5`
- Set **Lock out after how many forgot password attempts** to `5`
- Set **Amount of time a user is locked out** to `1 hour`

Wordfence also includes a web application firewall (WAF) that blocks common WordPress attacks before they reach your site — enable it from the Wordfence dashboard.

---

## Use strong authentication

**Strong passwords** — use the password generator built into WordPress. A password like `myfavoritecat` is trivial to crack; a random one like `8*kF$2nP!xq` is effectively unbreakable.

**Two-factor authentication (2FA)** — install **Wordfence Login Security** (free, from the same developer as Wordfence Security) or **Two Factor** (official WordPress plugin). Both support TOTP (Google Authenticator, Authy, etc.) and backup codes.

**Change the default admin username** — never use `admin`, `administrator`, `root`, or your domain name as the admin username. If you already have an `admin` user, create a new administrator account with a unique username, log in with it, and delete the old `admin` account.

---

## Hide WordPress version

Every WordPress install outputs its version number by default, making it easy for attackers to target known vulnerabilities. Remove it with your security plugin (Wordfence → All Options → **Hide WordPress version**) or by adding a filter:

```php
remove_action( 'wp_head', 'wp_generator' );
```

---

## Disable directory listing

If someone visits `https://yourdomain.com/wp-content/uploads/` directly, they should see a blank page or redirect — not a list of every file in the directory.

Add to `.htaccess`:

```apache
Options -Indexes
```

Arcline shared hosting has this enabled by default. Verify by visiting `https://yourdomain.com/wp-includes/` in your browser — you should see a 403 Forbidden, not a file list.

---

## Change the database table prefix

The default WordPress table prefix is `wp_`. Changing it to something random makes SQL injection attacks harder — the attacker has to guess your table names.

**For new installs:** change the prefix during installation when WordPress asks for it.

**For existing sites:** use the **Brozzme DB Prefix** plugin or do it manually (requires editing `wp-config.php` and renaming all database tables — not recommended unless you're comfortable with MySQL).

---

## Disable unused user enumeration

By default, visiting `https://yourdomain.com/?author=1` reveals the admin username in the URL or redirect. Attackers use this to collect usernames for brute-force attacks.

Block it with Wordfence (enabled by default) or add to your theme's `functions.php`:

```php
if ( ! is_admin() && isset( $_SERVER['QUERY_STRING'] ) ) {
    if ( preg_match( '/author=([0-9]*)/', $_SERVER['QUERY_STRING'] ) ) {
        wp_redirect( home_url() );
        exit;
    }
}
```

---

## Backup before you harden

Some security changes can break things. Before making any significant changes:

1. Take a full cPanel backup (cPanel → **Files → Backup**)
2. Export your database separately (see [Back Up and Restore a MySQL Database](/getting-started/mysql-backup/))
3. Test changes one at a time so you know which one caused a problem if something breaks

