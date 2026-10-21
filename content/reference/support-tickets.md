---
title: "How to Open a Support Ticket"
description: "How to contact Arcline support, what information to include, and expected response times."
section: reference
order: 1
---

# How to Open a Support Ticket

When you need help, opening a support ticket is the best way to get a response. This page explains how tickets work and what to include.

---

## Before opening a ticket

Check these resources first — your question may already be answered:

- **Knowledge base**: Browse the guides in this docs site
- **Status page**: Check [status.arcline.it](https://status.arcline.it) for any active incidents or scheduled maintenance
- **FAQ**: Visit [arcline.it/faq](https://arcline.it/faq) for common questions

---

## Opening a ticket

### Via the customer portal (recommended)

1. Log in to the [Arcline Customer Portal](https://portal.arcline.it)
2. Click **Support → New Ticket**
3. Select the relevant department:
   - **Billing** — invoices, plan changes, cancellations
   - **Technical** — alerts, firewall and remote access problems, backups and restores, anything in a system we built or monitor
4. Fill in the subject and description
5. Include as much detail as possible (see below)
6. Click **Submit**

You'll receive a confirmation email with the ticket number.

### Via email

Send an email to **support@arcline.it**. Include:

- Your account username or email
- A clear subject line
- A detailed description of the issue
- Any relevant screenshots or error messages

Your email will automatically create a ticket in the system and you'll receive a confirmation reply.

---

## What to include in a ticket

The more detail you provide, the faster we can help. Include:

| Detail | Example |
|--------|---------|
| **Your domain** | `example.com` |
| **What you're trying to do** | "I'm trying to connect to the office VPN from a new laptop" |
| **What's happening** | "WireGuard says the handshake never completes" |
| **What you've already tried** | "I re-imported the config and restarted the app" |
| **Relevant logs or errors** | Paste the exact error message, not a summary |
| **Recent changes** | "We switched internet providers yesterday" |

---

## Response times

How quickly we respond depends on your engagement:

| Engagement | Response time |
|------------|---------------|
| **Project** | Within 24 hours, for the length of the project |
| **Retainer** | 4–8 hours, during business hours and for emergencies |
| **Managed** | 1–2 hours, around the clock |

Business hours are Monday through Friday, Central Time. Your written agreement sets the exact terms. If you suspect a security incident, say so in the subject line.

---

## Ticket statuses

| Status | Meaning |
|--------|---------|
| **Open** | Received but not yet assigned |
| **In Progress** | Being worked on by a technician |
| **Awaiting Customer** | We need more information from you |
| **Resolved** | Issue has been resolved |
| **Closed** | Issue is complete and ticket is archived |

If a ticket is marked **Awaiting Customer**, please respond as soon as possible. Tickets that remain in this status for more than 7 days will be automatically closed.

---

## Providing SSH access

If your issue requires investigation on the server, you may be asked to provide SSH access. You have two options:

- **Temporary key**: Upload an SSH key and remove it after the issue is resolved
- **Add to known IPs**: We'll provide the IP we'll connect from — add it to your firewall's whitelist

Never share your account password. Arcline staff will never ask for your password.

---

## What's next

- [What SPF, DKIM, and DMARC actually do](/privacy/spf-dkim-dmarc/)
- [How to check if a host is self-hosted](/privacy/arcline-check-walkthrough/)

