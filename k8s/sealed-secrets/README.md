Sealed Secrets lets you store encrypted Kubernetes secrets in git. You create a regular Secret locally, encrypt it with the cluster's public key, and commit the encrypted SealedSecret manifest instead of the raw Secret. The controller running in the cluster can decrypt it, but the git repository only contains ciphertext.

## Install the controller

```bash
helm repo add sealed-secrets https://bitnami-labs.github.io/sealed-secrets
helm install sealed-secrets sealed-secrets/sealed-secrets -n kube-system
```

## Seal a secret

```bash
kubectl create secret generic hookdrop-secret \
  --from-literal=SOME_KEY=somevalue \
  --dry-run=client -o yaml | \
  kubeseal --format yaml > k8s/sealed-secrets/hookdrop-secret.yaml
```

## Apply the sealed secret

```bash
kubectl apply -f k8s/sealed-secrets/hookdrop-secret.yaml
```

The SealedSecret file is safe to commit to git. The raw Secret is not.
