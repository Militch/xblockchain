package backend

import (
	"errors"
	"time"
	"xblockchain/p2p"
)

type peer struct {
	p2pPeer *p2p.Peer
	version uint32
	head uint32
}

var MsgCodeVersion = uint32(3)
var GetBlockMsg = uint32(4)
var NewBlockMsg = uint32(5)

func newPeer(p *p2p.Peer) *peer {
	pt := &peer{
		p2pPeer: p,
	}
	return pt
}

func (p *peer) p2p() *p2p.Peer {
	return p.p2pPeer
}

func (p *peer) Handshake() error {
	go func() {
		if err := p2p.SendMsg(p.p2pPeer, MsgCodeVersion, nil); err != nil {
			return
		}
	}()
	for  {
		select {
		case msg := <-p.p2pPeer.GetProtocolMsgCh():
			msgCode := msg.Header.MsgCode
			switch msgCode {
			case MsgCodeVersion:
				if err := p2p.SendMsg(p.p2pPeer, MsgCodeVersion, nil); err != nil {
					return err
				}
				return nil
			}
		case <-time.After(3 * 60 * time.Second):
			return errors.New("time out")
		}
	}
}