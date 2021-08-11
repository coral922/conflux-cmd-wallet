package cmd

import (
	"cfxWorld/config"
	"cfxWorld/interface/cmdline"
	"cfxWorld/lib/util"
	"github.com/spf13/cobra"
)

var wallet = &cobra.Command{
	Use:   "wallet",
	Short: "wallet operation",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config.MustInitConfig()
		cmdline.MustInitWallet()
		cmdline.MustCheckStorage()
	},
}

var importKey = &cobra.Command{
	Use:   "import [PRIVATE_KEY] [ACCOUNT_NAME]",
	Short: "import the private key of your account",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		cmdline.MustCheckPassword(password)
		key := args[0]
		name := args[1]
		cmdline.ImportAccountByPrivateKey(name, key, password)
	},
}

var accountList = &cobra.Command{
	Use:   "list",
	Short: "list all accounts in your wallet",
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			cmdline.DetailedAccountList()
		} else {
			cmdline.AccountList()
		}
	},
}

var getAccount = &cobra.Command{
	Use:   "show [NAME]",
	Short: "show account info with giving name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmdline.DetailedAccount(args[0])
	},
}

var (
	accountAmount int
)
var createAccounts = &cobra.Command{
	Use:   "create [NAME]",
	Short: "create accounts",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmdline.MustCheckPassword(password)
		name := args[0]
		cmdline.CreateAccounts(name, password, accountAmount)
	},
}

var deleteAccount = &cobra.Command{
	Use:   "delete [ACCOUNT_NAME]",
	Short: "delete account",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmdline.MustCheckPassword(password)
		name := args[0]
		cmdline.DeleteAccount(name, password)
	},
}

var exportAccount = &cobra.Command{
	Use:   "export [ACCOUNT_NAME]",
	Short: "export the private key of account",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmdline.MustCheckPassword(password)
		name := args[0]
		cmdline.ExportPrivateKey(name, password)
	},
}

var exportCsv = &cobra.Command{
	Use:   "exportcsv [CSV_FILE]",
	Short: "export all private keys to csv file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmdline.MustCheckPassword(password)
		file := args[0]
		cmdline.ExportAllPrivateKeyToCSV(password, file)
	},
}

var setDefault = &cobra.Command{
	Use:   "default [ACCOUNT_NAME]",
	Short: "get or set your default account",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmdline.MustCheckPassword(password)
		if len(args) == 1 {
			cmdline.SetCurrentAccount(args[0])
		} else {
			cmdline.GetCurrentAccount()
		}
	},
}

var changeName = &cobra.Command{
	Use:   "rename [OLD_NAME] [NEW_NAME]",
	Short: "change your account name",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		cmdline.MustCheckPassword(password)

		cmdline.RenameAccount(args[0], args[1])
	},
}

var changePw = &cobra.Command{
	Use:   "changepw [NEW_PASSWORD]",
	Short: "change your password",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmdline.MustCheckPassword(password)

		cmdline.ChangePassword(args[0])
	},
}

var resetWallet = &cobra.Command{
	Use:   "reset",
	Short: "reset the service (all accounts will be deleted)",
	Run: func(cmd *cobra.Command, args []string) {
		cmdline.MustCheckPassword(password)

		if util.SecondConfirm("this operation will delete all accounts in your wallet, are you sure?") {
			cmdline.ResetWallet()
		}
	},
}

func init() {
	createAccounts.Flags().IntVarP(&accountAmount, "number", "n", 1, "account amount (default 1)")
}
