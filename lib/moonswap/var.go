package moonswap

import (
	conflux "github.com/Conflux-Chain/go-conflux-sdk"
	"github.com/coral922/moonswap-sdk-go/constants"
	moon "github.com/coral922/moonswap-sdk-go/entities"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"sync"
)

type (
	Pair  conflux.Contract
	Token conflux.Contract
)

const (
	SwapModeExactInput  = int(constants.ExactInput)
	SwapModeExactOutPut = int(constants.ExactOutput)
)

var (
	client *conflux.Client

	cfxClientOption = struct {
		nodeUrl string
		option  conflux.ClientOption
	}{}

	CUSDT, _ = moon.NewToken(1029, common.HexToAddress("0x8b8689C7F3014A4D86e4d1D0daAf74A47f5E0f27"), 18, "cUSDT", "conflux USDT")
	l        = sync.Mutex{}
)

func LazyLoad(nodeURL string, option ...conflux.ClientOption) {
	cfxClientOption.nodeUrl = nodeURL
	if len(option) > 0 {
		cfxClientOption.option = option[0]
	}
}

func cfx() *conflux.Client {
	l.Lock()
	defer l.Unlock()
	if client != nil {
		return client
	} else {
		c, err := conflux.NewClient(cfxClientOption.nodeUrl, cfxClientOption.option)
		if err != nil {
			log.Fatal(err)
		}
		client = c
		return client
	}
}

var (
	factory,
	route *conflux.Contract
)

func Factory() *conflux.Contract {
	if factory == nil {
		a := cfx().MustNewAddress(MoonFactoryAddress)
		con, err := cfx().GetContract([]byte(MoonFactoryABI), &a)
		if err != nil {
			log.Fatal(err)
		}
		factory = con
	}
	return factory
}

func Route() *conflux.Contract {
	if route == nil {
		a := cfx().MustNewAddress(MoonRouteAddress)
		con, err := cfx().GetContract([]byte(MoonRouteABI), &a)
		if err != nil {
			log.Fatal(err)
		}
		route = con
	}
	return route
}
