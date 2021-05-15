package kademlia

import (
	"xblockchain/uint256"
)
type Node struct {
	ID *uint256.UInt256
	IPAddr string
	Port uint
}

func (node *Node) Equals(target *Node) bool {
	if node == nil || target == nil {
		return false
	}
	return node.ID.Equals(target.ID)
}