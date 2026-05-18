# 🎉 Welcome to HookDrop Ansible Bootstrap!

> **TL;DR**: Run `make ansible-setup && make cluster-up-ansible` and your local cluster is ready in 15 minutes.

---

## 📋 What Is This?

Complete **Ansible automation** for your HookDrop local Kubernetes cluster bootstrap. Replaces the manual `setup-cluster.sh` bash script with **idempotent, modular, production-grade** automation.

---

## 🚀 Get Started Now (3 Steps)

```bash
# 1. Install Ansible (one-time)
make ansible-setup

# 2. Bootstrap your cluster
make cluster-up-ansible

# 3. Done! Verify:
kubectl cluster-info
```

**Time needed**: ~15 minutes ⏱️

---

## 📚 Documentation

Choose your path:

### 🏗️ Clean Up & Organize Files
→ Read **ANSIBLE_FILE_NAMES.md** (3 min read) **← START HERE!**
- Short, memorable file names
- Clean directory structure  
- How to reorganize

### 👤 I'm a Developer - Get Me Started!
→ Read **ANSIBLE_QUICK_START.md** (5 min read)
- Quick commands
- Common troubleshooting
- Next steps

### 🔍 I Need Full Details
→ Read **ANSIBLE_BOOTSTRAP.md** (10 min read)
- Complete setup guide
- All usage examples
- Advanced troubleshooting
- Customization options

### 🏗️ I'm Reviewing the Architecture
→ Read **ANSIBLE_IMPLEMENTATION.md** (10 min read)
- Implementation overview
- What each role does
- Comparison vs bash script
- Timeline expectations

### ✅ I Need to Verify Everything Works
→ Read **ANSIBLE_CHECKLIST.md** (5 min read)
- Verification procedures
- Success criteria
- Quick reference

### 📄 Just Show Me the Summary
→ Read **DELIVERY_SUMMARY.txt**
- Visual overview
- Quick reference
- All details in one place

---

## ✨ What You Get

### 🛠️ Tools Installed
- Docker, Go 1.22+, kubectl, helm, kind, argocd-cli, trivy, golangci-lint

### ☸️ Kubernetes Components
- Local kind cluster
- ArgoCD (GitOps)
- Kyverno (policies)
- Sealed Secrets (encryption)

### 🎯 Automation Features
- Idempotent (safe to run multiple times)
- Modular (run specific roles with `--tags`)
- Debuggable (verbose output support)
- Documented (comprehensive guides)
- Fast (~10-15 minutes)

---

## 📊 Files Created

**25 files total**:
- 1 main playbook
- 5 reusable Ansible roles (10 files)
- 3 configuration files
- 3 setup scripts
- 4 documentation files
- Makefile updates

---

## 🎮 Available Commands

```bash
make ansible-setup              # Install Ansible
make cluster-up-ansible         # Run bootstrap
make cluster-up-ansible-check   # Dry-run (see what would happen)
make ansible-status             # Check if installed
make cluster-down               # Tear down cluster
```

---

## 🔄 5 Ansible Roles Included

| Role | What It Does | Time |
|------|------------|------|
| **Prerequisites** | Install all tools | 2-5 min |
| **Kind Cluster** | Create local K8s | 2-3 min |
| **ArgoCD** | Install GitOps controller | 2-3 min |
| **Kyverno** | Install policy engine | 1-2 min |
| **Sealed Secrets** | Install secret encryption | 1-2 min |

**Total**: ~10-15 minutes

---

## 🆚 Why Ansible Over Bash?

| Feature | Bash | Ansible |
|---------|------|---------|
| Idempotent | ❌ | ✅ |
| Dry-run | ❌ | ✅ |
| Error recovery | ⚠️ | ✅ |
| Partial execution | ❌ | ✅ |
| Debugging | ⚠️ | ✅ |
| Extensibility | ⚠️ | ✅ |

---

## 🆘 Quick Help

**Ansible not found?**
```bash
make ansible-setup
```

**Cluster already exists?**
```bash
kind delete cluster --name hookdrop
make cluster-up-ansible
```

**Want to see what would happen first?**
```bash
make cluster-up-ansible-check
```

**Need ArgoCD password?**
```bash
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d
```

---

## 📖 How to Use This Guide

1. **First time?** → ANSIBLE_QUICK_START.md
2. **Need details?** → ANSIBLE_BOOTSTRAP.md  
3. **Reviewing design?** → ANSIBLE_IMPLEMENTATION.md
4. **Checking off items?** → ANSIBLE_CHECKLIST.md
5. **Need overview?** → DELIVERY_SUMMARY.txt

---

## ✅ Next Steps After Bootstrap

```bash
# 1. Build and load your app image
docker build -t hookdrop:local .
kind load docker-image hookdrop:local --name hookdrop

# 2. Apply ArgoCD Application
kubectl apply -f k8s/argocd/application.yaml

# 3. Access ArgoCD UI
kubectl port-forward svc/argocd-server -n argocd 8081:443 &

# 4. Open browser
open https://localhost:8081
# Username: admin
# Password: [from playbook output]
```

---

## 🎓 Key Concepts

**Idempotent**: Safe to run multiple times. Playbook detects current state and only makes necessary changes.

**Modular Roles**: Each role is independent. Use `--tags` to run only what you need.

**Dry-run Mode**: Use `--check` flag to see what would happen without making changes.

**Verbose Debugging**: Use `-vv` or `-vvv` for detailed execution logs.

---

## 📦 Installation Breakdown

### Prerequisites Role (~2-5 min)
Installs: Docker, Go, kubectl, helm, kind, argocd-cli, trivy, golangci-lint
- Auto-detects what's already installed
- Skips unnecessary installations
- Verifies each installation

### Kind Cluster Role (~2-3 min)
- Checks for existing cluster
- Auto-deletes if found (idempotent)
- Creates new cluster from config
- Waits for nodes ready

### ArgoCD Role (~2-3 min)
- Creates namespace
- Installs ArgoCD manifests
- Waits for deployments
- Extracts admin password

### Kyverno Role (~1-2 min)
- Installs policy engine
- Applies policies
- Shows applied policies

### Sealed Secrets Role (~1-2 min)
- Installs Helm repo
- Installs Sealed Secrets
- Backs up encryption key locally

---

## 🎯 Success Indicators

After running `make cluster-up-ansible`, you should see:
- ✅ All playbook tasks completed
- ✅ Cluster created: `kind get clusters`
- ✅ All components running: `kubectl get pods -A`
- ✅ ArgoCD password in output
- ✅ Sealed Secrets key backed up

---

## 🔧 Customization

Edit role defaults to customize:

- Prerequisites versions: `ansible_role_prerequisites_defaults.yaml`
- Kind cluster name: `ansible_role_kind_cluster_defaults.yaml`
- ArgoCD settings: `ansible_role_argocd_defaults.yaml`
- Kyverno policies: Applied from `k8s/kyverno/`
- Sealed Secrets: `ansible_role_sealed_secrets_defaults.yaml`

---

## 📁 File Structure

```
project-root/
├── ansible/playbooks/cluster-setup.yaml      # Main playbook
├── ansible_role_*_tasks.yaml            # Role tasks (5 files)
├── ansible_role_*_defaults.yaml         # Role config (5 files)
├── ansible/inventory/default.yaml     # Inventory
├── ansible.cfg                          # Settings
├── ansible-requirements.txt             # Dependencies
├── setup-ansible.sh                     # Setup script
├── Makefile (UPDATED)                   # New targets
└── ANSIBLE_*.md                         # Documentation (4 files)
```

---

## 🚀 Advanced Usage

```bash
# Run only prerequisites
ansible-playbook ansible_playbook_bootstrap.yaml --tags prerequisites

# Skip Kyverno
ansible-playbook ansible_playbook_bootstrap.yaml --skip-tags kyverno

# Verbose debugging
ansible-playbook ansible_playbook_bootstrap.yaml -vvv

# Dry-run
ansible-playbook ansible_playbook_bootstrap.yaml --check
```

---

## 💡 Pro Tips

1. **First time?** Use `--check` mode to see what will happen
2. **Debugging?** Add `-vv` or `-vvv` for verbose output
3. **Partial run?** Use `--tags` to run specific roles
4. **Need password?** Check playbook output, or use `kubectl` command
5. **Cluster issues?** Delete with `kind delete cluster --name hookdrop` and retry

---

## 🆘 Support

1. Check **ANSIBLE_QUICK_START.md** for quick answers
2. Review **ANSIBLE_BOOTSTRAP.md** for detailed docs
3. Run playbook with `-vvv` for debugging
4. Check specific role files for implementation details

---

## 📞 Questions?

- **"How do I use this?"** → ANSIBLE_QUICK_START.md
- **"How does it work?"** → ANSIBLE_BOOTSTRAP.md
- **"What was delivered?"** → DELIVERY_SUMMARY.txt
- **"How do I verify it?"** → ANSIBLE_CHECKLIST.md
- **"What's the architecture?"** → ANSIBLE_IMPLEMENTATION.md

---

## ✨ Ready?

```bash
make ansible-setup
make cluster-up-ansible
```

Your cluster will be ready in ~15 minutes! ☕

---

**Happy clustering!** 🚀
