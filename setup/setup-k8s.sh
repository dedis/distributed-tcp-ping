#!/usr/bin/env bash
pwd="$(pwd)"
. "${pwd}/setup/ip.sh"

# Build the remote install command
read -r -d '' cmd <<'EOF'
# remove old repo
sudo rm -f /etc/apt/sources.list.d/kubernetes.list &&

# install deps
sudo apt-get update &&
sudo apt-get install -y apt-transport-https ca-certificates curl gnupg &&

# prepare keyring
sudo mkdir -p /etc/apt/keyrings &&

# fetch the community key and dearmor it
curl -fsSL https://pkgs.k8s.io/core:/stable:/$( 
    curl -L -s https://dl.k8s.io/release/stable.txt \
    | sed -E 's/^(v[0-9]+\.[0-9]+)\..*$/\1/' 
)/deb/Release.key \
  | sudo gpg --batch --yes --dearmor \
    -o /etc/apt/keyrings/kubernetes-archive-keyring.gpg &&

# add the new pkgs.k8s.io repository
echo "deb [signed-by=/etc/apt/keyrings/kubernetes-archive-keyring.gpg] \
  https://pkgs.k8s.io/core:/stable:/\$(
    curl -L -s https://dl.k8s.io/release/stable.txt \
    | sed -E 's/^(v[0-9]+\.[0-9]+)\..*$/\1/' 
  )/deb/ /" \
  | sudo tee /etc/apt/sources.list.d/kubernetes.list &&

# install and hold the Kubernetes packages
sudo apt-get update &&
sudo apt-get install -y kubelet kubeadm kubectl &&
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