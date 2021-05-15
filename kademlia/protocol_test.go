package kademlia

import (
	"testing"
)

func TestProtocol_New(t *testing.T) {
	data := [] byte{
		// source
		1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
		// dest
		2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3,
		1, TypeFlagReqHello, 0,
	}
	p := NewProtocol(data)
	if p.HasTypeFlag(TypeFlagReqHello) && !p.HasTypeFlag(TypeFlagAckHello) {
		p.Header.TypeFlag = TypeFlagReqHello | TypeFlagAckHello
	}
	t.Fatal(p.ToBytes())
}
