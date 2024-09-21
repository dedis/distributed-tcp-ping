import sys

if len(sys.argv) < 3:
    exit("not enough arguments")

numReplicas = int(sys.argv[1])

replicaIPs = []


if len(sys.argv) < 2 + numReplicas:
    exit("not enough arguments")

argC = 2
for i in range(argC, argC + numReplicas, 1):
    replicaIPs.append(sys.argv[i])


def print_peers(replicaIPs):
    print("peers:")
    for j in range(1, 1 + numReplicas, 1):
        print("   - name: " + str(j))
        print("     address: " + str(replicaIPs[j - 1])+str(":")+str(10000))


print_peers(replicaIPs)
