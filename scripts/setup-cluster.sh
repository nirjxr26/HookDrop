#!/usr/bin/env bash
set -euo pipefail

kind create cluster --config k8s/kind/cluster.yaml

kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

kubectl wait --for=condition=available --timeout=120s deployment/argocd-server -n argocd

kubectl port-forward svc/argocd-server -n argocd 8081:443 &

# Install Kyverno
kubectl create -f https://github.com/kyverno/kyverno/releases/latest/download/install.yaml
kubectl wait --for=condition=available --timeout=120s deployment/kyverno -n kyverno

# Install Sealed Secrets controller
helm repo add sealed-secrets https://bitnami-labs.github.io/sealed-secrets
helm repo update
helm install sealed-secrets sealed-secrets/sealed-secrets \
	--namespace kube-system \
	--set fullnameOverride=sealed-secrets-controller

# Apply Kyverno policies
kubectl apply -f k8s/kyverno/

kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d

echo

kubectl apply -f k8s/argocd/application.yaml

echo "ArgoCD running at https://localhost:8081 — user: admin"
