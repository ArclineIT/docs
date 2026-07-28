---
title: "Migrate from GoDaddy"
description: "Step-by-step guide to moving your website and email from GoDaddy shared hosting to Arcline without downtime."
section: migrate
order: 1
---

# Migrate from GoDaddy

Moving from GoDaddy to Arcline involves three parts: transferring your files, moving your database, and finally switching DNS. Done in that order, you can complete the whole migration with zero downtime.

---

## Before you start

You'll need:
- GoDaddy cPanel login (or FTP credentials if you're on GoDaddy's non-cPanel shared hosting)
- Access to your Arcline cPanel account
- Your domain's registrar login (wherever you registered the domain — could be GoDaddy itself)
- About 30–60 minutes depending on site size

---

## Step 1 — Export your files from GoDaddy

### If GoDaddy uses cPanel

Log into GoDaddy's cPanel and create a full backup:

1. cPanel → **Files → Backup** (or **Backup Wizard**)
2. Choose **Full Account Backup** and download to your computer

The backup includes your `public_html` files and all databases.

### If GoDaddy uses their custom control panel (Websites + Marketing / Managed WordPress)

GoDaddy's Managed WordPress product doesn't give you direct file or database access. Options:

- Use the **Jetpack** or **UpdraftPlus** plugin to export a full backup (files + database) that you can import on Arcline
- Use the WordPress **Export** tool (Tools → Export) for content only, then reinstall plugins/themes manually
- Contact GoDaddy support and request a backup — they're required by some hosting terms to provide one on request

For standard GoDaddy shared hosting (cPanel-based), direct export is straightforward.

---

## Step 2 — Export your database from GoDaddy

If you used the full cPanel backup in Step 1, the database is already included. Otherwise:

1. GoDaddy cPanel → **Databases → phpMyAdmin**
2. Click your database name in the left sidebar
3. Click **Export → Quick → SQL → Go**

Save the `.sql` file.

For large databases: use **Custom** export, enable **Add DROP TABLE**, and export table-by-table if phpMyAdmin times out.

---

## Step 3 — Set up your site on Arcline

**Create the database:**

1. Arcline cPanel → **Databases → MySQL Databases**
2. Create a new database and a new user
3. Add the user to the database with **All Privileges**
4. Note the database name, username, and password

**Import the database:**

1. cPanel → phpMyAdmin → select the new database → **Import**
2. Upload the `.sql` file from GoDaddy
3. Click **Go**

For large databases, import via SSH instead (see [Back Up and Restore a MySQL Database](/getting-started/mysql-backup/)).

**Upload your files:**

Connect to your Arcline server via SFTP (see [Upload Files via SFTP](/getting-started/sftp/)) and upload your site files to `/home/username/public_html/`.

If you downloaded a full cPanel backup, extract it first and locate the `public_html` folder inside.

---

## Step 4 — Update your database connection settings

WordPress stores the database host, name, username, and password in `wp-config.php`. Update it to point to the new database you created in Step 3:

```php
define( 'DB_NAME',     'cpanelusername_newdbname' );
define( 'DB_USER',     'cpanelusername_newdbuser' );
define( 'DB_PASSWORD', 'your-new-password' );
define( 'DB_HOST',     'localhost' );
```

On GoDaddy, `DB_HOST` is sometimes a custom hostname. On Arcline it's always `localhost`.

For non-WordPress sites, update whatever config file stores your database credentials (`.env`, `config.php`, `database.yml`, etc.).

---

## Step 5 — Test before going live

Before switching DNS, preview your site on Arcline by editing your local `hosts` file.

Find your Arcline server IP in cPanel → **Server Information**.

**On macOS / Linux**, edit `/etc/hosts` (requires sudo):

```
123.456.789.0  yourdomain.com www.yourdomain.com
```

**On Windows**, edit `C:\Windows\System32\drivers\etc\hosts` as Administrator with the same line.

Now visit `https://yourdomain.com` in your browser — you'll see the Arcline version of your site even though DNS still points to GoDaddy. Check:
- Pages load correctly
- Images appear
- Forms work
- Login works (WordPress admin, etc.)
- SSL certificate is valid (if AutoSSL has run)

Remove the `hosts` entry when done testing.

---

## Step 6 — Switch DNS

When you're happy the site works on Arcline:

**If your domain is registered at GoDaddy:**

1. GoDaddy → My Products → Domains → click your domain
2. Scroll to **Nameservers → Change → I'll use my own nameservers**
3. Enter `ns1.arcline.it` and `ns2.arcline.it`
4. Save

**If your domain is registered elsewhere:**

Follow the nameserver update instructions for your registrar. See [Point Your Domain to Arcline](/getting-started/nameservers/) for provider-specific steps.

DNS propagation takes 15 minutes to a few hours. During this time, some visitors may still hit GoDaddy. Keep the GoDaddy site live until propagation is complete (at least 24 hours).

---

## Step 7 — Email migration

If you were using GoDaddy email (Workspace Email or Microsoft 365 through GoDaddy), you'll need to decide what to do with email:

**Switching to Arcline email:**
1. Create email accounts in Arcline cPanel → **Email → Email Accounts**
2. After nameservers switch, Arcline handles MX automatically
3. Export any email you want to keep from GoDaddy webmail before the cutover (GoDaddy Workspace Email → Settings → Export)

**Keeping GoDaddy email (or Microsoft 365):**
After switching nameservers to Arcline, add the original MX records back in cPanel → **Domains → Zone Editor**. GoDaddy's Workspace Email uses:

| Type | Priority | Value |
|---|---|---|
| MX | 0 | `mailstore1.secureserver.net` |
| MX | 10 | `smtp.secureserver.net` |

For Microsoft 365 through GoDaddy, the MX record is specific to your tenant — find it in your Microsoft 365 admin portal.

---

## Common GoDaddy-specific issues

**Site was on a subdomain (e.g., GoDaddy's temp URL)** — GoDaddy gives sites a staging URL while DNS isn't pointed. Make sure `siteurl` and `home` in WordPress → Settings → General (or `wp-config.php`) are set to your real domain, not GoDaddy's temp URL.

**GoDaddy's PHP version differs** — GoDaddy defaults vary. Check your PHP version in GoDaddy cPanel and match it in Arcline cPanel → **Software → Select PHP Version** to avoid compatibility issues.

**GoDaddy file permissions** — GoDaddy sometimes uses suEXEC with specific permission requirements. On Arcline, standard permissions (644 for files, 755 for directories) are correct.

**Email forwarding was set up at GoDaddy** — if you had forwarders, recreate them in Arcline cPanel → **Email → Forwarders** after switching nameservers.
