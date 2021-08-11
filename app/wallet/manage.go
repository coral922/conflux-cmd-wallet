package wallet

import (
	"cfxWorld/core"
	"cfxWorld/lib/standard"
	"cfxWorld/lib/util"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"log"
	"strings"
)

func (s *Service) ImportKey(name, priKey, password string) (*core.Account, error) {
	return s.am.Import(name, priKey, password)
}

func (s *Service) CurrentAccount() (*core.Account, error) {
	return s.am.GetCurrentAccount()
}

func (s *Service) SetCurrent(name string) error {
	return s.am.SetCurrentAccount(name)
}

func (s *Service) CreateAccount(name, password string) (*core.Account, error) {
	return s.am.Create(name, password)
}

func (s *Service) ChangeAccountName(oldName, newName string) error {
	return s.am.ChangeAccountName(oldName, newName)
}

func (s *Service) CreateAccountBatch(namePrefix, password string, amount int) []core.Account {
	res := make([]core.Account, 0)
	for i := 0; i < amount; i++ {
		name := namePrefix
		if amount > 1 {
			name = fmt.Sprintf("%s_%d", name, i)
		}
		account, err := s.CreateAccount(name, password)
		if err != nil {
			log.Println("ERR: create account failed : ", err)
			continue
		}
		res = append(res, *account)
	}
	return res
}

func (s *Service) AccountList() ([]core.Account, error) {
	return s.am.GetAll()
}

func (s *Service) GetBalance(symbol string, identifier string) (*hexutil.Big, error) {
	address, err := s.am.ParseStringToAddress(identifier, true)
	if err != nil {
		return nil, err
	}
	if strings.ToUpper(symbol) == standard.CFXToken {
		return s.c.GetBalance(*address)
	}
	return s.tm.GetBalance(symbol, *address)
}

func (s *Service) GetNftBalance(symbol string, identifier string) ([]*hexutil.Big, error) {
	address, err := s.am.ParseStringToAddress(identifier, true)
	if err != nil {
		return nil, err
	}
	return s.nm.GetBalance(symbol, *address)
}

func (s *Service) enrichAccountsInfo(accounts []core.Account) (core.DetailedAccountList, error) {
	var res core.DetailedAccountList
	CRC20List, err := s.tm.SortedFollowedTokenSymbolList()
	if err != nil {
		return res, err
	}
	res.CRC20List = CRC20List
	CRC1155List, err := s.nm.SortedFollowedTokenSymbolList()
	if err != nil {
		return res, err
	}
	res.CRC1155List = CRC1155List
	l := make([]core.DetailedAccount, 0)
	for _, a := range accounts {
		b, err := s.c.GetBalance(a.Address)
		if err != nil {
			log.Println(err)
		}
		da := core.DetailedAccount{
			Account:    a,
			CfxBalance: b,
		}
		da.TokenBalance = s.tm.GetBalanceOfFollowing(a.Address)
		da.NftBalance = s.nm.GetBalanceOfFollowing(a.Address)
		l = append(l, da)
	}
	res.List = l
	return res, nil
}

func (s *Service) DetailedAccountList() (core.DetailedAccountList, error) {
	var res core.DetailedAccountList
	accounts, err := s.AccountList()
	if err != nil {
		return res, err
	}
	return s.enrichAccountsInfo(accounts)
}

func (s *Service) DetailedAccount(name string) (core.DetailedAccountList, error) {
	var res core.DetailedAccountList
	a, err := s.am.GetByName(name)
	if err != nil {
		return res, err
	}
	if a == nil {
		return res, errors.New("account not exist")
	}
	return s.enrichAccountsInfo([]core.Account{*a})
}

func (s *Service) DeleteAccountByName(name string, password string) error {
	return s.am.DeleteByName(name, password)
}

func (s *Service) ExportPrivateKey(name, password string) (string, error) {
	account, err := s.am.GetByName(name)
	if err != nil {
		return "", err
	}
	if account == nil {
		return "", errors.New("account not exist")
	}
	return s.am.Export(account.Address, core.ToRealPW(password))
}

func (s *Service) ExportAllPrivateKeyToCSV(password, csvFile string) error {
	list, err := s.am.GetAll()
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return errors.New("no account found")
	}
	dat := make([][]string, 0)
	for _, a := range list {
		k, err := s.am.Export(a.Address, core.ToRealPW(password))
		if err != nil {
			return err
		}
		dat = append(dat, []string{a.Name, a.Address.String(), k})
	}
	return util.WriteCSV(csvFile, dat)
}

func (s *Service) ResetWallet() error {
	return s.am.Reset()
}
