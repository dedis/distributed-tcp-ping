package dummy

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Stats struct {
	in chan map[int]int64
}

func NewStats(in chan map[int]int64) *Stats {
	return &Stats{
		in: in,
	}
}

func (s *Stats) statsHandler(w http.ResponseWriter, r *http.Request) {
	select {
	case stats := <-s.in:
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(stats)
		if err != nil {
			http.Error(w, "Unable to encode JSON", http.StatusInternalServerError)
		}
	default:
		http.Error(w, "No stats available", http.StatusNotFound)
	}

}

func (s *Stats) Start() {
	http.HandleFunc("/stats", s.statsHandler)

	port := 8080 // Default port
	fmt.Printf("Starting http server on port %d...\n", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
