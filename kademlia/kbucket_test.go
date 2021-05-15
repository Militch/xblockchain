package kademlia

import (
	"testing"
	"xblockchain/uint256"
)

func TestKBucket_Push(t *testing.T) {
	kb := NewKBucket()
	kb.Push(&Node{
		ID: uint256.NewUInt256("0x01"),
	})
	if kb.IsEmpty() {
		t.Fatalf("Bucket nodes is empty")
	}
}

