# Tekton Pipeline

This directory contains a Kubernetes-native CI pipeline for HookDrop.

What it does:

- clones the repository
- runs `go test -race ./...`
- runs Trivy config scans against `helm/` and `k8s/`
- builds the container image into a tarball (no push)
- generates an SBOM and runs a Grype vulnerability scan against the SBOM/source

What it does not do yet:

- it does not update `helm/hookdrop/values-prod.yaml`
- it does not deploy directly to the cluster

ArgoCD still owns deployment. After Tekton pushes a new image, you still need to bump the Helm tag or use an image updater workflow.

## Manual setup

1. Install Tekton Pipelines in the cluster.
2. Install Tekton Triggers only if you want webhook or event-driven runs later.
3. Apply the manifests in this folder.
4. Edit `pipelinerun.yaml` with any pipeline params you want (image tag used for local naming only).
5. Run the PipelineRun and watch the logs in Tekton.

## Apply everything

```bash
kubectl apply -f k8s/tekton/namespace.yaml
kubectl apply -f k8s/tekton/serviceaccount.yaml
kubectl apply -f k8s/tekton/tasks.yaml
kubectl apply -f k8s/tekton/pipeline.yaml
kubectl create -f k8s/tekton/pipelinerun.yaml
```
