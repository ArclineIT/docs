---
title: "Migrate from Namecheap"
description: "Move your website and email from Namecheap shared hosting to Arcline — step by step, with zero downtime."
section: migrate
order: 4
---

# Migrate from Namecheap

Namecheap's shared hosting runs cPanel, which makes moving to Arcline straightforward. The process covers exporting your files and databases from Namecheap, setting them up on Arcline, testing, and then switching DNS.

---

## Before you start

- Namecheap cPanel login (reachable from your Namecheap dashboard → Hosting List → Manage → cPanel)
- Arcline cPanel account
- 30–90 minutes depending on site size
- Note your Namecheap PHP version: cPanel → **Software → Select PHP Version**

Namecheap is often both the host **and** the registrar — if your domain is registered there, you'll update nameservers in Namecheap's domain panel (covered in Step 6).

---

## Step 1 — Export your files from Namecheap

### Full cPanel backup

1. Namecheap cPanel → **Files → Backup** (or **Backup Wizard**)
2. Choose **Full Backup** and download to your local machine

This includes `public_html` and all databases in one archive.

### Manual export via SFTP

If you prefer to download files individually:

- Host: your server hostname or IP (find it in Namecheap cPanel → **Server Information**)
- Port: `21098` (Namecheap's default SFTP port) or `21` (FTP)
- Username and password: your cPanel credentials

Download `/home/username/public_html/` to your local machine.

---

## Step 2 — Export your database from Namecheap

1. Namecheap cPanel → **Databases → phpMyAdmin**
2. Click your database name in the left sidebar
3. Click **Export → Quick → SQL → Go**

Save the `.sql` file.

For large databases (over 50 MB), use **Custom** export and enable **Add DROP TABLE**, or export via SSH:

```bash
mysqldump -u cpaneluser_dbuser -p cpaneluser_dbname > namecheap_backup.sql
```

---

## Step 3 — Set up on Arcline

**Create the database:**

1. Arcline cPanel → **Databases → MySQL Databases**
2. Create a new database and a new user
3. Add the user to the database with **All Privileges**
4. Note the database name, username, and password — you'll need them for `wp-config.php`

**Import the database:**

1. cPanel → phpMyAdmin → select the new database → **Import**
2. Upload the `.sql` file from Namecheap
3. Click **Go**

For files over 50 MB, import via the command line instead (see [Back Up and Restore a MySQL Database](/getting-started/mysql-backup/)).

**Upload your files:**

Connect to Arcline via SFTP (see [Upload Files via SFTP](/getting-started/sftp/)) and upload your site to `/home/username/public_html/`.

---

## Step 4 — Update your database connection settings

WordPress stores database credentials in `wp-config.php`. Edit it to point to the new Arcline database:

```php
define( 'DB_NAME',     'arclineuser_newdbname' );
define( 'DB_USER',     'arclineuser_newdbuser' );
define( 'DB_PASSWORD', 'your-new-password' );
define( 'DB_HOST',     'localhost' );
```

On Namecheap, `DB_HOST` is typically `localhost` — same as Arcline. Change the other three values to match your new Arcline database.

For non-WordPress sites, update whatever config file stores your database credentials (`.env`, `config.php`, `database.yml`, etc.).

---

## Step 5 — Match the PHP version

Namecheap's default PHP version may differ from Arcline's. Check and match:

1. Namecheap cPanel → **Software → Select PHP Version** — note the version
2. Arcline cPanel → **Software → Select PHP Version** — set the same version
3. Also check the **Extensions** tab to enable any extensions your site needs (common ones: `mbstring`, `curl`, `gd`, `zip`, `intl`)

---

## Step 6 — Test before going live

Before switching DNS, preview your site on Arcline using your local `hosts` file.

Find your Arcline server IP in cPanel → **Server Information**.

**On macOS / Linux**, edit `/etc/hosts` (requires sudo):

```
YOUR_ARCLINE_IP  yourdomain.com www.yourdomain.com
```

**On Windows**, edit `C:\Windows\System32\drivers\etc\hosts` as Administrator with the same line.

Now visit `https://yourdomain.com` in your browser. Check:

- Pages load correctly
- Images appear
- WordPress admin (/wp-admin) works
- Contact forms submit
- Any dynamic features (shopping cart, search) work

Remove the `hosts` entry when done testing.

---

## Step 7 — Switch nameservers

Since Namecheap is likely both your host and registrar, you'll update nameservers in the Namecheap domain panel (not cPanel).

1. Log in to Namecheap → **Domain List**
2. Click **Manage** next to your domain
3. Under **Nameservers**, change the dropdown from **Namecheap BasicDNS** to **Custom DNS**
4. Enter:
   - `ns1.arcline.it`
   - `ns2.arcline.it`
5. Click the green checkmark to save

DNS propagation takes 15 minutes to a few hours. Keep Namecheap's hosting active during propagation — both servers will serve the site until DNS fully switches.

> If your domain is registered **elsewhere**, update nameservers at that registrar. See [Point Your Domain to Arcline](/getting-started/nameservers/) for steps by provider.

---

## Step 8 — Email migration

**Moving to Arcline email:**

1. Create email accounts in Arcline cPanel → **Email → Email Accounts**
2. After nameservers switch, Arcline's MX records activate automatically
3. Export any mail you want to keep from Namecheap webmail (usually Roundcube) before cancelling

**Keeping Namecheap email (Private Email):**

If you subscribed to Namecheap's Private Email service and want to keep it, add the MX records back in Arcline cPanel → **Domains → Zone Editor** after switching nameservers. Namecheap Private Email uses:

| Type  | Priority | Value                         |
|-------|----------|-------------------------------|
| MX    | 10       | `mx1.privateemail.com`        |
| MX    | 20       | `mx2.privateemail.com`        |

You'll also need the SPF record Namecheap provides in your Private Email settings.

---

## Namecheap-specific notes

**Namecheap's SSL certificates** — Namecheap sells PositiveSSL certificates. These don't transfer to Arcline. Arcline includes free AutoSSL (Let's Encrypt) which will issue automatically once your domain resolves to Arcline. Don't bother exporting Namecheap's SSL cert.

**Supersonic CDN** — Namecheap includes a basic CDN on some plans. Your site doesn't need it on Arcline; our infrastructure already includes caching and performance optimization by default.

**Namecheap's backup service** — Namecheap offers paid automated backups. You don't need this add-on on Arcline — use cPanel's built-in backup tools or set up a cron-based backup (see [Back Up and Restore a MySQL Database](/getting-started/mysql-backup/)).

**EasyWP (Namecheap managed WordPress)** — if you're on Namecheap's EasyWP platform instead of cPanel hosting, you can't export via cPanel. Use a WordPress backup plugin (UpdraftPlus or All-in-One WP Migration) to export your files and database, then import the backup on Arcline.

**.htaccess redirects** — if Namecheap added custom redirects in your `.htaccess` file (e.g., non-www to www), they'll come over with your files and work on Arcline without changes.

---

## After migration checklist

- [ ] Site loads at `https://yourdomain.com`
- [ ] WordPress admin accessible at `/wp-admin`
- [ ] Forms submit and send email
- [ ] SSL certificate valid (check cPanel → **SSL/TLS Status**)
- [ ] Email accounts created and working
- [ ] PHP version and extensions match
- [ ] DNS fully propagated (verify at dnschecker.org)
- [ ] Namecheap hosting account kept active for 48 hours as fallback

