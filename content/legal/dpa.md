# Data Processing Agreement (DPA)

**Arcline IT LLC**
Last updated: May 2026

This Data Processing Agreement ("DPA") forms part of the Master Service
Agreement ("MSA") between Arcline IT LLC ("Data Processor", "Arcline", "we",
"us") and the Customer ("Data Controller", "you").

---

## 1. Definitions

| Term | Definition |
|------|------------|
| **Controller** | The entity that determines the purposes and means of processing personal data |
| **Processor** | The entity that processes personal data on behalf of the Controller |
| **Data Subject** | An identified or identifiable natural person |
| **Personal Data** | Any information relating to an identified or identifiable natural person |
| **Processing** | Any operation performed on personal data (collection, storage, retrieval, transmission, deletion, etc.) |
| **GDPR** | Regulation (EU) 2016/679, the General Data Protection Regulation |
| **CCPA** | California Consumer Privacy Act, as amended |
| **Sub-processor** | A third party engaged by the Processor to process personal data |

---

## 2. Scope and Purpose

### 2.1 Application
This DPA applies whenever Arcline processes personal data on behalf of
Customer in the course of providing Services under the MSA.

### 2.2 Relationship
- **Customer** is the Data Controller
- **Arcline** is the Data Processor
- Customer retains full control over their personal data
- Arcline processes data only on Customer's documented instructions

### 2.3 Duration
This DPA remains in effect for as long as Arcline processes personal data on
behalf of Customer, plus the duration of any data retention obligations.

---

## 3. Description of Processing

### 3.1 Categories of Data Subjects
- Customer's employees, contractors, and agents
- Customer's end users and website visitors
- Individuals who communicate with Customer through their Arcline-hosted services

### 3.2 Categories of Personal Data
- Account information: name, email address, billing address, phone number
- Technical data: IP addresses, server access logs, browser user-agent strings
- Content data: files, databases, emails, and other content stored on Arcline
  infrastructure at Customer's direction
- Payment data: processed through Stripe (PCI-DSS compliant); Arcline does
  not store full credit card numbers

### 3.3 Special Categories of Data
Arcline does not intentionally process special categories of data (health
information, biometric data, political opinions, religious beliefs, etc.).
Customer agrees not to upload special category data to Arcline infrastructure
without additional contractual safeguards.

### 3.4 Processing Activities
- **Storage:** Customer data stored on Arcline's servers
- **Hosting:** Serving Customer's websites and applications to visitors
- **Backup:** Creating and maintaining backup copies for disaster recovery
- **Email:** Routing and storing email messages (where applicable)
- **Support:** Accessing data for troubleshooting and support purposes

---

## 4. Processor Obligations

### 4.1 Instructions
Arcline will process personal data only on documented instructions from
Customer, unless required to do otherwise by applicable law (in which case
Arcline will notify Customer of that legal requirement before processing,
unless prohibited by law).

### 4.2 Confidentiality
Arcline ensures that all personnel authorized to process personal data have
committed to confidentiality obligations.

### 4.3 Security
Arcline maintains appropriate technical and organizational security measures,
including:

**Technical Measures:**
- Encryption in transit (TLS 1.2+ for all services)
- Firewalls with default-deny rules
- Network segmentation (VLANs)
- Regular security patching
- Intrusion detection and prevention systems (Suricata)
- Access logging and monitoring
- Encrypted off-site backups

**Organizational Measures:**
- Access control based on least privilege
- Security training for personnel
- Incident response procedures
- Regular security assessments
- Vendor due diligence for sub-processors

### 4.4 Sub-processors
Customer authorizes Arcline to engage the following sub-processors:

| Sub-processor | Service | Location |
|---------------|---------|----------|
| Stripe, Inc. | Payment processing | United States |
| Let's Encrypt / ISRG | SSL certificate issuance | United States |
| GitLab B.V. | CI/CD and source control | United States/Europe |

Arcline will notify Customer of any intended changes to sub-processors.
Customer may object within 14 days. If reasonable objections cannot be
resolved, Customer may terminate the affected services.

### 4.5 Data Subject Rights
Arcline will assist Customer in responding to data subject requests under
applicable privacy laws, including:
- Right of access
- Right to rectification
- Right to erasure ("right to be forgotten")
- Right to restrict processing
- Right to data portability
- Right to object

Customer should forward any data subject requests they receive to
**privacy@arcline.it**. Arcline will respond within the timeframe required
by applicable law.

### 4.6 Data Breach Notification
In the event of a personal data breach, Arcline will:
1. Notify Customer within 72 hours of becoming aware of the breach
2. Provide details of the nature, scope, and impact of the breach
3. Describe measures taken to address the breach
4. Cooperate with Customer in notifying supervisory authorities and affected
   data subjects, where required

### 4.7 Data Protection Impact Assessments
Arcline will provide reasonable assistance to Customer in conducting data
protection impact assessments, where required by applicable law.

---

## 5. International Transfers

### 5.1 Data Location
Customer data is primarily stored on servers located in the United States.

### 5.2 Adequacy
For transfers of personal data from the European Economic Area (EEA),
Switzerland, or the United Kingdom to the United States, the parties agree
that the Standard Contractual Clauses (SCCs) approved by the European
Commission shall govern such transfers.

### 5.3 Alternative Mechanism
If the SCCs are deemed invalid or insufficient by a competent authority,
Arcline will implement an alternative lawful transfer mechanism.

---

## 6. Data Retention and Deletion

### 6.1 During the Term
Customer data is retained for the duration of the MSA or until Customer
requests deletion.

### 6.2 Upon Termination
Following termination of the MSA:
- **Active data:** Deleted within 30 days of termination
- **Backups:** Deleted within 90 days of termination
- **Access logs:** Anonymized or deleted within 12 months

### 6.3 Deletion Procedures
Data is securely deleted using:
- Secure file deletion (shred/overwrite) for files
- `DROP TABLE` for SQLite databases
- Cryptographic erasure for encrypted backups

### 6.4 Certificate of Deletion
Upon request, Arcline will provide a written certificate confirming that
Customer's data has been securely deleted.

---

## 7. Audit and Compliance

### 7.1 Right to Audit
Customer may request an audit of Arcline's data processing operations, at
Customer's expense, no more than once per 12-month period. Audits must:
- Be conducted during normal business hours
- Give at least 30 days notice
- Not unreasonably interfere with Arcline's operations
- Be performed by a mutually agreed independent auditor

### 7.2 Records of Processing
Arcline maintains written records of all processing activities conducted on
behalf of Customer, as required by Article 30 of the GDPR.

### 7.3 Compliance
Arcline will promptly notify Customer if any instruction from Customer
violates applicable data protection laws.

---

## 8. Liability

### 8.1 Liability Cap
Each party's liability under this DPA is subject to the limitations of
liability set forth in the MSA.

### 8.2 Direct Damages
Notwithstanding the general limitation above, either party may seek direct
damages for breaches of this DPA.

### 8.3 Regulatory Fines
Each party is responsible for administrative fines imposed on them by a
supervisory authority for their own violations of applicable data protection
law.

---

## 9. Governing Law

This DPA shall be governed by the same law as the MSA. Any dispute arising
from this DPA shall be resolved under the dispute resolution provisions of
the MSA.

---

## 10. Order of Precedence

In the event of any conflict or inconsistency between this DPA and the MSA,
this DPA shall prevail with respect to data processing matters.

---

*Questions about this DPA? Contact us at privacy@arcline.it*

