package backend

import (
	"time"
	"xblockchain/p2p"
)

type handler struct {
	handlerCallFn func(peer *p2p.Peer) error
	newPeerCh chan *peer
	peers map[string] *peer
}


func newHandler() (*handler,error) {
	h := &handler{
		newPeerCh: make(chan *peer, 1),
		peers: make(map[string] *peer),
	}
	h.handlerCallFn = h.handleNewPeer
	return h, nil
}

func (h *handler) handleNewPeer(p2p *p2p.Peer) error {
	p := newPeer(p2p)
	h.newPeerCh <- p
	return h.handle(p)
}

func (h *handler) handle(p *peer) error {
	if err := p.Handshake(); err != nil {
		return err
	}
	id := p.p2p().ID
	h.peers[id.Address] = p
	for {
		if err := h.handleMsg(p); err!= nil {
			return err
		}
	}
}

func  (h *handler) handleMsg(p *peer) error {
	msg := <-p.p2pPeer.GetProtocolMsgCh()
	msgCode := msg.Header.MsgCode
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
			go h.synchronise(h.basePeer())
		}
	}
}
func (h *handler) basePeer() *peer {
	return nil
}

func (h *handler) synchronise(p *peer) {
	if p == nil {
		return
	}
}
func (h *handler) Start() {
	go h.syncer()
}

