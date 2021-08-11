package cmdline

import (
	"cfxWorld/interface/ui"
	"fmt"
	"log"
)

func AccountList() {
	a, err := walletSvc.CurrentAccount()
	if err != nil {
		log.Fatal(err)
	}
	curName := ""
	if a != nil {
		curName = a.Name
	}
	as, err := walletSvc.AccountList()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ui.AccountListTable(curName, as))
}

func DetailedAccountList() {
	a, err := walletSvc.CurrentAccount()
	if err != nil {
		log.Fatal(err)
	}
	curName := ""
	if a != nil {
		curName = a.Name
	}
	as, err := walletSvc.DetailedAccountList()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ui.DetailedAccountListTable(curName, as))
}

func DetailedAccount(name string) {
	d, err := walletSvc.DetailedAccount(name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ui.DetailedAccountTable(d))
}

func ImportAccountByPrivateKey(name, key, password string) {
	if key == "" {
		log.Fatal("empty key")
	}
	account, err := walletSvc.ImportKey(name, key, password)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("your account has been successfully imported")
	fmt.Println(ui.AccountSimpleTable(account))
}

func CreateAccounts(namePrefix, password string, amount ...int) {
	if password == "" {
		log.Println("WARN: no password passed")
	}
	num := 1
	if len(amount) > 0 && amount[0] > 1 {
		num = amount[0]
	}
	accounts := walletSvc.CreateAccountBatch(namePrefix, password, num)
	log.Printf("%d accounts created \n", len(accounts))
	if len(accounts) > 0 {
		fmt.Println(ui.AccountsSimpleTable(accounts))
	}
}

func DeleteAccount(name string, password string) {
	err := walletSvc.DeleteAccountByName(name, password)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("account [%s] deleted", name)
}

func SetCurrentAccount(name string) {
	err := walletSvc.SetCurrent(name)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("default account has been set to [%s]", name)
}

func GetCurrentAccount() {
	a, err := walletSvc.CurrentAccount()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ui.AccountSimpleTable(a))
}

func RenameAccount(old, new string) {
	err := walletSvc.ChangeAccountName(old, new)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("account [%s] has been changed to [%s]", old, new)
}

func ExportPrivateKey(name string, password string) {
	str, err := walletSvc.ExportPrivateKey(name, password)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(str)
}

func ExportAllPrivateKeyToCSV(password string, csvFile string) {
	err := walletSvc.ExportAllPrivateKeyToCSV(password, csvFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("all private keys have been exported to [%s]", csvFile)
}

func ResetWallet() {
	err := walletSvc.ResetWallet()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("the wallet has been reset")
}

func ChangePassword(newPw string) {
	if newPw == "" {
		log.Fatal("empty or invalid password")
	}
	err := walletSvc.SetPassword(newPw)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("the password has been changed")
}
