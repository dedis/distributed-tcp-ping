#!/usr/bin/env bash
pwd="$(pwd)"
. "${pwd}/setup/ip.sh"

read -r -d '' cmd <<'EOF'
# remove any stale kubernetes.list
sudo rm -f /etc/apt/sources.list.d/kubernetes.list &&

# add the Google keyring and repository for Focal
sudo apt update &&
sudo apt install -y apt-transport-https ca-certificates curl gpg &&
sudo mkdir -p /etc/apt/keyrings &&
curl -fsSL https://packages.cloud.google.com/apt/doc/apt-key.gpg \
  | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-archive-keyring.gpg &&
echo "deb [signed-by=/etc/apt/keyrings/kubernetes-archive-keyring.gpg] \
  https://apt.kubernetes.io/ kubernetes-focal main" \
  | sudo tee /etc/apt/sources.list.d/kubernetes.list &&

# install Kubernetes components
sudo apt update &&
sudo apt install -y kubelet kubeadm kubectl &&
sudo apt-mark hold kubelet kubeadm kubectl
EOF

for host in "${replicas[@]}"; do
  echo "→ Installing Kubernetes on ${host}"
  ssh -i "${cert}" \
      -o StrictHostKeyChecking=no \
      -o UserKnownHostsFile=/dev/null \
      "${host}" \
      "${cmd}"
done

echo "✅ setup complete"
