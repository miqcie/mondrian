# Mondrian

> **⚠️ RESEARCH PROJECT - ARCHIVED JANUARY 2026**
>
> This repository represents a completed validation exercise for evidence-first Zero Trust architecture.
> Built in 2 days (September 2025), it successfully demonstrated the feasibility of the approach.
> Archived to focus on advisory/content work. See [RETROSPECTIVE.md](RETROSPECTIVE.md) for details.

---

> **Research Prototype:** Zero Trust validation exercise (not production-ready)

**The Problem:** Startups live below the security poverty line. They move fast, then lose enterprise deals when buyers demand SOC2 or ISO. Enterprise Zero Trust platforms cost more than a junior developer's salary.

**The Solution:** Open-source Zero Trust that blocks risky changes AND proves your controls actually ran. Built for engineers, not procurement departments.

**Current Release:** Policy enforcement CLI + GitHub Actions integration. First component of the larger Mondrian.dev Zero Trust OS vision.

## Quick Start

```bash
# Install
go install github.com/miqcie/mondrian/cmd/mondrian@latest

# Check your infrastructure
mondrian check

# Generate proof bundle
mondrian attest

# Verify evidence chain
mondrian verify
```

## Why This Matters

Enterprise buyers tighten security requirements every quarter. Startups that can't prove compliance lose deals. The current solution — Okta + VPNs + Drata + consultants — costs six figures and takes months.

Mondrian flips this: **cryptographic proofs instead of compliance theater**.

## What It Does Today

- 🛡️ **Blocks risky infrastructure** - Catches public S3 buckets, open security groups  
- 🔐 **Enforces deployment security** - Requires OIDC, prevents risky merges
- 📋 **Generates compliance proof** - SLSA-compatible attestations with tamper-evident chains
- ⚡ **Integrates with GitHub** - Works in your existing workflow, not against it
- 💰 **Costs $0** - Because security shouldn't be a luxury good

## The Bigger Vision

This CLI is the first piece of **Mondrian.dev** — the complete Zero Trust OS for developers. Coming soon: passkeys, zero-trust proxy, infra posture monitoring, and full compliance automation in one console.

## Installation

**From source:**
```bash
git clone https://github.com/miqcie/mondrian.git
cd mondrian
go build -o mondrian cmd/mondrian/main.go
```

**GitHub Action:**
```yaml
- uses: miqcie/mondrian-action@v1
  with:
    rules: 'iac,deploy,device'
```

## Historical Reference

This project's insights have been extracted as thought leadership content:

- **[Why I Built and Shelved a Zero Trust Tool](https://eagleridgeadvisory.com)** - Decision framework for CISOs (coming soon)
- **[The Security Poverty Line](https://eagleridgeadvisory.com)** - Economic analysis of the $200K compliance barrier (coming soon)
- **[Evidence-First vs Detection-Only Security](https://eagleridgeadvisory.com)** - Technical deep-dive with code examples (coming soon)

For advisory/consulting inquiries about compliance automation, evidence-first security, or build/buy decisions in the security space, visit [Eagle Ridge Advisory](https://eagleridgeadvisory.com) or [Caldris](https://caldris.io).

## License

Apache-2.0

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for details.