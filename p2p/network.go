package p2p

import (
	"arcoin/core"
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"sync"
)

const (
	Protocol      = "tcp"
	NodeVersion   = 1
	CommandLength = 12
)

type Node struct {
	Address     string
	Peers       map[string]bool
	Blockchain  *core.Blockchain
	Server      net.Listener
	Mutex       sync.Mutex
}

type Message struct {
	Command string
	Payload []byte
}

func StartNode(nodeID, minerAddress string) *Node {
	node := &Node{
		Address:    nodeID,
		Peers:      make(map[string]bool),
		Blockchain: core.NewBlockchain(),
	}
	
	go node.StartServer()
	node.SyncWithNetwork()
	return node
}

func (n *Node) StartServer() {
	server, err := net.Listen(Protocol, n.Address)
	if err != nil {
		panic(err)
	}
	defer server.Close()
	n.Server = server

	for {
		conn, err := server.Accept()
		if err != nil {
			continue
		}
		go n.HandleConnection(conn)
	}
}

func (n *Node) HandleConnection(conn net.Conn) {
	defer conn.Close()
	request := &Message{}
	gob.NewDecoder(conn).Decode(request)
	
	switch request.Command {
	case "version":
		n.HandleVersion(request)
	case "block":
		n.HandleBlock(request)
	case "tx":
		n.HandleTX(request)
	case "getblocks":
		n.HandleGetBlocks(request)
	case "inv":
		n.HandleInv(request)
	case "getdata":
		n.HandleGetData(request)
	}
}

func (n *Node) SendMessage(peer string, msg *Message) {
	conn, err := net.Dial(Protocol, peer)
	if err != nil {
		delete(n.Peers, peer)
		return
	}
	defer conn.Close()
	
	gob.NewEncoder(conn).Encode(msg)
}

// Blockchain synchronization logic
func (n *Node) SyncWithNetwork() {
	for peer := range n.Peers {
		n.SendGetBlocks(peer)
	}
}

func (n *Node) BroadcastBlock(block *core.Block) {
	for peer := range n.Peers {
		data := block.Serialize()
		msg := &Message{Command: "block", Payload: data}
		n.SendMessage(peer, msg)
	}
}

func (n *Node) BroadcastTX(tx *core.Transaction) {
	for peer := range n.Peers {
		data := tx.Serialize()
		msg := &Message{Command: "tx", Payload: data}
		n.SendMessage(peer, msg)
	}
}