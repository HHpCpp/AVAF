package blockchain

import (
	"fmt"
)

// Blockchain представляет собой цепочку блоков
type Blockchain struct {
	Chain    []Block
	Accounts map[string]Account // Хранилище аккаунтов
}

// NewBlockchain создает новый блокчейн с genesis блоком
func NewBlockchain() Blockchain {
	genesisBlock := NewBlock(0, []Transaction{}, "")
	accounts := make(map[string]Account)
	return Blockchain{
		Chain:    []Block{genesisBlock},
		Accounts: accounts,
	}
}

// AddBlock добавляет новый блок в блокчейн
func (bc *Blockchain) AddBlock(transactions []Transaction) {
	// Проверяем транзакции перед добавлением
	for _, tx := range transactions {
		if !bc.ValidateTransaction(tx) {
			fmt.Println("Invalid transaction:", tx)
			return
		}
	}

	// Обновляем балансы аккаунтов
	for _, tx := range transactions {
		sender := bc.Accounts[tx.Sender]
		recipient := bc.Accounts[tx.Recipient]

		sender.Balance -= tx.Amount
		recipient.Balance += tx.Amount

		bc.Accounts[tx.Sender] = sender
		bc.Accounts[tx.Recipient] = recipient
	}

	prevBlock := bc.Chain[len(bc.Chain)-1]
	newBlock := NewBlock(prevBlock.Index+1, transactions, prevBlock.Hash)
	bc.Chain = append(bc.Chain, newBlock)
}

// ValidateTransaction проверяет, что транзакция корректна
func (bc *Blockchain) ValidateTransaction(tx Transaction) bool {
	sender, exists := bc.Accounts[tx.Sender]
	if !exists {
		fmt.Println("Sender account does not exist:", tx.Sender)
		return false
	}

	if sender.Balance < tx.Amount {
		fmt.Println("Insufficient balance:", tx.Sender)
		return false
	}

	return true
}

// CreateAccount создает новый аккаунт
func (bc *Blockchain) CreateAccount(address string, balance float64) {
	if _, exists := bc.Accounts[address]; exists {
		fmt.Println("Account already exists:", address)
		return
	}
	bc.Accounts[address] = NewAccount(address, balance)
}

// IsValid проверяет валидность блокчейна
func (bc *Blockchain) IsValid() bool {
	for i := 1; i < len(bc.Chain); i++ {
		currentBlock := bc.Chain[i]
		prevBlock := bc.Chain[i-1]

		if currentBlock.Hash != currentBlock.CalculateHash() {
			return false
		}

		if currentBlock.PrevHash != prevBlock.Hash {
			return false
		}
	}
	return true
}
