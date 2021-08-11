package core

import (
	"cfxWorld/lib/crawler"
	"cfxWorld/lib/moonswap"
	"cfxWorld/lib/standard"
	"cfxWorld/lib/util"
	"errors"
	"fmt"
	"github.com/Conflux-Chain/go-conflux-sdk/types"
	"github.com/coral922/moonswap-sdk-go/constants"
	moon "github.com/coral922/moonswap-sdk-go/entities"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"go.uber.org/dig"
	"log"
	"math/big"
	"strings"
	"time"
)

type TxMgr struct {
	c  *Client
	tm *TokenMgr
	nm *NftMgr
}

type txmDep struct {
	dig.In
	C  *Client `name:"rw"`
	TM *TokenMgr
	NM *NftMgr
}

func NewTxMgr(dep txmDep) *TxMgr {
	return &TxMgr{
		c:  dep.C,
		tm: dep.TM,
		nm: dep.NM,
	}
}

func (m *TxMgr) TransferSupportedToken(tokenName string) bool {
	if strings.ToUpper(tokenName) != standard.CFXToken {
		token, err := m.tm.GetTokenBySymbol(tokenName)
		if err != nil {
			log.Println(err)
			return false
		}
		return token != nil && token.TransferType == "ERC20"
	}
	return true
}

func (m *TxMgr) TransferSupportedNft(tokenName string) bool {
	token, err := m.nm.GetTokenBySymbol(tokenName)
	if err != nil {
		log.Println(err)
		return false
	}
	return token != nil
}

func (m *TxMgr) SendTx(tokenName string,
	from types.Address, to types.Address, amount *hexutil.Big,
	txOpt ...TxOption) (*types.UnsignedTransaction, error) {
	if strings.ToUpper(tokenName) == standard.CFXToken {
		tx, err := m.c.CreateUnsignedTransaction(from, to, amount, nil)
		if err != nil {
			return nil, err
		}
		applyTxOpts(&(tx.UnsignedTransactionBase), txOpt...)
		if enough, need := crawler.IsCfxBalanceEnough(tx); !enough {
			return nil, fmt.Errorf("you don't have enough cfx, need: %s", util.AmountToFloatStr(need, 18))
		}
		return &tx, nil
	} else {
		ctr, err := m.CRC20Contract(tokenName)
		if err != nil {
			return nil, err
		}
		tx, err := ctr.CreateUnsignedTransaction(from, nil, standard.CRC20TransferMethod, to.MustGetCommonAddress(), (*big.Int)(amount))
		if err != nil {
			return nil, err
		}
		applyTxOpts(&(tx.UnsignedTransactionBase), txOpt...)
		if enough, need := crawler.IsCfxBalanceEnough(tx); !enough {
			return nil, fmt.Errorf("you don't have enough cfx, need: %s", util.AmountToFloatStr(need, 18))
		}
		return &tx, nil
	}
}

func (m *TxMgr) SendNftTx(tokenName string, tokenID *hexutil.Big,
	from types.Address, to types.Address,
	txOpt ...TxOption) (*types.UnsignedTransaction, error) {
	ctr, err := m.CRC1155Contract(tokenName)
	if err != nil {
		return nil, err
	}
	tx, err := ctr.CreateUnsignedTransaction(from,
		nil,
		standard.CRC1155TransferMethod,
		from.MustGetCommonAddress(),
		to.MustGetCommonAddress(),
		(*big.Int)(tokenID),
		big.NewInt(1),
		[]byte(""))
	if err != nil {
		return nil, err
	}
	applyTxOpts(&(tx.UnsignedTransactionBase), txOpt...)
	if enough, need := crawler.IsCfxBalanceEnough(tx); !enough {
		return nil, fmt.Errorf("you don't have enough cfx, need: %s", util.AmountToFloatStr(need, 18))
	}
	return &tx, nil
}

func (m *TxMgr) MakeSwap(mode int, inputToken, outputToken string, ExactAmount *big.Int) (*moon.Trade, error) {
	if mode != moonswap.SwapModeExactInput && mode != moonswap.SwapModeExactOutPut {
		return nil, errors.New("unknown swap mode")
	}
	var i, o *moon.Token
	if strings.ToUpper(inputToken) == standard.CFXToken {
		i = moon.WCFX[constants.Mainnet]
	} else {
		input, _ := m.tm.GetTokenBySymbol(inputToken)
		if input == nil || !input.SupportSwap {
			return nil, fmt.Errorf("token [%s] not exist or not swappable", inputToken)
		}
		i = input.MoonSwapToken()
	}
	if strings.ToUpper(outputToken) == standard.CFXToken {
		o = moon.WCFX[constants.Mainnet]
	} else {
		output, _ := m.tm.GetTokenBySymbol(outputToken)
		if output == nil || !output.SupportSwap {
			return nil, fmt.Errorf("token [%s] not exist or not swappable", outputToken)
		}
		o = output.MoonSwapToken()
	}
	if i.Equals(o) {
		return nil, errors.New("same token")
	}

	pairs, err := m.tm.pairs()
	if err != nil {
		return nil, err
	}
	var ts []*moon.Trade
	var errT error
	if mode == moonswap.SwapModeExactInput {
		tmIn, err := moon.NewTokenAmount(i, ExactAmount)
		if err != nil {
			return nil, err
		}
		ts, errT = moon.BestTradeExactIn(
			pairs,
			tmIn,
			o,
			nil,
			nil,
			nil,
			nil,
		)
	} else {
		tmOut, err := moon.NewTokenAmount(o, ExactAmount)
		if err != nil {
			return nil, err
		}
		ts, errT = moon.BestTradeExactOut(
			pairs,
			i,
			tmOut,
			nil,
			nil,
			nil,
			nil,
		)
	}
	if errT != nil {
		return nil, err
	}
	if len(ts) == 0 {
		return nil, errors.New("no route for the swap")
	}
	return ts[0], nil
}

func (m *TxMgr) SwapTx(addr types.Address, trade *moon.Trade, slipTolerance int, txOpt ...TxOption) (*types.UnsignedTransaction, error) {
	if trade.InputAmount().Token.Equals(trade.OutputAmount().Token) {
		return nil, errors.New("same token")
	}
	var method string
	var cfxAmount *hexutil.Big
	ctr := NewContract(moonswap.Route())
	var tx types.UnsignedTransaction
	var err error
	if trade.TradeType == constants.ExactInput {
		amountIn := trade.InputAmount().Raw()
		amountOutMinS, _ := trade.MinimumAmountOut(moon.NewPercent(big.NewInt(int64(slipTolerance)), big.NewInt(100)))
		amountOutMin := amountOutMinS.Raw()
		deadline := big.NewInt(time.Now().Add(5 * time.Minute).Unix())
		path := trade.Route.PathToAddress()
		params := []interface{}{
			amountIn,
			amountOutMin,
			path,
			addr.MustGetCommonAddress(),
			deadline,
		}
		method = "swapExactTokensForTokens"
		if trade.InputAmount().Token.Equals(moon.WCFX[constants.Mainnet]) {
			method = "swapExactCFXForTokens"
			cfxAmount = (*hexutil.Big)(trade.InputAmount().Raw())
			params = []interface{}{
				amountOutMin,
				path,
				addr.MustGetCommonAddress(),
				deadline,
			}
		}
		if trade.OutputAmount().Token.Equals(moon.WCFX[constants.Mainnet]) {
			method = "swapExactTokensForCFX"
		}
		tx, err = ctr.CreateUnsignedTransaction(addr, cfxAmount, method, params...)
		if err != nil {
			return nil, err
		}
	} else {
		amountOut := trade.OutputAmount().Raw()
		amountInMaxS, _ := trade.MaximumAmountIn(moon.NewPercent(big.NewInt(int64(slipTolerance)), big.NewInt(100)))
		amountInMax := amountInMaxS.Raw()
		deadline := big.NewInt(time.Now().Add(5 * time.Minute).Unix())
		path := trade.Route.PathToAddress()
		params := []interface{}{
			amountOut,
			amountInMax,
			path,
			addr.MustGetCommonAddress(),
			deadline,
		}
		method = "swapTokensForExactTokens"
		if trade.InputAmount().Token.Equals(moon.WCFX[constants.Mainnet]) {
			method = "swapCFXForExactTokens"
			cfxAmount = (*hexutil.Big)(trade.InputAmount().Raw())
			params = []interface{}{
				amountOut,
				path,
				addr.MustGetCommonAddress(),
				deadline,
			}
		}
		if trade.OutputAmount().Token.Equals(moon.WCFX[constants.Mainnet]) {
			method = "swapTokensForExactCFX"
		}
		tx, err = ctr.CreateUnsignedTransaction(addr, cfxAmount, method, params...)
		if err != nil {
			return nil, err
		}
	}

	applyTxOpts(&(tx.UnsignedTransactionBase), txOpt...)
	if enough, need := crawler.IsCfxBalanceEnough(tx); !enough {
		return nil, fmt.Errorf("you don't have enough cfx, need: %s", util.AmountToFloatStr(need, 18))
	}
	return &tx, nil
}

func (m *TxMgr) Do(tx *types.UnsignedTransaction) (types.Hash, error) {
	return m.c.SendTransaction(*tx)
}

func (m *TxMgr) CRC20Contract(symbol string) (*Contract, error) {
	token, err := m.tm.GetTokenBySymbol(symbol)
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, fmt.Errorf("token [%s] not exist", symbol)
	}
	ctr, err := m.c.GetContract([]byte(standard.CRC20BaseABI), &token.Address)
	if err != nil {
		return nil, err
	}
	return NewContract(ctr), nil
}

func (m *TxMgr) CRC1155Contract(symbol string) (*Contract, error) {
	token, err := m.nm.GetTokenBySymbol(symbol)
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, fmt.Errorf("nft token [%s] not exist", symbol)
	}
	ctr, err := m.c.GetContract([]byte(standard.CRC1155BaseABI), &token.Address)
	if err != nil {
		return nil, err
	}
	return NewContract(ctr), nil
}
