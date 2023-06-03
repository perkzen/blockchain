package main

import (
	"blockchain/pkg/blockchain"
	"blockchain/pkg/wallet"
	"fmt"
)

func main() {
	walletM := wallet.NewWallet()
	walletA := wallet.NewWallet()
	walletB := wallet.NewWallet()

	t := wallet.NewTransaction(walletA.PrivateKey(), walletA.PublicKey(), walletA.BlockchainAddress(), walletB.BlockchainAddress(), 1.0)
	bc := blockchain.InitBlockchain(walletM.BlockchainAddress())
	isAdd := bc.AddTransaction(walletA.BlockchainAddress(), walletB.BlockchainAddress(), 1.0, t.GenerateSignature(), walletA.PublicKey())
	fmt.Println(isAdd)
}
