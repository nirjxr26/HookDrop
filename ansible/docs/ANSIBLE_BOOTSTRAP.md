# HookDrop Ansible Local Cluster Bootstrap

This Ansible playbook automates the complete setup of a local Kubernetes cluster for HookDrop, replacing the manual `setup-cluster.sh` script.

## What It Does

The bootstrap playbook automates:

1. **Prerequisites Installation** - Docker, kubectl, helm, kind, argocd-cli, trivy, golangci-lint
2. **Kind Cluster Creation** - Creates a local Kubernetes cluster named "hookdrop"
3. **ArgoCD Installation** - GitOps deployment controller
4. **Kyverno Installation** - Policy engine for admission control
5. **Sealed Secrets Setup** - Secret encryption before git commit

## Prerequisites

- **Ansible 2.10+** - Install via `pip install ansible`
- **Python 3.8+** - For Ansible and kubernetes.core collection
- **Kubernetes Python Client** - Install via `pip install kubernetes`
- **macOS/Linux** - Currently configured for Unix-like systems
- **sudo access** - Some tasks require elevated privileges

## Installation

### 1. Install Ansible and Dependencies

```bash
# Create virtual environment
python3 -m venv ansible-env
source ansible-env/bin/activate  # On macOS/Linux
# or
ansible-env\Scripts\activate  # On Windows

# Install dependencies
pip install ansible kubernetes pyyaml
```

### 2. Install Ansible Collections

```bash
ansible-galaxy collection install kubernetes.core
```

## Usage

### Quick Start

```bash
# Run the complete bootstrap (organized layout)
ansible-playbook ansible/playbooks/cluster-setup.yaml

# Or with verbose output
ansible-playbook ansible_bootstrap_playbook.yaml -vv

# Or with debug output
ansible-playbook ansible_bootstrap_playbook.yaml -vvv
```

### Run Individual Roles

```bash
# Only prerequisites
ansible-playbook ansible_bootstrap_playbook.yaml --tags prerequisites

# Only kind cluster
ansible-playbook ansible_bootstrap_playbook.yaml --tags kind_cluster

# Only ArgoCD
ansible-playbook ansible_bootstrap_playbook.yaml --tags argocd

# Only Kyverno
ansible-playbook ansible_bootstrap_playbook.yaml --tags kyverno

# Only Sealed Secrets
ansible-playbook ansible_bootstrap_playbook.yaml --tags sealed_secrets
```

### Check Mode (Dry Run)

```bash
# See what would be done without making changes
ansible-playbook ansible_bootstrap_playbook.yaml --check
```

### Specific Hosts (localhost only in this case)

```bash
# Limit to specific host
ansible-playbook ansible_bootstrap_playbook.yaml -l localhost
```

## Playbook Structure

```
ansible/
├── ansible.cfg                        # Ansible configuration
├── ansible_bootstrap_playbook.yaml    # Main bootstrap playbook
├── ansible_inventory_localhost.yaml   # Inventory (localhost)
└── roles/
    ├── prerequisites/                 # Install tools
    │   ├── defaults/
    │   │   └── main.yaml             # Version defaults
    │   └── tasks/
    │       └── main.yaml             # Install tasks
    ├── kind_cluster/                  # Create local cluster
    │   ├── defaults/
    │   │   └── main.yaml
    │   └── tasks/
    │       └── main.yaml
    ├── argocd/                        # Install ArgoCD
    │   ├── defaults/
    │   │   └── main.yaml
    │   └── tasks/
    │       └── main.yaml
    ├── kyverno/                       # Install Kyverno policies
    │   ├── defaults/
    │   │   └── main.yaml
    │   └── tasks/
    │       └── main.yaml
    └── sealed_secrets/                # Install Sealed Secrets
        ├── defaults/
        │   └── main.yaml
        └── tasks/
            └── main.yaml
```

## What Happens After Bootstrap

After running the playbook, you'll have:

✅ **Local Kind cluster** named "hookdrop"  
✅ **ArgoCD** deployed and ready in `argocd` namespace  
✅ **Kyverno** policies enforcing security  
✅ **Sealed Secrets** for encrypted secret management  

### Next Steps

```bash
# 1. Apply the HookDrop Application manifest
kubectl apply -f k8s/argocd/application.yaml

# 2. Access ArgoCD UI (password printed in playbook output)
kubectl port-forward svc/argocd-server -n argocd 8081:443

# 3. Open browser to https://localhost:8081
# Username: admin
# Password: [shown in playbook output and stored]
```

## Troubleshooting

### Role-based Access Control (RBAC) Errors

If you get RBAC errors, ensure your kubeconfig is pointing to the right cluster:

```bash
kubectl config current-context
# Should show: kind-hookdrop
```

### Kind Cluster Already Exists

The playbook checks and deletes existing clusters before creating new ones:

```bash
# Manual cleanup if needed
kind delete cluster --name hookdrop
```

### ArgoCD Password Not Showing

If the password isn't captured:

```bash
kubectl -n argocd get secret argocd-initial-admin-secret \
  -o jsonpath="{.data.password}" | base64 -d
```

### Kyverno Policies Blocking Pods

This is intentional! Kyverno enforces:
- ❌ No `latest` image tags
- ❌ Missing resource limits

Update manifests to use specific tags and include `resources.limits`.

### Sealed Secrets Controller Not Ready

Wait a bit longer:

```bash
kubectl wait --for=condition=available --timeout=300s \
  deployment/sealed-secrets-controller -n kube-system
```

## Customization

### Change Kind Cluster Name

Edit `ansible_role_kind_cluster_defaults.yaml`:

```yaml
kind_cluster_name: "my-cluster"
```

### Change ArgoCD Port Forward

Edit `ansible_role_argocd_defaults.yaml`:

```yaml
argocd_port_forward: 9000  # Instead of 8081
```

### Skip Kyverno Policies

Comment out the Kyverno role in `ansible_bootstrap_playbook.yaml`:

```yaml
roles:
  - prerequisites
  - kind_cluster
  - argocd
  # - kyverno  # Skip this
  - sealed_secrets
```

## Advantages Over `setup-cluster.sh`

| Feature | Bash Script | Ansible |
|---------|-----------|---------|
| Idempotent | ❌ No | ✅ Yes |
| Error handling | ⚠️ Basic | ✅ Robust |
| Dry run support | ❌ No | ✅ Yes |
| Role-based execution | ❌ No | ✅ Yes |
| Cross-platform | ⚠️ Limited | ✅ macOS/Linux/Windows |
| Retry logic | ❌ No | ✅ Yes |
| Documentation | ⚠️ Inline | ✅ Separate |
| Extensibility | ⚠️ Hard | ✅ Easy |

## Maintenance

### Updating Tools

Tool versions are pinned in role defaults. To update:

```yaml
# In each role's defaults/main.yaml
kubectl_version: "latest"  # or specific version
helm_version: "v3.14.0"    # specific version
```

### Removing Resources

To tear down:

```bash
# Delete kind cluster
kind delete cluster --name hookdrop

# Cleanup files
rm -rf ~/.sealed-secrets ~/.kube/config-hookdrop
```

## Contributing

To extend the bootstrap:

1. Create a new role in `ansible/roles/`
2. Add defaults in `roles/new-role/defaults/main.yaml`
3. Add tasks in `roles/new-role/tasks/main.yaml`
4. Include in `ansible_bootstrap_playbook.yaml`
5. Document in this README

## Support

For issues:
- Check Ansible logs: `ansible-playbook ... -vvv`
- Review specific role files
- Consult [Kubernetes Ansible Collection](https://docs.ansible.com/ansible/latest/collections/kubernetes/core/index.html)
