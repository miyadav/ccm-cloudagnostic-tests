# Quick Start Guide

## Running Tests - Choose Your Method

### Method 1: Using Go Test (Recommended - No Extra Installation)
```bash
cd cmd/existing-ccm-test
go test -v --kubeconfig ~/.kube/config
```

### Method 2: Using Ginkgo CLI (Better Output, Requires Installation)
```bash
# Install Ginkgo CLI first
go install github.com/onsi/ginkgo/v2/ginkgo@latest

# Run tests
cd cmd/existing-ccm-test
ginkgo run -v -- --kubeconfig ~/.kube/config
```

### Method 3: Using Makefile (Requires Ginkgo CLI Installed)
```bash
# Uses default $HOME/.kube/config
make test-ginkgo-verbose

# Or specify custom path (use $HOME not ~)
make test-ginkgo-verbose KUBECONFIG=$HOME/.kube/config

# Or use absolute path
make test-ginkgo-verbose KUBECONFIG=/Users/yourname/.kube/config
```

---

## Common Mistakes to Avoid

### ❌ WRONG:
```bash
make test-ginkgo-verbose --kubeconfig ~/.kube/config
```
**Error:** `make: unrecognized option '--kubeconfig'`

**Why it fails:** Make doesn't understand `--kubeconfig`. It's a flag for the underlying test, not for make.

### ✅ CORRECT:
```bash
# Best: Let it use the default
make test-ginkgo-verbose

# Or: Use $HOME for path expansion
make test-ginkgo-verbose KUBECONFIG=$HOME/.kube/config

# Or: Use absolute path
make test-ginkgo-verbose KUBECONFIG=/Users/yourname/.kube/config
```
**How it works:** `KUBECONFIG=` is a Makefile variable. Use `$HOME` or absolute paths, not `~` (tilde won't expand properly).

---

## Quick Test Commands

### Run All Tests (Verbose)
```bash
# Using go test (easiest) - use $HOME or absolute path
cd cmd/existing-ccm-test && go test -v --kubeconfig $HOME/.kube/config

# Using ginkgo (requires installation)
cd cmd/existing-ccm-test && ginkgo run -v -- --kubeconfig ~/.kube/config

# Using makefile (requires ginkgo) - uses default or specify with $HOME
make test-ginkgo-verbose
# or
make test-ginkgo-verbose KUBECONFIG=$HOME/.kube/config
```

### Run With JUnit Output (for CI/CD)
```bash
# Using ginkgo
cd cmd/existing-ccm-test
ginkgo run --junit-report=../../test-results/junit.xml -- --kubeconfig ~/.kube/config

# Using makefile (uses default $HOME/.kube/config)
make test-ginkgo-junit
```

### Run Legacy E2E Test Runner
```bash
./bin/e2e-test-runner --provider existing --kubeconfig ~/.kube/config --suite all --verbose
```

---

## Installation

### Install Ginkgo CLI (Optional but Recommended)
```bash
go install github.com/onsi/ginkgo/v2/ginkgo@latest

# Verify installation
ginkgo version
```

### Build Test Runner
```bash
make build
# or
go build -o bin/e2e-test-runner ./cmd/e2e-test-runner/
```

---

## Makefile Variable Reference

All Ginkgo test targets accept these variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `KUBECONFIG` | `~/.kube/config` | Path to kubeconfig file |
| `FOCUS` | (none) | Run tests matching pattern |
| `SKIP` | (none) | Skip tests matching pattern |
| `LABELS` | (none) | Run tests with specific labels |
| `REPEAT` | (none) | Number of times to repeat tests |
| `ATTEMPTS` | (none) | Number of flake retry attempts |

### Examples:
```bash
# Run specific tests (uses default kubeconfig)
make test-ginkgo-focus FOCUS="LoadBalancer"

# Skip certain tests
make test-ginkgo-skip SKIP="slow"

# Run tests with specific labels
make test-ginkgo-labels LABELS="loadbalancer"

# Repeat tests 3 times
make test-ginkgo-repeat REPEAT=3

# Or specify custom kubeconfig with any target
make test-ginkgo-focus FOCUS="LoadBalancer" KUBECONFIG=$HOME/.kube/config
```

---

## Test Results Location

- **Test logs:** Console output
- **JUnit XML:** `test-results/junit.xml` (when using junit target)
- **Coverage:** `coverage.out` (when using coverage target)

---

## Troubleshooting

### "ginkgo: command not found"
**Solution:** Either:
1. Install ginkgo: `go install github.com/onsi/ginkgo/v2/ginkgo@latest`
2. Use `go test` instead: `cd cmd/existing-ccm-test && go test -v --kubeconfig ~/.kube/config`

### "make: unrecognized option"
**Solution:** Use `KUBECONFIG=` not `--kubeconfig`:
```bash
# Wrong
make test-ginkgo-verbose --kubeconfig ~/.kube/config

# Right - use default
make test-ginkgo-verbose

# Right - use $HOME or absolute path
make test-ginkgo-verbose KUBECONFIG=$HOME/.kube/config
```

### "stat ~/.kube/config: no such file or directory"
**Problem:** The tilde (`~`) is not being expanded

**Solution:** Use `$HOME` or absolute paths in Makefile commands:
```bash
# Wrong - tilde won't expand in make variables
make test-ginkgo-verbose KUBECONFIG=~/.kube/config

# Right - use $HOME
make test-ginkgo-verbose KUBECONFIG=$HOME/.kube/config

# Right - use absolute path
make test-ginkgo-verbose KUBECONFIG=/Users/yourname/.kube/config

# Best - use default (already set to $HOME/.kube/config)
make test-ginkgo-verbose
```

### Tests can't connect to cluster
**Solution:** Verify kubeconfig path:
```bash
# Test connection
kubectl get nodes --kubeconfig ~/.kube/config

# Use absolute path if needed
make test-ginkgo-verbose KUBECONFIG=/Users/yourname/.kube/config
```

---

## Next Steps

- Read full [README.md](README.md) for detailed documentation
- Check [docs/existing-ccm-testing.md](docs/existing-ccm-testing.md) for CCM testing guide
- See [docs/automatic-cleanup.md](docs/automatic-cleanup.md) for cleanup information

