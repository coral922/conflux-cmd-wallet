package core

import (
	"cfxWorld/lib/crawler"
	"cfxWorld/lib/standard"
	"cfxWorld/lib/storage"
	"cfxWorld/lib/util"
	"errors"
	"fmt"
	"github.com/Conflux-Chain/go-conflux-sdk/types"
	"github.com/Conflux-Chain/go-conflux-sdk/types/cfxaddress"
	"github.com/boltdb/bolt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"go.uber.org/dig"
	"log"
	"math/big"
	"sort"
	"sync"
)

type NftMgr struct {
	c  *Client
	db *storage.Storage
}

type nmDep struct {
	dig.In
	D  *storage.Storage
}

func NewNftMgr(dep nmDep) *NftMgr {
	return &NftMgr{
		db: dep.D,
	}
}

func (m *NftMgr) IsEmpty() (bool, error) {
	return m.db.HasBucket(NFTBucket)
}

func (m *NftMgr) Register(contractAddress string) (*CRC1155Token, error) {
	resp, err := crawler.TokenInfo(contractAddress)
	if err != nil {
		return nil, err
	}
	if resp.Symbol == "" {
		return nil, fmt.Errorf("no symbol for the token of %s", contractAddress)
	}
	if resp.TransferType != "ERC1155" {
		return nil, fmt.Errorf("the %s token of %s is not CRC1155", resp.Symbol, contractAddress)
	}
	if !resp.IsCustodianToken {
		log.Printf("WARN: the %s token of %s is not registerd in official website\n", resp.Symbol, contractAddress)
	}

	exist, err := m.HasSymbol(resp.Symbol)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, fmt.Errorf("token [%s] already exist", resp.Symbol)
	}
	t := CRC1155Token{
		Address:      cfxaddress.MustNewFromBase32(contractAddress),
		Name:         resp.Name,
		Symbol:       resp.Symbol,
	}
	err = m.db.SetStruct(NFTBucket, t.Symbol, &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (m *NftMgr) GetAll() ([]CRC1155Token, error) {
	strs, err := m.db.All(NFTBucket)
	if err != nil {
		return make([]CRC1155Token, 0), err
	}
	return CRC1155TokensFromJsonArr(strs), nil
}

func (m *NftMgr) GetFollowed() ([]CRC1155Token, error) {
	strs, err := m.db.All(NFTBucket, storage.JsonBoolAttrFilter("followed", true))
	if err != nil {
		return make([]CRC1155Token, 0), err
	}
	return CRC1155TokensFromJsonArr(strs), nil
}

func (m *NftMgr) GetTokenBySymbol(symbol string) (*CRC1155Token, error) {
	j, err := m.db.Get(NFTBucket, symbol)
	if err != nil {
		return nil, err
	}
	t := CRC1155TokenFromJson(j)
	return t, nil
}

func (m *NftMgr) SetTokenFollowState(symbol string, follow bool) error {
	token, err := m.GetTokenBySymbol(symbol)
	if err != nil {
		return err
	}
	if token == nil {
		return errors.New("token not exist")
	}
	token.Followed = follow
	return m.db.SetStruct(NFTBucket, token.Symbol, &token)
}

func (m *NftMgr) DeleteTokenByName(symbol string) error {
	ac, err := m.GetTokenBySymbol(symbol)
	if err != nil {
		return err
	}
	if ac == nil {
		return errors.New("token not exist")
	}
	return m.db.Del(NFTBucket, symbol)
}

func (m *NftMgr) Reset() error {
	return m.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(NFTBucket))
	})
}

func (m *NftMgr) HasSymbol(symbol string) (bool, error) {
	j, err := m.db.Get(NFTBucket, symbol)
	if err != nil {
		return false, err
	}
	return j != "", nil
}

func (m *NftMgr) SortedTokenSymbolList() ([]string, error) {
	ks, err := m.db.Keys(NFTBucket)
	if err != nil {
		return nil, err
	}
	sort.Strings(ks)
	return ks, nil
}

func (m *NftMgr) SortedFollowedTokenSymbolList() ([]string, error) {
	ks, err := m.db.Keys(NFTBucket, storage.JsonBoolAttrFilter("followed", true))
	if err != nil {
		return nil, err
	}
	sort.Strings(ks)
	return ks, nil
}

func (m *NftMgr) GetBalance(symbol string, addr types.Address) ([]*hexutil.Big, error) {
	token, err := m.GetTokenBySymbol(symbol)
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, fmt.Errorf("token [%s] not exist", symbol)
	}
	ctr, err := m.c.GetContract([]byte(standard.CRC1155BaseABI), &token.Address)
	if err != nil {
		return nil, err
	}
	b := make([]*big.Int, 0)
	err = ctr.Call(nil, &b, standard.CRC1155BalanceMethod, addr.MustGetCommonAddress())
	if err != nil {
		return nil, err
	}
	return util.BigSliceToHexSlice(b), nil
}

func (m *NftMgr) GetBalanceOfFollowing(addr types.Address) map[string][]*hexutil.Big {
	res := make(map[string][]*hexutil.Big)
	a := sync.Map{}
	tokens, err := m.GetFollowed()
	if err != nil {
		log.Println(err)
		return nil
	}
	if len(tokens) == 0 {
		return nil
	}
	wg := sync.WaitGroup{}
	wg.Add(len(tokens))
	for _, t := range tokens {
		token := t
		go func() {
			defer wg.Done()
			ctr, err := m.c.GetContract([]byte(standard.CRC1155BaseABI), &token.Address)
			if err != nil {
				log.Println(err)
			}
			b := make([]*big.Int, 0)
			err = ctr.Call(nil, &b, standard.CRC1155BalanceMethod, addr.MustGetCommonAddress())
			if err != nil {
				log.Println(err)
			}
			a.Store(token.Symbol, util.BigSliceToHexSlice(b))
		}()
	}
	wg.Wait()
	a.Range(func(key, value interface{}) bool {
		res[key.(string)] = value.([]*hexutil.Big)
		return true
	})
	return res
}
