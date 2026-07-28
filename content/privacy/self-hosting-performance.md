---
title: "Self-Hosting Without a CDN: Performance Tips"
description: "How to make your self-hosted site fast without relying on Cloudflare or other CDNs."
section: privacy
order: 3
---

# Self-Hosting Without a CDN: Performance Tips

A common objection to self-hosting is "but it won't be fast without a CDN." That's not true for most sites. With proper configuration, a well-tuned Nginx server on a decent VPS will load pages in under 200ms for visitors on the same continent — more than fast enough for a great user experience.

---

## Tip 1 — Enable HTTP/2 and HTTP/3

HTTP/2 multiplexes requests over a single connection. HTTP/3 (QUIC) reduces latency even further.

```nginx
server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    # HTTP/3 requires Nginx 1.25+ with quic support
    listen 443 quic reuseport;
}
```

Check that HTTP/2 is active by looking at the Chrome DevTools → Network tab — your requests should show `h2` as the protocol.

---

## Tip 2 — Aggressive static file caching

Serve static assets with long cache lifetimes and immutable headers so browsers never re-request them:

```nginx
location ~* \.(css|js|png|jpg|jpeg|gif|ico|svg|woff2?)$ {
    expires 1y;
    add_header Cache-Control "public, immutable";
    access_log off;
}
```

With immutable caching, the browser will never check if the file has changed until the user force-reloads.

---

## Tip 3 — Enable Gzip compression

Compress text-based responses:

```nginx
gzip on;
gzip_vary on;
gzip_proxied any;
gzip_comp_level 5;
gzip_min_length 256;
gzip_types
    text/plain
    text/css
    text/javascript
    application/javascript
    application/json
    application/xml
    image/svg+xml;
```

This can reduce HTML, CSS, and JS payloads by 60–80%.

---

## Tip 4 — Use a PHP opcode cache

For WordPress and other PHP sites, enable OPcache:

```ini
; /etc/php/8.3/cli/conf.d/10-opcache.ini
opcache.enable=1
opcache.memory_consumption=128
opcache.max_accelerated_files=10000
opcache.revalidate_freq=120
```

This keeps compiled PHP scripts in memory, avoiding recompilation on every request.

---

## Tip 5 — Install a WordPress caching plugin

For WordPress sites, use a caching plugin that generates static HTML:

- **WP Super Cache** or **W3 Total Cache**: Generate static HTML files served directly by Nginx
- **LiteSpeed Cache**: If you're running LiteSpeed (LSAPI mode)

See our [W3 Total Cache Configuration](/wordpress/w3-total-cache/) guide for detailed setup instructions without a CDN.

---

## Tip 6 — Tune Nginx worker settings

```nginx
worker_processes auto;
worker_rlimit_nofile 65535;

events {
    worker_connections 4096;
    use epoll;
    multi_accept on;
}
```

`auto` sets the worker count to the number of CPU cores. `worker_connections 4096` allows each worker to handle 4,000 concurrent connections.

---

## Tip 7 — Serve images efficiently

- **Use WebP** instead of JPEG/PNG (30% smaller). Convert with `cwebp`:
  ```bash
  sudo apt install webp -y
  cwebp -q 80 input.jpg -o output.webp
  ```
- **Use responsive images** with srcset so mobile devices don't download desktop-sized images
- **Lazy-load below-the-fold images** with `loading="lazy"` attribute:
  ```html
  <img src="photo.jpg" loading="lazy" alt="...">
  ```
- **Compress images** with a tool like ImageMagick or optipng before uploading

---

## Tip 8 — Optimize your database

For MySQL/MariaDB:
- Enable query cache (MySQL 5.7 and earlier)
- Run `mysqlcheck -o --all-databases` weekly
- Remove unused plugins and post revisions in WordPress
- Use an index on frequently queried columns

```sql
-- Clean up WordPress post revisions
DELETE FROM wp_posts WHERE post_type = 'revision';

-- Optimize tables
OPTIMIZE TABLE wp_posts, wp_postmeta;
```

---

## Tip 9 — Set proper buffer sizes

```nginx
client_body_buffer_size 128k;
client_max_body_size 50m;
client_header_buffer_size 1k;
large_client_header_buffers 4 8k;
output_buffers 32 32k;
postpone_output 1460;
```

These settings prevent Nginx from buffering more data than necessary per connection.

---

## Tip 10 — Choose a geographically close VPS location

Arcline's data center location matters. A VPS in Dallas will serve US visitors faster than one in Frankfurt. For international audiences, consider:

- One VPS in the US + one in Europe with a simple round-robin DNS
- Or just accept slightly higher latency for overseas visitors — a 200ms response time is still perfectly usable

---

## Benchmarking your setup

After making these changes, test your performance:

```bash
# Install siege for load testing
sudo apt install siege -y
siege -c 50 -t 60s https://example.com

# Or use curl to measure response time
curl -w "@curl-format.txt" -o /dev/null -s https://example.com
```

Create a format file (`curl-format.txt`):

```
    time_namelookup:  %{time_namelookup}s
       time_connect:  %{time_connect}s
    time_appconnect:  %{time_appconnect}s
   time_redirect:  %{time_redirect}s
time_pretransfer:  %{time_pretransfer}s
   time_starttransfer:  %{time_starttransfer}s
                     ----------
          time_total:  %{time_total}s
```

---

## Real-world results

A typical Arcline VPS running WordPress with:
- Nginx + PHP 8.3 FPM
- OPcache enabled
- W3 Total Cache (page cache + database cache)
- WebP images
- Gzip compression

Will serve pages in **150–300ms** to US visitors and handle **500+ concurrent users** on a $15/mo plan.

You don't need Cloudflare to be fast. You need a well-configured server.

---

## What's next

- [What SPF, DKIM, and DMARC actually do](/privacy/spf-dkim-dmarc/)
- [Why you shouldn't put Cloudflare in front of everything](/privacy/why-not-cloudflare/)
- [How to check if a host is self-hosted](/privacy/arcline-check-walkthrough/)

