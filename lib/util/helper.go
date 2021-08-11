package util

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
	"strings"
)

func InArrayStr(needle string, haystack []string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

func BoolToYesOrNo(b bool) string {
	if b {
		return "YES"
	}
	return "NO"
}

func BigSliceToHexSlice(h []*big.Int) []*hexutil.Big {
	res := make([]*hexutil.Big, 0)
	for _, b := range h {
		res = append(res, (*hexutil.Big)(b))
	}
	return res
}

func NftIDToStr(balance []*hexutil.Big) string {
	ids := make([]string, 0)
	for _, v := range balance {
		ids = append(ids, v.ToInt().Text(10))
	}
	return strings.Join(ids, ",")
}