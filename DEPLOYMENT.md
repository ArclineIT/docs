# Deployment Guide

> **From source to production** — how every Arcline service is deployed,
> managed, and maintained.

---

## Table of Contents

- [Overview](#overview)
- [Production Environment](#production-environment)
- [Deployment Methods](#deployment-methods)
- [Service Directory Layout](#service-directory-layout)
- [User & Permissions Model](#user--permissions-model)
- [Systemd Service Units](#systemd-service-units)
- [Nginx Reverse Proxy](#nginx-reverse-proxy)
- [Deploy by Service](#deploy-by-service)
- [Database Management](#database-management)
- [Backup & Restore](#backup--restore)
- [Rollback Procedures](#rollback-procedures)
- [Health Checks](#health-checks)
- [Troubleshooting](#troubleshooting)

---

## Overview

Every Arcline service deploys as a single static Go binary. The deployment
strategy depends on the service type:

| Type | Method | Examples |
|------|--------|---------|
| **Web service** | Docker container | portal, billing, docs |
| **Native binary** | Direct install on host | website (OpenBSD), uptime |
| **CLI tool** | Install to `$PATH` | check, audit, dns, migrate, vault |

All services share the same core deployment flow:

1. **Build** — Compile static Go binary with `CGO_ENABLED=0`
2. **Package** — Containerize (Docker) or ship raw binary
3. **Deploy** — Transfer to production host and start
4. **Verify** — Health check endpoint or process monitoring

---

## Production Environment

### Host Specifications

| Host | Role | OS | Location |
|------|------|----|----------|
| `srv01` | Primary web server | OpenBSD | Arcline datacenter |
| `srv02` | Application server | Linux (Alpine) | Arcline datacenter |
| `srv03` | Database & storage | Linux (Alpine) | Arcline datacenter |

### Network

- All services listen on **127.0.0.1** (localhost) only
- Nginx handles TLS termination and reverse proxy to local services
- Public access is through pfSense port forwarding (443 → Nginx)

### Runtime Dependencies

- **Docker** (for containerized services): `docker-ce`
- **Nginx**: `nginx` (OpenBSD: `nginx` package)
- **SQLite**: No runtime dependency (pure Go, file-based)
- **ca-certificates**: Required in all containers for TLS

---

## Deployment Methods

### Method 1: Docker Deployment (Web Services)

Used for: portal, billing, docs (dynamic), email (future)

```bash
# ── On production host ────────────────────────────────────────────────

# 1. Pull the latest image
docker login -u <user> -p <token> <registry>
docker pull <registry>/<project>:latest

# 2. Stop and remove existing container
docker stop <container-name> 2>/dev/null || true
docker rm   <container-name> 2>/dev/null || true

# 3. Start new container
docker run -d \
  --name <container-name> \
  --restart unless-stopped \
  -p 127.0.0.1:<host-port>:<container-port> \
  -v /opt/<service>/.env:/app/.env:ro \
  -v /opt/<service>/data:/app/data \
  <registry>/<project>:latest

# 4. Verify
docker ps | grep <container-name>
docker logs <container-name> --tail 20
```

### Method 2: Native Binary (OpenBSD)

Used for: website

```bash
# ── On build machine ──────────────────────────────────────────────────

# Cross-compile
GOOS=openbsd GOARCH=amd64 go build -ldflags="-s -w" -o binary-openbsd-amd64 .

# Deploy
scp binary-openbsd-amd64 srv01:/usr/local/bin/<binary>
rsync -av --delete static/ srv01:/var/www/<domain>/static/

# ── On production host ────────────────────────────────────────────────

# Restart service
doas rcctl restart <service>
```

### Method 3: Static Binary (CLI Tools)

Used for: uptime, status, check, audit, dns, migrate, vault

```bash
# ── On build machine ──────────────────────────────────────────────────

# Cross-compile for the target architecture
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o binary-linux-amd64 .

# Deploy
scp binary-linux-amd64 srv02:/usr/local/bin/<binary>
ssh srv02 "chmod 0755 /usr/local/bin/<binary>"

# Verify
ssh srv02 "<binary> version"
```

---

## Service Directory Layout

### Containerized Services

```
/opt/arcline-<service>/
├── .env              # Environment variables (arcline:arcline, 0640)
└── data/             # Persistent data directory
    └── <service>.db  # SQLite database (created at runtime)
```

### Native Services (OpenBSD)

```
/usr/local/bin/
├── arcline-web          # Go binary
└── rc.d/                # rc.d scripts (managed by rcctl)

/var/www/arclineit/
├── static/              # CSS, JS, images
├── templates/           # Go HTML templates
└── .env                 # Environment config
```

### CLI Tools

```
/usr/local/bin/
├── arcline-uptime
├── arcline-status
├── arcline-check
├── arcline-audit
├── arcline-dns
├── arcline-migrate
├── arcline-vault
└── arcline-email

/etc/arcline/
├── uptime.yaml          # Uptime monitor config
├── status.yaml          # Status page config
├── status.d/            # Modular status config directory
└── email.toml           # Email server config

/var/lib/arcline/
├── uptime.db            # Uptime monitoring database
└── vault/               # Vault secrets storage
```

---

## User & Permissions Model

### System User

```bash
# Create the arcline user (done once per host)
addgroup -S arcline
adduser -S -G arcline -h /opt arcline
```

### Directory Permissions

```bash
# Application directory
chown -R arcline:arcline /opt/arcline-<service>
chmod 0750 /opt/arcline-<service>

# Environment file (sensitive!)
chown arcline:arcline /opt/arcline-<service>/.env
chmod 0640 /opt/arcline-<service>/.env

# Database file (created by app, but ensure permissions)
chown arcline:arcline /opt/arcline-<service>/data/*.db
chmod 0640 /opt/arcline-<service>/data/*.db
```

### Docker Container User

All Docker containers run as the `arcline` non-root user:

```dockerfile
RUN addgroup -S arcline && adduser -S -G arcline arcline
USER arcline
```

---

## Systemd Service Units

### Template: `arcline-<service>.service`

```ini
[Unit]
Description=Arcline <Service Name> — <short description>
After=network.target

[Service]
Type=simple
User=arcline
Group=arcline
WorkingDirectory=/opt/arcline-<service>
EnvironmentFile=/opt/arcline-<service>/.env
ExecStart=/opt/arcline-<service>/<binary> [flags]
Restart=on-failure
RestartSec=5s

# Hardening
NoNewPrivileges=yes
PrivateTmp=yes
ProtectSystem=strict
ProtectHome=yes
ReadWritePaths=/opt/arcline-<service>

[Install]
WantedBy=multi-user.target
```

### Service Management

```bash
# Install the unit file
sudo cp arcline-<service>.service /etc/systemd/system/
sudo systemctl daemon-reload

# Enable on boot
sudo systemctl enable arcline-<service>

# Start / stop / restart / status
sudo systemctl start   arcline-<service>
sudo systemctl stop    arcline-<service>
sudo systemctl restart arcline-<service>
sudo systemctl status  arcline-<service>

# View logs
sudo journalctl -u arcline-<service> -f
```

### Hardening Options

| Option | Purpose |
|--------|---------|
| `NoNewPrivileges=yes` | Prevent privilege escalation via `suid` binaries |
| `PrivateTmp=yes` | Isolated `/tmp` for the service |
| `ProtectSystem=strict` | Read-only `/usr` and `/etc` |
| `ProtectHome=yes` | No access to `/home`, `/root` |
| `ReadWritePaths=...` | Only allow writes to the data directory |

**Reference implementation:** `portal/deploy/arcline-portal.service`

---

## Nginx Reverse Proxy

### Template: `nginx-<service>.conf`

```nginx
# HTTP → HTTPS redirect
server {
    listen 80;
    server_name <domain>;
    return 301 https://$host$request_uri;
}

# HTTPS server
server {
    listen 443 ssl;
    server_name <domain>;

    ssl_certificate     /etc/ssl/<domain>/fullchain.pem;
    ssl_certificate_key /etc/ssl/<domain>/privkey.pem;
    ssl_protocols       TLSv1.2 TLSv1.3;
    ssl_ciphers         HIGH:!aNULL:!MD5;

    # Security headers
    add_header Strict-Transport-Security "max-age=63072000; includeSubDomains" always;
    add_header X-Frame-Options DENY always;
    add_header X-Content-Type-Options nosniff always;

    location / {
        proxy_pass         http://127.0.0.1:<upstream-port>;
        proxy_set_header   Host              $host;
        proxy_set_header   X-Real-IP         $remote_addr;
        proxy_set_header   X-Forwarded-For   $proxy_add_x_forwarded_for;
        proxy_set_header   X-Forwarded-Proto $scheme;
        proxy_read_timeout 30s;
    }
}
```

### SSL Certificate Management

Certificates from **Let's Encrypt** (certbot) or **Step CA** (internal).

```bash
# External (Let's Encrypt)
certbot certonly --webroot -w /var/www/acme -d <domain>

# Internal (Step CA)
step certificate install /etc/step-ca/certs/intermediate_ca.crt
```

**Reference implementation:** `portal/deploy/nginx-portal.conf`

---

## Deploy by Service

### website — `arcline.it`

| Detail | Value |
|--------|-------|
| **Method** | Native binary on OpenBSD |
| **Binary** | `/usr/local/bin/arcline-web` |
| **Static files** | `/var/www/arclineit/static/` |
| **Service manager** | OpenBSD `rcctl` |
| **Deploy command** | `make deploy` |

```bash
# Quick deploy
make cross
scp arcline-web-openbsd-amd64 srv01:/usr/local/bin/arcline-web
rsync -av --delete static/ srv01:/var/www/arclineit/static/
ssh srv01 "rcctl restart arcline-web"
```

### billing — `client.arcline.it`

| Detail | Value |
|--------|-------|
| **Method** | Docker container |
| **Container name** | `arcline-billing` |
| **Internal port** | 8082 |
| **Data volume** | `/opt/arcline-billing/.env:/app/.env:ro` |
| **Database** | SQLite in container (ephemeral — backup on host) |

```bash
docker run -d \
  --name arcline-billing \
  --restart unless-stopped \
  -p 127.0.0.1:8082:8082 \
  -v /opt/arcline-billing/.env:/app/.env:ro \
  <registry>/arcline-billing:latest
```

### portal — `portal.arclineit.com`

| Detail | Value |
|--------|-------|
| **Method** | Docker container |
| **Container name** | `arcline-portal` |
| **Internal port** | 8082 |
| **Data volume** | `/opt/arcline-portal/.env:/app/.env:ro` |
| **Database** | SQLite in container (ephemeral — backup on host) |

```bash
docker run -d \
  --name arcline-portal \
  --restart unless-stopped \
  -p 127.0.0.1:8082:8082 \
  -v /opt/arcline-portal/.env:/app/.env:ro \
  <registry>/arcline-portal:latest
```

### git — `git.arcline.it`

| Detail | Value |
|--------|-------|
| **Method** | Docker container |
| **Container name** | `arcline-gitea` |
| **Internal ports** | 3000 (web), 22 (SSH) |
| **Data volume** | `/opt/arcline-gitea/data:/data` |
| **Database** | SQLite (or PostgreSQL for multi-instance) |

```bash
# Deploy with docker-compose
docker compose -f git/deploy/docker-compose.yml up -d

# Or manual
docker run -d \
  --name arcline-gitea \
  --restart unless-stopped \
  -p 127.0.0.1:3000:3000 \
  -p 22:22 \
  -v /opt/arcline-gitea/data:/data \
  -e DOMAIN=git.arcline.it \
  -e ROOT_URL=https://git.arcline.it \
  -e SSH_DOMAIN=git.arcline.it \
  gitea/gitea:latest

# Nginx config (redirect + proxy)
sudo cp git/deploy/nginx-git-redirect.conf /etc/nginx/sites-enabled/
sudo nginx -t && sudo systemctl reload nginx

# Verify
curl -I https://git.arcline.it
ssh -T _gitea@git.arcline.it
```

### uptime — Internal Monitor

| Detail | Value |
|--------|-------|
| **Method** | Static binary or Docker |
| **Binary** | `/usr/local/bin/arcline-uptime` |
| **Config** | `/etc/arcline/uptime.yaml` |
| **Database** | `/var/lib/arcline/uptime.db` |

```bash
# As a service
arcline-uptime start --config /etc/arcline/uptime.yaml

# Or via systemd
[Service]
ExecStart=/usr/local/bin/arcline-uptime start --config /etc/arcline/uptime.yaml
```

### status — `status.arclineit.com`

| Detail | Value |
|--------|-------|
| **Method** | Static HTML generation |
| **Binary** | `/usr/local/bin/arcline-status` |
| **Config** | `/etc/arcline/status.yaml` + `/etc/arcline/status.d/` |
| **Output** | `/var/www/status/` (served by Nginx) |

```bash
# Generate status page
arcline-status build --config /etc/arcline/status.yaml --out /var/www/status/

# Watch mode (regenerate on config change)
arcline-status build --config /etc/arcline/status.yaml --out /var/www/status/ --watch
```

### docs — `docs.arclineit.com`

| Detail | Value |
|--------|-------|
| **Method** | Docker container (dynamic) or Nginx (static) |
| **Container name** | `arcline-docs` |
| **Internal port** | 8080 |

```bash
# Dynamic variant (Go server with Markdown rendering)
docker run -d \
  --name arcline-docs \
  --restart unless-stopped \
  -p 127.0.0.1:8080:8080 \
  -v /opt/arcline-docs/.env:/app/.env:ro \
  <registry>/arcline-docs:latest

# Static variant (pre-built HTML via Dockerfile.static)
docker run -d \
  --name arcline-docs \
  --restart unless-stopped \
  -p 127.0.0.1:80:80 \
  <registry>/arcline-docs:latest-static
```

---

## Database Management

### Backup

```bash
# SQLite databases — simple file copy
cp /opt/arcline-<service>/data/<service>.db /backup/<service>-$(date +%Y%m%d).db

# Compress
gzip /backup/<service>-*.db
```

### Restore

```bash
# Stop the service first
docker stop arcline-<service>
# or
systemctl stop arcline-<service>

# Restore database
cp /backup/<service>-<date>.db.gz /opt/arcline-<service>/data/
gunzip /opt/arcline-<service>/data/<service>-<date>.db.gz
mv /opt/arcline-<service>/data/<service>-<date>.db /opt/arcline-<service>/data/<service>.db
chown arcline:arcline /opt/arcline-<service>/data/<service>.db

# Restart
docker start arcline-<service>
# or
systemctl start arcline-<service>
```

### Maintenance

```bash
# Vacuum SQLite database (reclaim space)
sqlite3 /opt/arcline-<service>/data/<service>.db "VACUUM;"

# Integrity check
sqlite3 /opt/arcline-<service>/data/<service>.db "PRAGMA integrity_check;"
```

---

## Backup & Restore

### What to Back Up

| Asset | Location | Frequency |
|-------|----------|-----------|
| SQLite databases | `/opt/arcline-*/data/*.db` | Daily |
| Environment files | `/opt/arcline-*/.env` | On change |
| SSL certificates | `/etc/ssl/` | On renewal |
| Nginx configs | `/etc/nginx/` | On change |
| Systemd units | `/etc/systemd/system/arcline-*.service` | On change |
| Status config | `/etc/arcline/` | On change |
| Uptime config | `/etc/arcline/uptime.yaml` | On change |

### Backup Script

```bash
#!/bin/sh
# /usr/local/bin/arcline-backup

BACKUP_DIR="/backup/arcline/$(date +%Y-%m-%d)"
mkdir -p "$BACKUP_DIR"

# Databases
for db in /opt/arcline-*/data/*.db; do
    cp "$db" "$BACKUP_DIR/"
done

# Environment files (mask secrets)
for env in /opt/arcline-*/.env; do
    cp "$env" "$BACKUP_DIR/$(basename $(dirname $env)).env"
done

# Configs
cp -r /etc/arcline "$BACKUP_DIR/"

# Compress
tar czf "$BACKUP_DIR.tar.gz" -C "$(dirname $BACKUP_DIR)" "$(basename $BACKUP_DIR)"
rm -rf "$BACKUP_DIR"

# Sync to backup host
rsync -av "$BACKUP_DIR.tar.gz" backup-host:/backups/arcline/
```

---

## Rollback Procedures

### Docker Rollback

```bash
# 1. Check previous image tags
docker images <registry>/<project>

# 2. Deploy a specific tag instead of latest
docker run -d \
  --name <container>-rollback \
  <registry>/<project>:<previous-sha>

# 3. Verify health
curl http://127.0.0.1:<port>/health

# 4. Swap if healthy
docker stop <container>
docker rm <container>
docker rename <container>-rollback <container>
```

### Native Binary Rollback

```bash
# 1. Keep previous binary versions
cp /usr/local/bin/<binary> /usr/local/bin/<binary>.bak

# 2. Restore previous version
cp /usr/local/bin/<binary>.bak /usr/local/bin/<binary>

# 3. Restart service
systemctl restart <service>
# or
rcctl restart <service>
```

---

## Health Checks

### Docker Health Check (in Dockerfile)

```dockerfile
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1
```

### Manual Health Check

```bash
# Web service
curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:<port>/health

# Process check
pgrep -x <binary>

# Docker check
docker ps --filter "name=<container>" --filter "status=running"
```

---

## Troubleshooting

### Service Won't Start

```bash
# Check systemd logs
journalctl -u arcline-<service> -f

# Check Docker logs
docker logs arcline-<service>

# Check binary directly
/opt/arcline-<service>/<binary> 2>&1

# Common issues:
# - Missing .env file → ExecStart fails
# - Port conflict → change port in .env
# - Permission denied → check chown/chmod
```

### Database Issues

```bash
# Check SQLite integrity
sqlite3 /opt/arcline-<service>/data/<service>.db "PRAGMA integrity_check;"
# Expected: "ok"

# Check disk space
df -h /opt/arcline-<service>/data/

# Check file permissions
ls -la /opt/arcline-<service>/data/<service>.db
# Expected: -rw-r----- arcline arcline
```

### Connection Refused

```bash
# Check if service is listening
ss -tlnp | grep <port>

# Check Nginx is running
systemctl status nginx

# Check pfSense NAT rules
# Verify port forwarding from WAN:443 → host:443
```

### Quick Recovery

```bash
# Universal restart sequence
docker stop arcline-<service> 2>/dev/null
docker rm arcline-<service> 2>/dev/null
docker pull <registry>/<project>:latest
docker run -d --restart unless-stopped \
  -p 127.0.0.1:<port>:<port> \
  -v /opt/arcline-<service>/.env:/app/.env:ro \
  <registry>/<project>:latest
```

