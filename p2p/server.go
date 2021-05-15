package p2p

import (
	"github.com/perlin-network/noise"
)

type Server struct {
	Addr string
	Node string
}



func (srv *Server) ListenAndServe() error {
	//node,err := noise.NewNode()
	//if err != nil {
	//	return err
	//}
	return nil
	//return s.Serve(node)
}

func (srv *Server) Serve(l *noise.Node) error {
	l.Handle(func(ctx noise.HandlerContext) error {
		return nil
	})
	return l.Listen()
}


