package server

import "blockchain/pkg/wallet"

func (s *Server) GetWallet() *wallet.Wallet {
	w, ok := s.cache[WALLET].(*wallet.Wallet)
	if !ok {
		w = wallet.NewWallet()
		s.cache[WALLET] = w
	}
	return w
}
