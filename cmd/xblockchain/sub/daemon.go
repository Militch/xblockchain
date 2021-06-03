package sub

import (
	"github.com/spf13/cobra"
	"xblockchain/backend"
	"xblockchain/node"
	"xblockchain/util"
)

var (
	daemonCmd = &cobra.Command{
		Use:   "daemon",
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
		P2PListenAddress: P2PListenAddress,
		P2PBootstraps: P2PBootstraps,
		RPCListenAddress: RPCListenAddress,
	}); err != nil {
		return err
	}
	if back, err = backend.NewBackend(stack, &backend.Opts{
		BlockDbPath: BlockDbPath,
		KeyStoragePath: KeyStoragePath,
	}); err != nil {
		return err
	}
	if err = util.StartNodeAndBackend(stack,back); err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(daemonCmd)
}