---
title: "Install Nginx + PHP-FPM + MySQL on Debian/Ubuntu"
description: "Set up a LEMP stack on your Arcline VPS for hosting PHP applications and WordPress."
section: vps
order: 2
---

# Install Nginx + PHP-FPM + MySQL

This guide walks through setting up a LEMP stack (Linux, Nginx, MySQL, PHP) on your Arcline VPS. This is the foundation for hosting WordPress, Laravel, and most PHP applications.

---

## Prerequisites

- A VPS provisioned through Arcline
- SSH access with sudo privileges (see [Initial VPS Setup](/vps/initial-setup/))
- A domain pointed to your VPS IP (see [Point Your Nameservers](/getting-started/nameservers/))

---

## Step 1 — Install Nginx

Nginx is the web server. It's lightweight, fast, and handles concurrent connections much better than Apache.

```bash
sudo apt update
sudo apt install nginx -y
```

Verify it's running:

```bash
sudo systemctl status nginx
```

Visit your VPS IP in a browser — you should see the default Nginx welcome page.

---

## Step 2 — Install MySQL (MariaDB)

MariaDB is a drop-in replacement for MySQL that's faster and more open:

```bash
sudo apt install mariadb-server mariadb-client -y
```

Run the security script:

```bash
sudo mysql_secure_installation
```

Follow the prompts:
- Set a root password
- Remove anonymous users: **Y**
- Disallow root login remotely: **Y**
- Remove test database: **Y**
- Reload privilege tables: **Y**

Verify the installation:

```bash
sudo mysql -u root -p
```

You should get a MariaDB prompt. Type `exit` to quit.

---

## Step 3 — Install PHP-FPM

PHP-FPM (FastCGI Process Manager) runs PHP scripts. Install the version appropriate for your needs:

```bash
# PHP 8.3 (recommended for most applications)
sudo apt install php8.3-fpm php8.3-mysql php8.3-curl php8.3-gd php8.3-mbstring php8.3-xml php8.3-xmlrpc php8.3-zip php8.3-intl php8.3-bcmath -y
```

For WordPress, the required extensions are: `mysql`, `curl`, `gd`, `mbstring`, `xml`, `zip`.

Verify PHP-FPM is running:

```bash
sudo systemctl status php8.3-fpm
```

---

## Step 4 — Configure Nginx for PHP

Create a site configuration file:

```bash
sudo nano /etc/nginx/sites-available/example.com
```

Replace `example.com` with your actual domain:

```nginx
server {
    listen 80;
    server_name example.com www.example.com;
    root /var/www/example.com;
    index index.php index.html;

    location / {
        try_files $uri $uri/ =404;
    }

    location ~ \.php$ {
        include snippets/fastcgi-php.conf;
        fastcgi_pass unix:/var/run/php/php8.3-fpm.sock;
    }

    location ~ /\.ht {
        deny all;
    }
}
```

Create the web root and enable the site:

```bash
sudo mkdir -p /var/www/example.com
sudo chown -R $USER:$USER /var/www/example.com
sudo ln -s /etc/nginx/sites-available/example.com /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

---

## Step 5 — Test PHP processing

Create a test file:

```bash
echo "<?php phpinfo();" > /var/www/example.com/info.php
```

Visit `http://example.com/info.php` in your browser — you should see the PHP information page.

Remove this file after testing — it exposes sensitive server information:

```bash
rm /var/www/example.com/info.php
```

---

## Step 6 — Set up SSL with Let's Encrypt

Install Certbot:

```bash
sudo apt install certbot python3-certbot-nginx -y
```

Obtain a certificate:

```bash
sudo certbot --nginx -d example.com -d www.example.com
```

Follow the prompts. Certbot will automatically modify your Nginx config to serve HTTPS and set up auto-renewal. Verify the renewal timer:

```bash
sudo systemctl status certbot.timer
```

---

## Directory structure summary

```
/var/www/example.com/        # Web root — your site files go here
├── index.php                # Main entry point
├── wp-admin/                # (WordPress) admin panel
├── wp-content/              # (WordPress) themes, plugins, uploads
└── ...                      # Your application files

/etc/nginx/
├── sites-available/         # All site configs (enabled via symlink)
│   └── example.com
└── sites-enabled/           # Active site configs
    └── example.com → ../sites-available/example.com
```

---

## Troubleshooting

**Nginx fails to start:** Check syntax with `sudo nginx -t`. Look at the error log: `sudo journalctl -u nginx`.

**PHP not processing:** Verify the PHP-FPM socket path matches in both your Nginx config and PHP-FPM pool config: `sudo nano /etc/php/8.3/fpm/pool.d/www.conf` — look for `listen =`.

**MySQL connection refused:** Make sure MySQL is running: `sudo systemctl status mariadb`. Check the socket: `sudo mysql -u root -p -S /var/run/mysqld/mysqld.sock`.

---

## What's next

- [Install WordPress](/wordpress/install-vps/) on your VPS
- [Deploy a Node.js app](/vps/nodejs-pm2/) with PM2
- [Set up automated backups](/vps/automated-backups/) with restic

