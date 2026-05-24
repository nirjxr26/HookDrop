<div align="center">

<h1>HookDrop</h1>

<p>A self-hosted webhook receiver with in-memory event buckets, live SSE streaming, health probes, and a production-style Kubernetes delivery stack.</p>

</div>

---

## What It Does

HookDrop is a small Go service that captures and inspects webhooks. POST to a named bucket, and the event lands in memory with its headers, body, source IP, and a per-event trace ID. GET the bucket to list events. Subscribe to `/stream` and you get a live SSE feed as new webhooks arrive.

Around the app, the repo is also a working example of how to wire a Go service into a real delivery stack — Helm chart, ArgoCD sync, Kyverno admission policies, and a full observability pipeline (Prometheus, Grafana, Loki, Tempo, OpenTelemetry).

---

## Tech Stack
 
- **Backend** — Go, `net/http`, Zerolog  
- **Telemetry** — OpenTelemetry, OTLP gRPC exporter  
- **Container** — Docker, distroless runtime image  
- **Kubernetes** — Helm, kind, ArgoCD, Kyverno, NetworkPolicy, HPA  
- **Observability** — Prometheus, Grafana, Loki, Tempo, OTel Collector  
- **Supply Chain** — Trivy, Cosign, Renovate
 

---

## Features

**Webhook capture**
- `POST /h/<bucket>` — store a webhook event with full headers and body
- `GET /h/<bucket>` — list stored events (latest 50 per bucket, in memory)
- `GET /h/<bucket>/stream` — live SSE stream; receives events as they arrive

**Health and tracing**
- `/healthz` and `/readyz` for liveness and readiness probes
- Every request gets a trace ID (picks up `X-Trace-Id` if provided, generates one otherwise)
- OpenTelemetry spans emitted per request when `OTEL_EXPORTER_OTLP_ENDPOINT` is set

**Platform hardening**
- Resource limits, HPA, ServiceMonitor, and PrometheusRule in the Helm chart
- Service account token auto-mount disabled
- NetworkPolicy scoped to app and ArgoCD namespaces with OTLP egress only
- Kyverno policy requiring a valid Cosign signature before pods are admitted

---

## Architecture

### Request path

```
HTTP client
    └─▶ main.go (ServeMux, requestLogger middleware, OTel span)
            ├─▶ POST /h/<bucket>   webhook.go → memory.go → SSE subscribers
            ├─▶ GET  /h/<bucket>   webhook.go ← memory.go
            ├─▶ GET  /h/<bucket>/stream  webhook.go (SSE fan-out)
            ├─▶ /healthz /readyz   health.go
            └─▶ /                  dashboard.go (static route listing)
```

Two trace IDs are in play: the middleware creates a transport-level trace per request; the webhook handler generates a separate per-event trace ID that goes into the stored payload and response. They're related, but not the same value.

### CI/CD Architecture
 
<div align="center">
<img
  src="./diagrams/hookdrop_architecture.png"
  alt="Pipeline Architecture"
/>
</div>

### Pipeline Overview — from code push to running pod

- Push or PR to `main` triggers `ci.yml`, which sets up Go 1.25.10, restores the module cache, compiles the binary, runs `go test -race`, and gates on `golangci-lint`.
- Trivy scans both the container image and the Helm chart configuration for CVEs before anything is pushed.
- On merge to `main`, CI configures AWS credentials, creates the ECR repository if it doesn't exist, builds the multi-stage Docker image (Alpine builder → distroless runtime), and pushes it to ECR tagged by commit SHA.
- Cosign signs the image keylessly against Sigstore immediately after push.
- *(Gap — not yet automated)* `reusable-build.yml` knows how to commit the updated image tag back into `helm/hookdrop/values.yaml` and push it to the deploy branch, which would close the GitOps loop. Nothing in the active CI calls it yet, so tag promotion is currently manual.
- ArgoCD watches the deploy branch, detects the updated Helm values, and syncs the `hookdrop` namespace automatically.
- Before the pod is admitted, the Kyverno admission webhook verifies the Cosign signature on the image. Unsigned images are rejected at the cluster boundary.
- Pod starts; if `OTEL_EXPORTER_OTLP_ENDPOINT` is set, the OTel collector receives traces over gRPC and forwards them to Tempo.




### Kubernetes and GitOps
 
```
setup-cluster.sh
    ├─▶ kind cluster (kind-config.yaml)
    ├─▶ ArgoCD install
    ├─▶ Kyverno install + policies
    └─▶ kubectl apply application.yaml
            └─▶ ArgoCD syncs helm/hookdrop → hookdrop namespace
 
setup-observability.sh
    ├─▶ kube-prometheus-stack (Prometheus + Grafana + Alertmanager)
    ├─▶ Loki
    ├─▶ Tempo
    └─▶ OTel Collector (OTLP 4317/4318 → forwards traces to Tempo)
```
 
The ArgoCD application uses `values-local.yaml` by default. Prod uses `values-prod.yaml` (different registry, ingress config).
 
The Helm chart ships a ServiceMonitor and a PrometheusRule, but the app has no `/metrics` endpoint. The ServiceMonitor scrapes `/healthz`. The PrometheusRule alerts on pod restarts. Custom app metrics aren't implemented yet.
 
---
## Quick Start

### Prerequisites

- Go, Docker, kind, kubectl, Helm

```bash
git clone https://github.com/nirjxr26/HookDrop.git
cd HookDrop
```

### Option 1 — Local Go

```bash
make dev          # run the server on :8080
make build
make test
make lint
make scan
```

### Option 2 — Docker

```bash
make docker-build
make docker-run
```

### Option 3 — Kubernetes with kind

```bash
make cluster-up
make observability-up
make docker-build
kind load docker-image hookdrop:local --name hookdrop
kubectl apply -f k8s/argocd/application.yaml
kubectl get pods -n hookdrop -w
```

**Port-forward to access cluster UIs:**

```bash
# ArgoCD
kubectl port-forward svc/argocd-server -n argocd 8081:443

# Grafana
kubectl port-forward -n observability svc/kube-prometheus-stack-grafana 3000:80

# Prometheus
kubectl port-forward -n observability svc/kube-prometheus-stack-prometheus 9090:9090

# Alertmanager
kubectl port-forward -n observability svc/kube-prometheus-stack-alertmanager 9093:9093
```

**Verify the app:**

```bash
kubectl port-forward svc/hookdrop-hookdrop -n hookdrop 8080:80

curl -X POST http://localhost:8080/h/test \
  -H "Content-Type: application/json" \
  -d '{"hello":"world"}'

curl http://localhost:8080/h/test
curl http://localhost:8080/h/test/stream
```

---

## Environment Variables

```env
PORT=8080
LOG_LEVEL=info
OTEL_SERVICE_NAME=hookdrop
OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector.observability.svc.cluster.local:4317
```

OTEL tracing is a no-op if `OTEL_EXPORTER_OTLP_ENDPOINT` is not set.

---

## Project Structure

```text
├── main.go                       # HTTP server, routing, graceful shutdown
├── telemetry.go                  # OTel bootstrap and request span middleware
├── handler/
│   ├── webhook.go                # POST/GET/SSE for bucket events
│   ├── health.go                 # /healthz and /readyz
│   └── dashboard.go              # static route listing at /
├── store/
│   └── memory.go                 # in-memory store with SSE fan-out (50 events/bucket)
├── helm/hookdrop/                # Helm chart (Deployment, Service, NetworkPolicy, monitoring)
├── k8s/
│   ├── argocd/application.yaml   # ArgoCD app pointing at the Helm chart
│   ├── kyverno/                  # Cosign signature enforcement
│   └── observability/            # OTel collector, Loki, Tempo configs
├── scripts/
│   ├── setup-cluster.sh          # kind + ArgoCD + Kyverno bootstrap
│   └── setup-observability.sh    # Prometheus stack + Loki + Tempo
├── Dockerfile                    # multi-stage: Alpine builder → distroless runtime
├── docker-compose.yml
└── kind-config.yaml
```

---

## Notes

- Events are in memory. A pod restart clears everything.
- SSE connections receive periodic keepalive comments so they survive proxies and load balancers that close idle connections.