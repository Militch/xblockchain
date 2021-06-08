package backend

import (
	"github.com/perlin-network/noise"
	"time"
	"xblockchain"
	"xblockchain/p2p"
)

type handler struct {
	handlerCallFn func(peer *p2p.Peer) error
	newPeerCh chan *peer
	peers map[noise.PublicKey] *peer
	blockchain *xblockchain.BlockChain
	version uint32
	network uint32
}


func newHandler(bc *xblockchain.BlockChain, pv uint32, nv uint32) (*handler,error) {
	h := &handler{
		newPeerCh: make(chan *peer, 1),
		peers: make(map[noise.PublicKey] *peer),
		blockchain: bc,
		version: pv,
		network: nv,
	}
	h.handlerCallFn = h.handleNewPeer
	return h, nil
}

func (h *handler) handleNewPeer(p2p *p2p.Peer) error {
	p := newPeer(p2p, h.version, h.network)
	h.newPeerCh <- p
	return h.handle(p)
}

func (h *handler) handle(p *peer) error {
	head := h.blockchain.GetLastBlockHash()
	if err := p.Handshake(head); err != nil {
		return err
	}
	id := p.p2p().ID
	h.peers[id.ID] = p
	defer delete(h.peers, id.ID)
	for {
		if err := h.handleMsg(p); err!= nil {
			return err
		}
	}
}

func  (h *handler) handleMsg(p *peer) error {
	msg := <-p.p2pPeer.GetProtocolMsgCh()
	msgCode := msg.Header.MsgCode
	//TODO: 这里处理P2P链路消息
	switch msgCode {
	case NewBlockMsg:
		// 处理区块广播
	case GetBlockMsg:
		// 获取本机区块

	}
	return nil
}

func (h *handler) syncer() {
	forceSync := time.Tick(10 * time.Second)
	for {
		select {
		case <-h.newPeerCh:
			if len(h.peers) < 5 {
				break
			}
			go h.synchronise(h.basePeer())
		case <-forceSync:
			// 强制同步
			go h.synchronise(h.basePeer())
		}
	}
}
func (h *handler) basePeer() *peer {
	for _, v := range h.peers {
		return v
	}
	return nil
}

func (h *handler) synchronise(p *peer) {
	if p == nil {
		return
	}
	//TODO: 这里处理链路同步
	// 1. 寻找共同父块
	// 2. 确定同步区间
	// 3. 下载区块头、区块体，并持久化存储

}
func (h *handler) Start() {
	go h.syncer()
}

