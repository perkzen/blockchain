package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
	"log"
)

type Wallet struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	address    string
}

func NewWallet() *Wallet {
	// 1. Create ECDSA private  and public key
	private, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	w := &Wallet{
		privateKey: private,
		publicKey:  &private.PublicKey,
	}

	// 2. Perform SHA-256 hashing on the public key
	h2 := sha256.New()
	h2.Write(w.publicKey.X.Bytes())
	h2.Write(w.publicKey.Y.Bytes())
	digest2 := h2.Sum(nil)

	// 3. Perform RIPEMD-160 hashing on the result of SHA-256
	h3 := ripemd160.New()
	h3.Write(digest2)
	digest3 := h3.Sum(nil)

	// 4. Add version byte in front of RIPEMD-160 hash (0x00 for Main Network)
	vd4 := make([]byte, 21)
	vd4[0] = 0x00
	copy(vd4[1:], digest3[:])

	// 5. Perform SHA-256 hash on the extended RPIEMD-160 result
	h5 := sha256.New()
	h5.Write(vd4)
	digest5 := h5.Sum(nil)

	// 6. Perform SHA-256 hash on result of the previous SHA-256 hash.
	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)

	// 7. Take the first 4 bytes of the second SHA-256 hash for checksum
	checksum := digest6[:4]

	// 8. Add the 4 checksum bytes from 7 at the end of extended RIPEMD-160 hash from 4
	dc8 := make([]byte, 25)
	copy(dc8[:21], vd4[:])
	copy(dc8[21:], checksum[:])

	// 9. Convert the result from a byte string into base58
	address := base58.Encode(dc8)
	w.address = address

	return w
}

func (w *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}

func (w *Wallet) PrivateKeyStr() string {
	return fmt.Sprintf("%x", w.privateKey.D.Bytes())
}

func (w *Wallet) PublicKey() *ecdsa.PublicKey {
	return w.publicKey
}

func (w *Wallet) PublicKeyStr() string {
	return fmt.Sprintf("%x%x", w.publicKey.X.Bytes(), w.publicKey.Y.Bytes())
}

func (w *Wallet) BlockchainAddress() string {
	return w.address
}

func (w *Wallet) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PrivateKey *ecdsa.PrivateKey `json:"private_key"`
		PublicKey  *ecdsa.PublicKey  `json:"public_key"`
		Address    string            `json:"blockchain_address"`
	}{
		PrivateKey: w.privateKey,
		PublicKey:  w.publicKey,
		Address:    w.BlockchainAddress(),
	})
}

// Doesn't work in GO 1.20 or more
//func (w *Wallet) Save() {
//	type WalletEncode struct {
//		PrivateKey ecdsa.PrivateKey
//		PublicKey  ecdsa.PublicKey
//		Address    string
//	}
//
//	var content bytes.Buffer
//
//	gob.Register(elliptic.P256())
//	encoder := gob.NewEncoder(&content)
//
//	err := encoder.Encode(WalletEncode{
//		PrivateKey: *w.privateKey,
//		PublicKey:  *w.publicKey,
//		Address:    w.BlockchainAddress(),
//	})
//
//	if err != nil {
//		log.Panic(err)
//	}
//	err = os.WriteFile("data/wallet.dat", content.Bytes(), 0644)
//	if err != nil {
//		log.Panic(err)
//	}
//
//}
//
//func Load() *Wallet {
//	var wallet struct {
//		PrivateKey ecdsa.PrivateKey
//		PublicKey  ecdsa.PublicKey
//		Address    string
//	}
//
//	content, err := os.ReadFile("data/wallet.dat")
//	if err != nil {
//		log.Panic(err)
//	}
//
//	gob.Register(elliptic.P256())
//	decoder := gob.NewDecoder(bytes.NewReader(content))
//	err = decoder.Decode(&wallet)
//	if err != nil {
//		log.Panic(err)
//	}
//
//	return &Wallet{
//		privateKey: &wallet.PrivateKey,
//		publicKey:  &wallet.PublicKey,
//		address:    wallet.Address,
//	}
//}
