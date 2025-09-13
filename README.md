# Mondrian

> Evidence-first Zero Trust runtime for startups

Mondrian blocks risky changes before merge or deploy and proves controls run with tamper-evident records.

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

## How It Works

1. **Policy**: Declarative rules across code, infrastructure, device, and identity
2. **Gateways**: CI/CD checks, pre-commit hooks, and runtime webhooks  
3. **Evidence**: Signed attestations with immutable log and query API

## MVP Features

- ✅ Infrastructure as Code policy checks (S3, Security Groups)
- ✅ Deploy policy enforcement (OIDC, two-person review)
- ✅ Device health verification (osquery integration)
- ✅ GitHub Actions integration
- ✅ SLSA-compatible attestation generation
- ✅ Offline proof verification

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

## License

Apache-2.0

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for details.