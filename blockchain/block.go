package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// Block представляет собой блок в блокчейне
type Block struct {
	Index        int
	Timestamp    string
	Transactions []Transaction // Список транзакций в блоке
	PrevHash     string
	Hash         string
}

// NewBlock создает новый блок
func NewBlock(index int, transactions []Transaction, prevHash string) Block {
	block := Block{
		Index:        index,
		Timestamp:    time.Now().String(),
		Transactions: transactions,
		PrevHash:     prevHash,
	}
	block.Hash = block.CalculateHash()
	return block
}

// CalculateHash вычисляет хеш блока
func (b *Block) CalculateHash() string {
	record := fmt.Sprintf("%d%s%v%s", b.Index, b.Timestamp, b.Transactions, b.PrevHash)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}
