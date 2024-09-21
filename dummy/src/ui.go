package dummy

import (
	"encoding/json"
	"github.com/rs/cors"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Metrics struct {
	Throughput float64 `json:"throughput"`
	Latency    float64 `json:"latency"`
	CPUUsage   float64 `json:"cpuUsage"`
	MemUsage   float64 `json:"memUsage"`
}

var (
	metrics Metrics
	mu      sync.Mutex
)

func ListenFrontEnd(name string, pr *Proxy) {

	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", handleMetrics)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Allow requests from any origin
		AllowedMethods: []string{"GET"},
	})

	handler := c.Handler(mux)

	go generateMetrics(pr)

	log.Println("Server starting on :" + name)
	log.Fatal(http.ListenAndServe(":"+name, handler))

}

func handleMetrics(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	//fmt.Printf("%v\n", r.Header)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

func generateMetrics(pr *Proxy) {
	for {
		mu.Lock()
		metrics = Metrics{
			Throughput: float64(pr.ui_stats.throughput),
			Latency:    float64(pr.ui_stats.latency),
			CPUUsage:   float64(pr.ui_stats.cpu),
			MemUsage:   float64(pr.ui_stats.mem),
		}
		mu.Unlock()
		time.Sleep(1 * time.Second)
	}
}

func DoUi(pr *Proxy) {
	ListenFrontEnd(strconv.FormatInt(pr.name*10000+200, 10), pr)
}
