package xblockchain

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"xblockchain/kademlia"
)

type NetServer struct {
	rt *kademlia.RoutingTable
}


func NewNetServer() *NetServer {
	rt := kademlia.NewRoutingTable()

	return &NetServer {
		rt: rt,
	}
}

func (netServer *NetServer) handleConnection(conn net.Conn) {
	for {
		buf := [1024]byte{}
		reader := bufio.NewReader(conn)
		n, err := reader.Read(buf[:])
		if err != nil {
			break
		}
		p := kademlia.NewProtocol(buf[:n])
		if p.HasTypeFlag(kademlia.TypeFlagReqHello) && !p.HasTypeFlag(kademlia.TypeFlagAckHello)  {
			// ping

			//netServer.rt.GetBucketFor()
		}else if p.HasTypeFlag(kademlia.TypeFlagReqHello | kademlia.TypeFlagAckHello) {
			//
		}
		resp := string(buf[:n])
		fmt.Printf("Read from conn success, data: %v\n", resp)
	}
	err := conn.Close()
	if err != nil {
		fmt.Printf("Close server connection errors: %v\n", err)
		return
	}
}


func (netServer *NetServer) Start(ip string, port string) {
	address := ip + ":" + port
	server, err := net.Listen("tcp4", address)
	if err != nil {
		fmt.Printf("Listen server errors: %v", err)
		os.Exit(1)
	}
	defer server.Close()
	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Printf("Create server connection errors: %s\n", err)
			os.Exit(1)
		}
		go netServer.handleConnection(conn)
	}
}