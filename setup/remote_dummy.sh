pwd=$(pwd)
. "${pwd}"/setup/ip.sh

rm nohup.out

dummy_path="/dummy/dummy"
config="/home/${username}/dummy/dedus-config.yaml"

local_output_path="logs/dummy/"
rm -r "${local_output_path}"; mkdir -p "${local_output_path}"

for index in "${!replicas[@]}";
do
  sshpass ssh "${replicas[${index}]}"  -i ${cert}  "pkill -f dummy; pkill -f dummy"
done

echo "Killed previously running instances"

echo "starting dummy replicas"

nohup ssh ${replica1}  -i ${cert}   ".${dummy_path} --config ${config} --name 1 --debugOn --debugLevel 0">"${local_output_path}"1.log &
nohup ssh ${replica2}  -i ${cert}   ".${dummy_path} --config ${config} --name 2 --debugOn --debugLevel 0">"${local_output_path}"2.log &
nohup ssh ${replica3}  -i ${cert}   ".${dummy_path} --config ${config} --name 3 --debugOn --debugLevel 0">"${local_output_path}"3.log &
nohup ssh ${replica4}  -i ${cert}   ".${dummy_path} --config ${config} --name 4 --debugOn --debugLevel 0">"${local_output_path}"4.log &
nohup ssh ${replica5}  -i ${cert}   ".${dummy_path} --config ${config} --name 5 --debugOn --debugLevel 0">"${local_output_path}"5.log &

echo "Started dummy replicas"

sleep  60

for index in "${!replicas[@]}";
do
  sshpass ssh "${replicas[${index}]}"  -i ${cert}  "pkill -f dummy; pkill -f dummy"
done

echo "Killed dummy replicas"