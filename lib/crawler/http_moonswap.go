package crawler

import (
	"cfxWorld/lib/util"
	"encoding/json"
	"errors"
	"math/big"
)

type priceAllResp struct {
	Code int     `json:"code"`
	Data []price `json:"data"`
}

type price struct {
	PriceUsd        string `json:"price_usd"`
	ContractAddress string `json:"contract_address"`
}

func PriceOfAllTokenFromMoonSwap() (map[string]*big.Float, error) {
	var res = make(map[string]*big.Float)
	resp, err := Get(PathTokenPrice)
	if err != nil {
		return nil, err
	}
	var r priceAllResp
	err = json.Unmarshal(resp, &r)
	if err != nil {
		return nil, err
	}
	if r.Code != 0 {
		return nil, errors.New(string(resp))
	}
	for _, a := range r.Data {
		res[a.ContractAddress] = util.TextToBigFloatIgnoreErr(a.PriceUsd)
	}
	return res, nil
}
