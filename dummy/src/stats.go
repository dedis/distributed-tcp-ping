package dummy

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Stats struct {
	pr *Proxy
}

func NewStats(pr *Proxy) *Stats {
	return &Stats{
		pr: pr,
	}
}

func (s *Stats) statsHandler(w http.ResponseWriter, r *http.Request) {
	stats := s.pr.GetRtt()

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(stats)
	if err != nil {
		http.Error(w, "Unable to encode JSON", http.StatusInternalServerError)
	}
	//fmt.Println(fmt.Sprintf("Stats published to http: %v", stats))

}

func (s *Stats) Start() {
	http.HandleFunc("/stats", s.statsHandler)

	port := 8080 // Default port
	fmt.Printf("Starting http server on port %d...\n", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
