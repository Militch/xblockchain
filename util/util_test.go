package util

import (
	"testing"
	"xblockchain/backend"
	"xblockchain/node"
)

func TestStartNodeAndBackend(t *testing.T) {
	var err error = nil
	var stack *node.Node = nil
	var back *backend.Backend = nil
	if stack, err = node.New(&node.Opts{
		P2PListenAddress: ":9001",
		RPCListenAddress: ":9002",
		P2PBootstraps: []string{},
	}); err != nil {
		t.Fatal(err)
	}
	if back, err = backend.NewBackend(stack, &backend.Opts{
		BlockDbPath: "./data0/blocks",
		KeyStoragePath: "./data0/keys",
		Version: uint32(0),
		Network: uint32(0),
	}); err != nil {
		t.Fatal(err)
	}
	if err = StartNodeAndBackend(stack,back); err != nil {
		t.Fatal(err)
	}
	select {}
}


func TestStartNodeAndBackend3(t *testing.T) {
	var err error = nil
	var stack *node.Node = nil
	var back *backend.Backend = nil
	if stack, err = node.New(&node.Opts{
		P2PListenAddress: ":9003",
		RPCListenAddress: ":9004",
		P2PBootstraps: []string{},
	}); err != nil {
		t.Fatal(err)
	}
	if back, err = backend.NewBackend(stack, &backend.Opts{
		BlockDbPath: "./data1/blocks",
		KeyStoragePath: "./data1/keys",
		Version: uint32(0),
		Network: uint32(0),
	}); err != nil {
		t.Fatal(err)
	}
	if err = StartNodeAndBackend(stack,back); err != nil {
		t.Fatal(err)
	}
	select {}
}