# GO Blockchain

### CLI

```bash
go run cmd/main.go start-node -port 3000
go run cmd/main.go send -to 0xa123b267c -amount 0.1 -node 3000
go run cmd/main.go balance -node 3000 
```

### Add in the future (maybe)
- [ ] Persistence
- [ ] Miner fees
- [ ] Max transaction in block (100 txs)
- [ ] Max number of coins
- [ ] Difficulty adjustment

### DONE
- [X] Proof of work
- [x] Wallet
- [x] Transactions
- [X] Peer to peer network
- [X] Longest chain rule 
- [X] UTXO model

