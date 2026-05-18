# 🚀 Ansible Local Cluster Bootstrap - Quick Start

## TL;DR - Get Started Now

```bash
# 1. Setup Ansible (one-time)
make ansible-setup

# 2. Bootstrap your local cluster
make cluster-up-ansible

# 3. Done! Your cluster is ready
```

---

## What You Get

✅ Local Kubernetes cluster (kind)  
✅ ArgoCD installed and ready  
✅ Kyverno policies for security  
✅ Sealed Secrets for safe secret management  
✅ All prerequisites (Docker, kubectl, helm, etc.)  

---

## Files Overview

| File | Purpose |
|------|---------|
| `ansible_playbook_bootstrap.yaml` | Main playbook - orchestrates everything |
| `ansible_role_*_tasks.yaml` | Individual role tasks |
| `ansible_role_*_defaults.yaml` | Role configuration/defaults |
| `ansible/inventory/default.yaml` | Inventory (just localhost) |
| `ansible.cfg` | Ansible configuration |
| `ansible-requirements.txt` | Python dependencies |
| `setup-ansible.sh` | Setup script for Ansible env |
| `ANSIBLE_BOOTSTRAP.md` | Detailed documentation |

---

## Available Make Commands

```bash
# Setup Ansible (one-time)
make ansible-setup

# Run bootstrap
make cluster-up-ansible

# Dry-run (see what would happen)
make cluster-up-ansible-check

# Check if Ansible is installed
make ansible-status

# Tear down cluster
make cluster-down

# Old way (bash script)
make cluster-up
```

---

## Manual Usage (If Not Using Make)

```bash
# Activate environment
source ansible-env/bin/activate

# Run playbook
ansible-playbook ansible_playbook_bootstrap.yaml

# Dry-run
ansible-playbook ansible_playbook_bootstrap.yaml --check

# Verbose output
ansible-playbook ansible_playbook_bootstrap.yaml -vv

# Run only specific role
ansible-playbook ansible_playbook_bootstrap.yaml --tags prerequisites
```

---

## Troubleshooting

### Ansible Not Found
```bash
make ansible-setup
```

### Permission Errors
```bash
# On macOS with Docker Desktop
sudo chown -R $(whoami):staff ~/.docker ~/.kube
```

### Cluster Already Exists
The playbook auto-detects and recreates it. Or manually:
```bash
kind delete cluster --name hookdrop
```

### ArgoCD Password
```bash
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d
```

---

## What Happens Under the Hood

```
1. Prerequisites Role
   ├─ Installs Docker
   ├─ Installs kubectl
   ├─ Installs helm
   ├─ Installs kind
   ├─ Installs argocd-cli
   ├─ Installs trivy
   └─ Installs golangci-lint

2. Kind Cluster Role
   ├─ Checks for existing cluster
   ├─ Deletes old cluster if found
   ├─ Creates new cluster
   └─ Waits for nodes to be ready

3. ArgoCD Role
   ├─ Creates argocd namespace
   ├─ Installs ArgoCD
   ├─ Waits for deployments
   └─ Saves password

4. Kyverno Role
   ├─ Installs Kyverno
   ├─ Waits for webhooks
   └─ Applies policies (no latest tags, require limits)

5. Sealed Secrets Role
   ├─ Adds Helm repo
   ├─ Installs Sealed Secrets
   ├─ Waits for controller
   ├─ Backs up encryption key locally
   └─ Displays setup info
```

---

## Next Steps After Bootstrap

```bash
# 1. Load local docker image into cluster
docker build -t hookdrop:local .
kind load docker-image hookdrop:local --name hookdrop

# 2. Apply ArgoCD Application
kubectl apply -f k8s/argocd/application.yaml

# 3. Port-forward ArgoCD
kubectl port-forward svc/argocd-server -n argocd 8081:443 &

# 4. Get ArgoCD password (shown in playbook output)
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d

# 5. Open browser
open https://localhost:8081
# Username: admin
# Password: [from step 4]
```

---

## Directory Structure (After Organization)

```
.
├── ansible.cfg
├── ansible_playbook_bootstrap.yaml
├── ansible_inventory_localhost.yaml
├── ansible_role_*/tasks.yaml
├── ansible_role_*/defaults.yaml
├── ansible-requirements.txt
├── setup-ansible.sh
├── ANSIBLE_BOOTSTRAP.md
├── ANSIBLE_QUICK_START.md (this file)
├── Makefile (updated with ansible targets)
└── ... (other project files)
```

**Pro Tip**: Run `organize-ansible.sh` to move everything into proper `ansible/` directory structure.

---

## Comparing: Bash vs Ansible

| Task | Bash (`scripts/setup-cluster.sh`) | Ansible |
|------|-----------------------------------|---------|
| Setup time | Manual - install tools first | Automated |
| Idempotent | No - will error if run twice | Yes - safe to run multiple times |
| Partial recovery | Hard - have to start over | Easy - just re-run |
| Customization | Edit script | Edit defaults in roles |
| Testing | No dry-run | `--check` mode |
| Documentation | Inline comments | Separate docs |
| Debugging | Shell output | Ansible verbosity levels |

---

## Performance

Expected timing:
- Prerequisites installation: **2-5 minutes** (depends on internet/OS)
- Kind cluster creation: **2-3 minutes**
- ArgoCD install: **2-3 minutes**
- Kyverno install: **1-2 minutes**
- Sealed Secrets: **1-2 minutes**

**Total: ~10-15 minutes** for full bootstrap

---

## Support

For detailed docs, see: `ANSIBLE_BOOTSTRAP.md`

For issues:
1. Run with verbosity: `ansible-playbook ... -vvv`
2. Check specific role tasks
3. Review error messages carefully
4. Consult Kubernetes documentation

---

## Notes

- Only supports macOS and Linux (bash-based systems)
- Requires Python 3.8+
- Playbook is idempotent (safe to run multiple times)
- All changes are logged and reversible
- Virtual environment is isolated from system Python
