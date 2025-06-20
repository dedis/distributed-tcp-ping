cert="/home/pasindut/.ssh/id_mrg-0"

replica1_name=10.0.1.2
replica1=pasindut@${replica1_name}

replica2_name=10.0.1.3
replica2=pasindut@${replica2_name}

replica3_name=10.0.1.4
replica3=pasindut@${replica3_name}

replica4_name=10.0.1.5
replica4=pasindut@${replica4_name}

username="pasindut"

replicas=(${replica1} ${replica2} ${replica3} ${replica4})
replica_names=(${replica1_name} ${replica2_name} ${replica3_name} ${replica4_name})
