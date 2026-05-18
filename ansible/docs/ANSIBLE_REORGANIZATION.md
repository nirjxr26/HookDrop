# 🏗️ Ansible Reorganization Guide

## What We're Doing

Converting from **flat file structure** to **clean organized structure** with **short, memorable names**.

---

## Old vs New Structure

### OLD (Current - Flat)
```
project-root/
├── ansible_playbook_bootstrap.yaml
├── ansible_role_prerequisites_tasks.yaml
├── ansible_role_prerequisites_defaults.yaml
├── ansible_role_kind_cluster_tasks.yaml
├── ansible_role_kind_cluster_defaults.yaml
├── ansible_role_argocd_tasks.yaml
├── ansible_role_argocd_defaults.yaml
├── ansible_role_kyverno_tasks.yaml
├── ansible_role_kyverno_defaults.yaml
├── ansible_role_sealed_secrets_tasks.yaml
├── ansible_role_sealed_secrets_defaults.yaml
├── ansible_inventory_localhost.yaml
└── ansible.cfg
```

### NEW (Organized)
```
project-root/
└── ansible/
    ├── ansible.cfg
    ├── inventory/
    │   └── default.yaml              (was: ansible_inventory_localhost.yaml)
    ├── playbooks/
    │   └── cluster-setup.yaml         (was: ansible_playbook_bootstrap.yaml)
    └── roles/
        ├── prep/                      (was: prerequisites)
        │   ├── tasks/main.yaml
        │   └── defaults/main.yaml
        ├── cluster/                   (was: kind_cluster)
        │   ├── tasks/main.yaml
        │   └── defaults/main.yaml
        ├── argocd/
        │   ├── tasks/main.yaml
        │   └── defaults/main.yaml
        ├── policy/                    (was: kyverno)
        │   ├── tasks/main.yaml
        │   └── defaults/main.yaml
        └── secrets/                   (was: sealed_secrets)
            ├── tasks/main.yaml
            └── defaults/main.yaml
```

---

## Naming Convention

| Old Name | New Name | Purpose |
|----------|----------|---------|
| prerequisites | **prep** | Install tools & prerequisites |
| kind_cluster | **cluster** | Create local Kubernetes cluster |
| argocd | **argocd** | GitOps deployment |
| kyverno | **policy** | Policy enforcement |
| sealed_secrets | **secrets** | Secret encryption |
| bootstrap | **cluster-setup** | Main playbook |
| localhost.yaml | **default.yaml** | Default inventory |

---

## How to Reorganize

### Option 1: Automatic (Recommended)
```bash
bash reorganize-ansible.sh
```

This script:
- Creates the directory structure
- Moves all files to new locations
- Renames files with short names
- Updates ansible.cfg paths
- Shows final structure

### Option 2: Manual Steps

```bash
# Create structure
mkdir -p ansible/playbooks
mkdir -p ansible/roles/{prep,cluster,argocd,policy,secrets}/{tasks,defaults}
mkdir -p ansible/inventory

# Move inventory
mv ansible_inventory_localhost.yaml ansible/inventory/default.yaml

# Move playbook (if needed)
# Original flat playbook: ansible_playbook_bootstrap.yaml
# Organized playbook: ansible/playbooks/cluster-setup.yaml

# Move prerequisites role (prep)
mv ansible_role_prerequisites_defaults.yaml ansible/roles/prep/defaults/main.yaml
mv ansible_role_prerequisites_tasks.yaml ansible/roles/prep/tasks/main.yaml

# Move kind cluster role (cluster)
mv ansible_role_kind_cluster_defaults.yaml ansible/roles/cluster/defaults/main.yaml
mv ansible_role_kind_cluster_tasks.yaml ansible/roles/cluster/tasks/main.yaml

# Move argocd role
mv ansible_role_argocd_defaults.yaml ansible/roles/argocd/defaults/main.yaml
mv ansible_role_argocd_tasks.yaml ansible/roles/argocd/tasks/main.yaml

# Move kyverno role (policy)
mv ansible_role_kyverno_defaults.yaml ansible/roles/policy/defaults/main.yaml
mv ansible_role_kyverno_tasks.yaml ansible/roles/policy/tasks/main.yaml

# Move sealed secrets role (secrets)
mv ansible_role_sealed_secrets_defaults.yaml ansible/roles/secrets/defaults/main.yaml
mv ansible_role_sealed_secrets_tasks.yaml ansible/roles/secrets/tasks/main.yaml

# Move config
mv ansible.cfg ansible/ansible.cfg
```

---

## Update ansible.cfg

After moving, update `ansible/ansible.cfg`:

```ini
[defaults]
inventory = inventory/default.yaml
host_key_checking = False
roles_path = ./roles
# ... rest stays the same
```

---

## Update Makefile

Update `Makefile` in project root:

```makefile
cluster-up-ansible:
	@command -v ansible-playbook >/dev/null 2>&1 || { bash setup-ansible.sh; }
	@source ansible-env/bin/activate 2>/dev/null || true && \
	  ansible-playbook ansible/playbooks/cluster-setup.yaml

cluster-up-ansible-check:
	@command -v ansible-playbook >/dev/null 2>&1 || { bash setup-ansible.sh; }
	@source ansible-env/bin/activate 2>/dev/null || true && \
	  ansible-playbook ansible/playbooks/cluster-setup.yaml --check

ansible-setup:
	bash setup-ansible.sh

ansible-status:
	@command -v ansible-playbook >/dev/null 2>&1 && echo "✓ Ansible ready" || echo "✗ Install: make ansible-setup"
```

---

## After Reorganization

### Run from project root:
```bash
cd /path/to/hookdrop

# Setup (one-time)
make ansible-setup

# Run bootstrap
make cluster-up-ansible

# Or dry-run
make cluster-up-ansible-check
```

### Or run from ansible directory:
```bash
cd ansible

# Run playbook
ansible-playbook playbooks/cluster-setup.yaml

# Dry-run
ansible-playbook playbooks/cluster-setup.yaml --check

# Only install tools
ansible-playbook playbooks/cluster-setup.yaml --tags prep

# Only create cluster
ansible-playbook playbooks/cluster-setup.yaml --tags cluster

# Skip policies
ansible-playbook playbooks/cluster-setup.yaml --skip-tags policy
```

---

## Tag Reference

Use with `--tags` or `--skip-tags`:

| Tag | Role | What It Does |
|-----|------|------------|
| `prep` | prep | Install Docker, kubectl, helm, kind, etc. |
| `cluster` | cluster | Create local kind cluster |
| `argocd` | argocd | Install ArgoCD |
| `policy` | policy | Install Kyverno policies |
| `secrets` | secrets | Install Sealed Secrets |

---

## File Organization Benefits

✅ **Cleaner directory** - Everything in one place
✅ **Standard pattern** - Follows Ansible best practices
✅ **Shorter names** - Easier to type and remember
✅ **Easy navigation** - Clear folder hierarchy
✅ **Git-friendly** - Proper structure for version control
✅ **Scalable** - Easy to add more roles/playbooks
✅ **IDE support** - Better editor/IDE recognition

---

## Before Reorganizing

**Commit your work!**

```bash
git add .
git commit -m "Add Ansible bootstrap automation (before reorganization)"
```

---

## After Reorganizing

```bash
# Remove old flat files (after verifying new structure)
rm -f ansible_*.yaml

# Commit reorganization
git add -A
git commit -m "Reorganize Ansible files into proper structure with short names"

# Verify it works
make cluster-up-ansible-check
make cluster-up-ansible
```

---

## Cleanup

After reorganization, remove old files:

```bash
# Old flat files (after verifying new structure works)
rm -f ansible_*.yaml
rm -f ansible_playbook_bootstrap.yaml
rm -f ansible_role_*.yaml
rm -f ansible_inventory_localhost.yaml
rm -f ansible.cfg  # If moved to ansible/ansible.cfg

# Clean old setup scripts
rm -f setup-ansible-dirs.sh  # No longer needed
rm -f organize-ansible.sh    # Redundant with reorganize-ansible.sh
```

---

## Verification

After reorganization, verify structure:

```bash
# Check directory structure
find ansible -type f | sort

# Should show:
# ansible/ansible.cfg
# ansible/inventory/default.yaml
# ansible/playbooks/cluster-setup.yaml
# ansible/roles/prep/tasks/main.yaml
# ansible/roles/prep/defaults/main.yaml
# ... (and other roles)
```

---

## Rollback (If Needed)

If something goes wrong:

```bash
# Git has your back!
git checkout HEAD -- ansible_*.yaml ansible.cfg

# Or if you didn't commit yet
git status  # See what's new/moved
```

---

## Documentation Updates

After reorganization, update:

1. **README.md** - Point to `ansible/README.md`
2. **ANSIBLE_QUICK_START.md** - Update file paths
3. **ANSIBLE_BOOTSTRAP.md** - Update examples
4. **ANSIBLE_IMPLEMENTATION.md** - Update file listing
5. **Makefile** - Already covered above

---

## Quick Reference: New Commands

```bash
# From project root
make ansible-setup                                    # Install Ansible
make cluster-up-ansible                              # Run bootstrap
make cluster-up-ansible-check                        # Dry-run

# From ansible/ directory
cd ansible
ansible-playbook playbooks/cluster-setup.yaml        # Run bootstrap
ansible-playbook playbooks/cluster-setup.yaml -vvv   # Debug
```

---

## Summary

**What changes:**
- File organization (flat → organized structure)
- Naming conventions (long → short)
- Directory structure (standard Ansible pattern)

**What stays the same:**
- All functionality
- All automation
- All role logic
- Makefile compatibility

---

**Ready to reorganize?**

```bash
bash reorganize-ansible.sh
```

Then verify:

```bash
make cluster-up-ansible-check
make cluster-up-ansible
```

Enjoy your clean Ansible setup! 🎉
