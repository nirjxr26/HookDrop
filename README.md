# HookDrop

POST a webhook to `/h/<bucket>`. It catches it. You inspect it.

That's the app. The interesting part is everything running behind it.

---

## What It Does

Five routes:

- `POST /h/<bucket>` — receives any webhook, stores it in memory
- `GET /h/<bucket>` — returns the last 50 events for that bucket
- `GET /h/<bucket>/stream` — SSE stream, new events pushed live
- `GET /healthz` / `GET /readyz` — liveness and readiness probes
- `GET /` — usage hints

Each stored event has a trace ID, the bucket name, method, headers, body, source IP, and timestamp. The store is in-memory only — no database, no Redis. 50 events per bucket, oldest dropped first.

---

## Why I Built This

I wanted a project small enough to understand completely but with enough moving parts to build a real DevOps pipeline around. A webhook receiver fits that. The app itself is ~300 lines of Go. Everything else — the CI, the GitOps loop, the security layer — is the actual point.

---

## Stack

**Frontend**
- Server-rendered dashboard in Go (`handler/dashboard.go`) — usage hints, bucket URLs
- SSE stream at `/h/<bucket>/stream` — new webhooks pushed live to the browser, no polling

**Backend**
- Go 1.22, stdlib `net/http` only (no Gin, no Echo)
- `zerolog` for structured JSON logging, trace ID on every request
- In-memory ring buffer, 50 events per bucket, thread-safe with `sync.RWMutex`
- Graceful shutdown on `SIGINT`/`SIGTERM`, 10 second drain window

**DevOps**
- Docker: multi-stage build, distroless final image, non-root UID 65532, read-only root filesystem
- GitHub Actions: lint → test → Trivy image scan → Trivy Helm config scan → build → push ECR → update `values-prod.yaml`
- AWS ECR: image registry, SHA-tagged on every push
- Helm 3: separate `values-local.yaml` and `values-prod.yaml`, HPA with properly set requests so autoscaling actually works
- ArgoCD: GitOps sync, auto-prune, self-heal — cluster state always matches git
- kind: local Kubernetes cluster, EKS-compatible manifests
- Kyverno: two admission policies — no `latest` tags, required resource limits, both enforced not just documented
- Bitnami Sealed Secrets: secrets encrypted before they touch git

---

## Local Setup

You need: Go 1.22, Docker, kind, kubectl, helm, argocd CLI, golangci-lint, trivy

```bash
# Just run the app
go run ./...
# POST http://localhost:8080/h/test
# GET  http://localhost:8080/h/test

# Or with Docker
make docker-run

# Full local cluster
make cluster-up
kubectl apply -f k8s/argocd/application.yaml
```

The local ArgoCD Application uses `values-local.yaml` which points to `hookdrop:local` with `pullPolicy: Never`. To get an image into the cluster:

```bash
docker build -t hookdrop:local .
kind load docker-image hookdrop:local --name hookdrop
```

No pull secrets, no ECR auth, no token expiry headaches locally. ECR is only in the CI/production path where it makes sense — EKS nodes pull via IAM automatically.

---

## CI/CD Flow

```
git push to main
  → lint + test
  → docker build
  → trivy image scan     (fails CI on HIGH/CRITICAL)
  → trivy config scan    (fails CI on Helm misconfigs)
  → push to ECR          (tagged sha-<commit>)
  → update values-prod.yaml image.tag
  → ArgoCD detects drift
  → syncs Helm chart to cluster
  → new pod rolls out
```

The cluster updates itself. Nothing manual after the push.

---

## Tekton

Tekton is now available as an in-cluster CI runner for this project. The manifests live in [k8s/tekton](k8s/tekton), and the pipeline does the following:

- clones the repo
- runs Go tests with race detection
- scans Helm and Kubernetes manifests with Trivy
 - builds the container image into a tarball (no push)
 - generates an SBOM and runs a Grype vulnerability scan

Tekton does not replace ArgoCD here. It produces the image; ArgoCD still handles deployment after the Helm tag is updated.

Manual flow:

1. Install Tekton Pipelines in the cluster.
2. Wait for the Tekton controller and webhook deployments to be available if you just installed Tekton.
3. Apply the Tekton manifests.
4. Edit [k8s/tekton/pipelinerun.yaml](k8s/tekton/pipelinerun.yaml) with your image tag (used only for local naming).
5. Run the PipelineRun.
6. Update the Helm image tag or use an image updater so ArgoCD can deploy the new image pushed by GitHub Actions.

If you want the shortest command path, use `make tekton-apply` and `make tekton-run`.
`make tekton-run` uses `kubectl create` because the PipelineRun manifest uses `generateName`.

## Security Choices

- Distroless image — no shell, no package manager, smaller CVE surface
- Non-root UID 65532, `readOnlyRootFilesystem: true`, all capabilities dropped, `seccompProfile: RuntimeDefault`
- Trivy runs on both the image and the Helm chart in CI, not just one
- Kyverno rejects pods at admission if they use `latest` tag or missing resource limits — enforced, not just documented
- Sealed Secrets encrypts anything sensitive before it touches git

---

## Two Deployment Profiles

| | Local (kind) | Production (EKS) |
|---|---|---|
| Image | `hookdrop:local` | ECR URI + SHA tag |
| Pull policy | `Never` | `Always` |
| Auth | None | IAM role on node |
| Ingress | Disabled | Enabled |
| Values file | `values-local.yaml` | `values-prod.yaml` |

Switch the ArgoCD Application to `values-prod.yaml` for the production path.

---

## Repo Layout

```
.
├── main.go
├── handler/
│   ├── webhook.go       # ingest, list, SSE stream
│   ├── health.go        # /healthz, /readyz
│   └── dashboard.go     # landing page
├── store/
│   └── memory.go        # thread-safe ring buffer + SSE fanout
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── helm/hookdrop/
│   ├── values.yaml
│   ├── values-local.yaml
│   ├── values-prod.yaml
│   └── templates/
├── k8s/
│   ├── argocd/application.yaml
│   ├── kind/cluster.yaml
│   ├── kyverno/
│   │   ├── no-latest-tag.yaml
│   │   └── require-resource-limits.yaml
│   ├── tekton/
│   │   ├── pipeline.yaml
│   │   ├── pipelinerun.yaml
│   │   ├── README.md
│   │   ├── serviceaccount.yaml
│   │   ├── namespace.yaml
│   │   └── tasks.yaml
│   └── sealed-secrets/README.md
└── scripts/
    └── setup-cluster.sh
```

---

## Troubleshooting

**`ImagePullBackOff` on kind**
The cluster is using the prod values file and trying to hit ECR. Check that `k8s/argocd/application.yaml` references `values-local.yaml`, not `values-prod.yaml`. Then load the local image: `kind load docker-image hookdrop:local --name hookdrop`.

**ArgoCD shows OutOfSync but won't sync**
Usually a Kyverno policy block. Run `kubectl describe pod <pod> -n hookdrop` and check events — it'll say which policy failed.

**Kyverno blocking your pod**
Either missing resource limits or using `latest` tag. Both are intentional. Fix the values file, not the policy.