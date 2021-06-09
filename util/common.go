package util

import (
	"github.com/sirupsen/logrus"
	"xblockchain/backend"
	"xblockchain/node"
)

func StartNodeAndBackend(node *node.Node, backend *backend.Backend) error {
	logrus.Info("正在启动节点...")
	if err := node.Start(); err != nil {
		return err
	}
	logrus.Info("节点启动完成，正在启动后台服务...")
	if err := backend.Start(); err != nil {
		return err
	}
	logrus.Info("后台服务启动完成!!!")
	return nil
}