package sub

import (
	"fmt"
	"log"
	"os"
	"xblockchain"
	"xblockchain/api"
	configs "xblockchain/cmd/config"
	"xblockchain/rpc"
	"xblockchain/storage/badger"

	"github.com/nictuku/dht"
	"github.com/spf13/cobra"
)

var (
	daemonCmd = &cobra.Command{
		Use: "daemon",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDaemon()
		},
	}
)

func runDaemon() error {

	tempConfig := configs.GetConfig() // 配置获取
	keyStorage := badger.New(tempConfig.Blockchain.Keys)
	defer func() {
		if err := keyStorage.Close(); err != nil {
			log.Fatalf("Key Storage close errors: %s", err)
		}
	}()
	blocksStorage := badger.New(tempConfig.Blockchain.Blocks)
	defer func() {
		if err := blocksStorage.Close(); err != nil {
			log.Fatalf("Blocks Storage close errors: %s", err)
		}
	}()

	txPendingPool := xblockchain.NewTXPendingPool(100)

	ws := xblockchain.NewWallets(keyStorage)
	gopt := xblockchain.DefaultGenesisBlockOpts()
	bc, err := xblockchain.NewBlockChain(gopt, blocksStorage)
	if err != nil {
		log.Fatalf("blockchain initail err: %s", err)
	}

	miner := xblockchain.NewMiner(bc, ws, txPendingPool)
	starter, err := rpc.NewServerStarter()
	if err != nil {
		log.Fatalf("RPC server starter initail err: %s", err)
		return err
	}
	chainApiHandler := &api.ChainAPIHandler{
		BlockChain: bc,
	}

	minerApiHandler := &api.MinerAPIHandler{
		Miner: miner,
	}

	walletApiHandler := &api.WalletsHandler{
		Wallets: ws,
	}

	txApiHandler := &api.TXAPIHandler{
		Wallets:       ws,
		BlockChain:    bc,
		TxPendingPool: txPendingPool,
	}
	err = starter.RegisterName("Chain", chainApiHandler)
	if err != nil {
		log.Fatalf("RPC service register error: %s", err)
		return err
	}

	err = starter.RegisterName("Wallet", walletApiHandler)
	if err != nil {
		log.Fatalf("RPC service register error: %s", err)
		return err
	}

	err = starter.RegisterName("Miner", minerApiHandler)
	if err != nil {
		log.Fatalf("RPC service register error: %s", err)
		return err
	}

	err = starter.RegisterName("Transaction", txApiHandler)
	if err != nil {
		log.Fatalf("RPC service register error: %s", err)
		return err
	}
	err = starter.Run()
	if err != nil {
		log.Fatalf("RPC server run error: %s", err)
		return err
	}

	//server := net.Default()
	//server.StartAndListen("0.0.0.0:198")
	log.Printf("RPC server running by listen :9005\n")
	RunP2PPeer()
	return nil
}

func RunP2PPeer() {
	d, err := dht.New(nil)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "New DHT error: %v", err)
		os.Exit(1)
	}
	if err = d.Start(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "DHT start error: %v", err)
		os.Exit(1)
	}
	go drainresults(d)
}

func drainresults(n *dht.DHT) {
	count := 0
	for r := range n.PeersRequestResults {
		for _, peers := range r {
			for _, x := range peers {
				//fmt.Printf("%d: %v\n", count, dht.DecodePeerAddress(x))
				log.Printf("join node: %s\n", dht.DecodePeerAddress(x))
				count++
				//if count >= 10 {
				//	log.Printf("RPC server running by listen :9005\n")
				//}
			}
		}
	}
}

func init() {
	rootCmd.AddCommand(daemonCmd)
}
