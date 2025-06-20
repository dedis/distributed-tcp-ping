#!/usr/bin/env bash
# Get current directory and load replicas[] & cert from setup/ip.sh
pwd="$(pwd)"
. "${pwd}/setup/ip.sh"

# Build the remote install command, with --batch --yes to prevent GPG from touching /dev/tty
read -r -d '' cmd <<'EOF'
# Remove any old kubernetes.list
sudo rm -f /etc/apt/sources.list.d/kubernetes.list &&

# Install prerequisites and add Google keyring
sudo apt update &&
sudo apt install -y apt-transport-https ca-certificates curl gpg &&
sudo mkdir -p /etc/apt/keyrings &&
curl -fsSL https://packages.cloud.google.com/apt/doc/apt-key.gpg \
  | sudo gpg --batch --yes --dearmor -o /etc/apt/keyrings/kubernetes-archive-keyring.gpg &&

# Add the Kubernetes repo for Ubuntu Focal
echo "deb [signed-by=/etc/apt/keyrings/kubernetes-archive-keyring.gpg] \
  https://apt.kubernetes.io/ kubernetes-focal main" \
  | sudo tee /etc/apt/sources.list.d/kubernetes.list &&

# Install and hold kubelet, kubeadm, kubectl
sudo apt update &&
sudo apt install -y kubelet kubeadm kubectl &&
sudo apt-mark hold kubelet kubeadm kubectl
EOF

# Loop over each replica host and run the install command
for host in "${replicas[@]}"; do
  echo "→ Installing Kubernetes on ${host}"
  ssh -i "${cert}" \
      -o StrictHostKeyChecking=no \
      -o UserKnownHostsFile=/dev/null \
      "${host}" \
      "${cmd}"
done

echo "✅ setup complete"
