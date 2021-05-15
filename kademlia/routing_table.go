package kademlia

import (
	"xblockchain/uint256"
)

const RoutingTableSize = 256

type RoutingTable struct {
	kBuckets [RoutingTableSize] *KBucket
	selfNode *Node
}

func NewRoutingTable() *RoutingTable {
	kbs := [RoutingTableSize] *KBucket {}
	for i:=0;i<RoutingTableSize;i++{
		kbs[i] = NewKBucket()
	}
	node := &Node{
		ID: uint256.NewUInt256("0x01"),
	}
	return &RoutingTable{
		kBuckets: kbs,
		selfNode: node,
	}
}

func (routingTable *RoutingTable) AddContact(node *Node)  {
	index := routingTable.GetBucketFor(node)
	if index == -1 {
		return
	}
	kb := routingTable.kBuckets[index]
	kb.Push(node)
}


func (routingTable *RoutingTable) GetBucketFor(node *Node) int  {
	if node == nil || routingTable.selfNode == nil{
		return -1
	}
	selfId := routingTable.selfNode.ID
	nodeId := node.ID
	distance := selfId.Xor(nodeId)
	distance.HexstrFull()
	return -1
}