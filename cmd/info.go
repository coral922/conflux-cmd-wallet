package cmd

import (
	"cfxWorld/config"
	"cfxWorld/interface/cmdline"
	"fmt"
	"github.com/spf13/cobra"
)

var openCfxScan = &cobra.Command{
	Use:   "scan",
	Short: "open the conflux scan website",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
		cmdline.OpenConfluxScan(args...)
	},
}

var info = &cobra.Command{
	Use:   "info",
	Short: "information operation",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config.MustInitConfig()
	},
}

var tokenInfo = &cobra.Command{
	Use:   "token [ADDRESS]",
	Short: "show the token detail with giving address",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmdline.TokenInfo(args[0])
	},
}

func init() {
}
