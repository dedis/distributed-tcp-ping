package dummy

import (
	"bufio"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

// Proxy is the main struct of the proxy
type Proxy struct {
	name        int // unique node id
	numReplicas int

	addrList        map[int]string // map with the IP:port address of each remote replica
	incomingReaders map[int]*bufio.Reader
	outgoingWriters map[int]*bufio.Writer
	mutexes         map[int]*sync.Mutex // the mutexes for each remote replica writer

	serverAddress string // listening address of self

	debugOn    bool // if turned on, the debug messages will be print on the console
	debugLevel int  // debug level

	rttLatency    map[int]int64     // map to store the round trip time of each replica
	sentTimestamp map[int]time.Time // map to store the sent time of each message

	incomingChan chan *RPCPairPeer // all incoming messages are sent to this channel

	startTime time.Time

	rpcTable  map[uint8]*RPCPair // map each RPC type (message type) to its unique number
	statsChan chan map[int]int64
	Stats     *Stats
}

// NewProxy creates a new proxy object

func NewProxy(name int64, replicas []Replica, debugOn bool, debugLevel int) *Proxy {

	pr := Proxy{
		name:            int(name),
		numReplicas:     len(replicas),
		addrList:        make(map[int]string),
		incomingReaders: make(map[int]*bufio.Reader),
		outgoingWriters: make(map[int]*bufio.Writer),
		mutexes:         make(map[int]*sync.Mutex),
		serverAddress:   "",
		debugOn:         debugOn,
		debugLevel:      debugLevel,
		rttLatency:      make(map[int]int64),
		sentTimestamp:   make(map[int]time.Time),
		incomingChan:    make(chan *RPCPairPeer, 1000),
		startTime:       time.Now(),
		rpcTable:        make(map[uint8]*RPCPair),
		statsChan:       make(chan map[int]int64, len(replicas)),
	}
	pr.Stats = NewStats(pr.statsChan)

	// initialize the addrList

	for i := 0; i < pr.numReplicas; i++ {
		intName, _ := strconv.Atoi(replicas[i].Name)

		pr.addrList[intName] = replicas[i].IP
		pr.mutexes[intName] = &sync.Mutex{}
		pr.rttLatency[intName] = 0
		pr.sentTimestamp[intName] = time.Now()

		if pr.name == intName {
			pr.serverAddress = replicas[i].IP
		}
	}

	pr.RegisterRPC(&Ping{}, GetRPCCodes().Ping)
	pr.RegisterRPC(&Pong{}, GetRPCCodes().Pong)

	rand.Seed(time.Now().UTC().UnixNano())

	pr.debug("initialed a new proxy "+strconv.Itoa(pr.name), 0)

	return &pr
}

func (n *Proxy) RegisterRPC(msgObj Serializable, code uint8) {
	n.rpcTable[code] = &RPCPair{Code: code, Obj: msgObj}
	n.debug("Registered RPC code "+strconv.Itoa(int(code)), 0)
}

/*
	the main loop of the proxy
*/

func (pr *Proxy) Run() {
	for id, _ := range pr.addrList {
		pr.sendPing(id)
	}

	for true {
		m_object := <-pr.incomingChan
		pr.debug(fmt.Sprintf("recevied message %v", m_object), 0)
		if m_object.RpcPair.Code == GetRPCCodes().Ping {
			pr.debug(fmt.Sprintf("received ping from %v", m_object.Peer), 0)
			pr.handlePing(m_object)
		} else if m_object.RpcPair.Code == GetRPCCodes().Pong {
			pr.debug(fmt.Sprintf("received pong from %v", m_object.Peer), 0)
			pr.handlePong(m_object)
		} else {
			panic("unknown message type")
		}
	}

}

// debug prints the message to console if the debug is turned on

func (pr *Proxy) debug(s string, i int) {
	if pr.debugOn && i >= pr.debugLevel {
		fmt.Printf("%s\n", s)
	}
}

func (pr *Proxy) sendPing(id int) {
	pr.debug(fmt.Sprintf("sending ping to %v", id), 0)
	pr.sendMessage(&RPCPairPeer{
		RpcPair: &RPCPair{
			Code: GetRPCCodes().Ping,
			Obj: &Ping{
				Id: int32(pr.name),
			},
		},
		Peer: id,
	})
	pr.sentTimestamp[id] = time.Now()
}

func (pr *Proxy) handlePing(object *RPCPairPeer) {
	sender := object.Peer
	pr.debug(fmt.Sprintf("received ping from %v", sender), 0)
	pr.sendPong(sender)
}

func (pr *Proxy) sendPong(sender int) {
	pr.sendMessage(&RPCPairPeer{
		RpcPair: &RPCPair{
			Code: GetRPCCodes().Pong,
			Obj:  &Pong{Id: int32(pr.name)},
		},
		Peer: sender,
	})
	pr.debug(fmt.Sprintf("sent pong to %v", sender), 0)
}

func (pr *Proxy) handlePong(object *RPCPairPeer) {
	pr.debug(fmt.Sprintf("received pong from %v", object.Peer), 0)
	sender := object.Peer
	pr.rttLatency[sender] = time.Since(pr.sentTimestamp[sender]).Microseconds()
	select {
	case pr.statsChan <- pr.rttLatency:
		pr.debug(fmt.Sprintf("sent rtt latency to stats %v", pr.rttLatency), 0)
	default:
		pr.debug(fmt.Sprintf("stats channel is full, dropping rtt latency %v", pr.rttLatency), 0)
	}
	pr.sendPing(sender)
}
