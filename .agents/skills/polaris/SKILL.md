---
name: polaris
description: Expertise in the Polaris (HookDrop) Go webhook receiver, Helm, Kind, ArgoCD, Kyverno, and OpenTelemetry stack.
---

# Polaris (HookDrop) Customization Skill

This skill helps the agent understand, modify, build, and deploy the Polaris (HookDrop) service.

## Core Responsibilities

1. **Go Webhook Handlers (`handler/`)**
   - [webhook.go](file:///C:/Users/Admin/Desktop/Desktop%20Backup/Projects/Polaris/handler/webhook.go) handles request routing: `/h/<bucket>`, `/h/<bucket>/stream`.
   - [dashboard.go](file:///C:/Users/Admin/Desktop/Desktop%20Backup/Projects/Polaris/handler/dashboard.go) renders the HTML homepage.
   - [health.go](file:///C:/Users/Admin/Desktop/Desktop%20Backup/Projects/Polaris/handler/health.go) provides `/healthz` and `/readyz` endpoints.

2. **In-Memory Store (`store/`)**
   - [memory.go](file:///C:/Users/Admin/Desktop/Desktop%20Backup/Projects/Polaris/store/memory.go) manages all buckets and events using an in-memory database with a read/write mutex (`sync.RWMutex`). It caps stored events at 50 per bucket.

3. **Infrastructure & Deployment**
   - Local cluster setup: `make cluster-up` (runs [setup-cluster.sh](file:///C:/Users/Admin/Desktop/Desktop%20Backup/Projects/Polaris/scripts/setup-cluster.sh) or [setup-cluster.ps1](file:///C:/Users/Admin/Desktop/Desktop%20Backup/Projects/Polaris/scripts/setup-cluster.ps1) targeting `docker-desktop`).
   - Monitoring stack: `make observability-up` (runs [setup-observability.sh](file:///C:/Users/Admin/Desktop/Desktop%20Backup/Projects/Polaris/scripts/setup-observability.sh)).
   - Building locally: `make build` or `make docker-build`.
   - Running tests: `make test`.

## Separation Context

- **Do NOT mix with AegisMesh IAM.** Polaris is a completely self-contained webhook-inspection workspace. Keep all telemetry, packaging (Helm), and code logic fully distinct and independent.
