#!/usr/bin/env bash

pwd="$(pwd)"
. "${pwd}/setup/ip.sh"

# Determine the current stable Kubernetes minor version (e.g., v1.30)
stable_minor_version=$(curl -sSL https://dl.k8s.io/release/stable.txt | sed -E 's/^(v[0-9]+\.[0-9]+)\..*$/\1/')

# Remote install command, with evaluated version
read -r -d '' cmd <<EOF
# Remove old Kubernetes source list
sudo rm -f /etc/apt/sources.list.d/kubernetes.list &&

# Install required packages
sudo apt-get update &&
sudo apt-get install -y apt-transport-https ca-certificates curl gnupg &&

# Create keyrings directory
sudo mkdir -p /etc/apt/keyrings &&

# Download and install Kubernetes APT key
curl -fsSL https://pkgs.k8s.io/core:/stable:/${stable_minor_version}/deb/Release.key \
  | sudo gpg --batch --yes --dearmor \
  -o /etc/apt/keyrings/kubernetes-archive-keyring.gpg &&

# Add Kubernetes APT repository
echo "deb [signed-by=/etc/apt/keyrings/kubernetes-archive-keyring.gpg] https://pkgs.k8s.io/core:/stable:/${stable_minor_version}/deb/ /" \
  | sudo tee /etc/apt/sources.list.d/kubernetes.list &&

# Install Kubernetes components and hold them
sudo apt-get update &&
sudo apt-get install -y kubelet kubeadm kubectl &&
sudo apt-mark hold kubelet kubeadm kubectl
EOF

${cmd}

# Execute on each replica
for host in "${replicas[@]}"; do
  echo "→ Installing Kubernetes on ${host}"
  ssh -i "${cert}" \
      -o StrictHostKeyChecking=no \
      -o UserKnownHostsFile=/dev/null \
      "${host}" \
      "${cmd}"
done

echo "✅ setup complete"