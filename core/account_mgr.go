package core

import (
	"cfxWorld/config"
	"cfxWorld/lib/storage"
	"cfxWorld/lib/util"
	"crypto/sha256"
	"errors"
	"fmt"
	conflux "github.com/Conflux-Chain/go-conflux-sdk"
	"github.com/Conflux-Chain/go-conflux-sdk/types"
	"github.com/Conflux-Chain/go-conflux-sdk/types/cfxaddress"
	"github.com/boltdb/bolt"
	"go.uber.org/dig"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strings"
	"time"
)

type AccountMgr struct {
	*conflux.AccountManager
	db *storage.Storage
}

type amDep struct {
	dig.In
	C *Client `name:"rw"`
	D *storage.Storage
}

func NewAccountMgr(dep amDep) (*AccountMgr, error) {
	if dep.C.AccountManager == nil {
		return nil, errors.New("no account manager found in client")
	}
	return &AccountMgr{
		AccountManager: dep.C.AccountManager.(*conflux.AccountManager),
		db:             dep.D,
	}, nil
}

func (a *AccountMgr) IsPasswordSet() (bool, error) {
	k, err := a.db.Get(PwBucket, PwKey)
	if err != nil {
		return false, err
	}
	return k != "", nil
}

func (a *AccountMgr) SetPassword(newPw string) error {
	p, err := bcrypt.GenerateFromPassword([]byte(newPw), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return a.db.Set(PwBucket, PwKey, string(p))
}

func (a *AccountMgr) CheckPassword(pw string) (bool, error) {
	k, err := a.db.Get(PwBucket, PwKey)
	if err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(k), []byte(pw))
	if err == nil {
		return true, nil
	}
	if !errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, err
	}
	return false, nil
}

func (a *AccountMgr) Reset() error {
	err := os.RemoveAll(config.C.Wallet.KeyStorePath)
	if err != nil {
		return err
	}
	err = a.db.Update(func(tx *bolt.Tx) error {
		err = tx.DeleteBucket([]byte(PwBucket))
		if err != nil {
			return err
		}
		return tx.DeleteBucket([]byte(AccountBucket))
	})
	return err
}

func (a *AccountMgr) GetByAddress(addr string) (*Account, error) {
	str, err := a.db.First(AccountBucket, storage.JsonStrAttrFilter("address", addr))
	if err != nil {
		return nil, err
	}
	return AccountFromJson(str), nil
}

func (a *AccountMgr) GetByName(name string) (*Account, error) {
	str, err := a.db.Get(AccountBucket, name)
	if err != nil {
		return nil, err
	}
	return AccountFromJson(str), nil
}

func (a *AccountMgr) HasName(name string) (bool, error) {
	str, err := a.db.Get(AccountBucket, name)
	if err != nil {
		return false, err
	}
	return str != "", nil
}

func (a *AccountMgr) GetAll() ([]Account, error) {
	strs, err := a.db.All(AccountBucket)
	if err != nil {
		return make([]Account, 0), err
	}
	return AccountsFromJsonArr(strs), nil
}

func (a *AccountMgr) GetAllAddress() ([]string, error) {
	res := make([]string, 0)
	as, err := a.GetAll()
	if err != nil {
		return res, err
	}
	for _, a := range as {
		res = append(res, a.Address.String())
	}
	return res, nil
}

func (a *AccountMgr) Create(name, password string) (*Account, error) {
	ac, err := a.HasName(name)
	if err != nil {
		return nil, err
	}
	if ac {
		return nil, errors.New("account name already exist")
	}
	addr, err := a.AccountManager.Create(ToRealPW(password))
	if err != nil {
		return nil, err
	}
	account := &Account{
		Name:      name,
		Address:   addr,
		Source:    SourceCreate,
		CreatedAt: time.Now(),
	}
	return account, a.db.SetStruct(AccountBucket, account.Name, account)
}

func (a *AccountMgr) Import(name, priKey, password string) (*Account, error) {
	exist, err := a.HasName(name)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, errors.New("account name already exist")
	}
	addr, err := a.AccountManager.ImportKey(priKey, ToRealPW(password))
	if err != nil {
		return nil, err
	}
	account := &Account{
		Name:      name,
		Address:   addr,
		Source:    SourceImport,
		CreatedAt: time.Now(),
	}
	err = a.db.SetStruct(AccountBucket, account.Name, account)
	if err != nil {
		return nil, err
	}
	return account, err
}

func (a *AccountMgr) DeleteByName(accountName, password string) error {
	ac, err := a.GetByName(accountName)
	if err != nil {
		return err
	}
	if ac == nil {
		return errors.New("account not exist")
	}
	err = a.AccountManager.Delete(ac.Address, ToRealPW(password))
	if err != nil {
		return err
	}
	return a.db.Del(AccountBucket, accountName)
}

func (a *AccountMgr) SetCurrentAccount(name string) error {
	ac, err := a.GetByName(name)
	if err != nil {
		return err
	}
	if ac == nil {
		return errors.New("account not exist")
	}
	return a.db.Set(StateBucket, CurrentAccountKey, ac.Address.String())
}

func (a *AccountMgr) GetCurrentAccount() (*Account, error) {
	str, err := a.db.Get(StateBucket, CurrentAccountKey)
	if err != nil {
		return nil, err
	}
	if str == "" {
		addr, err := a.AccountManager.GetDefault()
		if err != nil {
			return nil, err
		}
		str = addr.String()
	}
	ac, err := a.GetByAddress(str)
	if err != nil {
		return nil, err
	}
	return ac, nil
}

func (a *AccountMgr) ChangeAccountName(oldName, newName string) error {
	ac, err := a.GetByName(oldName)
	if err != nil {
		return err
	}
	if ac == nil {
		return errors.New("account not exist")
	}
	exist, err := a.HasName(newName)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("account name already exist")
	}
	ac.Name = newName
	err = a.db.SetStruct(AccountBucket, ac.Name, ac)
	if err != nil {
		return err
	}
	_ = a.db.Del(AccountBucket, oldName)
	return nil
}

func (a *AccountMgr) CheckStorage() error {
	accounts, err := a.GetAllAddress()
	if err != nil {
		return err
	}
	accountsInKS := a.AccountManager.List()
	if len(accounts) != len(accountsInKS) {
		return errors.New("check storage failed")
	}
	for _, account := range accountsInKS {
		if !util.InArrayStr(account.String(), accounts) {
			return errors.New("check storage failed")
		}
	}
	return nil
}

func (a *AccountMgr) ParseStringToAddress(account string, useDefault bool) (*types.Address, error) {
	if account == "" && useDefault {
		ca, err := a.GetCurrentAccount()
		if err != nil {
			return nil, err
		}
		if ca == nil {
			return nil, errors.New("no default account")
		}
		return &ca.Address, nil
	}
	if strings.HasPrefix(account, NamedAccountPrefix) {
		name := strings.TrimPrefix(account, NamedAccountPrefix)
		ac, err := a.GetByName(name)
		if err != nil || ac == nil {
			return nil, fmt.Errorf("account named [%s] not exist", name)
		}
		return &ac.Address, nil
	}
	addr, err := cfxaddress.NewFromBase32(account)
	if err != nil {
		return nil, err
	}
	return &addr, nil
}

func ToRealPW(_ string) string {
	phrase := "i love conflux !!"
	hash := sha256.New()
	hash.Write([]byte(phrase))
	return string(hash.Sum(nil))
}
