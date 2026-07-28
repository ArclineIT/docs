---
title: "Migrate from HostGator"
description: "How to move your website from HostGator shared hosting to Arcline, including WordPress and cPanel-to-cPanel transfers."
section: migrate
order: 5
---

# Migrate from HostGator

HostGator runs cPanel on its shared and reseller plans, so migrating to Arcline is familiar territory. You'll export your files and database from HostGator, set them up on Arcline, test, and then switch nameservers.

---

## Before you start

- HostGator customer portal login (at [portal.hostgator.com](https://portal.hostgator.com)) — you may need this to reach cPanel if single sign-on is enabled
- Arcline cPanel account
- 30–90 minutes depending on site size
- Note your HostGator PHP version: cPanel → **Software → MultiPHP Manager**

---

## Step 1 — Export your files from HostGator

### From cPanel

1. Log in to HostGator cPanel (via customer portal or directly at `https://yourdomain.com:2083`)
2. **Files → Backup** (or **Backup Wizard**)
3. Choose **Full Backup** and download to your local machine

This produces a single archive with your `public_html` directory and all databases.

### Via SFTP

If the cPanel backup is too large to download through the browser, use SFTP:

- Host: your server hostname or IP (find it in HostGator cPanel → **Server Information**)
- Port: `2222` (HostGator's default non-standard SSH/SFTP port)
- Username and password: your cPanel credentials

Download everything in `/home/username/public_html/` to your local machine.

---

## Step 2 — Export your database from HostGator

1. HostGator cPanel → **Databases → phpMyAdmin**
2. Click your database name in the left sidebar
3. **Export → Quick → SQL → Go**

Save the `.sql` file.

For databases over 50 MB, use **Custom** export and check **Add DROP TABLE**. If phpMyAdmin times out, export via SSH:

```bash
mysqldump -u cpaneluser_dbuser -p cpaneluser_dbname > hostgator_backup.sql
```

---

## Step 3 — Set up on Arcline

**Create the database:**

1. Arcline cPanel → **Databases → MySQL Databases**
2. Create a new database and a new user
3. Add the user to the database with **All Privileges**
4. Write down the database name, username, and password

**Import the database:**

1. cPanel → phpMyAdmin → select the new (empty) database → **Import**
2. Choose the `.sql` file from HostGator
3. Click **Go**

For large imports, use the command line (see [Back Up and Restore a MySQL Database](/getting-started/mysql-backup/)).

**Upload your files:**

Connect to Arcline via SFTP and upload everything to `/home/username/public_html/`. See [Upload Files via SFTP](/getting-started/sftp/).

---

## Step 4 — Update wp-config.php

Edit your `wp-config.php` to use the new Arcline database credentials:

```php
define( 'DB_NAME',     'arclineuser_newdbname' );
define( 'DB_USER',     'arclineuser_newdbuser' );
define( 'DB_PASSWORD', 'your-new-password' );
define( 'DB_HOST',     'localhost' );
```

HostGator typically uses `localhost` for `DB_HOST` — same as Arcline. Update the database name, username, and password to the values you created in Step 3.

For non-WordPress sites, update the equivalent database configuration file.

---

## Step 5 — Check PHP version and extensions

1. HostGator cPanel → **Software → MultiPHP Manager** or **Select PHP Version** — note the version
2. Arcline cPanel → **Software → Select PHP Version** — set the same version
3. Check the **Extensions** tab to match any custom extensions your site uses

HostGator supports PHP 5.6 through 8.x. If your site is running an older version (5.6 or 7.0–7.3), consider updating to PHP 8.1+ on Arcline — test thoroughly before switching DNS.

---

## Step 6 — Test before DNS cutover

Add Arcline's server IP to your local `hosts` file to preview the site before changing DNS.

Find your Arcline IP in cPanel → **Server Information**.

Add this line to your hosts file:

```
YOUR_ARCLINE_IP  yourdomain.com www.yourdomain.com
```

- **macOS / Linux:** `/etc/hosts` (requires sudo)
- **Windows:** `C:\Windows\System32\drivers\etc\hosts` (open as Administrator)

Visit `https://yourdomain.com` and verify:

- The homepage and key pages load
- Images and media work
- WordPress admin is accessible
- Contact forms, shopping cart, and other dynamic features work correctly

Remove the `hosts` entry after testing.

---

## Step 7 — Switch nameservers

**If your domain is registered at HostGator:**

1. Log in to HostGator customer portal → **Domains**
2. Select your domain and go to **Nameservers**
3. Choose **Custom Nameservers** and enter:
   - `ns1.arcline.it`
   - `ns2.arcline.it`
4. Save changes

**If your domain is registered elsewhere**, update nameservers at the registrar. See [Point Your Domain to Arcline](/getting-started/nameservers/) for provider-specific steps.

DNS propagation takes 15 minutes to a few hours. Keep HostGator active for at least 24 hours during propagation.

---

## Step 8 — Email migration

**Moving to Arcline email:**

1. Create your email accounts in Arcline cPanel → **Email → Email Accounts**
2. Once nameservers switch, Arcline handles MX records automatically
3. Export any mail you want to keep from HostGator webmail (Roundcube or Horde) before cancelling

**Keeping external email (Google Workspace, Microsoft 365):**

If you use an external email provider and only used HostGator for hosting, add your existing MX records back in Arcline cPanel → **Domains → Zone Editor** after switching nameservers.

---

## HostGator-specific notes

**SiteLock** — HostGator pushes SiteLock as an add-on for malware scanning. You don't need it on Arcline; standard WordPress security practices (keep everything updated, use strong passwords, limit login attempts) are sufficient. If you want extra scanning, use a free alternative like the Wordfence plugin.

**HostGator's custom caching** — HostGator sometimes uses server-level caching (Varnish or similar). Your site won't have this on Arcline, but standard WordPress caching plugins (W3 Total Cache, WP Super Cache) work well and give you more control.

**QuickInstall / Softaculous apps** — if you installed WordPress or other apps through HostGator's auto-installer, the installation itself is standard — it will migrate normally. The auto-installer is just a convenience tool; your site doesn't depend on it.

**HostGator's `public_html` / `www` difference** — HostGator sometimes uses `www` as a symlink to `public_html`. On Arcline, your site files go directly in `public_html`. Just upload everything there.

**Resource usage warnings** — HostGator's shared plans have strict inode and CPU limits. If you received warnings about resource usage, these don't apply on Arcline — our shared plans have different limits. Check your site's resource usage on Arcline's cPanel (left sidebar → **Resource Usage**).

**HostGator website builder** — if you built your site with HostGator's proprietary website builder (Gator Website Builder), it can't be exported as files. You'll need to rebuild the site on Arcline, or use a WordPress-based site instead. This guide only covers cPanel-based sites with standard file/database exports.

---

## After migration checklist

- [ ] Site loads correctly at `https://yourdomain.com`
- [ ] WordPress admin accessible at `/wp-admin`
- [ ] Forms submit and send email
- [ ] SSL certificate valid (AutoSSL in cPanel → **SSL/TLS Status**)
- [ ] PHP version matches what the site expects
- [ ] Email accounts created and working
- [ ] DNS propagated (check at dnschecker.org)
- [ ] HostGator account kept active for 48 hours as fallback
- [ ] HostGator account cancelled only after confirming everything works

