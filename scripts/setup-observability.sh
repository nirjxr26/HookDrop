#!/usr/bin/env bash
set -euo pipefail

kubectl create namespace observability --dry-run=client -o yaml | kubectl apply -f -

helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update

helm upgrade --install kube-prometheus-stack prometheus-community/kube-prometheus-stack \
  -n observability \
  -f k8s/observability/kube-prometheus-values.yaml

helm upgrade --install loki grafana/loki \
  -n observability \
  -f k8s/observability/loki-values.yaml \
  --set loki.useTestSchema=true

helm upgrade --install tempo grafana/tempo \
  -n observability \
  -f k8s/observability/tempo-values.yaml \
  --set loki.useTestSchema=true

kubectl apply -f k8s/observability/otel-collector-config.yaml
kubectl apply -f k8s/observability/otel-collector.yaml

echo "Observability stack installed."
echo "Grafana: kubectl port-forward -n observability svc/kube-prometheus-stack-grafana 3000:80"
echo "Prometheus: kubectl port-forward -n observability svc/kube-prometheus-stack-prometheus 9090:9090"
echo "Alertmanager: kubectl port-forward -n observability svc/kube-prometheus-stack-alertmanager 9093:9093"
