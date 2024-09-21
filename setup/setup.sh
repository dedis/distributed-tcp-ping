pwd=$(pwd)
. "${pwd}"/setup/ip.sh

rm ping-node/bin/ping-node
rm torture/bin/torture
echo "Removed old binaries"

mkdir -p temp

cp -r ping-node/                     temp
cp -r torture/                   temp
cp build.sh                      temp
cp go.mod                   temp
cp go.sum                   temp

zip -r temp.zip temp/
rm -r temp

kill_instances="pkill dummy ; pkill torture; pkill dummy ; pkill torture; pkill dummy ; pkill torture; pkill dummy ; pkill torture; pkill dummy ; pkill torture; pkill dummy ; pkill torture"

local_zip_path="temp.zip"
remote_home_path="/home/${username}/torture/"
reset_directory="rm -rf /home/${username}/torture; mkdir -p /home/${username}/torture/"
confs="sudo setcap cap_net_admin,cap_net_raw+ep $(which tc); sudo apt-get install libcap2-bin; sudo setcap cap_net_admin=eip  /usr/sbin/xtables-nft-multi; sudo modprobe ip_tables; sudo modprobe nfnetlink_queue"

go_command="sudo rm -rf /usr/local/go ; cd /home/${username}/torture/temp/ && wget https://go.dev/dl/go1.19.5.linux-amd64.tar.gz   && sudo tar -C /usr/local -xzf go1.19.5.linux-amd64.tar.gz && export PATH=\$PATH:/usr/local/go/bin && go version; export PATH=\$PATH:/usr/local/go/bin"
build_command="export PATH=\$PATH:/usr/local/go/bin; sudo apt-get update && sudo apt-get install -y pkg-config libnetfilter-queue-dev locales && sudo locale-gen de_CH.UTF-8 && sudo update-locale LANG=de_CH.UTF-8 LC_ALL=de_CH.UTF-8 ; cd /home/${username}/torture/temp/; go get -u github.com/AkihiroSuda/go-netfilter-queue; go get github.com/rs/cors; go get google.golang.org/protobuf; go get github.com/google/go-cmp; go get github.com/google/gopacket; /bin/bash build.sh"


command="sudo apt-get update; sudo apt-get install unzip;sudo apt-get install zip; cd /home/${username}/torture && unzip temp.zip"

for index in "${!replicas[@]}";
do
    echo "copying files to replica ${index}"
    sshpass ssh "${replicas[${index}]}" -i ${cert} "${reset_directory};${kill_instances}; ${confs}"
    scp -i ${cert} ${local_zip_path} "${replicas[${index}]}":${remote_home_path}
    sshpass ssh "${replicas[${index}]}" -i ${cert} "${command}; ${build_command};"
done

rm ${local_zip_path}

echo "setup complete"