package cmdline

import (
	"cfxWorld/lib/crawler"
	"cfxWorld/lib/util"
	"log"
)

func OpenConfluxScan(address ...string) {
	url := crawler.CfxScanBaseURL
	if len(address) > 0 && address[0] != "" {
		url = url + "/address/" + address[0]
	}
	err := util.OpenBrowser(url)
	if err != nil {
		log.Fatal(err)
	}
}

func TokenInfo(contractAddr string) {

}
