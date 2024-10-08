package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	dummy "ping-ping/dummy/src"
	"sort"
	"strings"
	"time"
)

// Stats represents the map[int]int64 structure from the server response
type Stats map[int]int64

func main() {

	configFile := flag.String("config", "dedis-config.yaml", "configuration file")

	flag.Parse()

	replicas, err := dummy.ReadYAMLNoModify(*configFile)
	if err != nil {
		panic(err.Error())
	}
	for i := 0; i < len(replicas); i++ {
		replica_name := replicas[i].Name
		replica_ip := strings.Split(replicas[i].IP, ":")[0]
		replicas[i] = dummy.Replica{Name: replica_name, IP: replica_ip + ":8080"}
	}

	for true {
		for i := 0; i < len(replicas); i++ {
			// Fetch the stats from each server
			stats, err := fetchStats(replicas[i].IP)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}

			keys := make([]int, 0, len(stats))
			for k := range stats {
				keys = append(keys, k)
			}

			sort.Ints(keys)

			fmt.Printf("\nStats from server %s -- ", replicas[i].Name)
			for _, k := range keys {
				fmt.Printf("%d: %d, ", k, stats[k])
			}
			fmt.Println()
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// fetchStats sends a request to the /stats endpoint of a given server and returns the stats
func fetchStats(server string) (Stats, error) {
	url := fmt.Sprintf("http://%s/stats", server)

	// Create an HTTP client with a timeout
	client := &http.Client{Timeout: 5 * time.Second}

	// Perform the GET request
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stats from server %s: %v", server, err)
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server %s returned non-OK status: %s", server, resp.Status)
	}

	// Decode the JSON response into the Stats map
	var stats Stats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, fmt.Errorf("failed to decode response from server %s: %v", server, err)
	}

	return stats, nil
}
