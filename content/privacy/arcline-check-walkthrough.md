---
title: "How to Check if Your Host is Self-Hosted"
description: "Use arcline-check to detect whether a website is truly self-hosted or routed through a CDN/cloud provider."
section: privacy
order: 2
---

# How to Check if Your Host is Self-Hosted

The `arcline-check` tool tells you whether a domain is truly self-hosted or routing through Cloudflare, Fastly, AWS CloudFront, or another CDN/cloud provider.

---

## What arcline-check does

When you run `arcline-check example.com`, it:

1. Resolves the domain to its IP address
2. Performs a reverse DNS lookup (PTR record)
3. Looks up the ASN (Autonomous System Number) and organization
4. Checks whether the IP falls within known CDN/cloud provider CIDR ranges
5. Inspects HTTP response headers for CDN signatures (CF-Ray, X-Served-By, etc.)
6. Produces a color-coded terminal report

This is especially useful for:
- Evaluating potential hosting providers during migration
- Verifying a host's claims about being "self-hosted"
- Understanding your own site's infrastructure

---

## Installation

```bash
# Download the latest binary
wget https://git.arcline.it/arcline/arcline-check/releases/latest/download/arcline-check-linux-amd64

# Make it executable
chmod +x arcline-check-linux-amd64

# Move to your PATH
sudo mv arcline-check-linux-amd64 /usr/local/bin/arcline-check
```

Or build from source:

```bash
git clone https://git.arcline.it/arcline/arcline-check.git
cd arcline-check
go build -o arcline-check .
```

---

## Basic usage

```bash
arcline-check arcline.it
```

Example output:

```
  domain      arcline.it
  resolved    203.0.113.42
  rdns        server1.arclineit.com
  asn         AS64496  Example ISP
  org         Example ISP LLC

  [OK]  not behind a known CDN
  [OK]  no Cloudflare headers detected
  [OK]  IP not in AWS/GCP/Azure ranges
```

---

## Checking a site behind Cloudflare

```bash
arcline-check example-cloudflare-site.com
```

Example output:

```
  domain      example-cloudflare-site.com
  resolved    104.16.42.42
  rdns        104.16.42.42  (no PTR)
  asn         AS13335  Cloudflare, Inc.
  org         Cloudflare

  [FAIL]  behind Cloudflare (CDN)
  [FAIL]  CF-Ray header detected
  [FAIL]  IP in Cloudflare CIDR range
```

---

## JSON output for scripting

```bash
arcline-check example.com --json
```

```json
{
  "domain": "example.com",
  "ip": "203.0.113.42",
  "rdns": "server1.arclineit.com",
  "asn": "AS64496",
  "org": "Example ISP LLC",
  "cdn_detected": false,
  "headers": {
    "server": "nginx/1.24.0"
  }
}
```

---

## Watch mode for DNS migration monitoring

During a DNS migration, use `--watch` to see when propagation completes:

```bash
arcline-check example.com --watch 30
```

This re-checks every 30 seconds. When the IP changes from the old provider to the new one, the output updates in place.

---

## Interpreting the results

| Indicator | What it means |
|-----------|---------------|
| **IP in Cloudflare range** | The site is behind Cloudflare's proxy (orange cloud) |
| **CF-Ray header** | Cloudflare is terminating the connection |
| **IP in AWS/GCP/Azure range** | The server is a cloud VM, not self-hosted hardware |
| **IP in a residential/business ISP range** | Likely self-hosted (on-premises or colocated) |
| **PTR matches domain** | Good operational practice — the host configured rDNS |
| **No CDN detected** | Traffic goes directly to the origin server |

---

## Limitations

- A CDN-detected result doesn't always mean bad hosting — some providers use CDNs for legitimate DDoS protection
- Arcline doesn't use any CDN by default, but customers are free to add one if they choose
- The tool can't detect every possible CDN or proxy — new providers are added regularly
- If a site uses Cloudflare spectrum or TCP tunnels, it may appear self-hosted even though Cloudflare is involved

---

## What's next

- [Self-hosting without a CDN: performance tips](/privacy/self-hosting-performance/)
- [Why you shouldn't put Cloudflare in front of everything](/privacy/why-not-cloudflare/)
- [What SPF, DKIM, and DMARC actually do](/privacy/spf-dkim-dmarc/)

