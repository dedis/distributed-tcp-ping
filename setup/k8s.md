### in master
sudo kubeadm init --apiserver-advertise-address=10.0.1.1 --pod-network-cidr=192.168.0.0/16
mkdir -p $HOME/.kube
sudo cp /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
kubectl apply -f https://raw.githubusercontent.com/projectcalico/calico/v3.27.0/manifests/calico.yaml

# run in workers, seperately
sudo kubeadm join 10.0.1.1:6443 --token d28asf.7y9nt8s4bck8nvqg \
	--discovery-token-ca-cert-hash sha256:b270e271fed40f8973722d4eca90a2d449bf79eae5397dfef01ceb8440b879df
	
# in master 
kubectl get pods -n kube-system -o wide
kubectl get nodes