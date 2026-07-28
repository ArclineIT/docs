---
title: "Deploy a Static Site with Nginx"
description: "Host a static HTML site, Hugo, or Jekyll site on your Arcline VPS with Nginx."
section: vps
order: 3
---

# Deploy a Static Site with Nginx

Static sites are fast, secure, and simple to host. This guide covers deploying plain HTML, Hugo, and Jekyll sites on your Arcline VPS.

---

## Prerequisites

- A VPS with Nginx installed (see [Nginx + PHP-FPM + MySQL](/vps/nginx-php-mysql/) for setup)
- A domain pointed to your VPS IP
- Your static site files (HTML, CSS, JS) ready to upload

---

## Step 1 — Create the site directory

```bash
sudo mkdir -p /var/www/example.com
sudo chown -R $USER:$USER /var/www/example.com
```

---

## Step 2 — Upload your site files

### Via SFTP (FileZilla, Cyberduck)

Connect to your VPS:

```
Host: your.vps.ip.address
User: yourname
Port: 22
Protocol: SFTP
```

Upload files to `/var/www/example.com/`.

### Via rsync (command line)

```bash
rsync -avz --delete ./_site/ yourname@your.vps.ip.address:/var/www/example.com/
```

The `--delete` flag removes remote files that no longer exist locally — perfect for rebuilds.

### Via SCP

```bash
scp -r ./my-site/* yourname@your.vps.ip.address:/var/www/example.com/
```

---

## Step 3 — Configure Nginx

Create a site configuration:

```bash
sudo nano /etc/nginx/sites-available/example.com
```

```nginx
server {
    listen 80;
    server_name example.com www.example.com;
    root /var/www/example.com;
    index index.html;

    # Gzip static assets
    location ~* \.(css|js|png|jpg|jpeg|gif|ico|svg|woff2?)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
        access_log off;
    }

    # HTML files — shorter cache
    location ~* \.html$ {
        expires 1h;
        add_header Cache-Control "public, must-revalidate";
    }

    location / {
        try_files $uri $uri/ =404;
    }
}
```

Enable the site:

```bash
sudo ln -s /etc/nginx/sites-available/example.com /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

---

## Step 4 — Set up SSL

```bash
sudo certbot --nginx -d example.com -d www.example.com
```

---

## Deploying a Hugo site

If you build your site with Hugo locally:

1. Generate the site: `hugo`
2. The output goes to the `public/` directory
3. Upload `public/` to your server

Or **build directly on the VPS** (for CI/CD style deployment):

```bash
# Install Hugo on the VPS
sudo apt install hugo -y
# Or download a specific version from GitHub releases

# Clone your repo
git clone https://git.arcline.it/yourname/your-site.git /var/www/example.com-source

# Build
cd /var/www/example.com-source
hugo -d /var/www/example.com
```

---

## Deploying a Jekyll site

Jekyll requires Ruby. Install it on the VPS:

```bash
sudo apt install ruby-full build-essential -y
sudo gem install jekyll bundler
```

Clone and build:

```bash
git clone https://git.arcline.it/yourname/your-site.git /var/www/example.com-source
cd /var/www/example.com-source
bundle install
jekyll build -d /var/www/example.com
```

---

## Security headers for static sites

Add this to your Nginx config inside the `server` block for recommended security headers:

```nginx
add_header X-Frame-Options "SAMEORIGIN" always;
add_header X-Content-Type-Options "nosniff" always;
add_header Referrer-Policy "strict-origin-when-cross-origin" always;
add_header Permissions-Policy "camera=(), microphone=(), geolocation=()";
```

Test and reload:

```bash
sudo nginx -t && sudo systemctl reload nginx
```

---

## What's next

- [Deploy a Node.js app](/vps/nodejs-pm2/) with PM2 + Nginx
- [Deploy a Go binary](/vps/go-systemd/) as a systemd service
- [Set up automated backups](/vps/automated-backups/) with restic

