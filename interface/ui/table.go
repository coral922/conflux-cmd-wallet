package ui

import (
	"cfxWorld/app/tx"
	"cfxWorld/core"
	"cfxWorld/lib/util"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
	"strings"
)

var (
	accountListHeader = []string{"NAME", "ADDRESS", "IMPORT TYPE", "CREATE TIME"}
	tokenListHeader   = []string{"FOLLOWED", "NAME", "SYMBOL", "ADDRESS", "SUPPORT TRANSFER", "SUPPORT SWAP"}
	nftListHeader     = []string{"FOLLOWED", "NAME", "SYMBOL", "ADDRESS"}

	detailedTokenListHeader = []string{"FOLLOWED", "NAME", "SYMBOL", "ADDRESS", "PRICE", "SUPPORT TRANSFER", "SUPPORT SWAP", "HAS SPONSOR", "TOTAL SUPPLY"}

	txResultHeader   = []string{"RESULT", "FROM", "TO", "TOKEN", "TOKEN ID", "VALUE", "GAS USED", "STORAGE", "HAS SPONSOR", "MESSAGE", "VIEW"}
	swapResultHeader = []string{"RESULT", "FROM", "SENT", "GOT", "GAS USED", "STORAGE", "HAS SPONSOR", "MESSAGE", "VIEW"}
)

func AccountListTable(curAccountName string, list []core.Account) string {
	var data = make([][]string, 0)
	for _, a := range list {
		cur := ""
		if a.Name == curAccountName {
			cur = "*"
		}
		row := []string{cur, a.Name, a.Address.String(), a.Source, a.CreatedAt.Format("2006-01-02 15:04:05")}
		data = append(data, row)
	}
	return util.BasicTable(append([]string{""}, accountListHeader...), data, nil)
}

func DetailedAccountListTable(curAccountName string, list core.DetailedAccountList) string {
	var data = make([][]string, 0)
	header := append([]string{""}, detailedAccountHeader(list)...)
	for _, a := range list.List {
		cur := ""
		if a.Name == curAccountName {
			cur = "*"
		}
		row := []string{
			cur,
			a.Name,
			a.Address.String(),
			a.Source,
			a.CreatedAt.Format("2006-01-02 15:04:05"),
			util.AmountToFloatStr(a.CfxBalance),
		}
		for _, t := range list.CRC20List {
			row = append(row, util.AmountToFloatStr(a.TokenBalance[t]))
		}
		for _, t := range list.CRC1155List {
			row = append(row, util.NftIDToStr(a.NftBalance[t]))
		}
		data = append(data, row)
	}
	return util.BasicTable(header, data, nil)
}

func DetailedAccountTable(list core.DetailedAccountList) string {
	var data = make([][]string, 0)
	for _, a := range list.List {
		row := []string{
			a.Name,
			a.Address.String(),
			a.Source,
			a.CreatedAt.Format("2006-01-02 15:04:05"),
			util.AmountToFloatStr(a.CfxBalance),
		}
		for _, t := range list.CRC20List {
			row = append(row, util.AmountToFloatStr(a.TokenBalance[t]))
		}
		for _, t := range list.CRC1155List {
			row = append(row, util.NftIDToStr(a.NftBalance[t]))
		}
		data = append(data, row)
	}
	return util.VerticalTable(detailedAccountHeader(list), data)
}

func detailedAccountHeader(list core.DetailedAccountList) []string {
	header := append([]string{"CFX"}, list.CRC20List...)
	nft := make([]string, 0)
	for _, v := range list.CRC1155List {
		nft = append(nft, v+"(NFT)")
	}
	header = append(header, nft...)
	header = append(accountListHeader, header...)
	return header
}

func AccountSimpleTable(account *core.Account) string {
	return util.SimpleTable([][]string{{account.Name, account.Address.String()}})
}

func AccountsSimpleTable(accounts []core.Account) string {
	d := make([][]string, 0)
	for _, a := range accounts {
		row := []string{a.Name, a.Address.String()}
		d = append(d, row)
	}
	return util.SimpleTable(d)
}

func TokenListTable(list []core.CRC20Token) string {
	var data = make([][]string, 0)
	for _, a := range list {
		var f string
		if a.Followed {
			f = "*"
		}
		row := []string{f, a.Name, a.Symbol, a.Address.String(), util.BoolToYesOrNo(a.TransferType == "ERC20"), util.BoolToYesOrNo(a.SupportSwap)}
		data = append(data, row)
	}
	return util.BasicTable(tokenListHeader, data, nil)
}

func NftListTable(list []core.CRC1155Token) string {
	var data = make([][]string, 0)
	for _, a := range list {
		var f string
		if a.Followed {
			f = "*"
		}
		row := []string{f, a.Name, a.Symbol, a.Address.String()}
		data = append(data, row)
	}
	return util.BasicTable(nftListHeader, data, nil)
}

func DetailedTokenListTable(list []core.DetailedCRC20Token) string {
	var data = make([][]string, 0)
	for _, a := range list {
		var f string
		if a.Followed {
			f = "*"
		}
		sponsor := "NO SPONSOR"
		gasCover := "GAS"
		storageCover := "STORAGE"
		cover := make([]string, 0)
		if a.GasFree {
			cover = append(cover, gasCover)
		}
		if a.StorageFree {
			cover = append(cover, storageCover)
		}
		if len(cover) > 0 {
			sponsor = strings.Join(cover, ",")
		}
		row := []string{f, a.Name, a.Symbol, a.Address.String(), a.PriceUSD, util.BoolToYesOrNo(a.TransferType == "ERC20"), util.BoolToYesOrNo(a.SupportSwap),
			sponsor, util.AmountToFloatStr(a.TotalSupply, 2)}
		data = append(data, row)
	}
	return util.BasicTable(detailedTokenListHeader, data, nil)
}

func TokenSimpleTable(t *core.CRC20Token) string {
	return util.SimpleTable([][]string{{t.Name, t.Symbol, t.Address.String()}})
}

func NftSimpleTable(t *core.CRC1155Token) string {
	return util.SimpleTable([][]string{{t.Name, t.Symbol, t.Address.String()}})
}

func TxResultTable(rec *tx.TransferReceipt) string {
	var data = make([][]string, 0)
	result := "SUCCESS"
	sponsor := "NO SPONSOR"
	gasCover := "GAS"
	storageCover := "STORAGE"
	cover := make([]string, 0)
	msg := ""
	if rec.TxExecErrorMsg != nil {
		result = "FAILED"
		msg = *rec.TxExecErrorMsg
	}
	if rec.GasCoveredBySponsor {
		cover = append(cover, gasCover)
	}
	if rec.StorageCoveredBySponsor {
		cover = append(cover, storageCover)
	}
	if len(cover) > 0 {
		sponsor = strings.Join(cover, ",")
	}
	data = append(data, []string{
		result,
		rec.From.String(),
		rec.To.String(),
		rec.Token,
		rec.NftTokenID,
		util.AmountToFloatStr(rec.Value, 5),
		rec.GasFee.ToInt().Text(10) + "drip" + fmt.Sprintf("(%s cfx)", util.AmountToFloatStr(rec.GasFee, 15)),
		fmt.Sprintf("%ddrip", uint64(rec.StorageCollateralized)) + fmt.Sprintf("(%s cfx)",
			util.AmountToFloatStr((*hexutil.Big)(new(big.Int).SetUint64(uint64(rec.StorageCollateralized))), 15)),
		sponsor,
		msg,
		core.TxURL(rec.TransactionHash),
	})
	return util.VerticalTable(txResultHeader, data)
}

func SwapResultTable(rec *tx.SwapReceipt) string {
	var data = make([][]string, 0)
	result := "SUCCESS"
	sponsor := "NO SPONSOR"
	gasCover := "GAS"
	storageCover := "STORAGE"
	cover := make([]string, 0)
	msg := ""
	if rec.TxExecErrorMsg != nil {
		result = "FAILED"
		msg = *rec.TxExecErrorMsg
	}
	if rec.GasCoveredBySponsor {
		cover = append(cover, gasCover)
	}
	if rec.StorageCoveredBySponsor {
		cover = append(cover, storageCover)
	}
	if len(cover) > 0 {
		sponsor = strings.Join(cover, ",")
	}
	data = append(data, []string{
		result,
		rec.From.String(),
		util.AmountToFloatStr(rec.SendAmount, 15) + " " + rec.SendToken,
		util.AmountToFloatStr(rec.GetAmount, 15) + " " + rec.GetToken,
		rec.GasFee.ToInt().Text(10) + "drip" + fmt.Sprintf("(%s cfx)", util.AmountToFloatStr(rec.GasFee, 15)),
		fmt.Sprintf("%ddrip", uint64(rec.StorageCollateralized)) + fmt.Sprintf("(%s cfx)",
			util.AmountToFloatStr((*hexutil.Big)(new(big.Int).SetUint64(uint64(rec.StorageCollateralized))), 15)),
		sponsor,
		msg,
		core.TxURL(rec.TransactionHash),
	})
	return util.VerticalTable(swapResultHeader, data)
}
