package backend

import (
	"encoding/json"
	"errors"
	"github.com/perlin-network/noise"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
	"xblockchain"
	"xblockchain/p2p"
	"xblockchain/uint256"
)

var MaxHashFetch = uint64(512)

type hashPack struct {
	p *peer
	hashes []uint256.UInt256
}
type blockPack struct {
	p *peer
	blocks []xblockchain.Block
}

type handler struct {
	handlerCallFn func(peer *p2p.Peer) error
	newPeerCh chan *peer
	hashPackCh chan hashPack
	blockPackCh chan blockPack
	peers map[noise.PublicKey] *peer
	blockchain *xblockchain.BlockChain
	syncLock sync.Mutex
	fetchHashesLock sync.Mutex
	fetchBlocksLock sync.Mutex
	processLock sync.Mutex
	version uint32
	network uint32
}


func newHandler(bc *xblockchain.BlockChain, pv uint32, nv uint32) (*handler,error) {
	h := &handler{
		newPeerCh: make(chan *peer, 1),
		hashPackCh: make(chan hashPack),
		blockPackCh: make(chan blockPack),
		peers: make(map[noise.PublicKey] *peer),
		blockchain: bc,
		version: pv,
		network: nv,
	}
	h.handlerCallFn = h.handleNewPeer
	return h, nil
}

func (h *handler) handleNewPeer(p2p *p2p.Peer) error {
	p := newPeer(p2p, h.version, h.network)
	h.newPeerCh <- p
	return h.handle(p)
}

func (h *handler) handle(p *peer) error {
	var err error = nil
	var head *xblockchain.Block = nil
	if head, err = h.blockchain.GetHeadBlock(); err != nil {
		return err
	}
	if err = p.Handshake(head.Hash,head.Height); err != nil {
		return err
	}
	logrus.Infof("handshake success, peer.height: %d, p.head: %s", p.height, p.head.Hex())
	id := p.p2p().ID
	h.peers[id.ID] = p
	defer delete(h.peers, id.ID)
	for {
		if err = h.handleMsg(p); err!= nil {
			return err
		}
	}
}

func  (h *handler) handleMsg(p *peer) error {
	msg := <-p.p2pPeer.GetProtocolMsgCh()
	msgCode := msg.Header.MsgCode
	switch msgCode {
	case NewBlockMsg:
		// 处理区块广播
	case GetBlockHashesFromNumberMsg:
		// 获取本地区块 Hash 列表
		bodyBs := msg.Body
		var data *getBlockHashesFromNumberData = nil
		if err := json.Unmarshal(bodyBs,&data); err != nil {
			logrus.Warnf("handle GetBlockHashesFromNumberMsg msg err: %s", err)
			return err
		}
		hashes := h.blockchain.GetBlockHashes(int(data.From), int(data.Count))
		// 发送本地hash值
		if err := p.SendBlockHashes(hashes); err != nil {
			logrus.Warnf("send block hashes data err: %s", err)
			return err
		}
	case BlockHashesMsg:
		// 接受区块 hash 列表消息
		bodyBs := msg.Body
		var data []uint256.UInt256 = nil
		if err := json.Unmarshal(bodyBs,&data); err != nil {
			logrus.Warnf("handle BlockHashesMsg msg err: %s", err)
			return err
		}
		h.hashPackCh <- hashPack{
			p: p,
			hashes: data,
		}
	case GetBlocksMsg:
		// 处理获取区块列表请求
		bodyBs := msg.Body
		var data []uint256.UInt256 = nil
		if err := json.Unmarshal(bodyBs,&data); err != nil {
			logrus.Warnf("handle GetBlocksMsg msg err: %s", err)
			return err
		}
		blocks := make([]xblockchain.Block, 0)
		for _, hash := range data {
			if block,err := h.blockchain.GetBlockByHash(&hash); err == nil && block != nil {
				blocks = append(blocks, *block)
			}
		}
		if err := p.SendBlocks(blocks); err != nil {
			logrus.Warnf("send blocks data err: %s", err)
			return err
		}
	case BlocksMsg:
		// 接受区块列表消息
		bodyBs := msg.Body
		var data []xblockchain.Block = nil
		if err := json.Unmarshal(bodyBs,&data); err != nil {
			logrus.Warnf("handle BlocksMsg msg err: %s", err)
			return err
		}
		h.blockPackCh <- blockPack{
			p: p,
			blocks: data,
		}
	}
	return nil
}

func (h *handler) syncer() {
	forceSync := time.Tick(10 * time.Second)
	for {
		select {
		case <-h.newPeerCh:
			if len(h.peers) < 5 {
				break
			}
			go h.synchronise(h.basePeer())
		case <-forceSync:
			// 强制同步
			go h.synchronise(h.basePeer())
		}
	}
}
func (h *handler) basePeer() *peer {
	var (
		bestPeer *peer = nil
		baseHeight uint64 = 0
	)
	for _, v := range h.peers {
		if ph := v.height; ph > baseHeight {
			bestPeer = v
			baseHeight = ph
		}
	}
	return bestPeer
}

func (h *handler) synchronise(p *peer) {
	h.syncLock.Lock()
	defer h.syncLock.Unlock()
	logrus.Warnf("正在进行同步....")
	if p == nil {
		logrus.Warnf("未找到合适的同步链路")
		return
	}
	var number uint64
	var err error
	if number, err = h.findAncestor(p); err != nil {
		return
	}
	logrus.Infof("获取到公共区块高度: %d", number)
	go func() {
		if err = h.fetchHashes(p, number + 1); err != nil {
			logrus.Warn("fetch hashes err")
		}
	}()
	go func() {
		if err = h.fetchBlocks(p); err != nil {
			logrus.Warn("fetch blocks err")
		}
	}()
}
// 寻找公共区块高度
func (h *handler) findAncestor(p *peer) (uint64,error)  {
	var err error = nil
	var headBlock *xblockchain.Block = nil
	if headBlock, err = h.blockchain.GetHeadBlock(); err != nil {
		return 0,err
	}
	head := int(headBlock.Height)
	from := head - int(MaxHashFetch)
	if from < 0 {
		from = 0
	}
	logrus.Infof("寻找固定高度区间: [%d, %d]", from, MaxHashFetch)
	// 获取区块hash列表
	if err = p.RequestHashesFromNumber(uint64(from), MaxHashFetch); err != nil {
		return 0,err
	}
	number := uint64(0)
	haveHash := *uint256.NewUInt256Zero()
	// 阻塞接收pack消息
	loop:
	for {
		select {
		case pack := <-h.hashPackCh:
			if pack.p != p {
				break
			}
			hashes := pack.hashes
			if len(hashes) == 0 {
				return 0, errors.New("empty hashes")
			}
			for i,hash := range hashes {
				if h.hashBlock(hash) {
					continue
				}
				// 记录高度与hash值
				number = uint64(from) + uint64(i)
				haveHash = hash
				break loop
			}
		}
	}

	if !haveHash.IsZero() {
		return number, nil
	}
	logrus.Infof("未找到固定区间值，继续遍历查找...")
	// 如果未找到固定区间值，遍历所有区块，二分查找
	left := 0
	right := int(MaxHashFetch) + 1
	for left < right {
		logrus.Infof("正在遍历查找高度区间: [%d, %d]", left, right)
		mid := (left + right) / 2
		if err = p.RequestHashesFromNumber(uint64(mid), 1); err != nil {
			return 0, err
		}
		for {
			select {
			case pack := <-h.hashPackCh:
				if pack.p != p {
					break
				}
				hashes := pack.hashes
				if len(hashes) != 1 {
					return 0, nil
				}
				if h.hashBlock(hashes[0]) {
					left = mid + 1
				} else {
					right = mid
				}
			}
		}
	}
	return uint64(left) - 1, nil
}
// 寻找hash值是否在本地存在本地区块列表中
func (h *handler) hashBlock(hash uint256.UInt256) bool {
	var err error = nil
	var block *xblockchain.Block = nil
	if block, err = h.blockchain.GetBlockByHash(&hash); err != nil {
		return false
	}
	if block == nil {
		 return false
	}
	return true
}

func (h *handler) fetchHashes(p *peer, from uint64) error {
	h.fetchHashesLock.Lock()
	defer h.fetchHashesLock.Unlock()
	go func() {
		if err := p.RequestHashesFromNumber(from, MaxHashFetch); err != nil {
			logrus.Warn("request hashes err")
		}
	}()
	for {
		select {
		case pack := <-h.hashPackCh:
			if pack.p != p {
				break
			}
			hashes := pack.hashes
			if len(hashes) == 0 {
				return nil
			}
			for _, hash := range hashes {
				logrus.Infof("handle fetch hash: %s", hash.Hex())
			}
			go func() {
				if err := p.RequestBlocks(hashes); err != nil {
					logrus.Warn("request blocks err")
				}
			}()
		}
	}
}

func (h *handler) fetchBlocks(p *peer) error {
	h.fetchBlocksLock.Lock()
	defer h.fetchBlocksLock.Unlock()
	for {
		select {
		case pack := <-h.blockPackCh:
			if pack.p != p {
				break
			}
			blocks := pack.blocks
			if len(blocks) == 0 {
				return nil
			}
			go h.process(blocks)
		}
	}
}

func (h *handler) process(blocks []xblockchain.Block) {
	h.processLock.Lock()
	defer h.processLock.Unlock()
	coverRawBlocks := func(blocks []xblockchain.Block) []*xblockchain.Block {
		tmp := make([]*xblockchain.Block, 0)
		for _, block := range blocks {
			logrus.Infof("cover block: %s", block.Hash.Hex())
			b := &xblockchain.Block{
				Height: block.Height,
				Timestamp: block.Timestamp,
				HashPrevBlock: block.HashPrevBlock,
				Nonce: block.Nonce,
				Hash: block.Hash,
				Transactions: block.Transactions,
			}
			tmp = append(tmp, b)
		}
		return tmp
	}
	covered := coverRawBlocks(blocks)
	if index, err := h.blockchain.InsertBatchBlock(covered); err != nil {
		logrus.Warnf("process blocks[%d], hash: %s err: %s", index,blocks[index].Hash.Hex(), err)
	}
}



func (h *handler) Start() {
	go h.syncer()
}

