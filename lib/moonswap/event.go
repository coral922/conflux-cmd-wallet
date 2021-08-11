package moonswap

import (
	"cfxWorld/lib/standard"
	conflux "github.com/Conflux-Chain/go-conflux-sdk"
	"github.com/Conflux-Chain/go-conflux-sdk/types"
	"github.com/Conflux-Chain/go-conflux-sdk/types/cfxaddress"
	"github.com/coral922/moonswap-sdk-go/constants"
	moon "github.com/coral922/moonswap-sdk-go/entities"
	"log"
)

type SwapResult struct {
	Input, Output *moon.TokenAmount
}

func ParseSwapResult(trade *moon.Trade, tx *types.TransactionReceipt) *SwapResult {
	if len(tx.Logs) == 0 {
		return nil
	}
	var r SwapResult
	//assert same address
	followIn, followOut := tx.From, tx.From
	input, output := trade.InputAmount().Token.Address, trade.OutputAmount().Token.Address
	iContract, _ := NewToken(cfxaddress.MustNewFromCommon(input, cfxaddress.NetowrkTypeMainnetID))
	oContract, _ := NewToken(cfxaddress.MustNewFromCommon(output, cfxaddress.NetowrkTypeMainnetID))
	for _, l := range tx.Logs {
		if trade.InputAmount().Token.Equals(moon.WCFX[constants.Mainnet]) {
			followIn = standard.ZeroAddr
		}
		if trade.OutputAmount().Token.Equals(moon.WCFX[constants.Mainnet]) {
			followOut = standard.ZeroAddr
		}
		if l.Address.MustGetCommonAddress() == input &&
			l.Topics[0] == standard.CRC20TransferEventID &&
			*l.Topics[1].ToCommonHash() == followIn.MustGetCommonAddress().Hash() {
			//input event
			var e standard.CRC20TransferEvent
			err := (*conflux.Contract)(iContract).DecodeEvent(&e, "Transfer", l)
			if err != nil {
				log.Println(err)
				return nil
			}
			i, err := moon.NewTokenAmount(trade.InputAmount().Token, e.Value)
			if err != nil {
				log.Println(err)
				return nil
			}
			r.Input = i
		}
		if l.Address.MustGetCommonAddress() == output &&
			l.Topics[0] == standard.CRC20TransferEventID &&
			*l.Topics[2].ToCommonHash() == followOut.MustGetCommonAddress().Hash() {
			//output event
			var e standard.CRC20TransferEvent
			err := (*conflux.Contract)(oContract).DecodeEvent(&e, "Transfer", l)
			if err != nil {
				log.Println(err)
				return nil
			}
			o, err := moon.NewTokenAmount(trade.OutputAmount().Token, e.Value)
			if err != nil {
				log.Println(err)
				return nil
			}
			r.Output = o
		}
	}
	return &r
}