---
title: "How to Open a Support Ticket"
description: "How to contact Arcline support, what information to include, and expected response times."
section: reference
order: 5
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
   - **Technical** — server issues, website problems, email setup
   - **Migration** — moving from another provider to Arcline
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
| **What you're trying to do** | "I'm trying to install WordPress on my VPS" |
| **What's happening** | "I get a 500 error when visiting example.com/wp-admin" |
| **What you've already tried** | "I've checked that PHP is running and the database credentials are correct" |
| **Relevant logs or errors** | Paste the exact error message, not a summary |
| **Recent changes** | "I updated WordPress to 6.7 yesterday" |

---

## Response times

| Priority | Response time | When to use |
|----------|--------------|-------------|
| **Critical** | < 1 hour | Site is down, email is not working, security issue |
| **High** | < 4 hours | Site is slow but accessible, email is delayed |
| **Normal** | < 24 hours | Configuration help, billing questions, migration support |
| **Low** | < 72 hours | Feature requests, general questions, feedback |

Response times are based on Eastern Time (ET), Monday through Friday. Weekend and holiday responses may take longer.

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

- [Arcline nameservers and DNS records](/reference/nameservers/)
- [Supported PHP versions](/reference/php-versions/)
- [Resource limits by plan](/reference/resource-limits/)

