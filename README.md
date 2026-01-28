# Headscale/Tailscale Health Exporter

> A Prometheus exporter that monitors Headscale and Tailscale APIs, deployed to EKS with a full GitOps pipeline.

## What This Does

This exporter polls your Headscale or Tailscale control server and exposes metrics that Prometheus can scrape:

- Is the API reachable?
- How long do API calls take?
- How many nodes are online?
- When was each node last seen?

If you're running a Headscale server (self-hosted Tailscale control plane) or using Tailscale's API, this gives you visibility into the health of your mesh network.

## Why This Exists

**1. A similar exporter exists, but it has tradeoffs I didn't want:**
- Uses Alpine as the runtime image, which means CVEs, shell access, and unnecessary attack surface
- Helm-only deployment, which doesn't suit everyone

Mine uses Alpine only in the build stage. The final runtime image is scratch. No shell, no package manager, nothing to exploit.

**2. I'm a DevOps/Platform engineer learning Go.** The exporter demonstrates I can read and write code. But the real portfolio piece is the infrastructure around it. This project deploys to EKS with Terraform, GitOps, service mesh, observability, and automated certificate management.

## How It Works (High Level)

**The Application:**
1. On startup, loads config from YAML (provider URLs, API keys, polling interval)
2. Runs a background loop that polls each provider on a configurable interval
3. For each poll: checks API health, measures latency, fetches node list
4. Updates Prometheus metrics with results
5. Exposes `/metrics` endpoint for Prometheus to scrape


**The Platform:**

6. Deploys to EKS via ArgoCD (GitOps)
7. Gateway API routes external traffic
8. CertManager handles TLS certificates via Let's Encrypt
9. ExternalDNS updates Route 53 automatically
10. Prometheus scrapes metrics, Grafana visualises them

## Key Features

- **Multi-provider support:** Works with both Headscale and Tailscale APIs
- **Security-hardened image:** Scratch runtime, no shell, runs as non-root, approximately 7MB, zero CVEs
- **Deployment flexibility:** Raw Kubernetes manifests and Kustomize (no Helm required)
- **Full GitOps pipeline:** Push to main, ArgoCD syncs, cluster updates

## Technology Choices

**Application:**
- **Go:** Industry standard for cloud-native tooling, produces static binaries, strong Prometheus client library
- **Scratch base image:** No OS, no package manager, no shell. Minimal attack surface.

**Infrastructure:**
- **Terraform modules:** Reusable, auditable infrastructure for EKS, VPC, IAM
- **Gateway API:** Kubernetes-native ingress, more expressive than Ingress resources
- **ArgoCD:** GitOps. Cluster state always matches Git.
- **Kustomize:** Simpler than Helm, easier to audit, you see exactly what gets deployed
- **CertManager and Let's Encrypt:** Automated TLS, no manual certificate management
- **ExternalDNS:** DNS records update automatically when services change

## Architecture Overview

*[Diagram coming]*

## Getting Started (For Developers)

### Prerequisites
- Go 1.21+
- Docker (for building the image)
- A Headscale instance or Tailscale API key (for real testing)

### Installation
```bash
git clone https://github.com/[your-username]/headscale-exporter.git
cd headscale-exporter
go mod download
```

### Configuration
```yaml
# config.yaml
providers:
  - name: headscale
    url: http://localhost:8080
    api_key: your-api-key-here

interval: 30s
```

### Running Locally
```bash
go run cmd/main.go
# Metrics available at http://localhost:8080/metrics
```

### Building the Docker Image
```bash
docker build -t headscale-exporter .
```

## Project Status

**Current state:** In active development

**What's working:**
- Config loading from YAML
- Headscale provider (CheckAPI, ListNodes)
- API health and latency metrics
- Background polling loop
- Docker image (7MB, scratch base)

**What's planned:**

| Component | Status |
|-----------|--------|
| Kubernetes manifests | Planned |
| Kustomize overlays (dev/prod) | Planned |
| EKS cluster (Terraform) | Planned |
| Gateway API routing | Planned |
| CertManager and TLS | Planned |
| ExternalDNS | Planned |
| ArgoCD GitOps | Planned |
| Prometheus and Grafana | Planned |
| Grafana dashboard JSON | Planned |
| Graceful shutdown | Nice to have |

## Challenges and Learnings

- **Challenge:** Handling different error types (server errors vs connection failures). **Solution:** Designed a clear contract where `(bool, duration, error)` returns error only for transport failures.
- **Challenge:** Defer doesn't work in infinite loops. **Solution:** Manual `cancel()` calls for context timeouts inside the ticker loop.
- **Learning:** Context pattern in Go. The caller sets the timeout, the function obeys via `http.NewRequestWithContext`.


## Questions or Issues?

Open an issue or reach out on [LinkedIn](www.linkedin.com/in/elsadevops).

