pwd=$(pwd)
. "${pwd}"/setup/ip.sh

rm -r dummy/bin/dummy
/bin/bash build.sh
echo "Removed old binaries and built project"

firewall="sudo ufw disable"
kill_instances="pkill -f dummy ;  pkill -f stats"
remote_home_path="/home/${username}/dummy/"

for index in "${!replicas[@]}";
do
    echo "copying files to replica ${index}"
    sshpass ssh "${replicas[${index}]}" -i ${cert} "rm -r ${remote_home_path}; mkdir -p ${remote_home_path}; ${kill_instances}; ${firewall}"
    scp -i ${cert} "dummy/bin/dummy" "${replicas[${index}]}":${remote_home_path}
    scp -i ${cert} "stats/bin/stats" "${replicas[${index}]}":${remote_home_path}
    scp -i ${cert} "dedis-config.yaml" "${replicas[${index}]}":${remote_home_path}
done

echo "setup complete"