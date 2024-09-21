package dummy

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"
)

// start listening to the proxy tcp connections

func (pr *Proxy) NetworkInit() {
	pr.waitForConnections()
}

/*
	Listen on the server ports for new connections
*/

func (pr *Proxy) waitForConnections() {

	for i := 0; i < len(pr.serverAddress); i++ {

		go func(la string) {

			var b [4]byte
			bs := b[:4]

			listener, err := net.Listen("tcp", la)
			if err != nil {
				panic(err.Error())
			}
			fmt.Printf("Listening to messages on " + la + "\n")

			for true {
				conn, err := listener.Accept()
				if err != nil {
					panic(err.Error())
				}
				pr.debug("Received incoming tcp connection from someone, id yet not read", 12)
				if _, err := io.ReadFull(conn, bs); err != nil {
					panic(err.Error())
				}
				id := int32(binary.LittleEndian.Uint16(bs))
				pr.debug("Received incoming tcp connection from "+strconv.Itoa(int(id)), 12)

				go pr.connectionListener(bufio.NewReader(conn), id)
				pr.debug("Started listening to "+strconv.Itoa(int(id)), 12)
			}
		}(pr.serverAddress[i])
	}
}

/*
	listen to a given connection. Upon receiving any message, put it into the central buffer
*/

func (pr *Proxy) connectionListener(reader *bufio.Reader, id int32) {

	var err error = nil
	for true {
		obj := (&Message{}).New()
		if err = obj.Unmarshal(reader); err != nil {
			//pr.debug("Error while unmarshalling", 0)
			return
		}
		pr.incomingChan <- &ReceivedMessage{message: obj, sender: id}
	}
}

/*
	make TCP connections to other replicas
*/

func (pr *Proxy) ConnectToReplicas() {

	for id, addresses := range pr.addrList {
		for i := 0; i < len(addresses); i++ {
			var b [4]byte
			bs := b[:4]

			for true {
				conn, err := net.Dial("tcp", addresses[i])
				if err == nil {
					pr.outgoingWriters[id] = append(pr.outgoingWriters[id], bufio.NewWriter(conn))
					pr.mutexes[id] = append(pr.mutexes[id], &sync.Mutex{})
					binary.LittleEndian.PutUint16(bs, uint16(pr.name))
					_, err := conn.Write(bs)
					if err != nil {
						//pr.debug("Error connecting to client "+strconv.Itoa(int(id)), 0)
						panic(err)
					}
					pr.debug("Started outgoing tcp connection to "+addresses[i], 12)
					break
				} else {
					pr.debug("failed to connect to "+addresses[i], 12)
					time.Sleep(time.Duration(10) * time.Millisecond)
				}
			}
		}
	}

}

/*
	write a message to the wire
*/

func (pr *Proxy) sendMessage(peer int64, msg *Message) {

	pr.debug("sending message to  "+strconv.Itoa(int(peer)), 0)

	randomWriter := rand.Intn(len(pr.outgoingWriters[peer]) + 1)
	if randomWriter == len(pr.outgoingWriters[peer]) {
		randomWriter--
	}
	w := pr.outgoingWriters[peer][randomWriter]
	m := pr.mutexes[peer][randomWriter]

	m.Lock()
	err := msg.Marshal(w)
	if err != nil {
		pr.debug("Error while marshalling", 0)
		m.Unlock()
		return
	}
	err = w.Flush()
	if err != nil {
		pr.debug("Error while flushing", 0)
		m.Unlock()
		return
	}
	pr.debug("sent message to  "+strconv.Itoa(int(peer)), 12)
	m.Unlock()
}
