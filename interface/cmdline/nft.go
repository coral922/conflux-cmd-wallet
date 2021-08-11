package cmdline

import (
	"cfxWorld/interface/ui"
	"fmt"
	"log"
)

func NftList() {
	list, err := nftSvc.TokenList()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ui.NftListTable(list))
}

func AddNft(address string) {
	t, err := nftSvc.AddToken(address)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ui.NftSimpleTable(t))
}

func FollowNft(names []string) {
	for _, name := range names {
		err := nftSvc.SetTokenFollowState(name, true)
		if err != nil {
			log.Printf("nft [%s] err : %v \n", name, err)
		} else {
			log.Printf("nft [%s] followed \n", name)
		}
	}
}

func UnFollowNft(names []string) {
	for _, name := range names {
		err := nftSvc.SetTokenFollowState(name, false)
		if err != nil {
			log.Printf("nft [%s] err : %v \n", name, err)
		} else {
			log.Printf("nft [%s] unfollowed \n", name)
		}
	}
}

func DeleteNft(name string) {
	err := nftSvc.DeleteToken(name)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("nft [%s] deleted", name)
}

func ResetNft() {
	err := nftSvc.Reset()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("the nft manager has been reset")
}
