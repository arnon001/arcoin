package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/ripemd160"
)

const (
	ChecksumLength = 4
	Version        = byte(0x00)
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func NewWallet() *Wallet {
	private, public := newKeyPair()
	return &Wallet{private, public}
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, _ := ecdsa.GenerateKey(curve, rand.Reader)
	public := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, public
}

func (w Wallet) Address() string {
	pubKeyHash := HashPubKey(w.PublicKey)
	versionedPayload := append([]byte{Version}, pubKeyHash...)
	checksum := checksum(versionedPayload)
	fullPayload := append(versionedPayload, checksum...)
	return Base58Encode(fullPayload)
}

func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)
	RIPEMD160Hasher := ripemd160.New()
	RIPEMD160Hasher.Write(publicSHA256[:])
	return RIPEMD160Hasher.Sum(nil)
}

func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])
	return secondSHA[:ChecksumLength]
}

// Signature and Verification
func SignTransaction(privateKey *ecdsa.PrivateKey, txHash []byte) ([]byte, error) {
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, txHash)
	if err != nil {
		return nil, err
	}
	signature := append(r.Bytes(), s.Bytes()...)
	return signature, nil
}

func VerifySignature(pubKey, txHash, signature []byte) bool {
	r := big.Int{}
	s := big.Int{}
	sigLen := len(signature)
	r.SetBytes(signature[:sigLen/2])
	s.SetBytes(signature[sigLen/2:])

	x := big.Int{}
	y := big.Int{}
	keyLen := len(pubKey)
	x.SetBytes(pubKey[:keyLen/2])
	y.SetBytes(pubKey[keyLen/2:])

	curve := elliptic.P256()
	publicKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
	return ecdsa.Verify(&publicKey, txHash, &r, &s)
}