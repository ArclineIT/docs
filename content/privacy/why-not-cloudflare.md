---
title: "Why You Shouldn't Put Cloudflare in Front of Everything"
description: "The hidden costs of relying on Cloudflare: privacy concerns, vendor lock-in, and single points of failure."
section: privacy
order: 1
---

# Why You Shouldn't Put Cloudflare in Front of Everything

Cloudflare is popular — it's fast, free, and easy to set up. But the convenience comes with tradeoffs that matter, especially for a hosting provider that prides itself on self-hosting and privacy.

---

## The single point of failure

When you put Cloudflare in front of your origin server, Cloudflare becomes your infrastructure. If Cloudflare has an outage (and they've had several), your site goes down even if your server is perfectly healthy.

In 2022, a Cloudflare global outage took down 20% of the internet's top million sites. Customers couldn't access their own servers because DNS wasn't resolving. If you self-host on Arcline, your DNS and traffic should be under your control.

---

## Cloudflare terminates TLS

This is the most important point for privacy-conscious site owners. When you use Cloudflare's proxy (orange cloud), Cloudflare terminates the TLS connection. Traffic between Cloudflare and your origin server travels in plaintext unless you set up a separate origin certificate.

This means Cloudflare can:
- Read all traffic passing through their network
- Inject JavaScript into your pages (which they do for features like Rocket Loader and bot management)
- Modify your HTML, images, and CSS
- See every visitor's IP address, browser, and behavior

If you're hosting anything sensitive — even a personal blog — you're trusting Cloudflare with all of it.

---

## Vendor lock-in

Once you're deep in the Cloudflare ecosystem, leaving becomes painful:

- **DNS records**: You manage DNS inside Cloudflare, not on your authoritative nameservers
- **Page Rules**: Caching, redirects, and rewriting rules are configured in Cloudflare's interface
- **SSL certificates**: Cloudflare issues and manages certs for your domains
- **Workers**: Serverless functions that only run on Cloudflare's network
- **Analytics**: Historical data lives in Cloudflare, not on your servers
- **Email routing**: Cloudflare's email forwarding is a sticky dependency

Each of these features makes Cloudflare harder to walk away from. The free tier becomes a golden handcuff.

---

## Privacy implications

Cloudflare sits between your visitors and your server. This means:

- Every visitor's IP address is visible to Cloudflare (they claim not to sell it, but you're still sharing it)
- Cloudflare can build behavioral profiles across the millions of sites on their network
- Visitors on Cloudflare see "This site is powered by Cloudflare" banners and CAPTCHAs
- The more sites that use Cloudflare, the more centralized the internet becomes

For a hosting company that markets itself on privacy and independence, directing customers to Cloudflare undermines the entire value proposition.

---

## When Cloudflare makes sense

There are legitimate use cases:

- **DDoS protection**: Cloudflare's network can absorb attacks that would overwhelm a single server
- **Global CDN**: If most of your visitors are on another continent, Cloudflare's edge cache speeds up delivery
- **API protection**: Rate limiting and bot detection at the edge

But these should be deliberate choices for specific problems, not a default for every site.

---

## The approach we recommend

Arcline is built on the principle that you should control your own infrastructure:

- DNS is answered by nameservers you choose and can audit
- TLS terminates at your server — no third party sees your traffic
- Static sites are served directly by your web server, without an intermediary
- Monitoring and analytics are self-hosted (for example, Prometheus and Wazuh)

This means you don't get Cloudflare's global CDN or DDoS scrubbing by default. What you get is:

- Full control over your traffic
- No third party reading your visitors' data
- No vendor lock-in
- Complete architectural flexibility

For most small to medium sites, a single well-tuned server with a good caching strategy and [fail2ban](/security/fail2ban/) for security is more than adequate.

---

## What's next

- [How to check if a host is self-hosted](/privacy/arcline-check-walkthrough/) — use arcline-check
- [What SPF, DKIM, and DMARC actually do](/privacy/spf-dkim-dmarc/)

