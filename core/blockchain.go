package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"
)

type Block struct {
	Header       BlockHeader
	Transactions []*Transaction
}

type BlockHeader struct {
	Version      int32
	PrevHash     []byte
	MerkleRoot   []byte
	Timestamp    int64
	Bits         uint32
	Nonce        uint32
	Height       int32
}

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

type TXInput struct {
	Txid      []byte
	Vout      int
	Signature []byte
	PubKey    []byte
}

type TXOutput struct {
	Value      int
	PubKeyHash []byte
}

type Blockchain struct {
	Blocks     []*Block
	UTXO       UTXOSet
	Difficulty int
	Mempool    []*Transaction
}

type UTXOSet struct {
	Data map[string][]TXOutput
}

// ProofOfWork represents the mining context
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

const (
	MaxBlockSize    = 1_000_000 // 1MB blocks
	Subsidy         = 50        // Mining reward
	TargetBits      = 24        // Mining difficulty
	CoinbaseData    = "ArCoin Genesis"
	MiningInterval  = 10 * time.Minute
)

func NewBlockchain() *Blockchain {
	genesis := createGenesisBlock()
	utxo := UTXOSet{
		Data: make(map[string][]TXOutput),
	}
	return &Blockchain{
		Blocks:     []*Block{genesis},
		UTXO:       utxo,
		Difficulty: TargetBits,
	}
}

func createGenesisBlock() *Block {
	coinbase := NewCoinbaseTX("", CoinbaseData)
	return &Block{
		Header: BlockHeader{
			Version:    1,
			PrevHash:   []byte{},
			MerkleRoot: []byte{},
			Timestamp:  time.Now().Unix(),
			Bits:       TargetBits,
			Height:     0,
		},
		Transactions: []*Transaction{coinbase},
	}
}

func (bc *Blockchain) MineBlock(transactions []*Transaction) (*Block, error) {
	lastHash := bc.Blocks[len(bc.Blocks)-1].Header.PrevHash
	newBlock := &Block{
		Header: BlockHeader{
			Version:    1,
			PrevHash:   lastHash,
			Timestamp:  time.Now().Unix(),
			Bits:       TargetBits,
			Height:     bc.Blocks[len(bc.Blocks)-1].Header.Height + 1,
		},
		Transactions: transactions,
	}

	pow := NewProofOfWork(newBlock)
	nonce, hash := pow.Run()

	newBlock.Header.Nonce = nonce
	newBlock.Header.MerkleRoot = hash

	if !pow.Validate() {
		return nil, errors.New("proof of work validation failed")
	}

	bc.Blocks = append(bc.Blocks, newBlock)
	bc.UpdateUTXO(newBlock)
	return newBlock, nil
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-TargetBits))
	return &ProofOfWork{b, target}
}

func (pow *ProofOfWork) Run() (uint32, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := uint32(0)

	for nonce < MaxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		}
		nonce++
	}
	return nonce, hash[:]
}

func (pow *ProofOfWork) prepareData(nonce uint32) []byte {
	data := bytes.Join([][]byte{
		pow.block.Header.PrevHash,
		pow.block.Header.MerkleRoot,
		IntToHex(pow.block.Header.Timestamp),
		IntToHex(int64(TargetBits)),
		IntToHex(int64(nonce)),
	}, []byte{})
	return data
}

// UTXO Management
func (utxo *UTXOSet) Update(block *Block) {
	for _, tx := range block.Transactions {
		for _, in := range tx.Vin {
			if !in.IsCoinbase() {
				utxo.RemoveOutput(string(in.Txid), in.Vout)
			}
		}
		for i, out := range tx.Vout {
			utxo.AddOutput(string(tx.ID), i, out)
		}
	}
}

// Transaction Handling
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to %s", to)
	}
	
	txin := TXInput{
		Txid:      []byte{},
		Vout:      -1,
		Signature: []byte(data),
	}
	
	txout := NewTXOutput(Subsidy, to)
	
	tx := Transaction{
		ID:   nil,
		Vin:  []TXInput{txin},
		Vout: []TXOutput{*txout},
	}
	tx.ID = tx.Hash()
	return &tx
}

func (tx *Transaction) Hash() []byte {
	var hash [32]byte
	txCopy := *tx
	txCopy.ID = []byte{}
	hash = sha256.Sum256(txCopy.Serialize())
	return hash[:]
}

func (tx *Transaction) Serialize() []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(tx)
	if err != nil {
		panic(err)
	}
	return buff.Bytes()
}