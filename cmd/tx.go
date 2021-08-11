package cmd

import (
	"cfxWorld/config"
	"cfxWorld/interface/cmdline"
	"cfxWorld/lib/moonswap"
	"errors"
	"github.com/spf13/cobra"
	"strings"
)

var tx = &cobra.Command{
	Use:   "tx",
	Short: "transaction operation",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config.MustInitConfig()
		cmdline.MustInitWallet()
		cmdline.MustCheckStorage()
		cmdline.MustCheckPassword(password)
		cmdline.MustInitTx()
	},
}

var (
	from     string
)

var sendToken = &cobra.Command{
	Use:   "send",
	Short: "transfer token (e.g. send 100 cMOON to name:myaddr1)",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 4 {
			return errors.New("4 args expected")
		}
		if args[2] != "to" {
			return errors.New("wrong syntax, except ([amount] [token] to [address])")
		}
		if !strings.HasPrefix(args[3], "cfx:") && !strings.HasPrefix(args[3], "name:") {
			return errors.New("wrong address syntax, except cfx:xxx OR name:xxx")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		amount, token, to := args[0], args[1], args[3]
		cmdline.Send(token, from, to, amount, password)
	},
}

var sendNft = &cobra.Command{
	Use:   "sendnft",
	Short: "transfer nft token (e.g. sendnft conDragon [your_token_id] to name:myaddr1)",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 4 {
			return errors.New("4 args expected")
		}
		if args[2] != "to" {
			return errors.New("wrong syntax, except ([nft_name] [nft_token_id] to [address])")
		}
		if !strings.HasPrefix(args[3], "cfx:") && !strings.HasPrefix(args[3], "name:") {
			return errors.New("wrong address syntax, except cfx:xxx OR name:xxx")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		tokenTyp, tokenID, to := args[0], args[1], args[3]
		cmdline.SendNft(tokenTyp, tokenID, from, to, password)
	},
}

var (
	slip int
)
var swap = &cobra.Command{
	Use:   "swap",
	Short: "swap tokens for another. (e.g. swap 100 cMOON for cUSDT) (e.g. swap cMOON for 100 cUSDT)",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 4 {
			return errors.New("4 args expected")
		}
		if args[1] != "for" && args[2] != "for" {
			return errors.New("wrong syntax, except ([amount] [token1] for [token2]) OR ([token1] for [amount] [token2])")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var mode int
		var amount string
		var input, output string
		if args[2] == "for" {
			mode = moonswap.SwapModeExactInput
			amount = args[0]
			input, output = args[1], args[3]
		}
		if args[1] == "for" {
			mode = moonswap.SwapModeExactOutPut
			amount = args[2]
			input, output = args[0], args[3]
		}
		cmdline.Swap(mode, from, input, output, amount, password, slip)
	},
}

func init() {
	tx.PersistentFlags().StringVar(&from, "from", "", "from which address (default your current account), "+
		"use [name:] as prefix if you want use account name instead of address (e.g. name:myaccount1)")
	swap.Flags().IntVarP(&slip, "slippage", "s", 1, "slippage tolerance (percent) (default 1)")
}
