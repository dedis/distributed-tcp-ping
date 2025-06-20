### Start Kubernetes control plane and advertise it on the master IP
sudo kubeadm init --apiserver-advertise-address=10.0.1.1 --pod-network-cidr=192.168.0.0/16

