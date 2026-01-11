# Mondrian Retrospective: Why This Project is Archived

## The Vision (September 2025)

"PostHog for Zero Trust" - make enterprise security accessible to startups without $50K+ budgets. Built as a validation exercise during Startup Virginia's Idea Factory program.

The core insight: there's a quantifiable "security poverty line" (coined by Wendy Nather at Cisco in 2011) where startups below $200K in available compliance capital simply cannot access enterprise customers, regardless of product quality.

## What We Built (2 days)

- **Working CLI** with 3 security rules (S3 public buckets, open security groups, OIDC vs long-lived credentials)
- **GitHub Actions integration** that blocks risky PRs before deployment
- **DSSE attestation system** with tamper-evident signature chains
- **Complete evidence generation** (`internal/evidence/` package) with attestation, hash chains, and cryptographic signing
- Positioned as **SBIR/STTR research project** on the "security poverty line" economic barrier

The prototype proved the technical approach works: you can generate cryptographically signed, tamper-evident proof that security controls ran.

## What We Learned

### 1. Market Validation: The Problem is Real

**The "security poverty line"** (coined by Wendy Nather, Cisco, 2011) isn't abstract - it's quantifiable:

- **29% of organizations lose enterprise deals** due to missing compliance certifications (A-LIGN 2023 Benchmark Report)
- **$200K-400K first-year compliance costs** create an economic barrier:
  - SOC 2 certification: $35K-150K initially, $12K-60K annually
  - Internal engineering effort: 100-200 hours = ~$50K opportunity cost
  - Insurance premium penalty: $45K/year higher for uncompliant startups
- **96% of GRC teams struggle** with keeping up with regulations (Swimlane study)
- **41% without continuous compliance** report sales cycle slowdowns (Drata 2025)
- **60% of SME cyber victims** go out of business within 6 months

I experienced this firsthand at DeepWaterPoint & Associates. Compliance theatre is real - teams spend enormous time collecting evidence that doesn't meaningfully reduce risk.

### 2. Technical Validation: Evidence-First Works

The technical approach is sound:

- **Cryptographic attestations are feasible**: DSSE signatures + hash chains create tamper-evident audit trails
- **GitHub Actions integration works**: Can prove controls ran, not just detect violations
- **SLSA framework provides standards**: Supply chain attestation formats already exist
- **Precedents exist**: DocuSign (legal signatures), GitHub (signed commits), AWS CloudTrail (crypto-verified logs)

**BUT**: Vanta and Drata already solve 80% of evidence collection through SaaS integrations. The remaining 20% (cryptographic proof) requires answering: "Do auditors actually accept crypto attestations over screenshots?"

Customer validation showed: Technical feasibility ≠ business viability. Differentiation requires ongoing R&D investment, not just a weekend prototype.

### 3. Business Reality: Strengths Mismatch

This is where honest self-assessment mattered:

**Customer discovery paralysis:**
- Accepted to Startup Virginia's Idea Factory (pre-incubator program)
- Needed to interview 50+ security leaders to validate market demand
- Actually completed: ~0 interviews
- Struggled with imposter syndrome about reaching out with the idea
- The validation framework required activities outside my comfort zone

**Advisory/content model is better ROI:**
- My strengths: Translation, synthesis, fast prototyping, strategic thinking
- Not my strengths: Community building, ongoing product maintenance, sales outreach, sustained customer development
- Advisory work (Eagle Ridge, Caldris) plays to strengths without OSS maintenance burden
- Content creates leveraged impact - blog posts reach more people than a GitHub repo

**Claude Code augments capability, but has limits:**
- AI tools dramatically accelerate prototyping (4 hours to working CLI)
- But they don't replace customer discovery, community management, or strategic positioning
- I'm still early in my development experience - building one prototype ≠ maintaining an OSS ecosystem

## Why It's Archived (January 2026)

**Immediate trigger:** Shift in family priorities (health situation) created bandwidth constraints.

**Deeper reasons:**
- **Advisory/content positioning** creates more impact than OSS product for my skill set
- **SBIR/STTR research path** requires 12-18 month commitment I don't have
- **Vanta/Drata** already lowered the barrier from $400K to $200K - remaining gap is customer willingness to pay vs technical feasibility
- **Core insights** more valuable as blog content (reaches broader audience, establishes thought leadership)

## What Good Assumptions Were Made

✅ **Technical feasibility** - Evidence-first approach works, prototype proves it
✅ **Market framing** - "Security poverty line" resonates with founders and security leaders
✅ **Fast validation** - 4 hours proved concept viability without months of planning
✅ **Documentation** - Comprehensive project artifacts (PROJECT_REPORT.md, LAUNCH_KIT.md) preserved value
✅ **Platform choice** - GitHub Actions integration was right wedge (developers already use it)

## What Was Missing

❌ **Customer validation** - Needed 50+ interviews, completed ~0. Can't validate market demand without talking to buyers
❌ **Go-to-market fit** - Community building and developer evangelism aren't my strengths
❌ **Competitive moat analysis** - Should have asked: "Why haven't Vanta/Drata added crypto attestations if valuable?"
❌ **Honest ROI analysis upfront** - Built first, analyzed later. Should have reversed this order
❌ **Auditor validation** - Core question "Do auditors accept crypto signatures?" remained unanswered

## What I'd Do Differently

1. **Customer discovery FIRST** - Talk to 20 security leaders before writing any code. Validate demand, not just technical feasibility.

2. **Validate differentiation** - Specifically answer: "Why can't Vanta/Drata add cryptographic attestations?" If it's valuable, why hasn't it been done?

3. **Match work to strengths** - Consulting/content creation aligns with my translation and synthesis skills. OSS community-building doesn't.

4. **Be honest about bandwidth** - SBIR research requires 12-18 month focused commitment. Advisory model gives flexibility for life situations.

5. **Test the "why" before building the "what"** - Economic analysis (A16Z margin impact, compliance scaling trap) should precede technical validation.

## For Future Explorers

This repository demonstrates several valuable patterns:

**Fast MVP validation (4 hours to working prototype):**
- Cobra CLI framework for Go command structure
- Hard-coded policy rules for rapid iteration
- DSSE signatures for cryptographic proof
- GitHub Actions integration for real-world testing

**Evidence-first Zero Trust architecture (SLSA-compatible):**
- `internal/evidence/attestation.go` - SLSA attestation format
- `internal/evidence/chain.go` - Hash-chain integrity tracking
- `internal/evidence/signer.go` - DSSE signature generation
- Policy engine → attestation → signature → verification chain

**Build/buy/partner decision framework for CISOs:**
- When to prototype (validate technical feasibility)
- When to partner (Vanta/Drata solve 80%)
- When to advise (translation skills > product maintenance)

**When to archive instead of launch:**
- Technical validation ≠ business viability
- Customer discovery paralysis is a signal (wrong skill/passion fit)
- Advisory/content can extract value without product liability
- Honest assessment of strengths prevents sunk cost fallacy

## Related Content

These blog posts extract the core insights from this research project:

- [Why I Built and Shelved a Zero Trust Tool](#) - Eagle Ridge blog (decision framework for CISOs)
- [The Security Poverty Line](#) - Eagle Ridge blog (economic analysis with quantified data)
- [Evidence-First vs Detection-Only Security](#) - Eagle Ridge blog (technical deep-dive with code examples)

## For Recruiters/Advisory Prospects

This project demonstrates:

✅ **Fast prototyping** - 4 hours from concept to working GitHub Actions integration
✅ **Technical depth** - Implemented SLSA attestations, DSSE signatures, cryptographic verification
✅ **Market research** - Quantified $200K-400K compliance barrier with cited sources
✅ **Honest assessment** - Recognized skill mismatch and pivoted to advisory model
✅ **Strategic thinking** - Build/buy/partner analysis, not just "build everything"

**What this means for clients:** I can validate technical approaches quickly, understand compliance/security economics deeply, and make honest build vs buy recommendations. The value isn't in maintaining an OSS tool - it's in helping you make the right decision faster.

---

**Status:** Archived January 2026 as completed research/validation project.

**GitHub:** https://github.com/miqcie/mondrian (read-only archive)

**Contact:** For advisory/consulting inquiries about compliance automation, evidence-first security, or build/buy decisions in the security space, reach out through Eagle Ridge Advisory or Caldris.
