package backend

import (
	"encoding/json"
	"errors"
	"time"
	"xblockchain"
	"xblockchain/p2p"
	"xblockchain/uint256"
)

type peer struct {
	p2pPeer *p2p.Peer
	version uint32
	network uint32
	head uint256.UInt256
	height uint64
}

var MsgCodeVersion = uint32(3)
var GetBlockHashesFromNumberMsg = uint32(4)
var BlockHashesMsg = uint32(5)
var GetBlocksMsg = uint32(6)
var BlocksMsg = uint32(7)
var NewBlockMsg = uint32(8)

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
	Version uint32 `json:"version"`
	Network uint32 `json:"network"`
	Head *uint256.UInt256 `json:"head"`
	Height uint64 `json:"height"`
}

type getBlockHashesFromNumberData struct {
	From uint64 `json:"from"`
	Count uint64 `json:"count"`
}

func (p *peer) Handshake(head *uint256.UInt256, height uint64) error {
	go func() {
		if err := p2p.SendMsgJSONData(p.p2pPeer, MsgCodeVersion, &statusData{
			Version: p.version,
			Network: p.network,
			Head: head,
			Height: height,
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
				p.height = status.Height
				return nil
			}
		case <-time.After(3 * 60 * time.Second):
			return errors.New("time out")
		}
	}
}

// RequestHashesFromNumber 请求获取header
func (p *peer) RequestHashesFromNumber(from uint64, count uint64) error {
	if err := p2p.SendMsgJSONData(p.p2pPeer, GetBlockHashesFromNumberMsg, &getBlockHashesFromNumberData{
		From: from,
		Count: count,
	}); err != nil {
		return err
	}
	return nil
}


// SendBlockHashes 发送header数据
func (p *peer) SendBlockHashes(hashes []uint256.UInt256) error {
	if err := p2p.SendMsgJSONData(p.p2pPeer, BlockHashesMsg, &hashes); err != nil {
		return err
	}
	return nil
}

// RequestBlocks 请求获取区块列表
func (p *peer) RequestBlocks(hashes []uint256.UInt256) error {
	if err := p2p.SendMsgJSONData(p.p2pPeer, GetBlocksMsg, &hashes); err != nil {
		return err
	}
	return nil
}
// SendBlocks 发送区块列表
func (p *peer) SendBlocks(blocks []xblockchain.Block) error {
	if err := p2p.SendMsgJSONData(p.p2pPeer,BlocksMsg ,&blocks); err != nil {
		return err
	}
	return nil
}