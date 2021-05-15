package kademlia

const KBucketSize = 20

type KBucket struct {
	nodes [KBucketSize]*Node
	top int
}

func NewKBucket() *KBucket {
	return &KBucket{
		top: -1,
	}
}

func (bucket *KBucket) Push(node *Node) {
	currentIndex := bucket.IndexOf(node)
	if currentIndex == -1 && bucket.top < KBucketSize {
		bucket.top++
		bucket.nodes[bucket.top] = node
	}else if currentIndex == -1 && bucket.top >= KBucketSize {

	}else {
		bucket.Move2Top(node)
	}
}
func (bucket *KBucket) Pop() *Node {
	tmp := bucket.nodes[bucket.top]
	bucket.nodes[bucket.top] = nil
	bucket.top--
	return tmp
}
func (bucket *KBucket) IndexOf(node *Node) int {
	for i := 0;i<KBucketSize;i++ {
		self := bucket.nodes[i]
		if self == nil || node == nil {
			return -1
		}
		selfId := self.ID
		targetId := node.ID
		if targetId.Equals(selfId) {
			return i
		}
	}
	return -1
}

func (bucket *KBucket) Move2Top(node *Node)  {
	index := bucket.IndexOf(node)
	tmp := bucket.nodes[index]
	for i:=index;i<KBucketSize-1;i++ {
		bucket.nodes[i] = bucket.nodes[i+1]
	}
	bucket.nodes[bucket.top] = tmp
}
func (bucket *KBucket) IsEmpty() bool {
	return bucket.top == -1
}