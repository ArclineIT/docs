---
title: "Upload Files via SFTP"
description: "How to connect to your Arcline server with FileZilla or Cyberduck and transfer website files."
section: getting-started
order: 2
---

# Upload Files via SFTP

SFTP (SSH File Transfer Protocol) lets you upload, download, and manage files on your server using a graphical client. It's encrypted and the standard method for file transfers — FTP without TLS is not supported on Arcline servers.

## What you'll need

- Your server's hostname or IP address
- Your cPanel username and password (or an SSH key)
- SFTP port: **2222** (cPanel default)

---

## FileZilla

FileZilla is free, cross-platform, and works well for most file transfer tasks.

**Download:** [filezilla-project.org](https://filezilla-project.org) — use the regular FileZilla Client, not FileZilla Pro.

### Quick connect

In the toolbar at the top, fill in:

- **Host:** `sftp://your-server.arcline.it`
- **Username:** your cPanel username
- **Password:** your cPanel password
- **Port:** `2222`

Click **Quickconnect**. Accept the server fingerprint when prompted.

### Site Manager (recommended for saved connections)

Go to **File → Site Manager → New Site** and fill in:

| Field | Value |
|---|---|
| Protocol | SFTP – SSH File Transfer Protocol |
| Host | `your-server.arcline.it` |
| Port | `2222` |
| Logon Type | Normal |
| User | your cPanel username |
| Password | your cPanel password |

Click **Connect**. The site is saved for next time.

### Navigating your files

- **Left panel** — your local computer
- **Right panel** — your server

Your website files live in `/home/username/public_html/`. Navigate there on the right side before uploading.

To upload: drag files from the left panel to the right. To download: drag from right to left.

> **Before uploading a new version of your site**, download a copy of `public_html` first so you have a backup.

---

## Cyberduck

Cyberduck is a good option on macOS (also available on Windows).

**Download:** [cyberduck.io](https://cyberduck.io) — it's free, though they ask for a donation.

Click **Open Connection** and set:

| Field | Value |
|---|---|
| Protocol | SFTP (SSH File Transfer Protocol) |
| Server | `your-server.arcline.it` |
| Port | `2222` |
| Username | your cPanel username |
| Password | your cPanel password |

Click **Connect**. Bookmark the connection (⌘D on Mac) to save it.

---

## SSH key authentication

If you've set up SSH keys (see [Connect via SSH](/getting-started/ssh/)), both FileZilla and Cyberduck support key-based login:

**FileZilla:** In Site Manager, set Logon Type to **Key file** and point it to your private key file (`~/.ssh/id_ed25519`).

**Cyberduck:** In Open Connection, leave the password blank and check **Use Public Key Authentication**. Cyberduck will pick up keys from `~/.ssh/` automatically.

---

## Common paths

| What | Path |
|---|---|
| Your website root | `/home/username/public_html/` |
| WordPress uploads | `/home/username/public_html/wp-content/uploads/` |
| cPanel mail | `/home/username/mail/` |
| cPanel logs | `/home/username/logs/` |
| cPanel backups | `/home/username/backup/` |

---

## Troubleshooting

**Authentication failed** — double-check your username, password, and port (2222, not 22). Make sure you're using SFTP not FTP.

**Connection timed out** — a firewall or fail2ban may be blocking your IP. Try connecting from a different network to test, then contact support.

**Permission denied when uploading** — the target directory may be owned by a different user or have restrictive permissions. Check the permissions in FileZilla (right-click → File permissions) or via SSH: `ls -la /home/username/public_html/`.
