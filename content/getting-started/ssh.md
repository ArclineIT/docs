---
title: "Connect to Your Server via SSH"
description: "How to SSH into your Arcline server on Windows, macOS, and Linux — plus setting up SSH keys for passwordless login."
section: getting-started
order: 1
---

# Connect to Your Server via SSH

SSH (Secure Shell) is the standard way to access your Arcline server's command line. Once connected you can manage files, install software, restart services, and run anything you'd run locally.

## What you'll need

- Your server's hostname or IP address (from your welcome email)
- Your SSH username (typically the cPanel username)
- Your SSH password **or** an SSH key pair

---

## Connect on macOS or Linux

Open Terminal and run:

```bash
ssh username@your-server.arcline.it
```

Replace `username` with your cPanel username and `your-server.arcline.it` with your server hostname. You'll be prompted for your password on first login.

If your server runs SSH on a non-standard port (common with shared hosting — cPanel uses port **2222**):

```bash
ssh -p 2222 username@your-server.arcline.it
```

---

## Connect on Windows

Windows 10 and 11 include OpenSSH built in. Open **PowerShell** or **Windows Terminal** and run the same command:

```powershell
ssh username@your-server.arcline.it
```

Or with port 2222:

```powershell
ssh -p 2222 username@your-server.arcline.it
```

**Using PuTTY** — if you prefer a GUI client, download PuTTY from [putty.org](https://www.putty.org). Fill in:

- Host Name: `your-server.arcline.it`
- Port: `2222`
- Connection type: SSH

Click **Open**, accept the server fingerprint on first connection, and log in with your username and password.

---

## Set up SSH key authentication

SSH keys are more secure than passwords and let you log in without typing a password every time.

**Step 1 — Generate a key pair** (skip if you already have one at `~/.ssh/id_ed25519`):

```bash
ssh-keygen -t ed25519 -C "your@email.com"
```

Accept the default file location. Set a passphrase if you want extra protection.

**Step 2 — Copy your public key to the server:**

```bash
ssh-copy-id -p 2222 username@your-server.arcline.it
```

Or manually: log in, open `~/.ssh/authorized_keys` on the server, and paste the contents of your local `~/.ssh/id_ed25519.pub` file.

**Step 3 — Verify it works:**

```bash
ssh -p 2222 username@your-server.arcline.it
```

You should log in without being asked for a password.

**Via cPanel** — you can also manage SSH keys in **cPanel → Security → SSH Access**. Upload your public key there and cPanel will authorize it automatically.

---

## Simplify repeated connections

Add an entry to `~/.ssh/config` to avoid typing the full command every time:

```
Host arcline
    HostName your-server.arcline.it
    User username
    Port 2222
    IdentityFile ~/.ssh/id_ed25519
```

Now you can connect with just:

```bash
ssh arcline
```

---

## Common errors

**Connection refused** — the SSH service may be running on a different port, or a firewall is blocking the connection. Confirm the port with Arcline support.

**Permission denied (publickey)** — your key isn't authorized on the server. Check that `~/.ssh/authorized_keys` contains your public key and has permissions `600`.

**Host key verification failed** — the server's fingerprint changed (common after a reinstall). Remove the old entry: `ssh-keygen -R your-server.arcline.it` and reconnect.

**Connection timed out** — your IP may be temporarily blocked by fail2ban after too many failed attempts. Contact support to whitelist your IP.
