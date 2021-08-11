package nft

import (
	"cfxWorld/core"
	"log"
)

func (s *Service) InitBuiltInNft() error {
	exist, err := s.nm.IsEmpty()
	if err != nil {
		return err
	}
	if !exist {
		for _, a := range core.BuiltInNftAddress {
			_, err := s.nm.Register(a)
			if err != nil {
				log.Println(err)
			}
		}
	}
	return nil
}

func (s *Service) SetTokenFollowState(tokenName string, followed bool) error {
	return s.nm.SetTokenFollowState(tokenName, followed)
}

func (s *Service) AddToken(address string) (*core.CRC1155Token, error) {
	return s.nm.Register(address)
}

func (s *Service) TokenList() ([]core.CRC1155Token, error) {
	return s.nm.GetAll()
}

func (s *Service) DeleteToken(tokenName string) error {
	return s.nm.DeleteTokenByName(tokenName)
}

func (s *Service) Reset() error {
	return s.nm.Reset()
}
