package crawler

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

type TokenResp struct {
	Name             string `json:"name"`
	Symbol           string `json:"symbol"`
	TransferType     string `json:"transferType"`
	Granularity      int    `json:"granularity"`
	IsCustodianToken bool   `json:"isCustodianToken"`
	Decimals         int    `json:"decimals"`
	TotalSupply      string `json:"totalSupply"`
	HolderCount      int    `json:"holderCount"`
}

func TokenInfo(address string) (*TokenResp, error) {
	url := TokenInfoURL + address
	resp, err := Get(url)
	if err != nil {
		return nil, err
	}
	var r TokenResp
	err = json.Unmarshal(resp, &r)
	if err != nil {
		return nil, err
	}
	return &r, err
}

func ABI(address string) (string, error) {
	url := ContractInfoURL + address
	resp, err := Get(url, map[string][]string{
		"fields": {"standard"},
	})
	if err != nil {
		return "", err
	}
	return gjson.GetBytes(resp, "standard").Str, nil
}
