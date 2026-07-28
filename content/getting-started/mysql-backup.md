---
title: "Back Up and Restore a MySQL Database"
description: "How to export and import MySQL databases on Arcline shared hosting using phpMyAdmin and the command line."
section: getting-started
order: 3
---

# Back Up and Restore a MySQL Database

Regular database backups are the most important thing you can do to protect your site. A file backup without a database backup is useless for dynamic sites like WordPress.

---

## Export via phpMyAdmin (easiest)

phpMyAdmin is available in cPanel → **Databases → phpMyAdmin**.

1. In the left sidebar, click the database name you want to export.
2. Click the **Export** tab at the top.
3. Leave the method as **Quick** and format as **SQL**.
4. Click **Go**.

The browser will download a `.sql` file. Store it somewhere safe — not on the same server.

For large databases (over 50 MB), use the **Custom** export method and enable **Add DROP TABLE** so the import won't fail on existing tables.

---

## Export via command line (mysqldump)

SSH into your server and run:

```bash
mysqldump -u username -p database_name > backup_$(date +%Y%m%d).sql
```

You'll be prompted for the MySQL password. On cPanel servers, the MySQL username is prefixed with your cPanel username:

```bash
mysqldump -u cpanelusername_dbusername -p cpanelusername_dbname > backup.sql
```

Check your cPanel → **Databases → MySQL Databases** for the exact username and database name.

To compress the backup immediately:

```bash
mysqldump -u username -p database_name | gzip > backup_$(date +%Y%m%d).sql.gz
```

---

## Download your backup via SFTP

After exporting from the command line, download the `.sql` or `.sql.gz` file to your local machine using FileZilla or Cyberduck (see [Upload Files via SFTP](/getting-started/sftp/)).

The file will be in your home directory: `/home/username/backup.sql`.

---

## Import a database

### phpMyAdmin

1. In cPanel → phpMyAdmin, click the target database in the left sidebar.
2. Click the **Import** tab.
3. Click **Choose File** and select your `.sql` file.
4. Click **Go**.

phpMyAdmin has a file size limit (usually 50–100 MB). For larger databases, use the command line.

### Command line import

```bash
mysql -u username -p database_name < backup.sql
```

For a compressed backup:

```bash
gunzip < backup.sql.gz | mysql -u username -p database_name
```

> Make sure the target database already exists in cPanel before importing. Create it in cPanel → **Databases → MySQL Databases** if needed, and make sure the database user is assigned to it with all privileges.

---

## Automate daily backups with cron

SSH in and open your crontab:

```bash
crontab -e
```

Add this line to run a backup every day at 2 AM:

```
0 2 * * * mysqldump -u cpanelusername_dbuser -pYOURPASSWORD cpanelusername_dbname | gzip > ~/backups/db_$(date +\%Y\%m\%d).sql.gz
```

A few notes:
- Replace `YOURPASSWORD` with your actual MySQL password (no space after `-p`)
- `~/backups/` must exist: `mkdir ~/backups`
- The `\%` escaping is required in crontab

Also set up a cron to delete backups older than 30 days to avoid filling up disk space:

```
30 2 * * * find ~/backups -name "*.sql.gz" -mtime +30 -delete
```

cPanel also has a built-in cron interface under **Advanced → Cron Jobs** if you prefer a GUI.

---

## cPanel full backup

For a complete backup of files **and** databases together, use cPanel → **Files → Backup Wizard**. Choose **Back Up Your Website** and download the full backup. These can be large — only practical for occasional snapshots, not daily automation.
