$ErrorActionPreference = "Stop"

function Assert-LastExitCode {
    param([string]$Step)
    if ($LASTEXITCODE -ne 0) {
        throw "$Step failed with exit code $LASTEXITCODE"
    }
}

$currentContext = kubectl config current-context
Assert-LastExitCode "kubectl context check"
if ($currentContext -ne "docker-desktop") {
    throw "Current kubectl context is '$currentContext'. Switch to 'docker-desktop' first."
}

Write-Host "Creating ArgoCD namespace if missing..."
kubectl create namespace argocd --dry-run=client -o yaml | kubectl apply -f -
Assert-LastExitCode "ArgoCD namespace setup"

Write-Host "Installing ArgoCD..."
kubectl apply --server-side -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
Assert-LastExitCode "ArgoCD installation"

Write-Host "Waiting for ArgoCD deployments..."
kubectl wait --for=condition=available --timeout=300s deployment/argocd-server -n argocd
Assert-LastExitCode "Wait for argocd-server"
kubectl wait --for=condition=available --timeout=300s deployment/argocd-repo-server -n argocd
Assert-LastExitCode "Wait for argocd-repo-server"
kubectl wait --for=condition=available --timeout=300s deployment/argocd-applicationset-controller -n argocd
Assert-LastExitCode "Wait for argocd-applicationset-controller"

Write-Host "Installing Kyverno..."
kubectl apply --server-side -f https://github.com/kyverno/kyverno/releases/latest/download/install.yaml
Assert-LastExitCode "Kyverno installation"

Write-Host "Waiting for Kyverno controllers..."
kubectl wait --for=condition=available --timeout=300s deployment/kyverno-admission-controller -n kyverno
Assert-LastExitCode "Wait for kyverno-admission-controller"
kubectl wait --for=condition=available --timeout=300s deployment/kyverno-background-controller -n kyverno
Assert-LastExitCode "Wait for kyverno-background-controller"
kubectl wait --for=condition=available --timeout=300s deployment/kyverno-cleanup-controller -n kyverno
Assert-LastExitCode "Wait for kyverno-cleanup-controller"
kubectl wait --for=condition=available --timeout=300s deployment/kyverno-reports-controller -n kyverno
Assert-LastExitCode "Wait for kyverno-reports-controller"

Write-Host "Applying Kyverno policies..."
kubectl apply -f k8s/kyverno/
Assert-LastExitCode "Kyverno policies"

Write-Host "Deploying ArgoCD application..."
kubectl apply -f k8s/argocd/application.yaml
Assert-LastExitCode "ArgoCD application deployment"

Write-Host ""
Write-Host "ArgoCD admin password:"
$password = kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}"
$decodedPassword = [System.Text.Encoding]::UTF8.GetString([System.Convert]::FromBase64String($password))
Write-Host $decodedPassword

Write-Host ""
Write-Host "Run this manually to access ArgoCD:"
Write-Host "kubectl port-forward svc/argocd-server -n argocd 8081:443"
Write-Host ""
Write-Host "Open:"
Write-Host "https://localhost:8081"
Write-Host ""
Write-Host "Username: admin"
