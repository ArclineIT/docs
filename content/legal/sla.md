# Service Level Agreement (SLA)

**Arcline IT LLC**
Last updated: May 2026
Applies to: All Arcline shared hosting, WordPress hosting, and VPS hosting services.

---

## 1. Overview

This Service Level Agreement ("SLA") governs the availability and performance
of services provided by Arcline IT LLC ("Arcline", "we", "us") to you ("Customer").
This SLA is incorporated by reference into the Arcline Terms of Service and
Master Service Agreement.

Arcline operates self-hosted infrastructure on owned hardware. We do not use
AWS, Azure, GCP, Cloudflare, or any hyperscale cloud provider. We are
transparent about what this means for availability, and we believe honesty
about our capabilities serves our customers better than inflated promises.

---

## 2. Uptime Commitment

| Service Tier | Monthly Uptime Target | Annual Uptime Target |
|-------------|----------------------|---------------------|
| Shared Hosting | 99.0% | 98.5% |
| WordPress Hosting | 99.0% | 98.5% |
| VPS Hosting | 99.5% | 99.0% |
| Network (edge) | 99.5% | 99.0% |

**Calculated as:**
```
Uptime % = (Total Minutes in Month − Downtime Minutes) ÷ Total Minutes × 100
```

---

## 3. Service Credits

If we fail to meet the uptime target in a given calendar month, you may request
a service credit. Credits are applied to your next invoice and are not
redeemable for cash.

| Monthly Uptime | Credit |
|----------------|--------|
| Below target but ≥ 95% | 10% of monthly fee |
| 90% – 94.99% | 25% of monthly fee |
| Below 90% | 50% of monthly fee |

### How to Request a Credit

1. Email **support@arcline.it** with subject line "SLA Credit Request"
2. Include your account email and the month in question
3. We will verify our monitoring data and respond within 5 business days

Credits must be requested within 30 days of the end of the month in which
the downtime occurred.

**Limitation:** Total credits in any single month shall not exceed the
Customer's monthly service fee for the affected service.

---

## 4. Exclusions

The following are excluded from SLA calculations and are not eligible for
service credits:

### Scheduled Maintenance
- Announced maintenance windows with at least 48 hours notice
- Emergency security patches (notice as circumstances permit)
- We schedule maintenance during off-peak hours (midnight–6 AM ET) when possible

### Customer-Caused Downtime
- Configuration errors by the Customer
- Exceeding resource limits (CPU, RAM, disk I/O, bandwidth)
- Failure to maintain payment on the account (past-due accounts)
- Actions taken by the Customer that trigger abuse prevention mechanisms

### Force Majeure
- Natural disasters, war, civil unrest, pandemic, or other events beyond
  reasonable control
- Upstream network outages (BGP issues, transit provider failures) — we will
  work to mitigate but do not guarantee alternate routing
- Power outages at our facility — we maintain UPS and generator backup but
  do not guarantee 100% uptime

### Third-Party Services
- DNS propagation delays when changing nameservers
- SSL certificate issuance delays by Let's Encrypt or other CAs
- Email deliverability issues caused by receiver-side filtering

### VPS-Specific Exclusions
- Guest operating system crashes or misconfiguration
- Resource exhaustion caused by the Customer's workload
- Actions taken by the Customer's users

---

## 5. Monitoring

Arcline measures uptime using our internal monitoring system (`arcline-uptime`),
which performs checks from multiple locations at 1-minute intervals. Monitoring
results are available at [status.arclineit.com](https://status.arclineit.com).

We monitor:
- **HTTP services** — connection success, HTTP 200 status, response within 10s
- **TCP services** — successful TCP connection to service port
- **Network** — ping response from edge router
- **Infrastructure** — system load, disk usage, temperature sensors

In the event of a disagreement about uptime, Arcline's monitoring data shall
be the primary source. Customers are encouraged to run independent monitoring
and may submit their own monitoring data for consideration.

---

## 6. Maintenance Windows

### Routine Maintenance
- Typically performed Tuesday–Thursday between midnight and 6 AM ET
- Announced on [status.arclineit.com](https://status.arclineit.com) at least
  48 hours in advance
- Brief service interruptions (< 15 minutes) for routine updates

### Emergency Maintenance
- Security vulnerabilities (CVSS ≥ 7.0): patched within 24 hours
- Critical hardware failures: immediate intervention
- Notice provided via status page and email as time permits

---

## 7. Incident Response

| Severity | Definition | Initial Response | Update Frequency |
|----------|-----------|-----------------|------------------|
| **Critical** | Service unavailable, all customers affected | 30 minutes | Every 60 minutes |
| **Major** | Service degraded or partially unavailable | 1 hour | Every 2 hours |
| **Minor** | Isolated issue affecting few customers | 2 hours | Every 4 hours |
| **Maintenance** | Planned work with brief interruption | Per schedule | Per schedule |

Response times are measured from the time Arcline becomes aware of the issue
(automated alert or customer report), not from the start of the incident.

---

## 8. Support Response Times

| Priority | Channel | Target Response |
|----------|---------|-----------------|
| Emergency (service down) | Ticket + email | 30 minutes (business hours) |
| High (degraded service) | Ticket | 2 hours (business hours) |
| Normal (general inquiry) | Ticket | 24 hours |
| Low (feature request) | Ticket | 48 hours |

**Business hours:** Monday–Friday, 9 AM–6 PM ET.
**After hours:** Best-effort for Critical and Major incidents only.

---

## 9. Data Preservation

In the event of account suspension or termination:
- **Shared/WordPress hosting:** Data preserved for 14 days after suspension
- **VPS hosting:** Data preserved for 7 days after suspension
- **Final deletion:** Data is securely erased after the preservation period

We do not provide data recovery for accounts that have been deleted for more
than 30 days.

---

## 10. SLA Exceptions and Changes

Arcline reserves the right to modify this SLA with 30 days written notice to
active customers. Material changes will be emailed and posted to our status
page.

---

*Questions about this SLA? Contact us at [arcline.it/contact](https://arcline.it/contact)*

