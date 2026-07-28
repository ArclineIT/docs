---
title: "Supported PHP Versions"
description: "PHP versions available on Arcline shared hosting and VPS plans."
section: reference
order: 2
---

# Supported PHP Versions

Arcline supports multiple PHP versions to accommodate different application requirements.

---

## Shared hosting

On shared hosting plans, you can select your PHP version in cPanel under **Software → Select PHP Version**.

| Version | Status | End of Life | Notes |
|---------|--------|-------------|-------|
| 8.3 | Active | 2026-12-31 | Recommended for new sites |
| 8.2 | Active | 2025-12-31 | Good stability |
| 8.1 | Security only | 2024-11-30 | Upgrade recommended |
| 8.0 | End of life | 2023-11-26 | Not recommended |

---

## VPS plans

On VPS plans, you can install any PHP version from the system package manager or use a third-party repository like [Ondřej Surý's PPA](https://deb.sury.org/):

```bash
sudo apt install php8.3 php8.3-fpm php8.3-mysql
```

---

## Switching PHP versions

In cPanel, go to **Select PHP Version**, choose the version you want, and click **Set as current**. Most WordPress sites will work on any PHP 8.x version without issues.

After switching, check your site thoroughly:
- If you see blank pages or errors, your application may not support the PHP version you selected
- WordPress 5.6+ supports PHP 8.x. WordPress 6.0+ requires PHP 7.4+
- Older plugins may not be compatible with PHP 8.x — check with the plugin developer

---

## PHP extensions

The following extensions are available on shared hosting:

```
bcmath, calendar, Core, ctype, curl, date, dom, exif, fileinfo,
filter, ftp, gd, gettext, gmp, hash, iconv, imagick, intl, json,
ldap, libxml, mbstring, mysqli, mysqlnd, openssl, pcntl, pcre,
PDO, pdo_mysql, pdo_sqlite, Phar, posix, random, readline, redis,
Reflection, session, shmop, SimpleXML, sockets, sodium, SPL, sqlite3,
standard, sysvsem, sysvshm, tokenizer, xml, xmlreader, xmlwriter,
xsl, Zend OPcache, zip, zlib
```

If you need an extension that isn't listed, contact Arcline support.

