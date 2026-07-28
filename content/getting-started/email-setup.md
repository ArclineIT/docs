---
title: "Set Up Email on Your Domain"
description: "How to create email accounts in cPanel, configure MX records, and connect to Outlook, Apple Mail, or Thunderbird."
section: getting-started
order: 4
---

# Set Up Email on Your Domain

Arcline hosting includes email hosting for your domain. You can create addresses like `hello@yourdomain.com`, `support@yourdomain.com`, or anything else — and access them via webmail or a mail client.

---

## Create an email account

In cPanel, go to **Email → Email Accounts** and click **Create**.

Fill in:

- **Domain** — select your domain from the dropdown
- **Username** — the part before the `@` (e.g. `hello`)
- **Password** — use a strong password or click **Generate**
- **Storage** — set a mailbox quota (500 MB is fine for most users; set to unlimited only if needed)

Click **Create**. The account is ready immediately.

---

## Access webmail

Every Arcline email account comes with webmail. Visit:

```
https://your-server.arcline.it/webmail
```

Or use the shortcut `https://yourdomain.com/webmail`. Log in with your full email address (`hello@yourdomain.com`) and the password you set.

Roundcube is the default webmail client. It works in any browser and is good for occasional access.

---

## Connect to Outlook

1. Open Outlook and go to **File → Add Account**.
2. Enter your full email address and click **Advanced options → Let me set up my account manually**.
3. Choose **IMAP**.

Fill in the incoming mail settings:

| Setting | Value |
|---|---|
| Server | `mail.yourdomain.com` |
| Port | `993` |
| Encryption | SSL/TLS |
| Username | your full email address |
| Password | your email password |

Outgoing mail (SMTP):

| Setting | Value |
|---|---|
| Server | `mail.yourdomain.com` |
| Port | `465` |
| Encryption | SSL/TLS |
| Username | your full email address |
| Password | your email password |

Click **Next** then **Done**. Outlook will verify the settings and download your mail.

> If you see a certificate warning, make sure you're using `mail.yourdomain.com` (not `localhost` or an IP address). Your SSL certificate covers that hostname.

---

## Connect to Apple Mail

1. Open Mail and go to **Mail → Add Account → Other Mail Account**.
2. Enter your name, full email address, and password. Click **Sign In**.
3. Apple Mail will try to auto-configure. If it fails, enter the server details manually:

**Incoming:**
- IMAP server: `mail.yourdomain.com`
- Port: `993`
- Use SSL: yes

**Outgoing:**
- SMTP server: `mail.yourdomain.com`
- Port: `465`
- Use SSL: yes
- Authentication: Password

Click **Sign In**. Mail will sync your inbox.

---

## Connect to Thunderbird

1. Open Thunderbird and go to **Account Settings → Account Actions → Add Mail Account**.
2. Enter your name, email address, and password. Click **Configure manually**.

Fill in:

| | Incoming (IMAP) | Outgoing (SMTP) |
|---|---|---|
| Server | `mail.yourdomain.com` | `mail.yourdomain.com` |
| Port | `993` | `465` |
| SSL | SSL/TLS | SSL/TLS |
| Auth | Normal password | Normal password |

Click **Done**. Thunderbird will test and save the account.

---

## Mobile devices (iOS / Android)

Use the same IMAP/SMTP settings as above. On both iOS and Android:

1. Go to **Settings → Mail → Accounts → Add Account → Other**.
2. Select **Add Mail Account**.
3. Enter your name, email, and password — then tap **Next**.
4. If auto-configuration fails, enter the incoming and outgoing server details manually (IMAP port 993, SMTP port 465, both SSL).

---

## MX records

If your domain's DNS is managed at Arcline (using Arcline nameservers), MX records are set automatically — you don't need to change anything.

If your DNS is managed elsewhere (at your registrar, Cloudflare, etc.), add these MX records manually:

| Type | Priority | Value |
|---|---|---|
| MX | 0 | `mail.yourdomain.com` |

Then add an A record pointing `mail.yourdomain.com` to your server's IP address (found in cPanel → **Server Information**).

DNS changes take 1–48 hours to propagate. While propagating, email sent to your domain may be delayed or bounce.

---

## SPF and DKIM

SPF and DKIM tell other mail servers that your server is authorized to send email for your domain. Without them, your outgoing email is more likely to be marked as spam.

**If you're using Arcline nameservers**, SPF and DKIM records are added automatically by cPanel.

**If you're using external DNS**, you'll need to add them manually.

**SPF** — add a TXT record to your domain's DNS:

```
v=spf1 +a +mx +ip4:YOUR_SERVER_IP ~all
```

Replace `YOUR_SERVER_IP` with your server's IP address.

**DKIM** — in cPanel, go to **Email → Email Deliverability** and click **Manage** next to your domain. Copy the DKIM TXT record shown and add it to your DNS provider.

After adding both records, go back to **Email → Email Deliverability** and click **Repair** if cPanel flags any problems.

---

## Forwarders

To forward email from one address to another (e.g., `hello@yourdomain.com` to a Gmail account):

1. cPanel → **Email → Forwarders → Add Forwarder**
2. Set the source address and destination
3. Click **Add Forwarder**

You can forward to an external address without creating a full mailbox, which is handy for aliases.

---

## Troubleshooting

**Can't send email** — check that port 465 (SMTP SSL) isn't blocked by your ISP. Some ISPs block port 25. Try port 587 with STARTTLS as an alternative.

**Mail going to spam** — verify SPF and DKIM records are set correctly. Check cPanel → **Email → Email Deliverability** for warnings.

**Authentication failed in mail client** — use your full email address as the username, not just the local part. Confirm the password in cPanel → **Email → Email Accounts → Manage**.

**Certificate warning** — you must use `mail.yourdomain.com` as the server hostname. Using your server's hostname directly (e.g. `server42.arcline.it`) is fine too as long as it has a valid SSL certificate — check **Email Deliverability** in cPanel for the correct hostname.
