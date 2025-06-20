pwd=$(pwd)
. "${pwd}"/setup/ip.sh

cmd="sudo apt update; sudo apt install -y apt-transport-https ca-certificates curl gpg; sudo mkdir -p /etc/apt/keyrings; curl -fsSL https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-archive-keyring.gpg; echo "deb [signed-by=/etc/apt/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee /etc/apt/sources.list.d/kubernetes.list ; sudo apt update ;sudo apt install -y kubelet kubeadm kubectl; sudo apt-mark hold kubelet kubeadm kubectl"



for index in "${!replicas[@]}";
do
    sshpass ssh "${replicas[${index}]}" -i ${cert}  "${cmd}"
done

echo "setup complete"
