package cmd

import (
	"cfxWorld/config"
	"cfxWorld/interface/cmdline"
	"cfxWorld/lib/util"
	"github.com/spf13/cobra"
)

var nft = &cobra.Command{
	Use:   "nft",
	Short: "nft token operation",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config.MustInitConfig()
		cmdline.MustInitNft()
	},
}

var nftList = &cobra.Command{
	Use:   "list",
	Short: "list all nft token supported",
	Run: func(cmd *cobra.Command, args []string) {
		cmdline.NftList()
	},
}

var registerNft = &cobra.Command{
	Use:   "register [ADDRESS]",
	Short: "register new CRC1155 token",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		address := args[0]
		cmdline.AddNft(address)
	},
}

var deleteNft = &cobra.Command{
	Use:   "delete [NFT_TOKEN_SYMBOL]",
	Short: "remove nft token",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		cmdline.DeleteNft(name)
	},
}

var followNft = &cobra.Command{
	Use:   "follow [TOKEN_NAME1] [TOKEN_NAME2] ...",
	Short: "follow nft token (you can check them in wallet list)",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmdline.FollowNft(args)
	},
}

var unfollowNft = &cobra.Command{
	Use:   "unfollow [TOKEN_NAME1] [TOKEN_NAME2] ...",
	Short: "unfollow nft token",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmdline.UnFollowNft(args)
	},
}

var resetNft = &cobra.Command{
	Use:   "reset",
	Short: "remove all nft token",
	Run: func(cmd *cobra.Command, args []string) {
		if util.SecondConfirm("this operation will delete all CRC1155 token, are you sure?") {
			cmdline.ResetNft()
		}
	},
}

func init() {
}
