#!/usr/bin/env bash
set -euo pipefail

# Load replicas[] and cert from setup/ip.sh
pwd="$(pwd)"
. "${pwd}/setup/ip.sh"

# Define the remote install command
read -r -d '' cmd <<'EOF'
sudo apt update &&
sudo apt install -y apt-transport-https ca-certificates curl gpg &&
sudo mkdir -p /etc/apt/keyrings &&
curl -fsSL https://packages.cloud.google.com/apt/doc/apt-key.gpg |
  sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-archive-keyring.gpg &&
echo "deb [signed-by=/etc/apt/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main" |
  sudo tee /etc/apt/sources.list.d/kubernetes.list &&
sudo apt update &&
sudo apt install -y kubelet kubeadm kubectl &&
sudo apt-mark hold kubelet kubeadm kubectl
EOF

# Loop over all replicas and run the install
for host in "${replicas[@]}"; do
  echo "→ Installing Kubernetes on ${host}"
  ssh -i "${cert}" \
      -o StrictHostKeyChecking=no \
      -o UserKnownHostsFile=/dev/null \
      "${host}" \
      "${cmd}"
done

echo "✅ setup complete"
