---
title: "Migrate from Bluehost"
description: "How to move your website and email from Bluehost to Arcline, including WordPress sites and cPanel-to-cPanel transfers."
section: migrate
order: 2
---

# Migrate from Bluehost

Bluehost runs cPanel on its shared hosting plans, which makes migrating to Arcline straightforward — both platforms speak the same language. The process is: export from Bluehost, set up on Arcline, test, then switch nameservers.

---

## Before you start

- Bluehost cPanel access (log in at your Bluehost dashboard → **Advanced** to reach cPanel)
- Arcline cPanel account
- 30–90 minutes depending on database and file size
- Note your Bluehost PHP version: cPanel → **Software → Select PHP Version**

---

## Step 1 — Back up files from Bluehost

### Full cPanel backup

1. Bluehost cPanel → **Files → Backup** (or search for "Backup Wizard")
2. Choose **Full Backup** and download to your local machine

This gives you an archive with `public_html` and all databases.

### Manual file backup via SFTP

If you only want specific files, connect to Bluehost via SFTP:

- Host: `ftp.yourdomain.com` or your server hostname
- Port: `21` (FTP) or `22` (SFTP) — use SFTP if available
- Username/password: your cPanel credentials

Download `/home/username/public_html/` to your local machine.

---

## Step 2 — Export your database from Bluehost

1. Bluehost cPanel → **Databases → phpMyAdmin**
2. Click your database name in the left sidebar
3. Click **Export → Quick → SQL → Go**

For large databases, use **Custom** export and enable **Add DROP TABLE**. You can also export via SSH:

```bash
mysqldump -u cpanelusername_dbuser -p cpanelusername_dbname > bluehost_backup.sql
```

---

## Step 3 — Set up on Arcline

**Create the database:**
1. Arcline cPanel → **Databases → MySQL Databases**
2. Create a database and a database user
3. Assign the user to the database with **All Privileges**

**Import the database:**
1. cPanel → phpMyAdmin → select your new database → **Import**
2. Upload the `.sql` file from Bluehost
3. Click **Go**

For files over 50 MB, use the command line (see [Back Up and Restore a MySQL Database](/getting-started/mysql-backup/)).

**Upload your files:**
Connect to Arcline via SFTP and upload your site to `/home/username/public_html/`. See [Upload Files via SFTP](/getting-started/sftp/).

---

## Step 4 — Update wp-config.php

WordPress stores database credentials in `wp-config.php`. Update the values to match the new database you created:

```php
define( 'DB_NAME',     'newcpanel_dbname' );
define( 'DB_USER',     'newcpanel_dbuser' );
define( 'DB_PASSWORD', 'your-password' );
define( 'DB_HOST',     'localhost' );
```

Bluehost sometimes uses a custom `DB_HOST` value. On Arcline it is always `localhost`.

---

## Step 5 — Match the PHP version

Bluehost's default PHP version may differ from Arcline's. To prevent compatibility errors:

1. Check the PHP version in Bluehost cPanel → **Software → Select PHP Version**
2. Set the same version in Arcline cPanel → **Software → Select PHP Version**
3. Also check that any PHP extensions your site requires are enabled on Arcline (the **Extensions** tab in MultiPHP Manager)

WordPress 6.x requires PHP 7.4 minimum; PHP 8.1 or 8.2 is recommended.

---

## Step 6 — Test the site before switching DNS

Edit your local `hosts` file to preview the Arcline version without changing DNS.

Find your Arcline server IP in cPanel → **Server Information**.

Add this line to your hosts file:

```
YOUR_ARCLINE_IP  yourdomain.com www.yourdomain.com
```

**macOS / Linux:** `/etc/hosts` (requires sudo)
**Windows:** `C:\Windows\System32\drivers\etc\hosts` (open as Administrator)

Visit your domain in a browser and verify:
- Pages and images load
- WordPress admin works
- Forms and any dynamic features work
- SSL is active (AutoSSL issues a cert after DNS resolves, which it will once nameservers switch — for pre-launch testing self-signed is normal)

Remove the hosts entry once done.

---

## Step 7 — Switch nameservers

When the site checks out, update nameservers to Arcline. Bluehost is often both the host and the registrar.

**If your domain is registered at Bluehost:**

1. Log in to Bluehost → **Domains → My Domains**
2. Click **Manage** → **Nameservers** → select **Custom Nameservers**
3. Enter:
   - `ns1.arcline.it`
   - `ns2.arcline.it`
4. Save

**If your domain is registered elsewhere**, update nameservers at that registrar instead. See [Point Your Domain to Arcline](/getting-started/nameservers/).

Keep Bluehost live for 24–48 hours while DNS propagates. Don't cancel your Bluehost account until propagation is confirmed and you've verified email is working.

---

## Step 8 — Email

**Moving to Arcline email:**

1. Create accounts in Arcline cPanel → **Email → Email Accounts**
2. Once nameservers switch, Arcline's MX records activate automatically
3. Export any mail you want to keep from Bluehost webmail (Bluehost uses Roundcube or Horde) before cancelling

**Keeping external email (Google Workspace, Microsoft 365):**

After switching nameservers to Arcline, add your existing MX records back in cPanel → **Domains → Zone Editor**. You'll find the required MX records in your Google or Microsoft 365 admin panel.

---

## Bluehost-specific notes

**SiteLock or CodeGuard backups** — these are Bluehost add-ons. You don't need them on Arcline; standard cPanel backups are included.

**WordPress staging site** — Bluehost provides a staging environment. After migration, delete the staging version from Bluehost to avoid confusion during the transition.

**Bluehost's custom wp-config constants** — Bluehost adds some non-standard constants to `wp-config.php` (like `COOKIE_DOMAIN` or Bluehost-specific feature flags). These can be removed safely after migrating.

**"Account Suspension" page appearing** — if Bluehost's CDN or caching layer still serves your domain during propagation, visitors may briefly see a Bluehost holding page. This resolves once propagation is complete.

**File permissions** — Bluehost uses suPHP/suEXEC. Standard permissions (644 files, 755 directories) work on Arcline. If Bluehost set unusual permissions (e.g., 777), correct them on Arcline.

---

## After migration checklist

- [ ] Site loads at `https://yourdomain.com`
- [ ] WordPress admin accessible at `/wp-admin`
- [ ] Forms submit and send email
- [ ] SSL certificate shows valid in browser
- [ ] Email accounts created and working
- [ ] Old Bluehost account not yet cancelled (wait 48 hours)
- [ ] DNS fully propagated (check at dnschecker.org)
