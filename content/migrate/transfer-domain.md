---
title: "Transfer a Domain to Arcline"
description: "How to transfer your domain registration to Arcline's registrar service — what you need to unlock, what it costs, and what to expect during the transfer."
section: migrate
order: 7
---

# Transfer a Domain to Arcline

Domain transfers move your domain **registration** from one registrar to another. It's separate from moving your website hosting — you can host at Arcline while keeping your domain registered wherever it is now, or you can transfer the registration to Arcline so everything is in one place.

---

## Before you transfer

A domain transfer typically takes **5–7 days**. During that time:

- Your website and email **are not affected** if DNS records stay the same
- You can't update nameservers or DNS records at the old registrar
- The domain is locked for 60 days after the transfer completes (ICANN rule)

**Transfers are not free.** ICANN charges registries a transfer fee, which becomes the renewal cost. When you transfer a domain to Arcline, you pay for **one additional year of registration** added to your current expiration date. If your domain expires in 6 months, after the transfer it will expire in 18 months.

---

## Step 1 — Check eligibility

A domain must meet these requirements to transfer:

- **Registered more than 60 days ago** — ICANN prohibits transfers within 60 days of initial registration
- **Not transferred in the last 60 days** — same lock applies after a previous transfer
- **More than 15 days from expiration** — most registrars block transfers of domains about to expire
- **Not in redemption or pending delete status** — the domain must be active

You can check your domain's registration date and status via WHOIS:

```
whois yourdomain.com
```

Or use a web-based WHOIS lookup at [whois.icann.org](https://whois.icann.org).

---

## Step 2 — Prepare the domain at your current registrar

### Unlock the domain

Most registrars keep domains locked to prevent unauthorized transfers. Find the lock setting:

- **GoDaddy:** Domain Settings → Domain lock → toggle off
- **Namecheap:** Domain List → Manage → Sharing & Transfer → unlock
- **Google Domains:** DNS → Registration settings → unlock
- **Cloudflare:** Domain Registration → Manage → Configuration → unlock
- **Other registrars:** look for "Domain lock", "Registrar lock", "Transfer lock", or "Theft protection"

The lock may take a few minutes to a few hours to release.

### Get the authorization code (EPP code)

The authorization code (also called EPP code or transfer key) is a one-time code required to approve the transfer:

- **GoDaddy:** Domain Settings → Transfer domain away from GoDaddy → Get authorization code
- **Namecheap:** Domain List → Manage → Transfer Out → Request EPP code (sent to your email)
- **Google Domains:** DNS → Registration settings → Get authorization code
- **Cloudflare:** Domain Registration → Manage → Configuration → Request auth code
- **Other registrars:** look for "Transfer out", "Authorization code", "EPP code", or "Auth info"

The code is typically a string of 6–32 characters. Save it — you'll need it in the next step.

### Verify admin contact email

Transfer confirmations are sent to the domain's **admin contact email** listed in WHOIS. Make sure this email address:

- Is one you have access to
- Isn't an email address at the domain you're transferring (e.g., `admin@yourdomain.com`) — because if DNS changes, you might miss the confirmation email

If needed, update the admin contact email at your current registrar before starting the transfer.

### Disable WHOIS privacy (if required)

Some registrars require WHOIS privacy to be disabled during a transfer. Check your registrar's transfer documentation. If required:

- **GoDaddy:** Domain Settings → WHOIS privacy → turn off
- **Namecheap:** Domain List → Manage → Domain Privacy → disable
- **Others:** look for "WHOIS privacy", "Privacy protection", or "ID protection"

Turn it back on after the transfer completes.

---

## Step 3 — Initiate the transfer at Arcline

Contact Arcline support to start a domain transfer. Provide:

- The domain name
- The authorization code (EPP code) from Step 2
- Confirmation that the domain is unlocked

Arcline will initiate the transfer through our registrar. The cost is one year's renewal fee, which will be added to your Arcline invoice.

---

## Step 4 — Approve the transfer

After Arcline initiates the transfer, you'll receive two emails:

1. **From your current registrar** — asking for confirmation that you want to transfer the domain away. You must click the approval link within a few days, or the transfer cancels.
2. **From the registry** (Verisign for `.com`/`.net`, Public Interest Registry for `.org`, etc.) — confirming the transfer has been requested.

Follow the instructions in both emails. If you don't receive them within a few hours, check your spam folder and verify the admin contact email is correct.

Some registrars (GoDaddy, Namecheap) allow you to approve the transfer from their dashboard instead of waiting for the email. If that option is available, use it — it's faster.

---

## Step 5 — Wait for the transfer to complete

Once approved, the transfer takes **5 days by default**. The losing registrar can release the domain sooner, but most of them wait the full period as a security measure.

During this time:

- Your website and email **continue working normally** — the transfer only affects the registration, not DNS
- You can't change nameservers at your current registrar
- The domain shows a "pending transfer" status in WHOIS

If your domain's DNS is managed at your current registrar, your DNS records are copied as-is and continue pointing wherever they pointed before. Nothing breaks during a transfer.

---

## Step 6 — After the transfer completes

Once Arcline confirms the transfer is complete:

1. **Verify your DNS records** — log in to your Arcline account and check that your DNS zone is correctly set up with the records you need (A, MX, CNAME, TXT, etc.)
2. **Update nameservers** — if you want DNS managed by Arcline, make sure nameservers are set to `ns1.arcline.it` and `ns2.arcline.it`
3. **Re-enable WHOIS privacy** — turn privacy protection back on if you want it
4. **Check email** — if your email uses the domain (e.g., `you@yourdomain.com`), verify that MX records are present and mail is flowing

---

## Transferring from specific registrars

### GoDaddy

GoDaddy often tries to retain domains by making the transfer process confusing. The key steps:

1. Unlock: Domain Settings → scroll to **Domain lock** → toggle off. Wait 10 minutes.
2. Auth code: GoDaddy calls it the **Authorization Code**. They'll email the admin contact with a link to retrieve it.
3. After initiating: GoDaddy sends a confirmation email with a link to approve the transfer. If you ignore this email, the transfer **will still complete after 5 days** — GoDaddy changed this policy and now auto-approves transfers after the waiting period.

### Namecheap

1. Unlock: Domain List → Manage → Sharing & Transfer → unlock
2. Auth code: Request via **Transfer Out** tab → **Request EPP Code** (emailed to you)
3. Namecheap's confirmation email includes an "Approve Transfer" link — click it for a faster transfer

### Google Domains

1. Unlock: DNS → Registration settings → unlock
2. Auth code: **Get authorization code** in Registration settings
3. Google Domains generally processes transfers quickly

### Cloudflare (as registrar)

1. Unlock: Domain Registration → Manage → Configuration → unlock
2. Auth code: **Request auth code** in Configuration
3. Cloudflare automatically removes your DNS zone after a successful transfer, which can cause downtime. Before transferring away from Cloudflare, export your DNS records (Cloudflare dashboard → DNS → Export) so you can re-import them at Arcline.

---

## What if the transfer fails?

Common reasons for transfer failure and how to fix them:

| Problem | Fix |
|---|---|
| Domain is locked | Unlock at current registrar and try again |
| Wrong authorization code | Request a new one from the current registrar |
| Domain registered less than 60 days ago | Wait until 60 days have passed since registration |
| Transferred within last 60 days | Wait until 60 days have passed since the previous transfer |
| Domain expired | Renew at current registrar first, then transfer |
| Admin contact didn't approve | Check email and spam folder. Ask Arcline support to resend |
| WHOIS privacy blocking | Disable privacy protection at current registrar |

Contact Arcline support if a transfer fails unexpectedly — we can check the exact reason and help you fix it.

---

## Transfer vs. just changing nameservers

You don't need to transfer your domain to host at Arcline. The alternative is simply updating your nameservers to `ns1.arcline.it` and `ns2.arcline.it` at your current registrar (see [Point Your Domain to Arcline](/getting-started/nameservers/)).

**Transfer when:**
- You want all billing in one place (Arcline)
- Your current registrar has raised prices or has bad support
- You're consolidating after an acquisition (Google Domains → Squarespace, etc.)

**Just change nameservers when:**
- You like your current registrar
- The domain isn't eligible for transfer yet (within 60-day lock)
- You want to keep registration and hosting with different providers

