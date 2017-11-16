package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

type Blockchain struct {
	chain             []*Block
	curr_transactions []*TranX
}

type Block struct {
	Index        int      `json:"index"`
	Timestamp    int64    `json:"timestamp"`
	Transactions []*TranX `json:"transactions"`
	Proof        int      `json:"proof"`
	PreviousHash string   `json:"previous_hash"`
}

type TranX struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Amount    int    `json:"amount"`
}

func NewBlockchain() *Blockchain {
	bc := &Blockchain{chain: []*Block{}, curr_transactions: []*TranX{}}
	bc.chain = append(bc.chain, bc.NewBlock(100, "1"))
	return bc
}

func (bc *Blockchain) NewBlock(proof int, prev ...string) *Block {
	var block *Block
	if len(prev) > 0 {
		block = &Block{
			Index:        len(bc.chain) + 1,
			Timestamp:    time.Now().Unix(),
			Transactions: bc.curr_transactions,
			Proof:        proof,
			PreviousHash: prev[0],
		}
		bc.chain = append(bc.chain, block)
		return block
	}

	block = &Block{
		Index:        len(bc.chain) + 1,
		Timestamp:    time.Now().Unix(),
		Transactions: bc.curr_transactions,
		Proof:        proof,
		PreviousHash: bc.Hash(bc.LastBlock()),
	}
	bc.chain = append(bc.chain, block)
	return block
}

func (bc *Blockchain) NewTranX(sender, recipient string, amount int) int {
	bc.curr_transactions = append(bc.curr_transactions,
		&TranX{
			Sender:    sender,
			Recipient: recipient,
			Amount:    amount,
		})
	return bc.LastBlock().Index + 1
}

func (bc *Blockchain) Chain() []*Block {
	return bc.chain
}

func (bc *Blockchain) Hash(b *Block) string {
	block_bytes, _ := json.Marshal(b)
	arr := sha256.Sum256(block_bytes)
	return fmt.Sprintf("%x", arr)
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) ProofOfWork(last int) int {
	curr := 0
	for !bc.ValidProof(last, curr) {
		curr += 1
	}
	return curr
}

func (bc *Blockchain) ValidProof(last, curr int) bool {
	buf := new(bytes.Buffer)
	io.WriteString(buf, fmt.Sprintf("%v%v", last, curr))
	arr := sha256.Sum256(buf.Bytes())
	return strings.Contains(string(arr[:]), "0")
}
