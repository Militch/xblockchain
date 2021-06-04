package backend

import (
	"log"
	"xblockchain"
	"xblockchain/node"
	"xblockchain/p2p"
	"xblockchain/storage/badger"
)

type Backend struct {
	txPool chan string
	blockchain *xblockchain.BlockChain
	handler *handler
	blockDb *badger.Storage
	keysDb *badger.Storage
	p2pServer *p2p.Server
	txPendingPool *xblockchain.TXPendingPool
	wallets *xblockchain.Wallets
	miner *xblockchain.Miner
}
type Opts struct {
	BlockDbPath string
	KeyStoragePath string
	Version uint32
	Network uint32
}
func NewBackend(stack *node.Node, opts *Opts) (*Backend,error) {
	var err error = nil
	back := &Backend{
		txPool: make(chan string),
		p2pServer: stack.P2PServer(),
	}
	back.blockDb = badger.New(opts.BlockDbPath)
	back.keysDb = badger.New(opts.KeyStoragePath)
	back.txPendingPool = xblockchain.NewTXPendingPool(100)
	genesisOpts := xblockchain.DefaultGenesisBlockOpts()
	if back.blockchain ,err = xblockchain.NewBlockChain(
		genesisOpts, back.blockDb); err != nil {
		return nil, err
	}
	back.wallets = xblockchain.NewWallets(back.keysDb)
	back.miner = xblockchain.NewMiner(back.blockchain, back.wallets, back.txPendingPool)

	if err = stack.RegisterBackend(
		back.blockchain,back.miner, back.wallets,
		back.txPendingPool); err != nil {
		return nil, err
	}

	if back.handler, err = newHandler(back.blockchain,
		opts.Version, opts.Network); err != nil {
		return nil, err
	}
	callFn := back.handler.handlerCallFn
	back.p2pServer.PeerHandlerFn = callFn
	return back, nil
}

func (b *Backend) Start() error {
	b.handler.Start()
	return nil
}

func (b *Backend) close() {
	if err := b.blockDb.Close(); err != nil {
		log.Fatalf("Blocks Storage close errors: %s", err)
	}
	if err := b.keysDb.Close(); err != nil {
		log.Fatalf("Blocks Storage close errors: %s", err)
	}
}

