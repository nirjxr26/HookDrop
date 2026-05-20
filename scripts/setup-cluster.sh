#!/usr/bin/env bash


set -euo pipefail

# Create kind cluster
kind create cluster --config k8s/kind/cluster.yaml

# Create ArgoCD namespace if missing
kubectl create namespace argocd --dry-run=client -o yaml | kubectl apply -f -

# Install ArgoCD
kubectl apply --server-side -n argocd \
-f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# Wait for ArgoCD
kubectl wait --for=condition=available --timeout=300s \
deployment/argocd-server -n argocd

kubectl wait --for=condition=available --timeout=300s \
deployment/argocd-repo-server -n argocd

kubectl wait --for=condition=available --timeout=300s \
deployment/argocd-applicationset-controller -n argocd

# Install Kyverno
kubectl apply --server-side \
-f https://github.com/kyverno/kyverno/releases/latest/download/install.yaml

# Wait for Kyverno controllers
kubectl wait --for=condition=available --timeout=300s \
deployment/kyverno-admission-controller -n kyverno

kubectl wait --for=condition=available --timeout=300s \
deployment/kyverno-background-controller -n kyverno

kubectl wait --for=condition=available --timeout=300s \
deployment/kyverno-cleanup-controller -n kyverno

kubectl wait --for=condition=available --timeout=300s \
deployment/kyverno-reports-controller -n kyverno

# Apply Kyverno policies
kubectl apply -f k8s/kyverno/

# Deploy ArgoCD application
kubectl apply -f k8s/argocd/application.yaml

echo
echo "ArgoCD admin password:"
kubectl -n argocd get secret argocd-initial-admin-secret \
-o jsonpath="{.data.password}" | base64 -d

echo
echo
echo "Run this manually to access ArgoCD:"
echo "kubectl port-forward svc/argocd-server -n argocd 8081:443"
echo
echo "Open:"
echo "https://localhost:12345"
echo
echo "Username: admin"