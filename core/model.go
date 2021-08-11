package core

import (
	"encoding/json"
	"github.com/Conflux-Chain/go-conflux-sdk/types"
	"github.com/coral922/moonswap-sdk-go/constants"
	moon "github.com/coral922/moonswap-sdk-go/entities"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"log"
	"math/big"
	"sort"
	"time"
)

const (
	SourceImport = "import"
	SourceCreate = "create"
)

type Account struct {
	Name      string        `json:"name"`
	Address   types.Address `json:"address"`
	Source    string        `json:"source"`
	CreatedAt time.Time     `json:"created_at"`
}

func AccountFromJson(j string) *Account {
	if j == "" {
		return nil
	}
	var a Account
	err := json.Unmarshal([]byte(j), &a)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &a
}

func AccountsFromJsonArr(ja []string) []Account {
	res := make([]Account, 0)
	for _, j := range ja {
		var a Account
		err := json.Unmarshal([]byte(j), &a)
		if err != nil {
			log.Println(err)
		} else {
			res = append(res, a)
		}
	}
	return res
}

type DetailedAccount struct {
	Account
	CfxBalance   *hexutil.Big
	TokenBalance map[string]*hexutil.Big
	TokenStaked  map[string]map[string]*hexutil.Big
	NftBalance   map[string][]*hexutil.Big
	//Nft          interface{}
}

type DetailedAccountList struct {
	CRC20List   []string
	CRC1155List []string
	List        []DetailedAccount
}

type CRC20Token struct {
	Address      types.Address `json:"address"`
	Name         string        `json:"name"`
	Symbol       string        `json:"symbol"`
	Granularity  int           `json:"granularity"`
	Decimals     int           `json:"decimals"`
	TransferType string        `json:"transfer_type"`
	SupportSwap  bool          `json:"support_swap"`
	OfficialCert bool          `json:"official_cert"`
	Followed     bool          `json:"followed"`
}

func (c *CRC20Token) MoonSwapToken() *moon.Token {
	t, err := moon.NewToken(constants.Mainnet, c.Address.MustGetCommonAddress(), c.Decimals, c.Symbol, c.Name)
	if err != nil {
		log.Println(err)
	}
	return t
}

type PairInfo struct {
	Address types.Address `json:"address"`
	Symbol0 string        `json:"token_0"`
	Symbol1 string        `json:"token_1"`
	Token0  *CRC20Token   `json:"-"`
	Token1  *CRC20Token   `json:"-"`
}

func (c *PairInfo) Key() string {
	return c.Address.String()
}

func (c *PairInfo) MoonSwapPair(amount0, amount1 *big.Int) *moon.Pair {
	t0Amount, err := moon.NewTokenAmount(c.Token0.MoonSwapToken(), amount0)
	if err != nil {
		log.Println(err)
	}
	t1Amount, err := moon.NewTokenAmount(c.Token1.MoonSwapToken(), amount1)
	if err != nil {
		log.Println(err)
	}
	pair, err := moon.NewPair(t0Amount, t1Amount)
	if err != nil {
		log.Println(err)
	}
	return pair
}

func PairInfoFromJson(j string) *PairInfo {
	if j == "" {
		return nil
	}
	var a PairInfo
	err := json.Unmarshal([]byte(j), &a)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &a
}

func PairInfosFromJsonArr(ja []string) []PairInfo {
	res := make([]PairInfo, 0)
	for _, j := range ja {
		var a PairInfo
		err := json.Unmarshal([]byte(j), &a)
		if err != nil {
			log.Println(err)
		} else {
			res = append(res, a)
		}
	}
	return res
}

type CRC1155Token struct {
	Address  types.Address `json:"address"`
	Name     string        `json:"name"`
	Symbol   string        `json:"symbol"`
	Followed bool          `json:"followed"`
}

type DetailedCRC20Token struct {
	CRC20Token
	TotalSupply *hexutil.Big
	PriceUSD    string
	GasFree     bool
	StorageFree bool
}

func CRC20TokenFromJson(j string) *CRC20Token {
	if j == "" {
		return nil
	}
	var a CRC20Token
	err := json.Unmarshal([]byte(j), &a)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &a
}

func CRC20TokensFromJsonArr(ja []string) []CRC20Token {
	res := make([]CRC20Token, 0)
	for _, j := range ja {
		var a CRC20Token
		err := json.Unmarshal([]byte(j), &a)
		if err != nil {
			log.Println(err)
		} else {
			res = append(res, a)
		}
	}
	//sort
	sort.SliceStable(res, func(i, j int) bool {
		return res[i].Followed && !res[j].Followed
	})
	return res
}

func CRC1155TokenFromJson(j string) *CRC1155Token {
	if j == "" {
		return nil
	}
	var a CRC1155Token
	err := json.Unmarshal([]byte(j), &a)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &a
}

func CRC1155TokensFromJsonArr(ja []string) []CRC1155Token {
	res := make([]CRC1155Token, 0)
	for _, j := range ja {
		var a CRC1155Token
		err := json.Unmarshal([]byte(j), &a)
		if err != nil {
			log.Println(err)
		} else {
			res = append(res, a)
		}
	}
	//sort
	sort.SliceStable(res, func(i, j int) bool {
		return res[i].Followed && !res[j].Followed
	})
	return res
}
