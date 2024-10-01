cert="/home/tennage/Pictures/pasindu"

replica1_name=10.156.33.141
replica1=pasindu@${replica1_name}

replica2_name=10.156.33.142
replica2=pasindu@${replica2_name}

replica3_name=10.156.33.143
replica3=pasindu@${replica3_name}

replica4_name=10.156.33.144
replica4=pasindu@${replica4_name}

replica5_name=10.156.33.145
replica5=pasindu@${replica5_name}

replica6_name=10.156.33.146
replica6=pasindu@${replica6_name}

replicas=(${replica1} ${replica2} ${replica3} ${replica4} ${replica5} ${replica6})
replica_names=(${replica1_name} ${replica2_name} ${replica3_name} ${replica4_name} ${replica5_name} ${replica6_name})

username="pasindu"