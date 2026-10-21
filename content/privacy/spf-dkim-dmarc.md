---
title: "What SPF, DKIM, and DMARC Actually Do"
description: "An explanation of email authentication standards and how to set them up on Arcline."
section: privacy
order: 4
---

# What SPF, DKIM, and DMARC Actually Do

Email authentication standards (SPF, DKIM, DMARC) prevent spammers from sending email that looks like it comes from your domain. Without them, your outgoing email is more likely to land in spam folders — or worse, be used to impersonate you.

---

## SPF — Sender Policy Framework

**What it does**: SPF publishes a list of IP addresses that are authorized to send email for your domain.

**How it works**:
1. A receiving mail server receives a message claiming to be from `@yourdomain.com`
2. The receiving server looks up your SPF record at `yourdomain.com`
3. If the sending IP is in the SPF record, the message passes. If not, it's subject to the server's spam policy.

**An SPF record looks like**:

```
v=spf1 mx ip4:203.0.113.42 ~all
```

This means: "Allow mail from your MX servers and from IP 203.0.113.42. Soft-fail everything else."

**Setting it up on Arcline**:

In your Arcline cPanel, go to **Domains → Zone Editor** and add a TXT record:

| Type | Name | Value |
|------|------|-------|
| TXT | `@` | `v=spf1 mx include:arcline.it ~all` |

The `include:arcline.it` will pull in Arcline's sending IPs automatically.

---

## DKIM — DomainKeys Identified Mail

**What it does**: DKIM signs your outgoing email with a cryptographic signature. The receiving server verifies the signature by looking up your public key in DNS.

**How it works**:
1. Arcline's mail server signs your outgoing message with a private key
2. The receiving server looks up your public DKIM key at `selector._domainkey.yourdomain.com`
3. It decrypts the signature and verifies the message wasn't tampered with in transit

**Setting it up on Arcline**:

Arcline cPanel automatically generates DKIM keys. To verify they're set up:

1. cPanel → **Email → Email Deliverability**
2. Find your domain and click **Manage**
3. You should see a green status for DKIM

The DNS record is automatically added. It looks like:

| Type | Name | Value |
|------|------|-------|
| TXT | `arcline._domainkey` | `v=DKIM1; h=sha256; p=MIGfMA0GCSqGSIb4DQEBAQUAA4GNADCBiQKBgQ...` |

---

## DMARC — Domain-based Message Authentication, Reporting & Conformance

**What it does**: DMARC tells receiving servers what to do when a message fails SPF or DKIM checks. It also provides reports so you can see who's sending email on your behalf.

**How it works**:
1. A message arrives claiming to be from your domain
2. The receiving server checks SPF and DKIM
3. If both pass, the DMARC policy doesn't matter
4. If one or both fail, the server follows your DMARC policy:
   - `none` — Take no action (just report)
   - `quarantine` — Mark as spam
   - `reject` — Reject the message outright

**Setting it up on Arcline**:

Start with `p=none` to see who's sending email for your domain without blocking anything. After a few weeks, review the reports and tighten to `p=quarantine`. After confirming all legitimate email is authenticated, move to `p=reject`.

| Type | Name | Value |
|------|------|-------|
| TXT | `_dmarc` | `v=DMARC1; p=none; rua=mailto:dmarc@yourdomain.com` |

The `rua` field specifies where you want to receive aggregate DMARC reports. Tools like [dmarcian.com](https://dmarcian.com) or [Postmark's DMARC tool](https://dmarc.postmarkapp.com) can help you parse them.

---

## All three records together

For `yourdomain.com`, your DNS zone should have these TXT records:

```
yourdomain.com.                TXT "v=spf1 mx include:arcline.it ~all"
arcline._domainkey.yourdomain.com.  TXT "v=DKIM1; h=sha256; p=..."
_dmarc.yourdomain.com.         TXT "v=DMARC1; p=quarantine; rua=mailto:dmarc@yourdomain.com"
```

---

## Testing your setup

Use one of these free tools to verify everything is working:

```bash
# Command line
dig TXT yourdomain.com +short
dig TXT arcline._domainkey.yourdomain.com +short
dig TXT _dmarc.yourdomain.com +short
```

Or visit:
- [MXToolbox](https://mxtoolbox.com/diagnostic.aspx) — enter your domain
- [Mail-Tester](https://www.mail-tester.com) — send a test email to the address shown

---

## Common issues

| Symptom | Likely cause |
|---------|--------------|
| SPF passes but emails still go to spam | Missing or misconfigured DKIM |
| DKIM passes but emails go to spam | Missing DMARC policy |
| DMARC reports show IPs you don't recognize | Someone is spoofing your domain — set `p=reject` |
| DMARC reports show a legitimate service failing | Add the service's IPs to your SPF record |
| Automated reports from third-party services | Their infrastructure needs to be included in your SPF record — contact them for their SPF include |

---

## What's next

- [Why you shouldn't put Cloudflare in front of everything](/privacy/why-not-cloudflare/)
- [Self-hosting without a CDN: performance tips](/privacy/self-hosting-performance/)
- [Set up email on Arcline](/getting-started/email-setup/)

