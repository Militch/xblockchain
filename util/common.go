package util

import (
	"xblockchain/backend"
	"xblockchain/node"
)

func StartNodeAndBackend(node *node.Node, backend *backend.Backend) error {
	if err := node.Start(); err != nil {
		return err
	}
	if err := backend.Start(); err != nil {
		return err
	}
	return nil
}