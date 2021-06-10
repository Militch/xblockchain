package sub

import (
	"xblockchain/backend"
	"xblockchain/node"
	"xblockchain/util"

	"github.com/spf13/cobra"
)

var (
	daemonCmd = &cobra.Command{
		Use: "daemon",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDaemon()
		},
	}
)

func runDaemon() error {
	var err error = nil
	var stack *node.Node = nil
	var back *backend.Backend = nil
	if stack, err = node.New(&node.Opts{
		P2PListenAddress: GetConfigSub().Network.P2PListenAddress,
		P2PBootstraps:    GetConfigSub().Network.P2PBootstraps,
		RPCListenAddress: GetConfigSub().Network.RPCListenAddress,
	}); err != nil {
		return err
	}
	if back, err = backend.NewBackend(stack, &backend.Opts{
		BlockDbPath:    GetConfigSub().Blockchain.BlockDbPath,
		KeyStoragePath: GetConfigSub().Blockchain.KeyStoragePath,
		Version:        GetConfigSub().Network.ProtocolVersion,
		Network:        GetConfigSub().Network.NetworkID,
	}); err != nil {
		return err
	}
	if err = util.StartNodeAndBackend(stack, back); err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(daemonCmd)
}
