package blockchain

import (
	"blockchain/pkg/wallet"
	"fmt"
	"testing"
)

func isCoinbase(tx *Tx) bool {
	return len(tx.TxInputs) == 1 && len(tx.TxInputs[0].ID) == 0 && tx.TxInputs[0].OutIdx == -1
}

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}

func TestCreateGenesisBlock(t *testing.T) {
	chain := InitBlockchain("", 3000)
	genesis := CreateGenesisBlock(chain.Address)
	hash := fmt.Sprintf("%x", []byte{})
	if genesis.PrevHash != hash {
		t.Error("Hashes do not equal")
	}
	if !isCoinbase(genesis.Transactions[0]) || len(genesis.Transactions) != 1 {
		t.Error("Genesis block should contain only 1 coinbase transaction")
	}
}

func TestBlockchain_Hashes(t *testing.T) {
	chain := InitBlockchain("", 3000)
	chain.MineBlock()
	hash := chain.Blocks[0].Hash()
	if hash != chain.Blocks[1].PrevHash {
		t.Error("Hashes should be equal")
	}

}

func TestInitBlockchain(t *testing.T) {
	chain := InitBlockchain("", 3000)
	if len(chain.Blocks) < 1 || len(chain.Blocks) > 1 {
		t.Error("Chain should have only 1 block on initiation")
	}
}

func TestBlockchain_AddBlock(t *testing.T) {
	chain := InitBlockchain("", 3000)
	chain.AddBlock()
	if len(chain.Blocks) <= 1 {
		t.Error("Chain should have more than 1 block")
	}
}

func TestBlockchain_FindUnspentTxs(t *testing.T) {
	walletA := wallet.NewWallet()
	chain := InitBlockchain(walletA.BlockchainAddress(), 3000)
	chain.MineBlock()
	utxo := chain.FindUnspentTxs(walletA.BlockchainAddress())
	if len(utxo) != 2 {
		t.Error("Wallet A should have 2 unspent transactions (from genesis block and new block reward)")
	}
}

func TestBlockchain_FindUTXO(t *testing.T) {
	walletA := wallet.NewWallet()
	walletB := wallet.NewWallet()
	chain := InitBlockchain(walletA.BlockchainAddress(), 3000)
	chain.MineBlock()
	utxo := chain.FindUTXO(walletA.BlockchainAddress())
	if len(utxo) != 2 {
		t.Error("Wallet A should have 2 unspent transactions (from genesis block and new block reward)")
	}
	utxo = chain.FindUTXO(walletB.BlockchainAddress())
	if len(utxo) != 0 {
		t.Error("Wallet B should have 0 unspent transactions")
	}
}

func TestBlockchain_FindSpendableOutputs1(t *testing.T) {
	walletA := wallet.NewWallet()
	walletB := wallet.NewWallet()
	chain := InitBlockchain(walletA.BlockchainAddress(), 3000)

	// WalletA has 0.1 coins
	// WalletA -> 0.1 coins -> WalletB
	// WalletA mines block gets 0.1 coins
	// WalletA has 0.1 coins
	// WalletB has 0.1 coins

	tx := NewTransaction(walletA.BlockchainAddress(), walletB.BlockchainAddress(), 0.1, chain)
	chain.AddTransaction(walletA.BlockchainAddress(), walletB.BlockchainAddress(), 0.1, tx.GenerateSignature(walletA.PrivateKey()), walletA.PublicKey())
	chain.MineBlock()

	//amount5, _ := chain.FindSpendableOutputs(walletA.BlockchainAddress(), 0.3)
	amount6, _ := chain.FindSpendableOutputs(walletB.BlockchainAddress(), 0.3)
	if amount6 != 0.1 {
		t.Error("Amount in B should be 0.1 but is", amount6)
	}
	//if amount5 != 0.1 {
	//	t.Error("Amount in A should be 0.1 but is", amount5)
	//}

}

func TestBlockchain_FindSpendableOutputs2(t *testing.T) {
	walletA := wallet.NewWallet()
	chain := InitBlockchain(walletA.BlockchainAddress(), 3000)
	chain.MineBlock()
	amount1, _ := chain.FindSpendableOutputs(walletA.BlockchainAddress(), 0.1)
	if amount1 != 0.1 {
		t.Error("Amount should be 0.1")
	}
	amount2, _ := chain.FindSpendableOutputs(walletA.BlockchainAddress(), 0.2)
	if amount2 != 0.2 {
		t.Error("Amount should be 0.2")
	}
}

func TestBlockchain_FindSpendableOutputs3(t *testing.T) {
	walletA := wallet.NewWallet()
	chain := InitBlockchain(walletA.BlockchainAddress(), 3000)
	chain.MineBlock()

	amount3, _ := chain.FindSpendableOutputs(walletA.BlockchainAddress(), 0.3)
	if amount3 == 0.3 {
		t.Error("Amount should not be 0.3 because there is only 0.2 unspent")
	}

	walletB := wallet.NewWallet()
	amount4, _ := chain.FindSpendableOutputs(walletB.BlockchainAddress(), 0.1)
	if amount4 != 0 {
		t.Error("Amount should be 0")
	}
}

func TestBlockchain_AddTransaction(t *testing.T) {
	walletA := wallet.NewWallet()
	walletB := wallet.NewWallet()
	chain := InitBlockchain(walletA.BlockchainAddress(), 3000)
	tx := NewTransaction(walletA.BlockchainAddress(), walletB.BlockchainAddress(), 0.1, chain)
	isAdd := chain.AddTransaction(walletA.BlockchainAddress(), walletB.BlockchainAddress(), 0.1, tx.GenerateSignature(walletA.PrivateKey()), walletA.PublicKey())
	chain.MineBlock()

	if !isAdd {
		t.Error("Tx should be added to tx pool")
	}

	assertPanic(t, func() {
		tx2 := NewTransaction(walletA.BlockchainAddress(), walletB.BlockchainAddress(), 0.3, chain)
		chain.AddTransaction(walletA.BlockchainAddress(), walletB.BlockchainAddress(), 0.3, tx2.GenerateSignature(walletA.PrivateKey()), walletA.PublicKey())
	})

}

func TestBlockchain_ClearPool(t *testing.T) {
	walletA := wallet.NewWallet()
	walletB := wallet.NewWallet()
	chain := InitBlockchain(walletA.BlockchainAddress(), 3000)
	chain.MineBlock()
	tx := NewTransaction(walletA.BlockchainAddress(), walletB.BlockchainAddress(), 0.1, chain)
	isAdd := chain.AddTransaction(walletA.BlockchainAddress(), walletB.BlockchainAddress(), 0.1, tx.GenerateSignature(walletA.PrivateKey()), walletA.PublicKey())

	if !isAdd {
		t.Error("Tx should be added to tx pool")
	}

	chain.AddBlock()
	if len(chain.Blocks) <= 1 {
		t.Error("Chain should have more than 1 block")
	}
	if len(chain.TxPool) >= 1 {
		t.Error("Transaction pool should be empty")
	}
}

func TestBlockchain_Mining(t *testing.T) {
	miner := wallet.NewWallet()
	chain := InitBlockchain(miner.BlockchainAddress(), 3000)
	chain.MineBlock()
	lastBlock := chain.lastBlock()

	if len(chain.Blocks) <= 1 {
		t.Error("Chain should have more than 1 block")
	}
	if len(chain.TxPool) >= 1 {
		t.Error("Transaction pool should be empty")
	}
	if len(lastBlock.Transactions) != 1 {
		t.Error("Block should have 1 transaction in it")
	}
	if isCoinbase(lastBlock.Transactions[0]) == false {
		t.Error("Should be a coinbase transaction")
	}

	walletB := wallet.NewWallet()
	tx := NewTransaction(miner.BlockchainAddress(), walletB.BlockchainAddress(), 0.1, chain)
	chain.AddTransaction(miner.BlockchainAddress(), walletB.BlockchainAddress(), 0.1, tx.GenerateSignature(miner.PrivateKey()), miner.PublicKey())
	chain.MineBlock()

	if len(chain.Blocks) <= 2 {
		t.Error("Chain should have more than 2 block")
	}
	lastBlock = chain.lastBlock()
	if len(lastBlock.Transactions) != 2 {
		t.Error("Block should have 2 transaction in it")
	}
	if isCoinbase(lastBlock.Transactions[0]) == false {
		t.Error("Should be a coinbase transaction")
	}
}

func TestBlockchain_VerifyTxSignature(t *testing.T) {
	walletA := wallet.NewWallet()
	walletB := wallet.NewWallet()
	chain := InitBlockchain(walletA.BlockchainAddress(), 3000)
	tx := NewTransaction(walletA.BlockchainAddress(), walletB.BlockchainAddress(), 0.1, chain)
	signature := tx.GenerateSignature(walletA.PrivateKey())
	isValid := chain.VerifyTxSignature(walletA.PublicKey(), signature, tx)

	if !isValid {
		t.Error("Transaction signature should be valid")
	}

	isValid = chain.VerifyTxSignature(walletB.PublicKey(), signature, tx)

	if isValid {
		t.Error("Transaction signature should not be valid")
	}
}
