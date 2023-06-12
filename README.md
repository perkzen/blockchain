# GO Blockchain

### CLI

```bash
go run cmd/main.go start-node -port 3000
go run cmd/main.go send -to 0xa123b267c -amount 0.1 -node 3000
go run cmd/main.go balance -node 3000 
```


### DONE

- [X] Proof of work
- [x] Wallet
- [x] Transactions
- [X] Peer to peer network
- [X] Longest chain rule

### Add in the future (maybe)
- [ ] Use sockets for peer to peer network
- [ ] Merkle trees
- [ ] Persistence
- [ ] Miner fees
- [ ] Max transaction in block
- [ ] remove node on disconnect
- [ ] save wallet to file