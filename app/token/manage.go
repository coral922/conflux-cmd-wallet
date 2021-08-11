package token

import (
	"cfxWorld/core"
	"cfxWorld/lib/crawler"
	"cfxWorld/lib/moonswap"
	"github.com/ethereum/go-ethereum/common"
	"log"
)

func (s *Service) SyncTokenAndPairFromMoonSwap() error {
	tokens := make(map[common.Address]*core.CRC20Token)
	getOrAddTokenFunc := func(token *moonswap.Token) *core.CRC20Token {
		addr := token.Address.MustGetCommonAddress()
		if tokens[addr] != nil {
			return tokens[addr]
		}
		name, err := token.Name()
		symbol, err := token.Symbol()
		decimals, err := token.Decimals()
		if err != nil {
			log.Println("load token info error :", err)
			return nil
		}
		scanInfo, err := crawler.TokenInfo(token.Address.String())
		if err != nil {
			scanInfo = new(crawler.TokenResp)
			log.Println("get cfxScan info error :", err)
		}
		t := &core.CRC20Token{
			Address:      *token.Address,
			Name:         name,
			Symbol:       symbol,
			Granularity:  scanInfo.Granularity,
			Decimals:     int(decimals),
			TransferType: "ERC20",
			SupportSwap:  true,
			OfficialCert: scanInfo.IsCustodianToken,
		}
		err = s.tm.UpsertToken(t)
		if err != nil {
			log.Println("storage error :", err)
		}
		log.Printf("sync Token [%s(%s)] . \n", t.Name, t.Symbol)
		tokens[addr] = t
		return t
	}
	count, err := moonswap.GetPairCount()
	if err != nil {
		return err
	}
	for i := 0; i < count; i++ {
		pair, err := moonswap.GetPairByIndex(i)
		if err != nil {
			log.Println("load pair error :", err)
			continue
		}
		tokens, err := pair.Tokens()
		if err != nil {
			log.Println("load token error :", err)
		}
		t0 := getOrAddTokenFunc(tokens[0])
		t1 := getOrAddTokenFunc(tokens[1])
		if t0 == nil || t1 == nil {
			continue
		}
		err = s.tm.UpdatePair(&core.PairInfo{
			Address: *pair.Address,
			Symbol0: t0.Symbol,
			Symbol1: t1.Symbol,
		})
		if err != nil {
			log.Println("store pair error :", err)
		}
	}
	return nil
}

func (s *Service) AddToken(address string) (*core.CRC20Token, error) {
	return s.tm.Register(address)
}

func (s *Service) SetTokenFollowState(tokenName string, followed bool) error {
	return s.tm.SetTokenFollowState(tokenName, followed)
}

func (s *Service) TokenList() ([]core.CRC20Token, error) {
	return s.tm.GetAll()
}

func (s *Service) DetailedTokenList() ([]core.DetailedCRC20Token, error) {
	return s.tm.GetAllDetailed()
}

func (s *Service) DeleteToken(tokenName string) error {
	return s.tm.DeleteTokenByName(tokenName)
}

func (s *Service) Reset() error {
	return s.tm.Reset()
}
