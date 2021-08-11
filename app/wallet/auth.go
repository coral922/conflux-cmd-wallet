package wallet

func (s *Service) IsPasswordSet() (bool, error) {
	return s.am.IsPasswordSet()
}

func (s *Service) SetPassword(newPw string) error {
	return s.am.SetPassword(newPw)
}

func (s *Service) CheckPassword(pw string) (bool, error) {
	return s.am.CheckPassword(pw)
}

func (s *Service) CheckStorage() error {
	return s.am.CheckStorage()
}
