package cmdline

import (
	"cfxWorld/app/nft"
	"cfxWorld/app/token"
	"cfxWorld/app/tx"
	"cfxWorld/app/wallet"
	"cfxWorld/container"
	"fmt"
	"log"
)

var (
	walletSvc *wallet.Service
	tokenSvc  *token.Service
	txSvc     *tx.Service
	nftSvc    *nft.Service
)

func MustInitWallet() {
	err := container.Ctn.Invoke(func(s *wallet.Service) {
		walletSvc = s
	})
	if err != nil {
		log.Println(err)
		log.Fatal("MustInitWallet failed")
	}
}

func MustInitTx() {
	err := container.Ctn.Invoke(func(s *tx.Service) {
		txSvc = s
	})
	if err != nil {
		log.Println(err)
		log.Fatal("MustInitTx failed")
	}
}

func MustInitToken() {
	err := container.Ctn.Invoke(func(s *token.Service) {
		tokenSvc = s
	})
	if err != nil {
		log.Println(err)
		log.Fatal("MustInitToken failed")
	}
}

func MustInitNft() {
	err := container.Ctn.Invoke(func(s *nft.Service) {
		nftSvc = s
	})
	if err != nil {
		log.Println(err)
		log.Fatal("MustInitNft failed")
	}

}

func MustCheckStorage() {
	err := walletSvc.CheckStorage()
	if err != nil {
		log.Fatal(err)
	}
}

func MustCheckPassword(pw string) {
	isSet, err := walletSvc.IsPasswordSet()
	if err != nil {
		log.Fatal(err)
	}
	if !isSet {
		//_ = walletSvc.ResetWallet()
		var pwd string
		fmt.Println("No password set, please set your password: ")
		for {
			_, _ = fmt.Scanln(&pwd)
			if pwd == "" {
				fmt.Println("Empty input, please input again: ")
			} else {
				err := walletSvc.SetPassword(pwd)
				if err != nil {
					log.Fatal(err)
				}
				log.Fatal("Password has been set, please run your previous command with <-p> arg")
			}
		}
	}

	if pw == "" {
		var input string
		fmt.Println("Please input your password: ")
		_, _ = fmt.Scanln(&input)
		pw = input
	}

	pass, err := walletSvc.CheckPassword(pw)
	if err != nil {
		log.Fatal(err)
	}
	if !pass {
		log.Fatal("Empty or incorrect password")
	}
}
