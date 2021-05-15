package xblockchain

import (
	"fmt"
	"log"
	"time"
)

type Miner struct {
	BlockChain *BlockChain
	Wallets *Wallets
	running bool
	TxPendingPool *TXPendingPool
	selflock bool
	utxolock bool
}

func NewMiner(
	blockChain *BlockChain,
	wallets *Wallets,
	pool *TXPendingPool) *Miner {
	return &Miner{
		BlockChain: blockChain,
		Wallets: wallets,
		TxPendingPool: pool,
	}
}

func (m *Miner) Run() error {
	err := m.PreRun()
	if err != nil {
		return err
	}
	if m.running {
		return fmt.Errorf("miner is running, unable to start again")
	}
	m.running = true
	m.selflock = true
	go m.selfRun()
	go m.run()
	return nil
}

func (m *Miner) run() {
	for m.running {
		pendingTxs := m.TxPendingPool.PopAll()
		m.utxolock = true
		for m.selflock {
			time.Sleep(1000 * time.Millisecond)
		}
		m.utxolock = true
		txs := make([]*Transaction,0)
		addr := m.Wallets.GetDefault()
		if addr == "" {
			log.Printf("not found default address!!!!")
			continue
		}
		coinbaseTx,err := NewCoinBaseTransaction(addr,"")
		if err != nil {
			log.Printf("create txs err: %v\n", err)
			return
		}
		txs = append(txs, coinbaseTx)
		txs = append(txs, pendingTxs...)
		//block, err := m.BlockChain.AddBlock(txs)
		_, err = m.BlockChain.AddBlock(txs)
		if err != nil {
			log.Printf("Miner block errors: %v\n", err)
		}
		m.utxolock = false
		//log.Printf("Miner tx block success, height: %d, hash: %s", block.Height, block.Hash.Hexstr(true))
	}
}


func (m *Miner) selfRun() {
	for m.running {
		for m.utxolock {
			time.Sleep(1000 * time.Millisecond)
		}
		m.selflock = true
		txs := make([]*Transaction,0)
		addr := m.Wallets.GetDefault()
		if addr == "" {
			log.Printf("not found default address!!!!")
			continue
		}
		coinbaseTx,err := NewCoinBaseTransaction(addr,"")
		if err != nil {
			log.Printf("create txs err: %v\n", err)
			return
		}
		txs = append(txs, coinbaseTx)
		//block, err := m.BlockChain.AddBlock(txs)

		_, err = m.BlockChain.AddBlock(txs)
		if err != nil {
			log.Printf("Miner block errors: %v\n", err)
		}
		m.selflock = false

		time.Sleep(10000 * time.Millisecond)
		//log.Printf("Miner zore block success, height: %d, hash: %s", block.Height, block.Hash.Hexstr(true))
	}
}

func (m *Miner) PreRun() error {
	addr := m.Wallets.GetDefault()
	if addr == "" {
		return fmt.Errorf("not found default address")
	}
	return nil
}
func (m *Miner) Stop() error {
	if !m.running {
		return fmt.Errorf("miner is not running, unable to stop miner")
	}
	m.running = false
	return nil
}