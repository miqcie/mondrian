# Mondrian Launch Kit
*Copy-paste ready messages for community launch*

---

## 🟢 WEEK 1: Soft Launch (Technical Communities)

### Monday: r/golang Post

**Title:** [Show r/golang] Built a Zero Trust CLI in Go - weekend project, looking for feedback

**Body:**
Hey r/golang!

Just shipped a weekend project I think this community might find interesting. It's a Zero Trust policy enforcement CLI built with Cobra.

**The Problem:** Startups need security but enterprise Zero Trust tools cost $50K+. So they ship risky configs, fail compliance, and lose enterprise deals.

**The Solution:** OSS Zero Trust CLI that blocks risky changes and generates cryptographic proof that your controls actually ran.

**What's built:**
✅ Policy engine with 3 security rules (S3 buckets, security groups, OIDC)
✅ GitHub Action integration for blocking bad PRs  
✅ SLSA-compatible attestation generation
✅ Hash-chain evidence integrity with DSSE signatures

**Tech stack:** Go + Cobra CLI + ECDSA signing + in-toto attestations

Repository: https://github.com/miqcie/mondrian

The CLI is clean, fast, and the attestation system is pretty elegant. This is my first serious Go project - curious what the community thinks about the architecture.

Also wondering: what security rules would be most valuable for Go teams? Database credentials in code? Missing TLS configs?

Built this because I believe security tools should help engineers ship faster, not slower.

P.S. - This CLI is the first piece of a bigger "Zero Trust OS for developers" vision, but wanted to start with something that works today.

---

### Tuesday: Dev.to Article

**Title:** Building an Evidence-First Zero Trust Platform in Go - Weekend Project

**Intro:**
Enterprise Zero Trust platforms cost $50K+. Startups need security but can't afford the luxury tax. So I built an alternative in a weekend using Go.

This is the story of Mondrian - an OSS Zero Trust CLI that proves your security controls actually ran.

**[Full article structure - expand with technical details, code snippets, architecture decisions]**

---

### Wednesday: Twitter Thread

🧵 1/8 Hot take: Zero Trust has a startup tax problem.

Enterprise vendors charge $50K+ because they can. Startups pay it because they must (compliance, customers, investors).

There's a better way. Built it this weekend. 🧵

🧵 2/8 Meet Mondrian: "PostHog for Zero Trust"

The insight: Security tools should make engineers more productive, not less.

Block risky changes ✅
Generate compliance proof ✅  
Cost $0 instead of $50K ✅

github.com/miqcie/mondrian

🧵 3/8 The economics are broken in security.

Vendors optimize for buyers (CISOs with budgets) not users (engineers shipping code).

Result: Tools that satisfy procurement but frustrate practitioners.

🧵 4/8 Mondrian does three things well:
1. Catches problems before production (S3 buckets, security groups, OIDC)
2. Plugs into GitHub Actions (because you have enough dashboards)  
3. Generates cryptographic proof (SLSA attestations + hash chains)

Built in Go because CLIs should be fast.

🧵 5/8 The "weekend project" economics:
- Day 1: Core rules + GitHub integration  
- Day 2: DSSE signatures + attestation chains
- Day 3: This launch thread

Total investment: ~20 hours
Enterprise alternative: $50K+ per year

🧵 6/8 Inspired by @slsa_framework @in_toto_io @sigstore but focused on runtime enforcement.

Think: OPA meets GitHub Actions with tamper-evident receipts.

The kind of tool engineers would build for engineers.

🧵 7/8 Early feedback welcome:
- What security rules would save you actual time?
- Is "evidence-first" the right framing?
- Would your startup use this?

First OSS project. Learning in public. 📚

🧵 8/8 The goal: Prove that security tools can be both effective AND affordable.

Zero Trust shouldn't be a luxury good.

What do you think? Worth building or am I solving the wrong problem?

#ZeroTrust #OpenSource #StartupLife

---

### Thursday: Hacker News

**Title:** Show HN: Mondrian – Evidence-first Zero Trust that doesn't bankrupt startups

**Body:**
Weekend project alert: Built an alternative to $50K+ Zero Trust platforms.

The market failure is obvious: Enterprise security vendors optimize for procurement departments, not engineering teams. Result? Tools that cost more than junior developers and integrate about as well.

Mondrian flips this. Evidence-first Zero Trust that:
- Prevents risky deployments (the prevention is worth the cure)
- Generates tamper-proof compliance records (auditors need evidence, not promises)
- Costs nothing but your time to set up

Tech: Go CLI + GitHub Actions + SLSA attestations + hash-chains for integrity.

Two days from frustrated tweet to working prototype. The code is straightforward because the problem is straightforward: make security help instead of hurt.

Repository: https://github.com/miqcie/mondrian

Curious what the HN crowd thinks about this approach.

---

## 🟡 WEEK 2: Broader Launch (Market Communities)

### Monday: r/devops

**Title:** [Show r/devops] Zero Trust for startups that doesn't cost $50K - built it in a weekend

**Body:**
Fellow r/devops folks,

Here's the thing nobody talks about: Zero Trust has a startup tax. Teleport wants $50K+. Zscaler wants your firstborn. Meanwhile, you're trying to ship features.

So I built something different. Call it "PostHog for Zero Trust."

**What it does:**
- Blocks risky changes before they hit production (like a good bouncer)
- Proves your security controls actually ran (with cryptographic receipts)
- Plugs into your existing GitHub workflow (because you have enough tools)

**What it costs:**
Zero. It's OSS. Because security shouldn't be a luxury good.

**What's working today:**
✅ CLI catches S3 public buckets, open security groups, missing OIDC
✅ GitHub Action blocks bad PRs automatically
✅ SLSA-style attestations prove compliance (auditors love this)
✅ Evidence chains that can't be tampered with

Look, this took me 2 days to build. The hard part isn't the tech—it's remembering that security tools should help engineers, not burden them.

Repository: https://github.com/miqcie/mondrian

What am I missing? What rules would actually save you time? 

P.S. - Yes, it's written in Go. Because life's too short for slow CLIs.

---

### Tuesday: r/cybersecurity

**Title:** Built an OSS Zero Trust platform focused on evidence generation - looking for security expert feedback

**Body:**
Hey r/cybersecurity,

Security engineer here who got tired of expensive Zero Trust platforms that generate more theater than actual security.

Built Mondrian - an evidence-first Zero Trust platform that focuses on cryptographic proof over compliance dashboards.

**Key insight:** It's not enough to block risky changes. You need tamper-evident proof that security controls actually executed.

**What it does:**
- Policy enforcement at CI/CD time (catch problems early)
- SLSA-compatible attestation generation (industry standard)
- Hash-chain evidence integrity (Byzantine fault tolerant)
- GitHub Actions integration (fits existing workflows)

**Security architecture:**
- DSSE (Dead Simple Signing Envelope) signatures
- ECDSA with ephemeral keys (production will use OIDC)
- In-toto attestation format compatibility
- Offline verification capabilities

Repository: https://github.com/miqcie/mondrian

Curious what the security community thinks:
1. Is "evidence-first" the right approach vs traditional network-centric Zero Trust?
2. What attack vectors am I missing in the attestation chain?
3. Would this be useful for compliance audits in your experience?

This is targeting startups who need security but can't afford enterprise platforms. Trying to prove that good security doesn't have to be expensive.

---

### Wednesday: r/startups 

**Title:** Built an OSS alternative to $50K Zero Trust platforms - feedback from founders?

**Body:**
Fellow startup folks,

Anyone else frustrated by the "security tax" when trying to sell to enterprise customers?

Enterprise buyers want SOC2, Zero Trust, compliance proof. Fair enough. But the solutions cost $50K+ and take months to implement. Meanwhile, you're trying to get to product-market fit.

So I built an alternative: Mondrian - open-source Zero Trust that proves your security controls actually ran.

**The economics:**
- Enterprise Zero Trust: $50K+/year + months of integration
- Mondrian: $0 + weekend setup

**What it does:**
- Blocks risky infrastructure changes before deployment
- Generates cryptographic proof for auditors (no more screenshots)
- Integrates with GitHub Actions (works with your existing workflow)

**Why it matters:**
Every startup hits this wall. You move fast, then enterprise buyers ask for compliance proof. The current path is expensive and slow. This flips the script - security that helps you move faster, not slower.

Repository: https://github.com/miqcie/mondrian

Questions for the community:
1. Have you lost enterprise deals due to security/compliance requirements?
2. What's your current security stack cost?
3. Would something like this help you sell upmarket faster?

This is my first OSS project. Learning what actual startup founders need vs what vendors think they need.

---

### Thursday: LinkedIn Post

**Professional LinkedIn Post:**

🚀 Just shipped my first open-source project: Mondrian - Zero Trust security that doesn't cost $50K+

After seeing too many startups lose enterprise deals due to security requirements they can't afford to meet, I built an alternative.

**The Problem:**
• Enterprise Zero Trust platforms cost more than a junior developer's salary  
• Startups need compliance proof but can't afford enterprise tools
• Current solutions optimize for procurement, not engineering teams

**The Solution:**
• Evidence-first Zero Trust that blocks risky changes AND proves controls ran
• Cryptographic attestations replace compliance screenshots  
• GitHub Actions integration works with existing workflows
• Open source = $0 cost

**Early Results:**
✅ CLI with security policy enforcement  
✅ Tamper-evident compliance records
✅ SLSA-compatible attestation generation
✅ Weekend implementation timeline

This is what happens when you build security tools FOR engineers instead of procurement departments.

Repository: https://github.com/miqcie/mondrian

Curious: What security challenges are blocking your team from selling upmarket?

#ZeroTrust #Cybersecurity #StartupLife #OpenSource

---

## 📋 Launch Checklist

**Before posting:**
- [ ] Test all GitHub links work
- [ ] Verify demo examples run correctly  
- [ ] Check repository is public and polished
- [ ] Enable GitHub Discussions (✅ Done)
- [ ] Prepare to respond quickly to comments

**After posting:**
- [ ] Monitor comments and respond within 2 hours
- [ ] Track engagement metrics (upvotes, stars, comments)
- [ ] Create GitHub issues for feature requests
- [ ] Thank early adopters and contributors

**Timing:**
- Reddit: Post 8-10 AM ET (when devs check Reddit)
- HN: Post 8-10 AM PT (prime HN time)  
- Twitter: Post 9 AM PT (max engagement)
- LinkedIn: Post 12-2 PM ET (lunch browsing)

---

## 📊 Success Metrics

**Week 1 Goals:**
- 50+ GitHub stars
- 100+ Reddit upvotes total
- 10+ quality comments/feedback
- 5+ feature requests

**Week 2 Goals:**  
- 100+ GitHub stars
- 3+ external contributors interested
- 1000+ Twitter impressions
- 5+ companies expressing interest

---

*Ready to copy-paste and launch! Good luck! 🚀*