package dummy

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"time"
)

// start listening to the proxy tcp connections

func (pr *Proxy) NetworkInit() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		pr.waitForConnections()
		wg.Done()
	}()

	go func() {
		pr.ConnectToReplicas()
		wg.Done()
	}()
	wg.Wait()
	pr.debug("Network initialized", 0)
}

/*
	Listen on the server ports for new connections
*/

func (pr *Proxy) waitForConnections() {

	counter := 0

	var b [4]byte
	bs := b[:4]

	listener, err := net.Listen("tcp", pr.serverAddress)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Listening to messages on " + pr.serverAddress + "\n")

	for counter < pr.numReplicas {
		conn, err := listener.Accept()
		if err != nil {
			panic(err.Error())
		}
		pr.debug("Received incoming tcp connection from someone, id yet not read", 0)
		if _, err := io.ReadFull(conn, bs); err != nil {
			panic(err.Error())
		}
		id := int32(binary.LittleEndian.Uint16(bs))
		pr.debug("Received incoming tcp connection from "+strconv.Itoa(int(id)), 0)
		pr.incomingReaders[int(id)] = bufio.NewReader(conn)
		go pr.connectionListener(bufio.NewReader(conn), id)
		pr.debug("Started listening to "+strconv.Itoa(int(id)), 0)
		counter++
	}

}

/*
	listen to a given connection. Upon receiving any message, put it into the central buffer
*/

func (pr *Proxy) connectionListener(reader *bufio.Reader, id int32) {

	var msgType uint8
	var err error = nil

	for true {
		if msgType, err = reader.ReadByte(); err != nil {
			pr.debug("Error while reading message code: connection broken from "+strconv.Itoa(int(id))+fmt.Sprintf(" %v", err.Error()), 3)
			return
		}
		if rpair, present := pr.rpcTable[msgType]; present {
			obj := rpair.Obj.New()
			if err = obj.Unmarshal(reader); err != nil {
				pr.debug("Error while unmarshalling from "+strconv.Itoa(int(id))+fmt.Sprintf(" %v", err.Error()), 3)
				return
			}
			pr.incomingChan <- &RPCPairPeer{
				RpcPair: &RPCPair{
					Code: msgType,
					Obj:  obj,
				},
				Peer: int(id),
			}
			pr.debug("Pushed a message from "+strconv.Itoa(int(id)), 0)
		} else {
			pr.debug("Error received unknown message type from "+strconv.Itoa(int(id)), 3)
			return
		}
	}
}

/*
	make TCP connections to other replicas
*/

func (pr *Proxy) ConnectToReplicas() {

	for id, addresses := range pr.addrList {
		var b [4]byte
		bs := b[:4]

		for true {
			conn, err := net.Dial("tcp", addresses)
			if err == nil {
				pr.outgoingWriters[id] = bufio.NewWriter(conn)
				pr.mutexes[id] = &sync.Mutex{}
				binary.LittleEndian.PutUint16(bs, uint16(pr.name))
				_, err := conn.Write(bs)
				if err != nil {
					panic(err)
				}
				pr.debug("Started outgoing tcp connection to "+addresses, 0)
				break
			} else {
				pr.debug("failed to connect to "+addresses, 0)
				time.Sleep(time.Duration(100) * time.Millisecond)
			}
		}

	}

}

/*
	write a message to the wire
*/

func (pr *Proxy) sendMessage(msg *RPCPairPeer) {

	peer := int32(msg.Peer)
	messageCode := msg.RpcPair.Code
	message := msg.RpcPair.Obj

	pr.debug("sending message to  "+strconv.Itoa(int(peer)), 0)

	w := pr.outgoingWriters[int(peer)]
	m := pr.mutexes[int(peer)]

	m.Lock()
	err := w.WriteByte(messageCode)
	if err != nil {
		pr.debug("Error writing message code byte:"+err.Error(), 0)
		m.Unlock()
		return
	}

	err = message.Marshal(w)
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
	pr.debug("sent message to  "+strconv.Itoa(int(peer)), 0)
	m.Unlock()
}
