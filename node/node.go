package node

import (
	"context"
	"errors"
	"fmt"
	"github.com/perlin-network/noise"
	"github.com/perlin-network/noise/kademlia"
	log "github.com/sirupsen/logrus"
	"net"
	"strconv"
	"strings"
	"time"
	"xblockchain/p2p/protocol"
)

type Node struct {
	p2pListenAddr string
	bootstrap string
}
func closehold(node *noise.Node) {
	if err := node.Close(); err !=nil{
		return
	}
}
func NewNode(listenAddr string, bootstrap string) *Node {
	n := &Node{
		p2pListenAddr: listenAddr,
		bootstrap: bootstrap,
	}
	return n
}

func (n *Node) Start() error {
	node,err := createLocalNode(n.p2pListenAddr)
	if err != nil {
		return err
	}
	defer closehold(node)
	node.RegisterMessage(protocol.Message{}, protocol.UnmarshalMessage)
	node.Handle(n.handle)
	overlay := n.createKademliaOverlay()
	node.Bind(overlay.Protocol())
	err = node.Listen()
	if err != nil {
		return err
	}
	log.Infof("start node: %s", node.ID())

	bootstrap(node, n.bootstrap)

	return nil
}

func (n *Node) createKademliaOverlay() *kademlia.Protocol {
	events := kademlia.Events{
		OnPeerAdmitted: n.handleOverlayOnPeerAdmitted,
		OnPeerEvicted: n.handleOverlayOnPeerEvicted,
	}
	return kademlia.New(kademlia.WithProtocolEvents(events))
}

func (n *Node) handleOverlayOnPeerAdmitted(id noise.ID) {
	fmt.Printf("Learned about a new peer %s(%s).\n", id.Address, id.ID.String()[:8])
}

func (n *Node) handleOverlayOnPeerEvicted(id noise.ID) {
	fmt.Printf("Forgotten a peer %s(%s).\n", id.Address, id.ID.String()[:8])
}


func createLocalNode(listenAddr string) (*noise.Node,error) {
	listenAddrSp := strings.Split(listenAddr,":")
	netIp := net.ParseIP(listenAddrSp[0])
	netPort, err := strconv.ParseInt(listenAddrSp[1],10, 16)
	if err != nil {
		return nil, err
	}
	nPort := uint16(netPort)
	node, err := noise.NewNode(
		noise.WithNodeBindHost(netIp),
		noise.WithNodeBindPort(nPort),
	)
	if err != nil {
		return nil, err
	}
	return node, err
}


func bootstrap(node *noise.Node, addresses ...string) {
	for _, addr := range addresses {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		_, err := node.Ping(ctx, addr)
		cancel()
		if err != nil {
			fmt.Printf("Failed to ping bootstrap node (%s). Skipping... [error: %s]\n", addr, err)
			continue
		}
	}
}

func (n *Node) handle(ctx noise.HandlerContext) error {
	if ctx.IsRequest() {
		return nil
	}
	obj, err := ctx.DecodeMessage()
	if err != nil {
		return err
	}
	data, ok := obj.(protocol.Message)
	if !ok {
		return errors.New("parse data err")
	}
	return n.handleProtocol(&data)
}


func (n *Node) handleProtocol(message *protocol.Message) error {
	header := message.Header
	t := header.Type
	if t == protocol.VersionT {
		return n.handleVersionMessage(message.Body)
	}
	return nil
}

func (n *Node) handleVersionMessage(data []byte) error {
	log.Infof("hello, request comming: %v", data)
	return nil
}

