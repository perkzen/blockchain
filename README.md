# GoBlockchain

### Description

**GoBlockchain** is a basic implementation of a blockchain using the Go programming language. The project aims to demonstrate the fundamental concepts behind blockchain technology and how they can be implemented in a decentralized and secure manner.

### Features

- **Blockchain Structure**: GoBlockchain utilizes a chain of blocks to store and organize transactional data. Each block contains a list of transactions, a timestamp, and a unique identifier called a hash. The blocks are linked together using cryptographic hashes, ensuring the integrity and immutability of the data.


- **UTXO Model**: GoBlockchain implements the UTXO (Unspent Transaction Output) model, which is widely used in blockchain systems. Instead of tracking account balances, transactions in the blockchain are represented as inputs and outputs. Each output is an unspent transaction that can be used as an input in subsequent transactions. The UTXO model provides an efficient way to verify transaction validity and ensure that funds are not double-spent.

- **Proof-of-Work Consensus**: The project implements a proof-of-work (PoW) consensus mechanism, similar to that used by Bitcoin. Miners in the network compete to solve computationally intensive puzzles to validate and append new blocks to the blockchain. This consensus algorithm ensures the security and trustworthiness of the blockchain.


- **Transaction Validation**: GoBlockchain validates the integrity of transactions using cryptographic techniques. Each transaction is digitally signed, ensuring its authenticity and preventing tampering. Transaction validation is an essential step in maintaining the security and validity of the blockchain.


- **Peer-to-Peer Networking**: The project incorporates a peer-to-peer network architecture, enabling nodes to communicate and share information. Nodes can join the network, propagate  blocks, and synchronize their local copy of the blockchain. The peer-to-peer networking layer facilitates a decentralized and distributed environment.


- **Command-Line Interface (CLI)**: GoBlockchain provides a simple command-line interface to interact with the blockchain.

[//]: # (Users can create wallets, send transactions, mine blocks, and view the blockchain's current state. )



### CLI Tool

```bash
go run cmd/main.go start-node -port 3000
go run cmd/main.go send -to 0xa123b267c -amount 0.1 -node 3000
go run cmd/main.go balance -node 3000 
```

### Server Endpoints
 `GET` http://localhost/

- Retrieves blockchain 

`GET` http://localhost/nodes

- Retries all connected nodes

`GET` http://localhost/wallet

- Retries wallet address and balance

`POST` http://localhost/transacions
- Create adds new transactions to blockchains the transaction pool



### Add in the future (maybe)
- [ ] Persistence
- [ ] Miner fees
- [ ] Max transaction in block (100 txs)
- [ ] Max number of coins
- [ ] Difficulty adjustment
- [ ] Support for multiple wallets
- [ ] add CLI commands to view chain

