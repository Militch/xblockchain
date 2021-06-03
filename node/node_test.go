package node

import (
	"os"
	"os/signal"
	"testing"
)

func TestNode_Start(t *testing.T) {
	n,err := New()
	if err != nil {
		t.Fatal(err)
	}
	err = n.Start()
	if err != nil {
		t.Fatal(err)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}