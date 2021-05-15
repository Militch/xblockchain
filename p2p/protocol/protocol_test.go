package protocol

import (
	"bytes"
	"fmt"
	"testing"
)

func TestProtocol_Marshal(t *testing.T) {
	data := []byte("hello")
	p := NewMessage(VersionT, 0, data)
	ps := p.Marshal()
	if len(ps) != int(p.Header.Length) {
		t.Fatal(fmt.Errorf("marshal protocol err"))
	}
}

func TestUnmarshalProtocol(t *testing.T) {
	data := []byte("hello")
	protocol := NewMessage(VersionAckT, 0, data)
	protocolBytes := protocol.Marshal()
	got,err := UnmarshalMessage(protocolBytes)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(protocol.Body, got.Body) != 0 {
		t.Fatal(fmt.Errorf("unmarshal protocol err"))
	}
}