package cmdline

import (
	"cfxWorld/app/tx"
	"cfxWorld/interface/ui"
	"cfxWorld/lib/moonswap"
	"cfxWorld/lib/util"
	"fmt"
	moon "github.com/coral922/moonswap-sdk-go/entities"
	"log"
	"math/big"
	"strings"
)

func Send(token, from, to, amount, password string) {
	b, err := walletSvc.GetBalance(token, from)
	if err != nil {
		log.Println("ERROR when check balance: ", err)
	}
	bStr := util.AmountToFloatStr(b, 18)
	fmt.Printf("Your %s Balance: %s\n", token, bStr)
	if amount == "all" {
		amount = bStr
	}
	fmt.Printf("You will send %s %s to %s \n", amount, token, to)
	if util.SecondConfirm() {
		rec, err := txSvc.Send(token, from, to, amount, password)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(ui.TxResultTable(rec))
	}
}

func SendNft(token, tokenID, from, to, password string) {
	b, err := walletSvc.GetNftBalance(token, from)
	if err != nil {
		log.Println("ERROR when check balance: ", err)
	}
	fmt.Printf("Your %s Balance: %s\n", token, util.NftIDToStr(b))
	fmt.Printf("You will send %s [%s] to %s \n", token, tokenID, to)
	if util.SecondConfirm() {
		rec, err := txSvc.SendNft(token, tokenID, from, to, password)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(ui.TxResultTable(rec))
	}
}

func Swap(mode int, addr, input, output, amount, password string, slipTolerance int) {
	b, err := walletSvc.GetBalance(input, addr)
	if err != nil {
		log.Println("ERROR when check balance", err)
	}
	bStr := util.AmountToFloatStr(b, 18)
	fmt.Printf("Your %s Balance: %s\n", input, bStr)

	if amount == "all" && mode == moonswap.SwapModeExactInput {
		amount = bStr
	}
	trade, err := txSvc.MakeSwap(mode, input, output, amount)
	if err != nil {
		log.Fatal(err)
	}
	if slipTolerance <= 0 {
		slipTolerance = tx.DefaultSlippageTolerance
	}
	if mode == moonswap.SwapModeExactInput {
		amountOutMinS, err := trade.MinimumAmountOut(moon.NewPercent(big.NewInt(int64(slipTolerance)), big.NewInt(100)))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("You will send %s %s \n", trade.InputAmount().ToSignificant(10), trade.InputAmount().Symbol)
		fmt.Printf("You will receive %s %s in general\n", trade.OutputAmount().ToSignificant(10), trade.OutputAmount().Symbol)
		fmt.Printf("You will receive %s %s at least \n", amountOutMinS.ToSignificant(10), trade.OutputAmount().Symbol)
	}
	if mode == moonswap.SwapModeExactOutPut {
		amountInMaxS, err := trade.MaximumAmountIn(moon.NewPercent(big.NewInt(int64(slipTolerance)), big.NewInt(100)))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("You will receive %s %s \n", trade.OutputAmount().ToSignificant(10), trade.OutputAmount().Symbol)
		fmt.Printf("You will send %s %s in general\n", trade.InputAmount().ToSignificant(10), trade.InputAmount().Symbol)
		fmt.Printf("You will send %s %s at most \n", amountInMaxS.ToSignificant(10), trade.InputAmount().Symbol)
	}

	tmp := make([]string, 0)
	for _, p := range trade.Route.Path {
		tmp = append(tmp, p.Symbol)
	}
	fmt.Printf("Token swap path: %s\n", strings.Join(tmp, " -> "))

	if util.SecondConfirm() {
		rec, err := txSvc.DoSwap(addr, trade, password, slipTolerance)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(ui.SwapResultTable(rec))
	}
}
