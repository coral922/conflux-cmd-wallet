package core

import (
	"cfxWorld/lib/crawler"
	"cfxWorld/lib/moonswap"
	"cfxWorld/lib/standard"
	"cfxWorld/lib/storage"
	"cfxWorld/lib/util"
	"errors"
	"fmt"
	"github.com/Conflux-Chain/go-conflux-sdk/types"
	"github.com/Conflux-Chain/go-conflux-sdk/types/cfxaddress"
	"github.com/boltdb/bolt"
	moon "github.com/coral922/moonswap-sdk-go/entities"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"go.uber.org/dig"
	"log"
	"math/big"
	"sort"
	"sync"
)

type TokenMgr struct {
	db    *storage.Storage
	cache sync.Map
	mu    sync.Mutex
}

type tmDep struct {
	dig.In
	D *storage.Storage
}

func NewTokenMgr(dep tmDep) *TokenMgr {
	return &TokenMgr{
		db:    dep.D,
		cache: sync.Map{},
		mu:    sync.Mutex{},
	}
}

func (m *TokenMgr) IsEmpty() (bool, error) {
	return m.db.HasBucket(TokenBucket)
}

func (m *TokenMgr) Register(tokenAddress string) (*CRC20Token, error) {
	//fetch data via conflux scan instead of contract
	resp, err := crawler.TokenInfo(tokenAddress)
	if err != nil {
		return nil, err
	}
	if resp.TransferType != "ERC20" {
		log.Printf("WARN: the %s token of %s is not CRC20\n", resp.Symbol, tokenAddress)
	}
	if !resp.IsCustodianToken {
		log.Printf("WARN: the %s token of %s is not registerd in official website\n", resp.Symbol, tokenAddress)
	}
	exist, err := m.HasSymbol(resp.Symbol)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, fmt.Errorf("token [%s] already exist", resp.Symbol)
	}
	t := CRC20Token{
		Address:      cfxaddress.MustNewFromBase32(tokenAddress),
		Name:         resp.Name,
		Symbol:       resp.Symbol,
		Granularity:  resp.Granularity,
		Decimals:     resp.Decimals,
		TransferType: resp.TransferType,
		OfficialCert: resp.IsCustodianToken,
		SupportSwap:  false,
		Followed:     true,
	}
	err = m.db.SetStruct(TokenBucket, t.Symbol, &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (m *TokenMgr) AddToken(token *CRC20Token) error {
	exist, err := m.HasSymbol(token.Symbol)
	if err != nil {
		return err
	}
	if exist {
		return fmt.Errorf("token [%s] already exist", token.Symbol)
	}
	return m.db.SetStruct(TokenBucket, token.Symbol, token)
}

func (m *TokenMgr) UpsertToken(token *CRC20Token) error {
	old, err := m.GetTokenBySymbol(token.Symbol)
	if err != nil {
		return err
	}
	if old != nil {
		token.Followed = old.Followed
	}
	return m.db.SetStruct(TokenBucket, token.Symbol, token)
}

func (m *TokenMgr) GetAll() ([]CRC20Token, error) {
	strs, err := m.db.All(TokenBucket)
	if err != nil {
		return make([]CRC20Token, 0), err
	}
	return CRC20TokensFromJsonArr(strs), nil
}

func (m *TokenMgr) GetAllPairInfo() ([]PairInfo, error) {
	ret := make([]PairInfo, 0)
	strs, err := m.db.All(PairBucket)
	if err != nil {
		return ret, err
	}
	for _, p := range PairInfosFromJsonArr(strs) {
		t0, err := m.GetTokenBySymbol(p.Symbol0)
		t1, err := m.GetTokenBySymbol(p.Symbol1)
		if err != nil {
			log.Println(err)
			continue
		}
		p.Token0 = t0
		p.Token1 = t1
		ret = append(ret, p)
	}
	return ret, nil
}

func (m *TokenMgr) GetFollowed() ([]CRC20Token, error) {
	strs, err := m.db.All(TokenBucket, storage.JsonBoolAttrFilter("followed", true))
	if err != nil {
		return make([]CRC20Token, 0), err
	}
	return CRC20TokensFromJsonArr(strs), nil
}

func (m *TokenMgr) GetAllDetailed() ([]DetailedCRC20Token, error) {
	tokens, err := m.GetAll()
	if err != nil {
		return make([]DetailedCRC20Token, 0), err
	}
	res := make([]DetailedCRC20Token, len(tokens))
	wg := sync.WaitGroup{}
	wg.Add(len(tokens))
	for i := 0; i < len(tokens); i++ {
		t := tokens[i]
		index := i
		go func() {
			defer wg.Done()
			dt := DetailedCRC20Token{
				CRC20Token: t,
			}
			sponsor, err := crawler.SponsorOf(t.Address)
			if err != nil {
				log.Println(err)
			}
			dt.GasFree = sponsor.Gas
			dt.StorageFree = sponsor.Collateral
			total, err := crawler.CRC20TotalSupply(t.Address)
			if err != nil {
				log.Println(err)
			}
			dt.TotalSupply = total
			dt.PriceUSD = "-"
			price, err := m.GetAccurateUSDPrice(t.Symbol)
			if err != nil {
				//log.Println(err)
			} else {
				dt.PriceUSD = "$" + price.ToSignificant(6)
			}
			res[index] = dt
		}()
	}
	wg.Wait()
	return res, nil
}

func (m *TokenMgr) pairs() ([]*moon.Pair, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, exist := m.cache.Load("pair")
	if exist {
		return v.([]*moon.Pair), nil
	}
	ret := make([]*moon.Pair, 0)
	pi, err := m.GetAllPairInfo()
	if err != nil {
		return ret, err
	}
	wg := sync.WaitGroup{}
	wg.Add(len(pi))
	ret = make([]*moon.Pair, len(pi))
	for i := 0; i < len(pi); i++ {
		j := i
		go func() {
			defer wg.Done()
			ctr, err := moonswap.NewPair(pi[j].Address)
			if err != nil {
				log.Println(err)
				return
			}
			rs, err := ctr.Reserves()
			if err != nil {
				log.Println(err)
				return
			}
			ret[j] = pi[j].MoonSwapPair(rs[0], rs[1])
		}()
	}
	wg.Wait()
	m.cache.Store("pair", ret)
	return ret, nil
}

func (m *TokenMgr) GetAccurateUSDPrice(symbol string) (*moon.Price, error) {
	token, err := m.GetTokenBySymbol(symbol)
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, errors.New("token not exist")
	}
	pairs, err := m.pairs()
	if err != nil {
		return nil, err
	}
	tm, err := moon.NewTokenAmount(token.MoonSwapToken(), util.CoinNum(1))
	if err != nil {
		return nil, err
	}
	ts, err := moon.BestTradeExactIn(pairs, tm, moonswap.CUSDT, nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	if len(ts) == 0 {
		return nil, errors.New("no route")
	}
	return ts[0].ExecutionPrice, nil
}

//func (m *TokenMgr) GetDetailedTokenBySymbol(symbol string) (*DetailedCRC20Token, error) {
//	t, err := m.GetTokenBySymbol(symbol)
//	if err != nil {
//		return nil, err
//	}
//	if t == nil {
//		return nil, nil
//	}
//	dt := &DetailedCRC20Token{
//		CRC20Token: *t,
//	}
//	sponsor, err := m.c.SponsorOf(t.Address)
//	if err != nil {
//		log.Println(err)
//	}
//	dt.GasFree = sponsor.Gas
//	dt.StorageFree = sponsor.Collateral
//	resp, err := crawler.TokenInfo(t.Address.String())
//	if err != nil {
//		log.Println(err)
//	}
//	b := &big.Int{}
//	err = b.UnmarshalText([]byte(resp.TotalSupply))
//	if err != nil {
//		log.Println(err)
//	}
//	dt.TotalSupply = (*hexutil.Big)(b)
//	return dt, nil
//}

func (m *TokenMgr) GetTokenBySymbol(symbol string) (*CRC20Token, error) {
	j, err := m.db.Get(TokenBucket, symbol)
	if err != nil {
		return nil, err
	}
	t := CRC20TokenFromJson(j)
	return t, nil
}

func (m *TokenMgr) SetTokenFollowState(symbol string, follow bool) error {
	token, err := m.GetTokenBySymbol(symbol)
	if err != nil {
		return err
	}
	if token == nil {
		return errors.New("token not exist")
	}
	token.Followed = follow
	return m.db.SetStruct(TokenBucket, token.Symbol, &token)
}

func (m *TokenMgr) DeleteTokenByName(symbol string) error {
	token, err := m.GetTokenBySymbol(symbol)
	if err != nil {
		return err
	}
	if token == nil {
		return errors.New("token not exist")
	}
	return m.db.Del(TokenBucket, symbol)
}

func (m *TokenMgr) Reset() error {
	return m.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(TokenBucket))
	})
}

func (m *TokenMgr) HasSymbol(symbol string) (bool, error) {
	j, err := m.db.Get(TokenBucket, symbol)
	if err != nil {
		return false, err
	}
	return j != "", nil
}

func (m *TokenMgr) SortedTokenSymbolList() ([]string, error) {
	ks, err := m.db.Keys(TokenBucket)
	if err != nil {
		return nil, err
	}
	sort.Strings(ks)
	return ks, nil
}

func (m *TokenMgr) SortedFollowedTokenSymbolList() ([]string, error) {
	ks, err := m.db.Keys(TokenBucket, storage.JsonBoolAttrFilter("followed", true))
	if err != nil {
		return nil, err
	}
	sort.Strings(ks)
	return ks, nil
}

func (m *TokenMgr) GetBalance(symbol string, addr types.Address) (*hexutil.Big, error) {
	token, err := m.GetTokenBySymbol(symbol)
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, fmt.Errorf("token [%s] not exist", symbol)
	}
	ctr, err := crawler.Contract([]byte(standard.CRC20BaseABI), token.Address)
	if err != nil {
		return nil, err
	}
	var b *big.Int
	err = ctr.Call(nil, &b, standard.CRC20BalanceMethod, addr.MustGetCommonAddress())
	if err != nil {
		return nil, err
	}
	return (*hexutil.Big)(b), nil
}

func (m *TokenMgr) UpdatePair(info *PairInfo) (err error) {
	return m.db.SetStruct(PairBucket, info.Key(), info)
}

func (m *TokenMgr) GetBalanceOfFollowing(addr types.Address) map[string]*hexutil.Big {
	res := make(map[string]*hexutil.Big)
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
			ctr, err := crawler.Contract([]byte(standard.CRC20BaseABI), token.Address)
			if err != nil {
				log.Println(err)
			}
			var b *big.Int
			err = ctr.Call(nil, &b, standard.CRC20BalanceMethod, addr.MustGetCommonAddress())
			if err != nil {
				log.Println(err)
			}
			a.Store(token.Symbol, (*hexutil.Big)(b))
		}()
	}
	wg.Wait()
	a.Range(func(key, value interface{}) bool {
		res[key.(string)] = value.(*hexutil.Big)
		return true
	})
	return res
}
