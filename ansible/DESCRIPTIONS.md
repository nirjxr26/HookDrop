# Ansible Organized Files - Short Descriptions

- `ansible.cfg`: Ansible config; updated inventory path and roles_path set to `./roles`.
- `inventory/default.yaml`: Localhost inventory (ansible_connection: local).
- `requirements.txt`: Python package pins required to run Ansible automation.
- `playbooks/cluster-setup.yaml`: Main orchestration playbook (bootstraps kind, ArgoCD, Kyverno, Sealed Secrets).

Roles (standard layout: `roles/<name>/{defaults,tests,tasks}`)
- `roles/prep/defaults/main.yaml`: Defaults for prerequisites (tool version pins).
- `roles/prep/tasks/main.yaml`: Install/verify Docker, kubectl, helm, kind, argocd CLI, trivy, golangci-lint.
- `roles/cluster/defaults/main.yaml`: Defaults for kind cluster (name, config path, kubeconfig path).
- `roles/cluster/tasks/main.yaml`: Create or recreate kind cluster and verify node readiness.
- `roles/argocd/defaults/main.yaml`: ArgoCD defaults (namespace, manifest URL, port-forward port).
- `roles/argocd/tasks/main.yaml`: Install ArgoCD, wait for deployments, extract admin password.
- `roles/policy/defaults/main.yaml`: Kyverno defaults (namespace, manifest URL, policies path).
- `roles/policy/tasks/main.yaml`: Install Kyverno, apply policies from `k8s/kyverno/`.
- `roles/secrets/defaults/main.yaml`: Sealed Secrets defaults (helm chart/repo, namespace).
- `roles/secrets/tasks/main.yaml`: Add helm repo, install Sealed Secrets, back up public key locally.

Notes:
- These files are copies of the flat `ansible_*` files moved into a standard Ansible layout.
- I left original docs in the project root (e.g., `ANSIBLE_BOOTSTRAP.md`, `ANSIBLE_QUICK_START.md`) to preserve history.
- To finish reorganization: run `make ansible-reorganize` or remove the flat files after verification.
