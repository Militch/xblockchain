package node

import (
	"log"
	"xblockchain"
	"xblockchain/api"
	"xblockchain/p2p"
	"xblockchain/rpc"

	"github.com/sirupsen/logrus"
)

type Node struct {
	*Opts
	p2pServer  *p2p.Server
	RPCStarter *rpc.ServerStarter
}

type Opts struct {
	P2PListenAddress string
	P2PBootstraps    []string
	RPCListenAddress string
}

func New(opts *Opts) (*Node, error) {
	n := &Node{
		Opts: opts,
		p2pServer: &p2p.Server{
			ListenAddr:     opts.P2PListenAddress,
			BootstrapNodes: opts.P2PBootstraps,
		},
	}
	var err error = nil
	if n.RPCStarter, err = rpc.NewServerStarter(); err != nil {
		return nil, err
	}
	return n, nil
}

func (n *Node) Start() error {
	if err := n.p2pServer.Start(); err != nil {
		return err
	}
	go func() {
		if err := n.RPCStarter.Run(); err != nil {
			logrus.Warnf("启动 RPC ERR: %s", err)
		}
	}()
	return nil
}

func (n *Node) RegisterBackend(
	bc *xblockchain.BlockChain,
	miner *xblockchain.Miner,
	wallets *xblockchain.Wallets,
	txPendingPool *xblockchain.TXPendingPool) error {
	chainApiHandler := &api.ChainAPIHandler{
		BlockChain: bc,
	}

	minerApiHandler := &api.MinerAPIHandler{
		Miner: miner,
	}

	walletApiHandler := &api.WalletsHandler{
		Wallets: wallets,
	}

	txApiHandler := &api.TXAPIHandler{
		Wallets:       wallets,
		BlockChain:    bc,
		TxPendingPool: txPendingPool,
	}
	starter := n.RPCStarter
	if err := starter.RegisterName("Chain", chainApiHandler); err != nil {
		log.Fatalf("RPC service register error: %s", err)
		return err
	}
	if err := starter.RegisterName("Wallet", walletApiHandler); err != nil {
		log.Fatalf("RPC service register error: %s", err)
		return err
	}
	if err := starter.RegisterName("Miner", minerApiHandler); err != nil {
		log.Fatalf("RPC service register error: %s", err)
		return err
	}
	if err := starter.RegisterName("Transaction", txApiHandler); err != nil {
		log.Fatalf("RPC service register error: %s", err)
		return err
	}
	return nil
}

func (n *Node) P2PServer() *p2p.Server {
	return n.p2pServer
}
