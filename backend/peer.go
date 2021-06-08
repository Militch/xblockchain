package backend

import (
	"encoding/json"
	"errors"
	"time"
	"xblockchain/p2p"
	"xblockchain/uint256"
)

type peer struct {
	p2pPeer *p2p.Peer
	version uint32
	network uint32
	head uint256.UInt256
}

var MsgCodeVersion = uint32(3)
var GetBlockMsg = uint32(4)
var NewBlockMsg = uint32(5)

func newPeer(p *p2p.Peer, version uint32, network uint32) *peer {
	pt := &peer{
		p2pPeer: p,
		version: version,
		network: network,
	}
	return pt
}

func (p *peer) p2p() *p2p.Peer {
	return p.p2pPeer
}

type statusData struct {
	Version uint32
	Network uint32
	Head *uint256.UInt256
}

func (p *peer) Handshake(head *uint256.UInt256) error {
	go func() {
		if err := p2p.SendMsgJSONData(p.p2pPeer, MsgCodeVersion, &statusData{
			Version: p.version,
			Network: p.network,
			Head: head,
		}); err != nil {
			return
		}
	}()
	for  {
		select {
		case msg := <-p.p2pPeer.GetProtocolMsgCh():
			msgCode := msg.Header.MsgCode
			switch msgCode {
			case MsgCodeVersion:
				bodyBs := msg.Body
				var status *statusData = nil
				if err := json.Unmarshal(bodyBs,&status); err != nil {
					return errors.New("error")
				}
				if status.Version != p.version {
					return errors.New("error")
				}
				if status.Network != p.network {
					return errors.New("error")
				}
				p.head = *status.Head
				return nil
			}
		case <-time.After(3 * 60 * time.Second):
			return errors.New("time out")
		}
	}
}