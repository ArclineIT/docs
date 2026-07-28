---
title: "Initial VPS Setup (Debian/Ubuntu)"
description: "First steps after provisioning a new VPS: users, SSH keys, firewall, and system updates."
section: vps
order: 1
---

# Initial VPS Setup

This guide walks you through the first steps after provisioning a new Arcline VPS. You'll create a non-root user, harden SSH, set up a firewall, and apply system updates.

---

## Before you begin

You'll receive your VPS login credentials from Arcline after provisioning. Your initial login is as `root` via SSH.

```
ssh root@your.vps.ip.address
```

If you're on macOS or Linux, the SSH client is built in. On Windows, use PowerShell, Windows Terminal, or WSL.

---

## Step 1 — Create a non-root user

Working as root for daily tasks is risky. Create an administrative user:

```bash
adduser yourname
```

Follow the prompts to set a strong password. Then add the user to the `sudo` group:

```bash
usermod -aG sudo yourname
```

For Debian, the sudo group may be named differently. Verify with `groups yourname` — if you see `sudo`, you're set.

---

## Step 2 — Copy your SSH key

From your local machine (not the VPS), copy your SSH public key to the new user:

```bash
ssh-copy-id yourname@your.vps.ip.address
```

If `ssh-copy-id` isn't available, manually create the `.ssh` directory and `authorized_keys` file:

```bash
# On the VPS, as your new user:
mkdir -p ~/.ssh
chmod 700 ~/.ssh
# Edit this file and paste your public key
nano ~/.ssh/authorized_keys
chmod 600 ~/.ssh/authorized_keys
```

Test that key-based login works from a new terminal:

```bash
ssh yourname@your.vps.ip.address
```

If you can log in without a password prompt, proceed.

---

## Step 3 — Harden SSH

Edit the SSH server configuration:

```bash
sudo nano /etc/ssh/sshd_config
```

Make the following changes:

```
PermitRootLogin no
PasswordAuthentication no
PubkeyAuthentication yes
Port 22
```

If you changed the SSH port, note it — you'll need it in firewall rules below.

Restart SSH:

```bash
sudo systemctl restart sshd
```

Before closing your current session, open a **second terminal** and verify you can still log in as your new user. If something went wrong, you still have the root session to fix it.

---

## Step 4 — Set up the firewall (UFW)

UFW (Uncomplicated Firewall) is the easiest way to manage iptables rules on Ubuntu/Debian.

First, allow SSH so you don't lock yourself out:

```bash
sudo ufw allow ssh
```

If you changed the SSH port:

```bash
sudo ufw allow 2222/tcp   # replace 2222 with your port
```

For a web server, allow HTTP and HTTPS:

```bash
sudo ufw allow http
sudo ufw allow https
```

Enable the firewall:

```bash
sudo ufw enable
```

Check the status:

```bash
sudo ufw status verbose
```

Default deny on incoming, allow on outgoing is the correct policy. Only the ports you explicitly opened should be listed.

---

## Step 5 — Apply system updates

Keep the system current:

```bash
sudo apt update
sudo apt upgrade -y
```

Enable automatic security updates:

```bash
sudo apt install unattended-upgrades -y
sudo dpkg-reconfigure --priority=low unattended-upgrades
```

Select **Yes** when prompted about automatically installing security updates.

---

## Step 6 — Set the timezone and hostname

Set the correct timezone:

```bash
sudo timedatectl set-timezone America/New_York   # or your timezone
```

Verify with `timedatectl`.

Set a descriptive hostname:

```bash
sudo hostnamectl set-hostname myserver
```

Add it to `/etc/hosts`:

```bash
echo "127.0.1.1 myserver" | sudo tee -a /etc/hosts
```

---

## Step 7 — Install essential tools

A few packages you'll want on every server:

```bash
sudo apt install -y curl wget git htop net-tools ufw fail2ban
```

Fail2ban will be configured in a dedicated guide. For now it runs with sensible defaults.

---

## What's next

- [Install Nginx + PHP-FPM + MySQL](/vps/nginx-php-mysql/) for a LEMP stack
- [Deploy a static site](/vps/static-site/) with Nginx
- [Set up fail2ban](/vps/fail2ban/) for SSH brute-force protection

