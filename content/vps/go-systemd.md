---
title: "Deploy a Go Binary as a Systemd Service"
description: "Run a Go application as a background service on your Arcline VPS with systemd."
section: vps
order: 5
---

# Deploy a Go Binary as a Systemd Service

Go compiles to a single static binary — no runtime, no dependencies, no package manager. This makes it ideal for running as a systemd service on your Arcline VPS.

---

## Prerequisites

- A VPS with SSH access
- A Go binary compiled for Linux amd64 (or arm64 if using an ARM VPS)

---

## Step 1 — Build your Go binary

On your local machine, cross-compile for your target VPS:

```bash
# For Linux amd64 (most common)
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o myapp

# For Linux arm64 (e.g., Raspberry Pi)
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o myapp
```

The `CGO_ENABLED=0` flag ensures a fully static binary with no external library dependencies.

---

## Step 2 — Upload the binary

```bash
scp myapp yourname@your.vps.ip.address:/tmp/
```

On the VPS, move it to its final location:

```bash
sudo mkdir -p /opt/myapp
sudo mv /tmp/myapp /opt/myapp/
sudo chmod +x /opt/myapp/myapp
```

---

## Step 3 — Create a systemd service file

```bash
sudo nano /etc/systemd/system/myapp.service
```

```ini
[Unit]
Description=My Go Application
After=network.target
Wants=network-online.target

[Service]
Type=simple
User=yourname
Group=yourname
WorkingDirectory=/opt/myapp
ExecStart=/opt/myapp/myapp
Restart=always
RestartSec=5
EnvironmentFile=-/opt/myapp/.env

# Security hardening
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/myapp

[Install]
WantedBy=multi-user.target
```

---

## Step 4 — Create an environment file

If your app reads configuration from environment variables:

```bash
sudo nano /opt/myapp/.env
```

```
PORT=8080
DATABASE_PATH=/opt/myapp/data.db
LOG_LEVEL=info
```

Secure the file:

```bash
sudo chmod 600 /opt/myapp/.env
sudo chown yourname:yourname /opt/myapp/.env
```

---

## Step 5 — Start and enable the service

```bash
sudo systemctl daemon-reload
sudo systemctl start myapp
sudo systemctl enable myapp   # starts on boot
```

Check the status:

```bash
sudo systemctl status myapp
```

---

## Step 6 — Set up Nginx reverse proxy (if it's a web app)

If your Go app serves HTTP on a port like `8080`, put Nginx in front:

```nginx
server {
    listen 80;
    server_name api.example.com;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Enable SSL with Certbot:

```bash
sudo certbot --nginx -d api.example.com
```

---

## Managing the service

| Command | Description |
|---------|-------------|
| `sudo systemctl start myapp` | Start the service |
| `sudo systemctl stop myapp` | Stop the service |
| `sudo systemctl restart myapp` | Restart the service |
| `sudo systemctl status myapp` | Show status and recent logs |
| `sudo systemctl enable myapp` | Enable auto-start on boot |
| `sudo systemctl disable myapp` | Disable auto-start |
| `journalctl -u myapp -f` | Follow live logs |

---

## Updating the binary

```bash
# Build new version locally
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o myapp

# Upload
scp myapp yourname@your.vps.ip.address:/tmp/

# On the VPS
sudo systemctl stop myapp
sudo cp /tmp/myapp /opt/myapp/
sudo systemctl start myapp
sudo systemctl status myapp
```

---

## Logging

Your Go app's stdout and stderr are automatically captured by systemd's journal. View them with:

```bash
journalctl -u myapp -f
```

For persistent log files, your app can write to a file, or you can configure systemd to forward logs to syslog:

```bash
sudo mkdir -p /var/log/myapp
sudo chown yourname:yourname /var/log/myapp
```

Then in your service file, add:

```
StandardOutput=append:/var/log/myapp/stdout.log
StandardError=append:/var/log/myapp/stderr.log
```

---

## What's next

- [Set up automated backups](/vps/automated-backups/) with restic
- [Install fail2ban](/vps/fail2ban/) for SSH brute-force protection

