# HookDrop

HookDrop is a webhook receiver with a production-style DevOps platform around it.

## What is implemented

1. Observability (three pillars)
- Metrics: Prometheus + Grafana (`kube-prometheus-stack`)
- Logs: Loki
- Traces: OpenTelemetry in app + OTel Collector + Tempo

2. Ops maturity
- Alertmanager installed through `kube-prometheus-stack`
- `PrometheusRule` for restart alerts in Helm chart

3. Supply chain security
- Trivy image + config scanning in CI
- Cosign keyless image signing in CI
- Kyverno signature verification policy for ECR images

4. Zero trust controls
- NetworkPolicy for HookDrop pods
- ServiceAccount hardening (`automountServiceAccountToken: false`)
- Kyverno enforce policies:
  - no `latest` tags
  - resource limits required
  - runAsNonRoot required

5. Dependency discipline
- Renovate config (`renovate.json`) for Go, Dockerfile, GitHub Actions, Helm values updates

## CI/CD flow

- GitHub Actions (`.github/workflows/ci.yml`)
  - build, test, lint
  - Trivy scans
  - build + push image to ECR
  - Cosign sign image (keyless)
  - update `helm/hookdrop/values-prod.yaml` image tag
- ArgoCD watches git and syncs cluster

## Run everything (fresh local)

From project root:

```bash
# 0) (Optional) wipe old local cluster
kind delete cluster --name hookdrop

# 1) Create kind + ArgoCD + Kyverno policies
make cluster-up

# 2) Install observability stack
make observability-up

# 3) Build local app image
make docker-build

# 4) Load image into kind
kind load docker-image hookdrop:local --name hookdrop

# 5) Apply ArgoCD application (local values)
kubectl apply -f k8s/argocd/application.yaml

# 6) Watch deployment
kubectl get pods -n hookdrop -w
```

## Access UIs

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

ArgoCD password:

```bash
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d && echo
```

## Verify observability

```bash
# Send sample traffic
kubectl port-forward svc/hookdrop-hookdrop -n hookdrop 8080:80
curl -X POST http://localhost:8080/h/test -H "Content-Type: application/json" -d '{"hello":"world"}'

# Check observability namespace
kubectl get pods -n observability
```

## New files added

- `k8s/observability/*` (Prometheus/Grafana/Loki/Tempo/OTel setup values/manifests)
- `scripts/setup-observability.sh`
- `helm/hookdrop/templates/servicemonitor.yaml`
- `helm/hookdrop/templates/prometheusrule.yaml`
- `helm/hookdrop/templates/networkpolicy.yaml`
- `k8s/kyverno/verify-cosign-signature.yaml`
- `renovate.json`
- `telemetry.go` (OTel app instrumentation)

## Important note

The Cosign verify policy currently matches this ECR repo and GitHub workflow identity:
- `654654364687.dkr.ecr.us-east-1.amazonaws.com/hookdrop:*`
- `https://github.com/nirjxr26/HookDrop/.github/workflows/ci.yml@refs/heads/main`

If those change, update `k8s/kyverno/verify-cosign-signature.yaml`.
