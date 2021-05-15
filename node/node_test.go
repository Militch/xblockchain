package node

import "testing"

//func TestListenAndServe(t *testing.T) {
//	err := ListenAndServe("0.0.0.0:9001")
//	if err != nil {
//		t.Fatal(err)
//	}
//}

func TestNode_ListenAndServe(t *testing.T) {
	n := NewNode("0.0.0.0:9001","")
	err := n.ListenAndServe()
	if err != nil {
		t.Fatal(err)
	}
}