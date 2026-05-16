# HookDrop

HookDrop is a mock webhook receiver written in Go. It accepts webhook traffic, stores requests in memory, exposes a live SSE feed per bucket, and includes a full DevOps pipeline around the app.

## Architecture

The Go service exposes `POST /h/:id` to ingest webhooks, `GET /h/:id` to read stored events, `GET /h/:id/stream` for SSE updates, plus `GET /healthz` and `GET /readyz` for health checks. Events live in an in-memory ring buffer with a max of 50 events per bucket, and each request is logged as structured JSON with `zerolog`.

The delivery path is GitHub Actions to AWS ECR to ArgoCD to a kind-based Kubernetes cluster. CI builds, tests, scans, and publishes an image, then ArgoCD syncs the Helm chart into the target cluster.

## DevOps Stack

| Tool | Purpose | Where it fits |
|---|---|---|
| Go | Application runtime | Webhook receiver and SSE server |
| Docker | Container build | Multi-stage image build |
| Distroless | Final image base | Minimal runtime image with no shell |
| GitHub Actions | CI pipeline | Build, test, lint, scan, publish |
| Trivy | Security scanning | Image scan and Helm config scan |
| AWS ECR | Container registry | Stores published images |
| Helm 3 | Kubernetes packaging | Deploys the app manifests |
| ArgoCD | GitOps deployment | Syncs Helm chart to cluster |
| kind | Local cluster | Runs the GitOps setup locally |
| Kyverno | Policy enforcement | Blocks bad pod specs |
| Sealed Secrets | Secret handling | Encrypted secrets in git |
| zerolog | Structured logging | JSON request/event logs |
| HPA | Autoscaling | Scales the Deployment on CPU |

## Getting Started

### Prerequisites

Go 1.25.10, Docker, kind, kubectl, helm, argocd CLI, kubeseal, golangci-lint, trivy.

### Run locally

```bash
go run ./...
```

### Run in Docker

```bash
make docker-run
```

### Full cluster setup

```bash
make cluster-up
kubectl apply -f k8s/argocd/application.yaml
argocd app get hookdrop
```

## CI/CD Pipeline

Pushes to `main` run lint, tests, Trivy image scanning, Trivy Helm config scanning, build and push the image to ECR, and then update `helm/hookdrop/values-prod.yaml` with the new tag. ArgoCD watches the repo, detects the change, and syncs the cluster automatically.

## Security Choices

- Distroless final image, so there is no shell or package manager in the runtime image.
- Non-root container user with UID 65532.
- `readOnlyRootFilesystem` and dropped capabilities.
- Trivy scans both the image and the Helm manifests in CI.
- Sealed Secrets keeps encrypted secret manifests safe to commit.
- Kyverno blocks `latest` tags and requires resource limits in the `hookdrop` namespace.

## Deploying to EKS

Everything here is cluster-agnostic. For local kind development, use `values-local.yaml` (local image). For EKS/production, set your ECR URI in `values-prod.yaml`, update the `repoURL` in `k8s/argocd/application.yaml` (or use a prod-specific ArgoCD app manifest), and point your kubeconfig at the EKS cluster.

## Project Structure

```text
.
├── .gitignore # Ignore local build outputs and env files.
├── Dockerfile # Multi-stage container build.
├── Makefile # Local, CI, Docker, and cluster commands.
├── README.md # Project overview and setup guide.
├── docker-compose.yml # Local container runtime composition.
├── go.mod # Go module definition.
├── go.sum # Go dependency checksums.
├── main.go # HTTP server entrypoint and graceful shutdown.
├── handler/ # HTTP handlers.
│   ├── dashboard.go # Basic landing page handler.
│   ├── health.go # Health and readiness handlers.
│   └── webhook.go # Webhook ingest, list, and SSE handlers.
├── helm/ # Helm chart for Kubernetes deployment.
│   └── hookdrop/ # HookDrop chart.
│       ├── Chart.yaml # Chart metadata.
│       ├── values-prod.yaml # Production overrides.
│       ├── values.yaml # Default chart values.
│       └── templates/ # Kubernetes manifests.
│           ├── _helpers.tpl # Shared name and label helpers.
│           ├── deployment.yaml # Application Deployment.
│           ├── hpa.yaml # HorizontalPodAutoscaler.
│           ├── ingress.yaml # Optional ingress.
│           ├── service.yaml # ClusterIP service.
│           └── serviceaccount.yaml # Service account and IRSA annotation hook.
├── k8s/ # GitOps and policy manifests.
│   ├── argocd/ # ArgoCD application definition.
│   │   └── application.yaml # ArgoCD Application resource.
│   ├── kind/ # Local cluster config.
│   │   └── cluster.yaml # kind cluster definition.
│   ├── kyverno/ # Policy-as-code.
│   │   ├── no-latest-tag.yaml # Blocks latest image tags.
│   │   └── require-resource-limits.yaml # Requires CPU and memory limits.
│   └── sealed-secrets/ # Sealed Secrets instructions.
│       └── README.md # How to seal and apply secrets.
├── scripts/ # Automation scripts.
│   └── setup-cluster.sh # Bootstraps kind, ArgoCD, Kyverno, and Sealed Secrets.
└── store/ # In-memory persistence.
    └── memory.go # Ring buffer and SSE subscriber store.
```
