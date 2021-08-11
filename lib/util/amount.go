package util

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
)

const (
	ratio       = 1e18
	defaultPrec = 2
)

func AmountToFloatStr(amount *hexutil.Big, prec ...int) string {
	if amount == nil {
		return "unknown"
	}
	r := big.Rat{}
	p := defaultPrec
	if len(prec) > 0 {
		p = prec[0]
	}
	return r.SetFrac(amount.ToInt(), big.NewInt(ratio)).FloatString(p)
}

func Float64ToAmount(f float64) *hexutil.Big {
	if f < 0 {
		return nil
	}
	var bf big.Float
	amount := bf.Mul(big.NewFloat(f), big.NewFloat(ratio))
	bi, _ := amount.Int(nil)
	return (*hexutil.Big)(bi)
}

func FloatStrToAmount(f string) (*hexutil.Big, error) {
	var bf big.Float
	err := bf.UnmarshalText([]byte(f))
	if err != nil {
		return nil, err
	}
	amount := bf.Mul(&bf, big.NewFloat(ratio))
	bi, _ := amount.Int(nil)
	return (*hexutil.Big)(bi), nil
}

func IntStrToBig(s string) (*hexutil.Big, error) {
	b := new(big.Int)
	err := b.UnmarshalText([]byte(s))
	if err != nil {
		return nil, err
	}
	return (*hexutil.Big)(b), nil
}

func TextToBigFloatIgnoreErr(s string) *big.Float {
	b := big.NewFloat(0)
	_ = b.UnmarshalText([]byte(s))
	return b
}

func CoinNum(value int64) *big.Int {
	return new(big.Int).Mul(big.NewInt(value), big.NewInt(ratio))
}