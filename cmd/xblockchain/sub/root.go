package sub

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	cfgFile     string
	rootCmd = &cobra.Command{
		Use:   "fixcoin",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}
			return nil
		},
	}
	getCommand = &cobra.Command{
		Use:   "get <command> [flags]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	rootCmd.AddCommand(getCommand)
}


