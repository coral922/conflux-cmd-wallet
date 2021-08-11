package cmd

import (
	"cfxWorld/config"
	"cfxWorld/interface/cmdline"
	"cfxWorld/lib/util"
	"github.com/spf13/cobra"
)

var token = &cobra.Command{
	Use:   "token",
	Short: "token operation",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config.MustInitConfig()
		cmdline.MustInitToken()
	},
}

var tokenList = &cobra.Command{
	Use:   "list",
	Short: "list all token supported",
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			cmdline.DetailedTokenList()
		} else {
			cmdline.TokenList()
		}
	},
}

var registerToken = &cobra.Command{
	Use:   "register [ADDRESS]",
	Short: "register new CRC20 token",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		address := args[0]
		cmdline.AddToken(address)
	},
}

var deleteToken = &cobra.Command{
	Use:   "delete [TOKEN_NAME]",
	Short: "remove token",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		cmdline.DeleteToken(name)
	},
}

var followToken = &cobra.Command{
	Use:   "follow [TOKEN_NAME1] [TOKEN_NAME2] ...",
	Short: "follow tokens (you can check them in wallet list)",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmdline.FollowTokens(args)
	},
}

var unFollowToken = &cobra.Command{
	Use:   "unfollow [TOKEN_NAME1] [TOKEN_NAME2] ...",
	Short: "unfollow tokens",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmdline.UnFollowTokens(args)
	},
}

var resetToken = &cobra.Command{
	Use:   "reset",
	Short: "remove all token",
	Run: func(cmd *cobra.Command, args []string) {
		if util.SecondConfirm("this operation will delete all CRC20 token, are you sure?") {
			cmdline.ResetToken()
		}
	},
}

func init() {
}
