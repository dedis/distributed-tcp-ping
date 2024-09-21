package dummy

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

/*
	config.go defines the structs and methods to pass the configuration file, that contains the IP:ports of each dummy replica
*/

type ReplicaInstance struct {
	Name  string
	IP    string
	PORTS []string
}

// InstanceConfig describes the set of replicas
type InstanceConfig struct {
	Peers []ReplicaInstance
}

// NewInstanceConfig loads an instance configuration from given file
func NewInstanceConfig(fname string, name int64) (*InstanceConfig, error) {
	cfg := InstanceConfig{
		Peers: []ReplicaInstance{},
	}

	file, err := os.Open(fname)
	if err != nil {
		panic(err.Error())
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err.Error())
		}
	}(file)

	var lines []string

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Iterate over each line
	for scanner.Scan() {
		// Append each line to the slice
		lines = append(lines, scanner.Text())
	}

	// Check for any errors encountered during scanning
	if err := scanner.Err(); err != nil {
		panic(err.Error())
	}

	// Iterate over each line
	for _, line := range lines {
		// Split the line by the space character
		parts := strings.Split(line, " ")

		// Create a new ReplicaInstance
		peer := ReplicaInstance{
			Name:  parts[0],
			IP:    parts[1],
			PORTS: parts[2:],
		}

		// Append the new ReplicaInstance to the configuration
		cfg.Peers = append(cfg.Peers, peer)
	}

	// set the self ip to 0.0.0.0
	cfg = configureSelfIP(cfg, name)
	return &cfg, nil
}

/*
	Replace the IP of my self to 0.0.0.0
*/

func configureSelfIP(cfg InstanceConfig, name int64) InstanceConfig {
	for i := 0; i < len(cfg.Peers); i++ {
		if cfg.Peers[i].Name == strconv.Itoa(int(name)) {
			cfg.Peers[i].IP = "0.0.0.0"
			return cfg
		}
	}
	return cfg
}
