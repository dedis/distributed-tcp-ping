### Start Kubernetes control plane and advertise it on the master IP
sudo kubeadm init --apiserver-advertise-address=10.0.1.1 --pod-network-cidr=192.168.0.0/16


# Create local kubeconfig for kubectl to access the cluster
mkdir -p $HOME/.kube
sudo cp /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

# Download Calico pod network plugin (required for inter-pod communication)
curl -O https://raw.githubusercontent.com/projectcalico/calico/v3.26.1/manifests/calico.yaml

# Apply it to the cluster
kubectl apply -f calico.yaml


# Temporarily (takes effect immediately)
sudo sysctl -w net.ipv4.ip_forward=1

# Persist across reboots
echo "net.ipv4.ip_forward=1" | sudo tee /etc/sysctl.d/99-kubernetes-ip-forward.conf
sudo sysctl --system

