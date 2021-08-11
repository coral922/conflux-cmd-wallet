package cmd

import (
	"cfxWorld/config"
	"cfxWorld/interface/cmdline"
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var (
	password string
	verbose  bool
)

var rootCmd = &cobra.Command{
	Use:   "cfxWorld",
	Short: "Welcome to conflux world!",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config.MustInitConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
		logo := figure.NewFigure("CONFLUX", "doom", true)
		logo.Print()

		fmt.Println("")

		err := cmd.Help()
		if err != nil {
			log.Fatal(err)
		}
	},
}

var update = &cobra.Command{
	Use:   "update",
	Short: "update crc20 and crc1155 token via data of <moonswap>",
	Run: func(cmd *cobra.Command, args []string) {
		cmdline.MustInitToken()
		cmdline.MustInitNft()
		cmdline.SyncTokenAndPair()
		//TODO: sync crc1155
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&config.ConfPath, "config", "c", "app", "config file path, default ./app.yaml")

	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "your password")

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose information")

	rootCmd.AddCommand(
		update,
		openCfxScan,
		info,
		wallet,
		token,
		nft,
		tx,
	)

	info.AddCommand(
		tokenInfo,
	)

	wallet.AddCommand(
		importKey,
		accountList,
		createAccounts,
		deleteAccount,
		resetWallet,
		changePw,
		exportAccount,
		exportCsv,
		changeName,
		setDefault,
		getAccount,
	)

	token.AddCommand(
		tokenList,
		registerToken,
		deleteToken,
		resetToken,
		followToken,
		unFollowToken,
	)

	nft.AddCommand(
		nftList,
		registerNft,
		deleteNft,
		resetNft,
		followNft,
		unfollowNft,
	)

	tx.AddCommand(
		sendToken,
		sendNft,
		swap,
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
