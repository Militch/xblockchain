package p2p

import (
	"github.com/perlin-network/noise"
	"github.com/sirupsen/logrus"
	"time"
)

type Peer struct {
	wc chan targetWrite
	ID noise.ID
	messageCh chan Message
	protocolMsgCh chan Message
}
type ProtocolHandler interface {
	RunProtocol(peer *Peer) error
}

func newPeer(wc chan targetWrite, id noise.ID) *Peer {
	p := &Peer{
		wc: wc,
		ID: id,
		messageCh: make(chan Message),
		protocolMsgCh: make(chan Message),
	}
	return p
}

func (p *Peer) run(callFn func(peer *Peer) error) error {
	go p.readLoop()
	go p.pingLoop()
	go func() {
		if err := callFn(p); err != nil{
			logrus.Error(err)
		}
	}()
	return nil
}

func (p *Peer) readLoop() {
	for {
		msg := <- p.messageCh
		if err := p.handle(msg); err != nil {
			logrus.Error(err)
		}
	}
}

func (p *Peer) handle(msg Message)  error {
	switch msg.Header.MsgCode {
	case MsgCodePing:
		logrus.Infof("pong")
		if err := SendMsg(p, MsgCodePong,nil); err != nil {
			return err
		}
	default:
		p.protocolMsgCh <- msg
	}
	return nil
}

func (p *Peer) pingLoop() {
	for range time.Tick(10000 * time.Millisecond) {
		logrus.Infof("peer ping loop target: %s", p.ID.Address)
		if err := SendMsg(p, MsgCodePing,nil); err != nil {
			logrus.Error(err)
		}
	}
}

func (p *Peer) sendData(data []byte) {
	tw := &targetWrite{
		target: p.ID,
		data: data,
	}
	p.wc <- *tw
}


func (p *Peer) handleData(data []byte) error {
	msg := &Message{}
	if err := UnmarshalMsg(data, msg); err != nil {
		return err
	}
	p.messageCh <- *msg
	return nil
}

func (p *Peer) GetProtocolMsgCh() chan Message {
	return p.protocolMsgCh
}