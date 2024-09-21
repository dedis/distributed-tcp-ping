package dummy

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strconv"
	"strings"
)

/*
	config.go defines the structs and methods to pass the configuration file, that contains the IP:ports of each dummy replica
*/

type Replica struct {
	Name string
	IP   string
}

// Config represents the structure of the YAML file.
type Config struct {
	Peers []struct {
		Name    string `yaml:"name"`
		Address string `yaml:"address"`
	} `yaml:"peers"`
}

// generate a config object from the given file

func ReadYAML(fileName string, name int) ([]Replica, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err.Error())
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err.Error())
	}

	// Create an array of Replicas
	var replicas []Replica
	for _, peer := range config.Peers {
		if strconv.Itoa(name) == peer.Name {
			peer.Address = "0.0.0.0:" + GetPort(peer.Address)
		}

		replicas = append(replicas, Replica{
			Name: peer.Name,
			IP:   peer.Address,
		})
	}

	return replicas, nil
}

func GetPort(address string) string {
	parts := strings.Split(address, ":")
	if len(parts) != 2 {
		panic("invalid address format")
	}
	return parts[1]
}

// generate a config object from the given file without modifying the address

func ReadYAMLNoModify(fileName string) ([]Replica, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err.Error())
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err.Error())
	}

	var replicas []Replica
	for _, peer := range config.Peers {
		replicas = append(replicas, Replica{
			Name: peer.Name,
			IP:   peer.Address,
		})
	}
	return replicas, nil
}
