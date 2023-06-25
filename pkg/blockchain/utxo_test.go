package blockchain

import "testing"

func TestUTXO_AddOutput(t *testing.T) {
	utxo := NewUTXO()
	output := NewOutput(TxOutput{Value: 10.0, PublicKey: "public_key"}, 0, "txID", 0)
	utxo.AddOutput(output)
	if len(utxo.Outputs) != 1 {
		t.Errorf("Expected 1 output, got %d", len(utxo.Outputs))
	}
}

func TestUTXO_SpendOutput(t *testing.T) {
	utxo := NewUTXO()
	output := NewOutput(TxOutput{Value: 10.0, PublicKey: "public_key"}, 0, "txID", 0)
	utxo.AddOutput(output)
	utxo.SpendOutput(output.Hash())
	if !utxo.IsSpent[output.Hash()] {
		t.Errorf("Expected hash to be spent")
	}
}

func TestUTXO_FindSpendableOutputs(t *testing.T) {
	utxo := NewUTXO()
	output := NewOutput(TxOutput{Value: 10.0, PublicKey: "public_key"}, 0, "txID", 0)
	utxo.AddOutput(output)
	acc, unspentOut := utxo.FindSpendableOutputs("public_key", 10)
	if acc != 10.0 {
		t.Errorf("Expected 10.0, got %f", acc)
	}
	if len(unspentOut) != 1 {
		t.Errorf("Expected 1 output, got %d", len(unspentOut))
	}

	utxo.SpendOutput(output.Hash())
	acc, unspentOut = utxo.FindSpendableOutputs("public_key", 10)
	if acc != 0.0 {
		t.Errorf("Expected 0.0, got %f", acc)
	}
	if len(unspentOut) != 0 {
		t.Errorf("Expected 0 output, got %d", len(unspentOut))
	}
}
