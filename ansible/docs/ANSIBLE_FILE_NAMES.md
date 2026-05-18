# 📋 Ansible File Organization Reference

## Quick Reorganize

```bash
# Reorganize flat files to clean structure with short names
make ansible-reorganize

# Verify it worked
make cluster-up-ansible-check

# Run bootstrap
make cluster-up-ansible
```

---

## Name Mapping (Old → New)

| Old Name | New Name | Path | Purpose |
|----------|----------|------|---------|
| prerequisites | **prep** | `ansible/roles/prep/` | Install tools |
| kind_cluster | **cluster** | `ansible/roles/cluster/` | Create K8s |
| argocd | **argocd** | `ansible/roles/argocd/` | GitOps |
| kyverno | **policy** | `ansible/roles/policy/` | Policies |
| sealed_secrets | **secrets** | `ansible/roles/secrets/` | Encryption |
| bootstrap | **cluster-setup** | `ansible/playbooks/cluster-setup.yaml` | Main playbook |
| localhost.yaml | **default.yaml** | `ansible/inventory/default.yaml` | Inventory |

---

## New Directory Structure

```
ansible/
├── ansible.cfg                      Configuration
├── inventory/
│   └── default.yaml                 Inventory (localhost)
├── playbooks/
│   └── cluster-setup.yaml          Main playbook
└── roles/
    ├── prep/                        Install tools
    │   ├── tasks/main.yaml
    │   └── defaults/main.yaml
    ├── cluster/                     Create cluster
    │   ├── tasks/main.yaml
    │   └── defaults/main.yaml
    ├── argocd/                      Install ArgoCD
    │   ├── tasks/main.yaml
    │   └── defaults/main.yaml
    ├── policy/                      Install Kyverno
    │   ├── tasks/main.yaml
    │   └── defaults/main.yaml
    └── secrets/                     Install Sealed Secrets
        ├── tasks/main.yaml
        └── defaults/main.yaml
```

---

## Usage After Reorganization

### From project root:
```bash
make ansible-setup              # One-time setup
make cluster-up-ansible         # Run bootstrap
make cluster-up-ansible-check   # Dry-run
make ansible-reorganize         # (If re-organizing)
```

### From ansible/ directory:
```bash
cd ansible

# Run full bootstrap
ansible-playbook playbooks/cluster-setup.yaml

# Dry-run
ansible-playbook playbooks/cluster-setup.yaml --check

# By role tags
ansible-playbook playbooks/cluster-setup.yaml --tags prep
ansible-playbook playbooks/cluster-setup.yaml --tags cluster
ansible-playbook playbooks/cluster-setup.yaml --tags argocd
ansible-playbook playbooks/cluster-setup.yaml --tags policy
ansible-playbook playbooks/cluster-setup.yaml --tags secrets

# Verbose
ansible-playbook playbooks/cluster-setup.yaml -vvv
```

---

## Tag Quick Reference

| Tag | Role | What It Does |
|-----|------|------------|
| `prep` | prep | Docker, kubectl, helm, kind, tools |
| `cluster` | cluster | Create local Kubernetes cluster |
| `argocd` | argocd | Install ArgoCD controller |
| `policy` | policy | Install Kyverno policies |
| `secrets` | secrets | Install Sealed Secrets |

Example:
```bash
# Only install tools
ansible-playbook playbooks/cluster-setup.yaml --tags prep

# Skip policies
ansible-playbook playbooks/cluster-setup.yaml --skip-tags policy
```

---

## Files to Delete After Reorganization

After verifying the new structure works:

```bash
# Old flat files (no longer needed)
rm -f ansible_*.yaml
rm -f ansible.cfg

# Old setup scripts
rm -f organize-ansible.sh
rm -f setup-ansible-dirs.sh
```

---

## Before & After Examples

### BEFORE (Flat)
```bash
ansible-playbook ansible_playbook_bootstrap.yaml
ansible-playbook ansible_playbook_bootstrap.yaml --tags prerequisites
```

### AFTER (Organized)
```bash
ansible-playbook ansible/playbooks/cluster-setup.yaml
ansible-playbook ansible/playbooks/cluster-setup.yaml --tags prep
```

---

## Why These Short Names?

| Name | Reason |
|------|--------|
| `prep` | Short for "prepare/prerequisites" |
| `cluster` | The main action: create a cluster |
| `argocd` | Clear component name |
| `policy` | What Kyverno does: enforce policies |
| `secrets` | What Sealed Secrets do: manage secrets |
| `cluster-setup` | The playbook sets up a cluster |
| `default.yaml` | Default inventory (standard Ansible) |

---

## Rollback (If Needed)

```bash
# Undo with git
git checkout HEAD -- .

# Or manually move files back:
# ansible/playbooks/cluster-setup.yaml → ansible_playbook_bootstrap.yaml
# ansible/inventory/default.yaml → ansible_inventory_localhost.yaml
# etc.
```

---

## Git Commit After Reorganization

```bash
git add -A
git commit -m "Reorganize Ansible files: flat → organized structure with short names

- Move files to ansible/ directory
- Rename roles: prerequisites→prep, kind_cluster→cluster, kyverno→policy, sealed_secrets→secrets
- Rename playbook: bootstrap→cluster-setup
- Rename inventory: localhost.yaml→default.yaml
- Update Makefile to support both structures
- All functionality unchanged"
```

---

## Comparison

| Aspect | Old (Flat) | New (Organized) |
|--------|-----------|-----------------|
| Files in root | 12+ flat files | All in ansible/ |
| Role names | Long (prerequisites) | Short (prep) |
| Playbook | Long (ansible_playbook_bootstrap.yaml) | Short (cluster-setup.yaml) |
| Structure | Non-standard | Standard Ansible pattern |
| IDE support | ⚠️ Hard to navigate | ✅ Clear hierarchy |
| Scalability | 🟡 Limited | ✅ Easy to extend |
| Git clarity | 🟡 Messy | ✅ Clean |

---

## Status Check

```bash
# Check which structure you have
make ansible-status

# Output: 
# ✓ Ansible installed
# ✓ Organized structure detected (ansible/playbooks/cluster-setup.yaml)
# or
# ⚠ Flat structure detected (ansible_playbook_bootstrap.yaml) - Run 'make ansible-reorganize'
```

---

**Ready to reorganize?**

```bash
make ansible-reorganize
```

Clean, organized, easy to remember! ✅
