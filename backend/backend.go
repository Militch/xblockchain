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
	blockChain *xblockchain.BlockChain
	miner *xblockchain.Miner
}
type Opts struct {
	BlockDbPath string
	KeyStoragePath string
}
func NewBackend(stack *node.Node, opts *Opts) (*Backend,error) {
	var err error = nil
	backend := &Backend{
		txPool: make(chan string),
		p2pServer: stack.P2PServer(),
	}
	backend.blockDb = badger.New(opts.BlockDbPath)
	backend.keysDb = badger.New(opts.KeyStoragePath)
	backend.txPendingPool = xblockchain.NewTXPendingPool(100)
	genesisOpts := xblockchain.DefaultGenesisBlockOpts()
	if backend.blockChain ,err = xblockchain.NewBlockChain(
		genesisOpts, backend.blockDb); err != nil {
		return nil, err
	}
	backend.wallets = xblockchain.NewWallets(backend.keysDb)
	backend.miner = xblockchain.NewMiner(backend.blockchain, backend.wallets, backend.txPendingPool)

	if err = stack.RegisterBackend(
		backend.blockchain,backend.miner, backend.wallets,
		backend.txPendingPool); err != nil {
		return nil, err
	}

	if backend.handler, err = newHandler(); err != nil {
		return nil, err
	}
	callFn := backend.handler.handlerCallFn
	backend.p2pServer.PeerHandlerFn = callFn
	return backend, nil
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

