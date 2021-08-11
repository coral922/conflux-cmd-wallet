package cmdline

import (
	"cfxWorld/interface/ui"
	"fmt"
	"log"
)

func TokenList() {
	list, err := tokenSvc.TokenList()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ui.TokenListTable(list))
}

func DetailedTokenList() {
	list, err := tokenSvc.DetailedTokenList()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ui.DetailedTokenListTable(list))
}

func AddToken(address string) {
	t, err := tokenSvc.AddToken(address)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ui.TokenSimpleTable(t))
}

func SyncTokenAndPair() {
	err := tokenSvc.SyncTokenAndPairFromMoonSwap()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("sync CRC20 token finished")
}

func DeleteToken(name string) {
	err := tokenSvc.DeleteToken(name)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("token [%s] deleted", name)
}

func FollowTokens(names []string) {
	for _, name := range names {
		err := tokenSvc.SetTokenFollowState(name, true)
		if err != nil {
			log.Printf("token [%s] err : %v \n", name, err)
		} else {
			log.Printf("token [%s] followed \n", name)
		}
	}
}

func UnFollowTokens(names []string) {
	for _, name := range names {
		err := tokenSvc.SetTokenFollowState(name, false)
		if err != nil {
			log.Printf("token [%s] err : %v \n", name, err)
		} else {
			log.Printf("token [%s] unfollowed \n", name)
		}
	}
}

func ResetToken() {
	err := tokenSvc.Reset()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("the token manager has been reset")
}
