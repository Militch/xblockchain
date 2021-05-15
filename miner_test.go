package xblockchain

import (
	"fmt"
	"testing"
	"xblockchain/storage/badger"
)

func TestMiner_Run(t *testing.T) {
	keyStorage := badger.New("./data0/keys")
	defer func() {
		if err := keyStorage.Close(); err != nil {
			t.Fatalf("Sotrage close errors: %s", err)
		}
	}()
	ws := NewWallets(keyStorage)
	blocksStorage := badger.New("./data0/blocks")
	defer func() {
		if err := blocksStorage.Close(); err != nil {
			t.Fatalf("Sotrage close errors: %s", err)
		}
	}()
	gopt := DefaultGenesisBlockOpts()
	bc,err := NewBlockChain(gopt, blocksStorage)
	if err != nil {
		t.Fatal(err)
	}
	miner := NewMiner(bc,ws,nil)
	_ = miner
	//miner.Run()
}

func TestRunning(t *testing.T) {
	//var mu0 sync.Mutex
	//var mu1 sync.Mutex
	//pool := NewTXPendingPool(100)
	//mpool := make([]string,0)
	//mu.Lock()
	//mu0.Lock()
	//go func() {
	//	for i:=0;;i++{
	//		pool.Push([]byte(fmt.Sprintf("c%d",i)))
	//		time.Sleep(5 * time.Second)
	//	}
	//}()
	//go func() {
	//	for {
	//		msg := pool.Pop()
	//		fmt.Printf("i: %s\n", string(msg))
	//		time.Sleep(5 * time.Second)
	//	}
	//}()
	//wait := make(chan struct{})
	//<- wait
	//mu0.Lock()
	//var num *int
	//for i:=0;i<10;i++{
	//	in := i
	//	mu.Lock()
	//
	//}
	//mu.Lock()
	// 互斥锁
	//go func() {
	//	run1("c1")
	//	mu.Unlock()
	//}()
	//mu.Lock()
	//go func() {
	//	run1("c2")
	//	mu.Unlock()
	//}()
	//mu.Lock()

}

func TestName(t *testing.T) {
	//s := []byte("abcdef")
	//fmt.Printf("len(s): %d\n",len(s))
	//s = s[0:1]
	//fmt.Printf("len(s): %d\n",len(s))
	fmt.Printf("len(s): %v\n", 10<10)
}

func run0(c string, i int){
	fmt.Printf("run[%s]: %d\n",c,i)
}

func run1(c string){
	for i:=0;i<10;i++{
		fmt.Printf("run[%s]: %d\n",c,i)
	}
}