# HookDrop

POST a webhook to `/h/<bucket>`. It catches it. You inspect it.

That is the app. The interesting part is the platform around it.

## What It Does

- `POST /h/<bucket>`: receive a webhook event
- `GET /h/<bucket>`: fetch buffered events (last 50)
- `GET /h/<bucket>/stream`: live SSE stream
- `GET /healthz` and `GET /readyz`: probe endpoints
- `GET /`: quick usage page

Each event stores trace ID, bucket, method, headers, body, source IP, and timestamp.

## Stack

- Go service (`net/http` + `zerolog`)
- Docker multi-stage build with distroless runtime
- GitHub Actions CI/CD (build, test, lint, scan, push, GitOps tag bump)
- AWS ECR image registry
- Helm chart with environment split (`values-local.yaml`, `values-prod.yaml`)
- ArgoCD GitOps sync
- kind for local Kubernetes
- Kyverno admission policies

## Local Setup

Prereqs: Go 1.25.10, Docker, kind, kubectl, helm, argocd CLI, golangci-lint, trivy.

```bash
go run ./...
```

```bash
make docker-run
```

```bash
make cluster-up
kubectl apply -f k8s/argocd/application.yaml
```

Local image flow:

```bash
docker build -t hookdrop:local .
kind load docker-image hookdrop:local --name hookdrop
```

## CI/CD Flow

1. Push to `main`
2. GitHub Actions runs lint/test/build
3. Trivy scans image and Helm config
4. Image is pushed to ECR as `sha-<commit>`
5. CI updates `helm/hookdrop/values-prod.yaml` tag
6. ArgoCD detects drift and syncs cluster

## Security and Policy

- Distroless runtime
- Non-root container UID
- Read-only root filesystem
- Dropped Linux capabilities
- Trivy security scanning in CI
- Kyverno admission enforcement:
  - no `latest` image tags
  - required CPU/memory limits
  - required `runAsNonRoot: true`

## Deployment Profiles

- Local (`values-local.yaml`): `hookdrop:local`, no registry auth
- Production (`values-prod.yaml`): ECR image + SHA tag

## Repo Layout

- `main.go`, `handler/`, `store/`: app runtime
- `helm/hookdrop/`: chart + values
- `k8s/argocd/`: ArgoCD application
- `k8s/kind/`: kind cluster config
- `k8s/kyverno/`: policy manifests
- `scripts/setup-cluster.sh`: cluster bootstrap
