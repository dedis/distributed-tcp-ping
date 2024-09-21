package dummy

import (
	"encoding/binary"
	"google.golang.org/protobuf/proto"
	"io"
)

type RPCPair struct {
	Code uint8
	Obj  Serializable
}

type RPCPairPeer struct {
	RpcPair *RPCPair
	Peer    int
}

/*
	each message sent over the network should implement this interface
*/

type Serializable interface {
	Marshal(io.Writer) error
	Unmarshal(io.Reader) error
	New() Serializable
}

/*
	A struct that allocates a unique uint8 for each message type
*/

type MessageCode struct {
	Ping uint8
	Pong uint8
}

/*
	A static function which assigns a unique uint8 to each message type
*/

func GetRPCCodes() MessageCode {
	return MessageCode{
		Ping: 1,
		Pong: 2,
	}
}

func marshalMessage(wire io.Writer, m proto.Message) error {
	data, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	lengthWritten := len(data)
	var b [8]byte
	bs := b[:8]
	binary.LittleEndian.PutUint64(bs, uint64(lengthWritten))
	_, err = wire.Write(bs)
	if err != nil {
		return err
	}
	_, err = wire.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func unmarshalMessage(wire io.Reader, m proto.Message) error {
	var b [8]byte
	bs := b[:8]
	_, err := io.ReadFull(wire, bs)
	if err != nil {
		return err
	}
	numBytes := binary.LittleEndian.Uint64(bs)
	data := make([]byte, numBytes)
	length, err := io.ReadFull(wire, data)
	if err != nil {
		return err
	}
	err = proto.Unmarshal(data[:length], m)
	if err != nil {
		return err
	}
	return nil
}

func (t *Ping) Marshal(wire io.Writer) error {
	return marshalMessage(wire, t)
}

func (t *Ping) Unmarshal(wire io.Reader) error {
	return unmarshalMessage(wire, t)
}

func (t *Ping) New() Serializable {
	return new(Ping)
}

func (t *Pong) Marshal(wire io.Writer) error {
	return marshalMessage(wire, t)
}

func (t *Pong) Unmarshal(wire io.Reader) error {
	return unmarshalMessage(wire, t)
}

func (t *Pong) New() Serializable {
	return new(Pong)
}
