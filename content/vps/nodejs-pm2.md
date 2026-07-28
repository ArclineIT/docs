---
title: "Deploy a Node.js App with PM2 and Nginx"
description: "Run a Node.js application behind Nginx reverse proxy with PM2 process management."
section: vps
order: 4
---

# Deploy a Node.js App with PM2 and Nginx

This guide covers running a Node.js application on your Arcline VPS with PM2 for process management and Nginx as a reverse proxy.

---

## Prerequisites

- A VPS with Nginx installed
- A Node.js application (Express, Koa, Fastify, or similar)
- SSH access to your VPS

---

## Step 1 — Install Node.js

Install Node.js from the official NodeSource repository (recommended over the system package manager):

```bash
# Node.js 22.x (LTS)
curl -fsSL https://deb.nodesource.com/setup_22.x | sudo -E bash -
sudo apt install nodejs -y
```

Verify:

```bash
node --version
npm --version
```

---

## Step 2 — Install PM2

PM2 is a process manager that keeps your Node.js app running, handles logging, and provides monitoring.

```bash
sudo npm install -g pm2
```

---

## Step 3 — Upload your application

Create a directory for your app:

```bash
sudo mkdir -p /var/www/example.com
sudo chown -R $USER:$USER /var/www/example.com
```

Upload your application files via SFTP, rsync, or git:

```bash
cd /var/www/example.com
git clone https://git.arcline.it/yourname/your-app.git .
npm install --production
```

---

## Step 4 — Start the app with PM2

```bash
pm2 start app.js --name example-app
```

Or if your app uses `npm start`:

```bash
pm2 start npm --name example-app -- start
```

Save the PM2 process list so it restarts on reboot:

```bash
pm2 save
pm2 startup
```

The `pm2 startup` command will output a command for you to run with sudo. Follow its instructions.

---

## Step 5 — Configure Nginx as a reverse proxy

Your Node.js app runs on a port like `3000` or `8080`. Nginx will sit in front of it, handling SSL and serving static assets directly.

Create an Nginx config:

```bash
sudo nano /etc/nginx/sites-available/example.com
```

```nginx
server {
    listen 80;
    server_name example.com www.example.com;

    # If your app serves static files from a public directory
    location /static/ {
        alias /var/www/example.com/public/;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    location / {
        proxy_pass http://127.0.0.1:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
}
```

Enable the site:

```bash
sudo ln -s /etc/nginx/sites-available/example.com /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

Set up SSL:

```bash
sudo certbot --nginx -d example.com -d www.example.com
```

---

## Step 6 — Environment variables

Create an environment file:

```bash
nano /var/www/example.com/.env
```

```
PORT=3000
NODE_ENV=production
DATABASE_URL=postgres://...
```

Update your PM2 process to load it:

```bash
pm2 delete example-app
pm2 start app.js --name example-app --env-file /var/www/example.com/.env
pm2 save
```

Or use a PM2 ecosystem file (`ecosystem.config.js`):

```javascript
module.exports = {
  apps: [{
    name: 'example-app',
    script: 'app.js',
    env_file: '/var/www/example.com/.env',
    instances: 2,
    exec_mode: 'cluster',
    max_memory_restart: '500M',
  }]
};
```

Then:

```bash
pm2 start ecosystem.config.js
```

---

## PM2 useful commands

| Command | Description |
|---------|-------------|
| `pm2 list` | List all processes |
| `pm2 logs` | Show live logs |
| `pm2 logs --lines 100` | Show last 100 lines |
| `pm2 monit` | Real-time CPU/memory monitor |
| `pm2 restart example-app` | Restart an app |
| `pm2 reload example-app` | Zero-downtime reload |
| `pm2 stop example-app` | Stop an app |
| `pm2 delete example-app` | Remove from PM2 |
| `pm2 save` | Save process list |
| `pm2 startup` | Generate startup script |

---

## WebSocket support

If your app uses WebSockets, ensure the upgrade headers are passed through. The config above already includes them. For socket.io, add these to your Nginx config:

```nginx
location /socket.io/ {
    proxy_pass http://127.0.0.1:3000;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_read_timeout 86400;
}
```

---

## What's next

- [Deploy a Go binary](/vps/go-systemd/) as a systemd service
- [Set up automated backups](/vps/automated-backups/) with restic
- [Install fail2ban](/vps/fail2ban/) for SSH protection

