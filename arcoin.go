package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Block represents a block in the ArCoin blockchain
type Block struct {
	Index     int
	Timestamp string
	Data      string
	PrevHash  string
	Hash      string
	Nonce     int
}

// Blockchain represents the ArCoin blockchain
type Blockchain struct {
	Blocks     []*Block
	Difficulty int
}

// CalculateHash calculates the hash of a block
func (b *Block) CalculateHash() string {
	record := strconv.Itoa(b.Index) + b.Timestamp + b.Data + b.PrevHash + strconv.Itoa(b.Nonce)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// CreateGenesisBlock creates the first block in the blockchain
func CreateGenesisBlock(difficulty int) *Block {
	genesisBlock := &Block{
		Index:     0,
		Timestamp: time.Now().String(),
		Data:      "ArCoin Genesis Block",
		PrevHash:  "",
		Nonce:     0,
	}
	genesisBlock.Hash = genesisBlock.CalculateHash()
	return genesisBlock
}

// NewBlockchain creates a new ArCoin blockchain
func NewBlockchain(difficulty int) *Blockchain {
	genesisBlock := CreateGenesisBlock(difficulty)
	return &Blockchain{
		Blocks:     []*Block{genesisBlock},
		Difficulty: difficulty,
	}
}

// ProofOfWork performs the mining process
func (b *Block) ProofOfWork(difficulty int) {
	for {
		hash := b.CalculateHash()
		if hash[:difficulty] == "0000" {
			b.Hash = hash
			return
		}
		b.Nonce++
	}
}

// AddBlock adds a new block to the blockchain
func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := &Block{
		Index:     prevBlock.Index + 1,
		Timestamp: time.Now().String(),
		Data:      data,
		PrevHash:  prevBlock.Hash,
		Nonce:     0,
	}

	newBlock.ProofOfWork(bc.Difficulty)
	bc.Blocks = append(bc.Blocks, newBlock)
}

// IsValid checks if the blockchain is valid
func (bc *Blockchain) IsValid() bool {
	for i := 1; i < len(bc.Blocks); i++ {
		currentBlock := bc.Blocks[i]
		prevBlock := bc.Blocks[i-1]

		if currentBlock.Hash != currentBlock.CalculateHash() {
			return false
		}

		if currentBlock.PrevHash != prevBlock.Hash {
			return false
		}
	}
	return true
}

// PrintBlockchain prints the blockchain in readable format
func (bc *Blockchain) PrintBlockchain() {
	for _, block := range bc.Blocks {
		blockJSON, _ := json.MarshalIndent(block, "", "  ")
		fmt.Println(string(blockJSON))
	}
}

func main() {
	// Initialize blockchain with difficulty 2 (number of leading zeros required)
	difficulty := 2
	arnCoin := NewBlockchain(difficulty)

	// Add some blocks
	arnCoin.AddBlock("First ARN Transaction")
	arnCoin.AddBlock("Second ARN Transaction")

	// Print blockchain
	arnCoin.PrintBlockchain()

	// Validate blockchain
	fmt.Println("\nBlockchain valid:", arnCoin.IsValid())

	// Try to tamper with blockchain
	arnCoin.Blocks[1].Data = "Tampered transaction"
	fmt.Println("After tampering - Blockchain valid:", arnCoin.IsValid())
}