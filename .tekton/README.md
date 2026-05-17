# Tekton Pipelines-as-Code

This directory contains the Tekton Pipelines-as-Code configuration for HookDrop.

## Overview

The `pipelinerun.yaml` file defines how Tekton automatically triggers pipeline runs based on GitHub events (push and pull requests to `main`).

## How It Works

**Automatic Triggers:**
- `on-event`: `[push, pull_request]` — trigger on both push and PR events
- `on-target-branch`: `[main]` — only trigger on the `main` branch
- `on-status`: All task statuses trigger pipeline execution (can be filtered to `[success]` if needed)

**Template Variables:**
- `{{ repo_url }}` — repository URL (injected by Tekton Pipelines-as-Code)
- `{{ revision }}` — git ref/commit SHA (injected by Tekton Pipelines-as-Code)
- `{{ timestamp }}` — unique timestamp for each run
- `{{ repo_name_slug }}` — sanitized repository name

**Workspace:**
- A 1Gi PVC (`hookdrop-source-pvc`) is created per namespace to cache repository data

**Pipeline Reference:**
- Points to `hookdrop-release` pipeline in the `hookdrop-ci` namespace
- Uses `tekton-hookdrop` service account for task execution

## Setup

1. Install Tekton Pipelines and Tekton Pipelines-as-Code controller in your cluster.
2. Configure GitHub integration with your Tekton Pipelines-as-Code controller.
3. Push or create a PR to trigger the pipeline automatically.

## Manual Trigger (Optional)

If you want to manually trigger a PipelineRun outside of GitHub events:

```bash
kubectl apply -f .tekton/pipelinerun.yaml
```

This will create a one-off run with placeholder values. In production, Tekton Pipelines-as-Code substitutes real values from the GitHub webhook.

## Debugging

Check PipelineRun status:
```bash
kubectl get pipelinerun -n hookdrop-ci
kubectl describe pipelinerun <name> -n hookdrop-ci
kubectl logs -l tekton.dev/pipelinerun=<name> -n hookdrop-ci --all-containers
```

Check the Tekton Pipelines-as-Code controller logs:
```bash
kubectl logs -n tekton-pipelines deployment/pipelines-as-code-controller -f
```
