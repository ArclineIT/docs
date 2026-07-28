---
title: "Install WordPress on Shared Hosting"
description: "Install WordPress on Arcline shared hosting with cPanel — from one-click installers to manual setup and WordPress Toolkit."
section: wordpress
order: 1
---

# Install WordPress on Shared Hosting

WordPress can be installed on Arcline shared hosting in a few minutes. This covers the one-click installer (fastest), the manual method (for full control), and WordPress Toolkit (for managing multiple sites).

---

## Option 1 — Softaculous one-click installer (recommended)

Softaculous is included with every Arcline cPanel account. It installs WordPress with one click and lets you choose the install location, admin credentials, and all settings upfront.

1. cPanel → **Software → WordPress Manager by Softaculous**
2. Click **Install** at the top
3. Fill in the form:

| Field | What to enter |
|---|---|
| Choose Protocol | `https://` (recommended) or `https://www.` |
| Choose Domain | your domain from the dropdown |
| In Directory | leave blank (for the root of your domain) or enter a subdirectory like `blog` |
| Site Name | your site's title (can be changed later) |
| Site Description | short tagline (optional) |
| Admin Username | pick a unique username — **not** `admin` |
| Admin Password | use a strong password or click the key icon to generate one |
| Admin Email | your email address |

Scroll down and click **Install**. WordPress is ready in under a minute.

> Do not use `admin` as your username — it's the most targeted name for brute-force attacks. Pick something unique.

---

## Option 2 — Manual install

A manual install gives you full control over file placement, database setup, and the initial configuration.

### Step 1 — Create a database

1. cPanel → **Databases → MySQL Databases**
2. Under **Create New Database**, enter a name (e.g., `wp_yoursite`) and click **Create Database**
3. Under **Add New User**, create a user with a strong password
4. Under **Add User To Database**, select the user and database, then click **Add**
5. Check **All Privileges** and click **Make Changes**

Write down the database name, username, and password.

### Step 2 — Upload WordPress

1. Download the latest WordPress `.zip` from [wordpress.org](https://wordpress.org/download/)
2. cPanel → **Files → File Manager**
3. Navigate to `public_html` (or the subdirectory where you want WordPress)
4. Click **Upload**, select the `.zip`, and wait for it to finish
5. In File Manager, select the uploaded `.zip` and click **Extract**
6. Move the extracted `wordpress/` folder contents into `public_html/` (or keep it in a subdirectory)
7. Delete the `.zip` file and the now-empty `wordpress/` folder

### Step 3 — Run the installer

1. Visit `https://yourdomain.com` in your browser
2. WordPress will detect no config file and show the setup screen
3. Click **Let's go!**, then enter your database details:

| Field | Value |
|---|---|
| Database Name | the database you created in Step 1 |
| Username | the database user you created |
| Password | the database user's password |
| Database Host | `localhost` |
| Table Prefix | `wp_` (default is fine, or change it for extra security) |

4. Click **Submit** → **Run the installation**
5. Fill in the site info: site title, admin username (not `admin`), password, and your email
6. Click **Install WordPress**

---

## Option 3 — WordPress Toolkit (for multiple sites)

WordPress Toolkit is included in cPanel and is ideal if you manage several WordPress sites on one account.

1. cPanel → **Software → WordPress Toolkit**
2. Click **Install**
3. Choose your domain, directory, and basic settings
4. Click **Install**

After installation, WordPress Toolkit lets you:
- Clone a site to a subdomain or another domain
- Create staging sites (copy production to a test area)
- Run WordPress, plugin, and theme updates from one screen
- Reset passwords and toggle debug mode without logging into WordPress

---

## After installation

### Enable SSL

Your site should use HTTPS immediately. cPanel's AutoSSL will issue a certificate within a few minutes of installation. If it doesn't:

1. cPanel → **Security → SSL/TLS Status**
2. Click **Run AutoSSL**
3. Wait a few minutes and refresh

Once the certificate is issued, install **Really Simple SSL** in WordPress to automatically redirect HTTP to HTTPS and fix mixed content.

### Update permalinks

WordPress's default permalink structure (`?p=123`) is bad for SEO and usability. Change it immediately:

1. WordPress admin → **Settings → Permalinks**
2. Choose **Post name** (the most common and SEO-friendly option)
3. Click **Save Changes**

If you see a 404 error after switching permalinks, WordPress couldn't write to `.htaccess`. Copy the code WordPress shows at the bottom of the Permalinks page and paste it manually into `.htaccess` in File Manager.

### Install essential plugins

Start with these free plugins (install after you've confirmed the site works):

- **Really Simple SSL** — handles HTTPS redirect and mixed content
- **Wordfence Security** — firewall, malware scanner, and login protection
- **UpdraftPlus** — automated backups to remote storage (Google Drive, Dropbox, etc.)

### Set up automated backups

Don't rely on manual backups. Configure UpdraftPlus to run daily backups to off-server storage. See [Back Up and Restore a MySQL Database](/getting-started/mysql-backup/) for database backup details.

---

## Installing in a subdirectory

If you want WordPress at `https://yourdomain.com/blog/` rather than the root of your domain:

1. In Softaculous or during manual setup, enter `blog` as the directory
2. WordPress files go into `public_html/blog/`
3. Your homepage at `yourdomain.com` can be a static site or an HTML landing page

This is a common setup when your main site is not WordPress (e.g., a static business site with a separate blog).

---

## Common setup issues

**Error establishing a database connection** — the database credentials in `wp-config.php` don't match what you created in cPanel. Double-check the database name, username, and password. Note that all three are prefixed with your cPanel username (e.g., `cpaneluser_wp_yoursite`).

**White screen after install** — a PHP error. Enable debugging temporarily by adding this to `wp-config.php`:

```php
define( 'WP_DEBUG', true );
define( 'WP_DEBUG_LOG', true );
```

Check `wp-content/debug.log` for the specific error. Remove or set to `false` after debugging.

**Can't upload files** — the `wp-content/uploads/` directory may have incorrect permissions. In File Manager, right-click the `uploads` folder and set permissions to **755**. If that doesn't fix it, try **775**.

**404 on all pages except home** — permalink rules aren't being applied. Go to Settings → Permalinks and click **Save Changes** twice (this forces WordPress to regenerate `.htaccess` rules). If it still doesn't work, check that your `.htaccess` file exists in `public_html/` and is writable.

