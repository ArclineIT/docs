---
title: "Set Up Fail2ban for SSH Brute-Force Protection"
description: "Protect your Arcline VPS from SSH brute-force attacks with fail2ban."
section: vps
order: 7
---

# Set Up Fail2ban for SSH Brute-Force Protection

Fail2ban monitors system logs for repeated failed login attempts and temporarily bans the offending IP addresses using the firewall. It's essential for any internet-facing server.

---

## Prerequisites

- A VPS with SSH access and sudo privileges
- UFW or iptables already installed (see [Initial VPS Setup](/vps/initial-setup/))

---

## Step 1 — Install fail2ban

```bash
sudo apt update
sudo apt install fail2ban -y
```

---

## Step 2 — Configure fail2ban for SSH

The default configuration file is `/etc/fail2ban/jail.conf`. Don't edit it directly — it gets overwritten on updates. Instead, create a local override:

```bash
sudo nano /etc/fail2ban/jail.local
```

```ini
[DEFAULT]
# Ban IPs for 1 hour after 5 failed attempts within 10 minutes
bantime = 3600
findtime = 600
maxretry = 5

# Send email alerts (optional)
# destemail = you@example.com
# action = %(action_mwl)s

[sshd]
enabled = true
port = ssh
logpath = %(sshd_log)s
```

If you changed your SSH port, specify it:

```ini
[sshd]
enabled = true
port = 2222
logpath = %(sshd_log)s
```

---

## Step 3 — Start fail2ban

```bash
sudo systemctl enable fail2ban
sudo systemctl start fail2ban
```

Check the status:

```bash
sudo systemctl status fail2ban
```

---

## Step 4 — Monitor banned IPs

View the SSH jail status:

```bash
sudo fail2ban-client status sshd
```

This shows the total bans and currently active bans.

View the ban log:

```bash
sudo tail -f /var/log/fail2ban.log
```

---

## Step 5 — Unban an IP

If you accidentally lock yourself out (you should have tested SSH key access before enabling, but just in case):

```bash
sudo fail2ban-client set sshd unbanip 203.0.113.42
```

Or from the console (if you still have a root session open):

```bash
sudo iptables -D f2b-sshd -s 203.0.113.42 -j DROP
```

---

## Step 6 — Additional jails (optional)

### Nginx

```ini
[nginx-http-auth]
enabled = true
logpath = /var/log/nginx/error.log
```

### Nginx bot protection (repeat offenders)

```ini
[nginx-botsearch]
enabled = true
logpath = /var/log/nginx/access.log
maxretry = 2
findtime = 86400
bantime = 86400
```

This bans IPs that hit common admin paths (wp-admin, etc.) that don't exist on your server.

### Wordpress

```ini
[wordpress]
enabled = true
filter = wordpress
logpath = /var/log/auth.log
```

You may need to create a custom filter for your specific application logs.

---

## Step 7 — Whitelist IPs

To exclude trusted IPs from bans (your office IP, for example):

```ini
[DEFAULT]
ignoreip = 127.0.0.1/8 ::1 203.0.113.100
```

---

## Permanent bans with recidive jail

Habitual offenders get progressively longer bans:

```ini
[recidive]
enabled = true
logpath = /var/log/fail2ban.log
maxretry = 3
findtime = 604800   # 1 week
bantime = 604800    # 1 week
```

An IP that triggers bans 3 times in a week gets banned for a week.

---

## Testing fail2ban

From a different machine (or after whitelisting your IP), intentionally fail SSH login a few times:

```bash
ssh nonexistent@your.vps.ip.address
```

After 5 failures, further attempts should hang or be refused. Check with:

```bash
sudo fail2ban-client status sshd
```

---

## Performance notes

Fail2ban uses minimal resources — typically under 50MB of RAM with a few jails enabled. It reads log files using Python's `pyinotify` (if available) or polls every second.

If you have high-traffic sites with aggressive bots, increase `findtime` and lower `maxretry` to catch them sooner:

```ini
[nginx-botsearch]
maxretry = 2
findtime = 3600
bantime = 86400
```

---

## What's next

- [Deploy a Go binary](/vps/go-systemd/) as a systemd service
- [Set up automated backups](/vps/automated-backups/) with restic

