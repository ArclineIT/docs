---
title: "Install WordPress on a VPS"
description: "Set up WordPress on your Arcline VPS from scratch — LAMP or LEMP stack with PHP, MySQL, and Nginx or Apache."
section: wordpress
order: 2
---

# Install WordPress on a VPS

Installing WordPress on a VPS gives you full control over the server configuration, performance tuning, and security. This guide covers both a LEMP stack (Linux, Nginx, MySQL, PHP-FPM) and a LAMP stack (Linux, Apache, MySQL, PHP).

If you haven't set up your VPS yet, start with [Initial VPS Setup](/vps/initial-setup/).

---

## Choose your stack

| Stack | Web Server | Best for |
|---|---|---|
| **LEMP** | Nginx + PHP-FPM | High traffic, static caching, modern setups |
| **LAMP** | Apache + PHP | Simpler `.htaccess` support, beginner-friendly |

This guide covers the LEMP stack (Nginx) as the primary setup with LAMP (Apache) notes where they differ.

---

## Step 1 — Install the stack

### LEMP (Nginx)

```bash
sudo apt update
sudo apt install -y nginx mysql-server php-fpm php-mysql php-curl php-gd \
  php-mbstring php-xml php-zip php-intl php-imagick unzip curl
```

### LAMP (Apache)

```bash
sudo apt update
sudo apt install -y apache2 mysql-server php libapache2-mod-php php-mysql \
  php-curl php-gd php-mbstring php-xml php-zip php-intl php-imagick unzip curl
```

### Secure MySQL and create the database

```bash
sudo mysql_secure_installation
```

Follow the prompts — set a root password, remove anonymous users, disallow remote root login, remove test databases, and reload privileges.

Now create the WordPress database and user:

```bash
sudo mysql -u root -p
```

```sql
CREATE DATABASE wordpress CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'wpuser'@'localhost' IDENTIFIED BY 'a-strong-password-here';
GRANT ALL PRIVILEGES ON wordpress.* TO 'wpuser'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

Replace `a-strong-password-here` with an actual strong password.

---

## Step 2 — Download and set up WordPress

```bash
cd /tmp
curl -O https://wordpress.org/latest.tar.gz
tar xzf latest.tar.gz
sudo mv wordpress /var/www/yourdomain.com
sudo chown -R www-data:www-data /var/www/yourdomain.com
```

---

## Step 3 — Configure Nginx (LEMP)

Create the Nginx site configuration:

```bash
sudo nano /etc/nginx/sites-available/yourdomain.com
```

```nginx
server {
    listen 80;
    server_name yourdomain.com www.yourdomain.com;
    root /var/www/yourdomain.com;
    index index.php index.html;

    location / {
        try_files $uri $uri/ /index.php?$args;
    }

    location ~ \.php$ {
        include snippets/fastcgi-php.conf;
        fastcgi_pass unix:/var/run/php/php8.1-fpm.sock;
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        include fastcgi_params;
    }

    location = /favicon.ico { access_log off; log_not_found off; }
    location = /robots.txt  { access_log off; log_not_found off; }

    # Block access to sensitive files
    location ~* /\.(?!well-known\/) { deny all; }
    location ~* /wp-config\.php { deny all; }
    location ~* /xmlrpc\.php { deny all; }

    # Cache static assets in the browser
    location ~* \.(css|js|ico|gif|jpg|jpeg|png|webp|svg|woff2?|ttf|otf|eot)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}
```

Enable the site and test the config:

```bash
sudo ln -s /etc/nginx/sites-available/yourdomain.com /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

Replace `php8.1-fpm.sock` with the PHP version you installed. Check your PHP version:

```bash
php -v
```

---

## Step 3 (alt) — Configure Apache (LAMP)

```bash
sudo nano /etc/apache2/sites-available/yourdomain.com.conf
```

```apache
<VirtualHost *:80>
    ServerName yourdomain.com
    ServerAlias www.yourdomain.com
    DocumentRoot /var/www/yourdomain.com

    <Directory /var/www/yourdomain.com>
        AllowOverride All
        Require all granted
    </Directory>
</VirtualHost>
```

Enable and restart:

```bash
sudo a2ensite yourdomain.com.conf
sudo a2enmod rewrite
sudo systemctl reload apache2
```

---

## Step 4 — Set up SSL with Let's Encrypt

```bash
sudo apt install -y certbot python3-certbot-nginx
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com
```

For Apache, use `python3-certbot-apache` instead.

Certbot modifies your Nginx/Apache config to add SSL automatically. Certificates renew automatically via a systemd timer — verify with:

```bash
sudo certbot renew --dry-run
```

---

## Step 5 — Complete the WordPress install

Visit `https://yourdomain.com` in your browser and complete the WordPress setup wizard. When prompted for database credentials, enter:

| Field | Value |
|---|---|
| Database Name | `wordpress` |
| Username | `wpuser` |
| Password | the password you set in Step 1 |
| Database Host | `localhost` |
| Table Prefix | `wp_` |

---

## Step 6 — Configure file permissions

WordPress needs write access to `wp-content/uploads` (for media) but the rest of your install should be read-only for the web server to prevent tampering.

```bash
sudo chown -R www-data:www-data /var/www/yourdomain.com
sudo find /var/www/yourdomain.com -type d -exec chmod 755 {} \;
sudo find /var/www/yourdomain.com -type f -exec chmod 644 {} \;
sudo chmod 640 /var/www/yourdomain.com/wp-config.php
```

---

## Step 7 — Configure PHP-FPM for WordPress

Edit the PHP-FPM pool config:

```bash
sudo nano /etc/php/8.1/fpm/pool.d/www.conf
```

Increase resource limits for a production WordPress site:

```ini
pm = dynamic
pm.max_children = 20
pm.start_servers = 5
pm.min_spare_servers = 5
pm.max_spare_servers = 10
pm.max_requests = 500
```

Adjust `pm.max_children` based on your VPS RAM: roughly `(available_RAM - 512MB) / 50MB` per child.

Apply changes:

```bash
sudo systemctl restart php8.1-fpm
```

---

## Step 8 — Set up a firewall

```bash
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
```

For stricter security, see [Set up fail2ban](/vps/fail2ban/).

---

## Step 9 — Configure WordPress cron properly

WordPress's default pseudo-cron runs on every page load, which is wasteful on a VPS. Replace it with a real system cron job:

```bash
sudo crontab -u www-data -e
```

Add:

```
*/5 * * * * /usr/bin/php /var/www/yourdomain.com/wp-cron.php > /dev/null 2>&1
```

Then disable WordPress pseudo-cron in `wp-config.php`:

```php
define( 'DISABLE_WP_CRON', true );
```

---

## Performance tuning

### PHP opcache

Uncomment and tweak in `/etc/php/8.1/fpm/php.ini`:

```ini
opcache.enable=1
opcache.memory_consumption=256
opcache.interned_strings_buffer=16
opcache.max_accelerated_files=10000
opcache.revalidate_freq=2
opcache.fast_shutdown=1
```

### MySQL tuning

For a 2 GB VPS, add to `/etc/mysql/mysql.conf.d/mysqld.cnf`:

```ini
innodb_buffer_pool_size = 512M
innodb_log_file_size = 128M
query_cache_type = 0
```

Run MySQLTuner after a few days of uptime for more specific recommendations:

```bash
sudo apt install mysqltuner
sudo mysqltuner
```

### WordPress object caching with Redis

If your VPS has enough RAM (2 GB+ free), Redis dramatically speeds up WordPress:

```bash
sudo apt install redis-server
sudo systemctl enable redis-server
```

Install the **Redis Object Cache** plugin in WordPress admin and click **Enable Object Cache**.

---

## Automating updates

Set up unattended security updates for the OS:

```bash
sudo apt install unattended-upgrades
sudo dpkg-reconfigure unattended-upgrades
```

WordPress core auto-updates are enabled by default for minor versions. For plugins and themes, enable auto-updates in WordPress admin → Plugins → toggle **Enable auto-updates** on each plugin you trust.

---

## Next steps

- [Configure W3 Total Cache](/wordpress/w3-total-cache/) for page caching and performance
- [WordPress security hardening](/wordpress/security/) — tighten file permissions, disable XML-RPC, and set up login protection
- [Automated backups with restic](/vps/restic-backups/) for off-site backups of your files and database

