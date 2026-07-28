---
title: "Migrate from SiteGround"
description: "How to transfer your website and database from SiteGround to Arcline, including handling SiteGround's custom tools and SuperCacher."
section: migrate
order: 3
---

# Migrate from SiteGround

SiteGround uses cPanel on older accounts and their own Site Tools panel on newer ones. The migration process is similar either way: export your files and database, set up on Arcline, test, then cut over DNS.

---

## Before you start

- SiteGround Site Tools or cPanel access
- Arcline cPanel account
- 30–90 minutes depending on site size
- Note any SiteGround-specific plugins installed (see [SiteGround-specific notes](#siteground-specific-notes))

---

## Step 1 — Export files from SiteGround

### Via Site Tools (newer SiteGround accounts)

SiteGround's Site Tools replaced cPanel on newer accounts. To export your files:

1. Site Tools → **Site → FTP Accounts** — create an FTP user if you don't have one
2. Connect with FileZilla or Cyberduck to `yourdomain.com` on port 22 (SFTP) or 21 (FTP)
3. Download your `public_html` folder

Or use SiteGround's Backup tool:

1. Site Tools → **Security → Backups**
2. Select the most recent backup date
3. Click **Restore** on **File System** — but choose **Download** instead of restore to download the archive

### Via cPanel (older SiteGround accounts)

cPanel → **Files → Backup Wizard → Full Backup → Download**

This creates a full archive including databases.

---

## Step 2 — Export your database

### Site Tools

1. Site Tools → **Site → MySQL** → select your database → click **phpMyAdmin**
2. Click the database name in the left sidebar
3. **Export → Quick → SQL → Go**

Or use the **Backup** tool → select the most recent backup → **Download** on the database section.

### cPanel

cPanel → **Databases → phpMyAdmin** → select database → **Export → Quick → SQL → Go**

For large databases, export via SSH:

```bash
mysqldump -u dbusername -p dbname > siteground_backup.sql
```

---

## Step 3 — Set up on Arcline

**Create the database:**
1. Arcline cPanel → **Databases → MySQL Databases**
2. Create a database, a user, and assign the user with **All Privileges**
3. Note all three: database name, username, password

**Import the database:**
1. cPanel → phpMyAdmin → select the new (empty) database → **Import**
2. Choose the `.sql` file from SiteGround
3. Click **Go**

For large files (over 50 MB), import via SSH (see [Back Up and Restore a MySQL Database](/getting-started/mysql-backup/)).

**Upload files:**
Connect via SFTP to Arcline and upload everything to `/home/username/public_html/`. See [Upload Files via SFTP](/getting-started/sftp/).

---

## Step 4 — Update database credentials

Edit `wp-config.php` (or your site's database config file) with the new credentials:

```php
define( 'DB_NAME',     'arclinecpanel_dbname' );
define( 'DB_USER',     'arclinecpanel_dbuser' );
define( 'DB_PASSWORD', 'your-new-password' );
define( 'DB_HOST',     'localhost' );
```

SiteGround sometimes uses a non-localhost `DB_HOST` (e.g., a private IP). On Arcline, it's always `localhost`.

---

## Step 5 — Check PHP version

1. SiteGround Site Tools → **Devs → PHP Manager** — note the PHP version
2. Arcline cPanel → **Software → Select PHP Version** — match it

Also enable any PHP extensions your site needs. Common ones: `mbstring`, `curl`, `gd`, `imagick`, `intl`, `zip`. Check the **Extensions** tab in MultiPHP Manager on Arcline.

---

## Step 6 — Remove SiteGround caching

SiteGround's SuperCacher (now called SG Optimizer) installs a caching plugin and server-side caching. Before testing on Arcline:

1. Deactivate the **SG Optimizer** plugin (or **SiteGround Optimizer**) — it will try to connect to SiteGround's infrastructure which won't be available on Arcline
2. Delete any `.htaccess` lines added by SG Optimizer that reference SiteGround-specific caching rules
3. Install a standard WordPress caching plugin instead (LiteSpeed Cache, W3 Total Cache, or WP Super Cache work well on Arcline)

If you leave SG Optimizer active, the plugin will disable itself gracefully on non-SiteGround servers, but it's cleaner to remove it.

---

## Step 7 — Test before DNS cutover

Add your Arcline server IP to your local `hosts` file:

```
YOUR_ARCLINE_IP  yourdomain.com www.yourdomain.com
```

**macOS / Linux:** `/etc/hosts` (requires sudo)
**Windows:** `C:\Windows\System32\drivers\etc\hosts` (as Administrator)

Visit your domain and check:
- Homepage and key pages load
- Images display
- WordPress admin works
- Contact forms submit
- Any e-commerce checkout works (if applicable)

Remove the `hosts` entry when done.

---

## Step 8 — Switch nameservers

**If your domain is registered at SiteGround** (they became a domain registrar):

1. SiteGround → **Domains** → select your domain
2. Go to **Nameservers** and select **Custom Nameservers**
3. Enter:
   - `ns1.arcline.it`
   - `ns2.arcline.it`
4. Save

**If registered elsewhere**, update nameservers at that registrar. See [Point Your Domain to Arcline](/getting-started/nameservers/).

Leave SiteGround running for 24–48 hours. Both servers will serve the site during propagation — that's fine, the content is the same.

---

## Step 9 — Email migration

**Moving to Arcline email:**
1. Create accounts in Arcline cPanel → **Email → Email Accounts**
2. Arcline's MX records activate after nameservers switch
3. Export existing email from SiteGround's webmail (Roundcube) before cancelling

**Keeping Google Workspace or another provider:**
After switching nameservers, re-add your external MX records in Arcline cPanel → **Domains → Zone Editor**.

---

## SiteGround-specific notes

**SG Optimizer / SuperCacher** — disable and remove before or shortly after migration. It won't cause harm on Arcline but serves no purpose.

**SiteGround staging** — if you used SiteGround's staging feature, the staging site lives on their servers. You can ignore it; migrate only the production site.

**Cloudflare via SiteGround** — SiteGround has a built-in Cloudflare integration. If your site was using it, you have a Cloudflare account attached to SiteGround. After migrating, you can either:
- Switch to Arcline nameservers directly (disables the SiteGround-managed Cloudflare)
- Keep Cloudflare: transfer the zone to your own Cloudflare account and update the A record to point to Arcline's IP

**SiteGround emails about resource limits** — if you were on SiteGround's entry-level plan, you may have hit CPU/memory limits. These limits don't carry over; Arcline has its own resource allocation.

**WordPress auto-updates** — SiteGround has WordPress auto-update settings in Site Tools. These won't apply on Arcline. Set up your own update preferences in WordPress → **Dashboard → Updates** or use a plugin like ManageWP.

---

## After migration checklist

- [ ] Site loads at `https://yourdomain.com`
- [ ] SSL certificate valid
- [ ] WordPress admin accessible
- [ ] SG Optimizer deactivated/removed
- [ ] Caching plugin configured for Arcline
- [ ] Email accounts working
- [ ] DNS propagated (verify at dnschecker.org)
- [ ] SiteGround account kept active for 48 hours as fallback
