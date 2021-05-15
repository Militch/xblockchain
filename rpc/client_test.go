package rpc

import (
	"testing"
)

type wallet struct {
	Address string `json:"address"`
}

func TestClient_CallMethod(t *testing.T) {
	cli := NewClient("http://localhost:9005")
	a  := make([]*wallet,0)
	err := cli.CallMethod(1,"Wallet.List",nil,&a)
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_CallMethod_New(t *testing.T) {
	cli := NewClient("http://localhost:9005")
	var a *string = nil
	err := cli.CallMethod(1,"Wallet.New",nil,&a)
	if err != nil {
		t.Fatal(err)
	}
}