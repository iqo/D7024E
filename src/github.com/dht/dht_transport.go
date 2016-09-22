package dht

import (
	"fmt"
	"json"
	"net"
)

//connect upd

type Msg struct {
	Key   string //värdet
	Src   string //från noden som kalla
	Dst   string //destinationsadress
	Bytes []byte //transport funktionen, msg.Bytes

}

type Transport struct {
	node        *DHTNode
	bindAddress string // rad 20, bindadress måste finnas.
}

func (transport *Transport) listen() {
	udpAddr, err := net.ResolveUDPAddr("udp", transport.bindAddress)
	conn, err := net.Listen("udp", udpAddr)
	defer conn.Close()
	dec := json.NewDecoder(conn)
	for {
		msg = Msg{}
		err = dec.Decode(&msg)
	}
}
