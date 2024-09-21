package dummy

import (
	"bufio"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type UI_Stats struct {
	throughput float32
	latency    float32
	cpu        float32
	mem        float32
}

// Proxy is the main struct of the proxy
type Proxy struct {
	name        int64 // unique node id
	numReplicas int

	addrList        map[int64][]string        // map with the IP:port address of each remote replica
	outgoingWriters map[int64][]*bufio.Writer // the ip ports of each remote replica
	mutexes         map[int64][]*sync.Mutex   // the mutexes for each remote replica writer

	serverAddress []string // listening address of self

	debugOn    bool // if turned on, the debug messages will be print on the console
	debugLevel int  // debug level

	lastStarTime    time.Time // last time when stats were printed
	receivedLatency []int64   // list of latencies of received messages since the lastStartTime

	sent sync.Map // to save the set of sent messages

	incomingChan chan *ReceivedMessage // all incoming messages are sent to this channel

	counter int64 // to create unique message index

	ui_stats UI_Stats

	startTime time.Time

	interArrivalTime int
}

// each time a new request is sent, a Request object is created and stored in sent

type Request struct {
	sentTime time.Time
}

// ReceivedMessage is a struct to store the received message and the sender id

type ReceivedMessage struct {
	message Serializable
	sender  int32
}

// NewProxy creates a new proxy object

func NewProxy(name int64, cfg InstanceConfig, debugOn bool, debugLevel int, interArrivalTime int) *Proxy {

	pr := Proxy{
		name:            name,
		numReplicas:     len(cfg.Peers),
		addrList:        make(map[int64][]string),
		outgoingWriters: make(map[int64][]*bufio.Writer),
		mutexes:         make(map[int64][]*sync.Mutex),
		serverAddress:   []string{},
		debugOn:         debugOn,
		debugLevel:      debugLevel,
		receivedLatency: make([]int64, 0),
		incomingChan:    make(chan *ReceivedMessage, 100000),
		counter:         0,
		ui_stats: UI_Stats{
			throughput: 0,
			latency:    0,
			cpu:        0,
			mem:        0,
		},
		interArrivalTime: interArrivalTime,
	}

	// initialize the addrList

	for i := 0; i < pr.numReplicas; i++ {
		intName, _ := strconv.Atoi(cfg.Peers[i].Name)
		addresses := make([]string, 0)
		for j := 0; j < len(cfg.Peers[i].PORTS); j++ {
			addresses = append(addresses, cfg.Peers[i].IP+":"+cfg.Peers[i].PORTS[j])
		}
		if pr.name != int64(intName) {
			pr.addrList[int64(intName)] = addresses
			pr.outgoingWriters[int64(intName)] = make([]*bufio.Writer, 0)
			pr.mutexes[int64(intName)] = make([]*sync.Mutex, 0)
		}
		if pr.name == int64(intName) {
			pr.serverAddress = addresses
		}
	}

	rand.Seed(time.Now().UTC().UnixNano())

	pr.debug("initialed a new proxy "+strconv.Itoa(int(pr.name)), 0)

	return &pr
}

/*
	the main loop of the proxy
*/

func (pr *Proxy) Run() {
	go func() {
		for true {
			m_object := <-pr.incomingChan
			pr.debug("Received message from "+string(m_object.sender), 12)
			pr.handleMessage(m_object.message.(*Message), m_object.sender)
		}
	}()

}

// debug prints the message to console if the debug is turned on

func (pr *Proxy) debug(s string, i int) {
	if pr.debugOn && i >= pr.debugLevel {
		fmt.Printf("%s\n", s)
	}
}
