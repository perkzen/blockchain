package blockchain

import (
	"blockchain/pkg/wallet"
	"math"
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
	genesis := chain.CreateGenesisBlock()
	hash := "GENESIS"
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
	if len(chain.Blocks) != 1 {
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

func TestBlockchain_Balance(t *testing.T) {
	walletA := wallet.NewWallet()
	chain := InitBlockchain(walletA.BlockchainAddress(), 3000)
	balanceA, _ := chain.UTXO.FindSpendableOutputs(walletA.BlockchainAddress(), math.MaxInt)
	if balanceA != COINBASE_REWARD {
		t.Error("Balance should be equal to COINBASE_REWARD but got", balanceA)
	}

	walletB := wallet.NewWallet()
	balanceB, _ := chain.UTXO.FindSpendableOutputs(walletB.BlockchainAddress(), math.MaxInt)
	if balanceB != 0 {
		t.Error("Balance should be equal to 0 but got", balanceB)
	}

	chain.MineBlock()
	balanceA, _ = chain.UTXO.FindSpendableOutputs(walletA.BlockchainAddress(), math.MaxInt)
	if balanceA != COINBASE_REWARD*2 {
		t.Error("Balance should be equal to COINBASE_REWARD*2 but got", balanceA)
	}
}

func TestBlockchain_AddTransaction(t *testing.T) {
	walletA := wallet.NewWallet()
	walletB := wallet.NewWallet()
	chain := InitBlockchain(walletA.BlockchainAddress(), 3000)

	balanceA, outsA1 := chain.UTXO.FindSpendableOutputs(walletA.BlockchainAddress(), math.MaxInt)
	if balanceA != 0.1 {
		t.Error("Balance should be equal to COINBASE_REWARD*2-0.1 but got", balanceA)
	}
	if len(outsA1) != 1 {
		t.Error("Outs should be equal to 1 but got", len(outsA1))
	}

	isAdd := chain.AddTransaction(walletA.BlockchainAddress(), walletB.BlockchainAddress(), 0.1, walletA.PrivateKey(), walletA.PublicKey())

	balanceA, outsA2 := chain.UTXO.FindSpendableOutputs(walletA.BlockchainAddress(), math.MaxInt)
	if balanceA != 0 {
		t.Error("Balance should be equal to COINBASE_REWARD*2-0.1 but got", balanceA)
	}
	if len(outsA2) != 0 {
		t.Error("Outs should be equal to 0 but got", len(outsA2))
	}

	if chain.UTXO.IsSpent[outsA1[0].Hash()] == false {
		t.Error("Output should be spent")
	}

	if !isAdd {
		t.Error("Tx should be added to tx pool")
	}

	balanceB, outsB := chain.UTXO.FindSpendableOutputs(walletB.BlockchainAddress(), math.MaxInt)
	if balanceB != 0.1 {
		t.Error("Balance should be equal to 0.1 but got", balanceB)
	}
	if len(outsB) != 1 {
		t.Error("Outs should be equal to 1 but got", len(outsB))
	}

	assertPanic(t, func() {
		chain.AddTransaction(walletA.BlockchainAddress(), walletB.BlockchainAddress(), 0.3, walletA.PrivateKey(), walletA.PublicKey())
	})

	chain.MineBlock()
	balanceA, _ = chain.UTXO.FindSpendableOutputs(walletA.BlockchainAddress(), math.MaxInt)
	if balanceA != COINBASE_REWARD*2-0.1 {
		t.Error("Balance should be equal to COINBASE_REWARD but got", balanceA)
	}
}

func TestBlockchain_ClearPool(t *testing.T) {
	walletA := wallet.NewWallet()
	walletB := wallet.NewWallet()
	chain := InitBlockchain(walletA.BlockchainAddress(), 3000)
	chain.MineBlock()
	isAdd := chain.AddTransaction(walletA.BlockchainAddress(), walletB.BlockchainAddress(), COINBASE_REWARD, walletA.PrivateKey(), walletA.PublicKey())

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
	chain.AddTransaction(miner.BlockchainAddress(), walletB.BlockchainAddress(), COINBASE_REWARD, miner.PrivateKey(), miner.PublicKey())
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

func TestBlockchain_UTXO(t *testing.T) {
	walletA := wallet.NewWallet()
	chain := InitBlockchain(walletA.BlockchainAddress(), 3000)
	isAdded := chain.AddTransaction(walletA.BlockchainAddress(), walletA.BlockchainAddress(), 0.1, walletA.PrivateKey(), walletA.PublicKey())

	if !isAdded {
		t.Error("Transaction should be added to tx pool")
	}

	//chain.MineBlock()

	amount, _ := chain.UTXO.FindSpendableOutputs(walletA.BlockchainAddress(), math.MaxInt)

	if amount != COINBASE_REWARD {
		t.Error("Balance should be equal to COINBASE_REWARD but got", amount)
	}

	if len(chain.UTXO.Outputs) != 2 {
		t.Error("UTXO should have 2 output ")
	}

	chain.AddBlock()
	if len(chain.UTXO.Outputs) != 3 {
		t.Error("UTXO should have 3 output ")
	}

}

func TestBlockchain_AddBlock2(t *testing.T) {
	walletA := wallet.NewWallet()
	chain := InitBlockchain(walletA.BlockchainAddress(), 3000)

	block := chain.AddBlock()

	pow := NewProofOfWork(block)
	valid := pow.IsValid(block.Nonce)

	if !valid {
		t.Error("Block should be valid")
	}

}
