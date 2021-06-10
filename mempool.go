package xblockchain

import (
	"fmt"
	"sync"
	"xblockchain/uint256"
)

type TxPool struct {
	blockchain *BlockChain
	pending map[uint256.UInt256] *Transaction
	mu sync.RWMutex
}

func NewTxPool(chain *BlockChain) *TxPool {
	txPool := &TxPool{
		blockchain: chain,
		pending: make(map[uint256.UInt256] *Transaction),
	}
	return txPool
}

func (pool *TxPool) AddTx(tx *Transaction) error {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	hash := *uint256.NewUInt256BS(tx.Hash())
	if pool.pending[hash] != nil {
		return fmt.Errorf("know transaction (%s)", hash.Hex())
	}
	if pool.blockchain.VerifyTransaction(tx) {
		return fmt.Errorf("verify transaction err, hash: %s\n", hash.Hex())
	}

	return nil
}


