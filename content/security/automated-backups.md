---
title: "Set Up Automated Backups with Restic"
description: "Automate encrypted off-site backups of a Linux server using restic."
section: security
order: 2
---

# Set Up Automated Backups with Restic

Restic is a fast, encrypted backup tool that supports local and remote storage backends (SFTP, S3, B2, rsync.net). This guide covers backing up a Linux server to a remote repository.

---

## Prerequisites

- A Linux server with sudo access
- A backup destination (SFTP server, Backblaze B2, or local storage)

---

## Step 1 — Install restic

```bash
sudo apt update
sudo apt install restic -y
```

Verify:

```bash
restic version
```

---

## Step 2 — Initialize a repository

### Option A: SFTP/SSH (recommended)

If you have SSH access to a backup server:

```bash
restic init --repo sftp:backup@backup-server:/var/backups/example-vps/
```

You'll be prompted for a repository password — this encrypts your backups. Store it in a password manager — if you lose it, you cannot recover your data.

### Option B: Backblaze B2

```bash
export B2_ACCOUNT_ID="your-application-key-id"
export B2_ACCOUNT_KEY="your-application-key"
restic init --repo b2:bucket-name:/example-vps
```

### Option C: Local directory

```bash
sudo mkdir -p /backups/example-vps-repo
restic init --repo /backups/example-vps-repo
```

---

## Step 3 — Create a backup script

```bash
sudo nano /usr/local/bin/backup.sh
```

```bash
#!/bin/bash
set -e

# Repository location and password
export RESTIC_REPOSITORY="sftp:backup@backup-server:/var/backups/example-vps/"
export RESTIC_PASSWORD="your-repo-password"

# Files and directories to back up
BACKUP_PATHS=(
    /var/www
    /etc/nginx
    /etc/letsencrypt
    /opt
    /home
)

# Directories to exclude
EXCLUDE_PATTERNS=(
    --exclude "/var/www/example.com/cache"
    --exclude "*.log"
)

echo "Starting backup at $(date)"

# Create the backup
restic backup "${EXCLUDE_PATTERNS[@]}" "${BACKUP_PATHS[@]}"

# Keep last 7 daily, 4 weekly, 6 monthly snapshots
restic forget --keep-daily 7 --keep-weekly 4 --keep-monthly 6 --prune

echo "Backup completed at $(date)"
```

Make it executable:

```bash
sudo chmod +x /usr/local/bin/backup.sh
```

---

## Step 4 — Test the backup

Run the backup manually:

```bash
sudo /usr/local/bin/backup.sh
```

List snapshots:

```bash
restic -r sftp:backup@backup-server:/var/backups/example-vps/ snapshots
```

> You'll need `RESTIC_PASSWORD` exported or passed via `--password-file` for any restic command.

---

## Step 5 — Schedule daily backups with systemd

Create a service file:

```bash
sudo nano /etc/systemd/system/restic-backup.service
```

```ini
[Unit]
Description=Restic backup
Wants=network-online.target
After=network-online.target

[Service]
Type=oneshot
ExecStart=/usr/local/bin/backup.sh
```

Create a timer:

```bash
sudo nano /etc/systemd/system/restic-backup.timer
```

```ini
[Unit]
Description=Daily restic backup

[Timer]
OnCalendar=daily
RandomizedDelaySec=3600
Persistent=true

[Install]
WantedBy=timers.target
```

Enable and start the timer:

```bash
sudo systemctl daemon-reload
sudo systemctl enable restic-backup.timer
sudo systemctl start restic-backup.timer
```

Verify:

```bash
sudo systemctl status restic-backup.timer
sudo systemctl list-timers | grep restic
```

---

## Step 6 — Restoring from a backup

List available snapshots:

```bash
restic -r sftp:backup@backup-server:/var/backups/example-vps/ snapshots
```

Restore the latest snapshot:

```bash
restic -r sftp:backup@backup-server:/var/backups/example-vps/ restore latest --target /tmp/restore
```

Or restore a specific snapshot by ID:

```bash
restic -r ... restore <snapshot-id> --target /tmp/restore
```

To restore only specific paths:

```bash
restic -r ... restore <snapshot-id> --target /tmp/restore --path /var/www
```

---

## Step 7 — Database backups

For MySQL databases, add a pre-backup dump step:

```bash
#!/bin/bash
set -e

# Dump all databases
mysqldump --all-databases --single-transaction --quick | gzip > /tmp/mysql-all.sql.gz

# Include the dump in the backup
restic backup --hostname example-vps /tmp/mysql-all.sql.gz "${BACKUP_PATHS[@]}"

rm /tmp/mysql-all.sql.gz
```

---

## Monitoring backups

Add a health check notification:

```bash
# After successful backup
curl -fsS -m 10 --retry 5 https://hc-ping.com/your-uuid

# Or on failure, notify via Discord/Slack webhook
curl -fsS -m 10 -X POST -H "Content-Type: application/json" \
  -d '{"content":"Backup failed on example-vps"}' \
  https://discord.com/api/webhooks/your-webhook-url
```

---

## What's next

- [Install fail2ban](/security/fail2ban/) for SSH brute-force protection
- [Harden WordPress](/security/wordpress/) if the server runs it

