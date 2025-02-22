package blockchain

import (
	"fmt"

	AA "github.com/HHpCpp/AVAF/accounts"
)

// Blockchain представляет собой цепочку блоков
type Blockchain struct {
	Chain          []Block
	AccountManager *AA.AccountManager // Используем AccountManager для управления аккаунтами
}

// NewBlockchain создает новый блокчейн с genesis блоком
func NewBlockchain() Blockchain {
	genesisBlock := NewBlock(0, []Transaction{}, "")
	accountManager := AA.NewAccountManager()
	return Blockchain{
		Chain:          []Block{genesisBlock},
		AccountManager: accountManager,
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
		sender, err := bc.AccountManager.GetAccount(tx.Sender)
		if err != nil {
			fmt.Println("Error getting sender account:", err)
			return
		}

		recipient, err := bc.AccountManager.GetAccount(tx.Recipient)
		if err != nil {
			fmt.Println("Error getting recipient account:", err)
			return
		}

		sender.Balance -= tx.Amount
		recipient.Balance += tx.Amount

		bc.AccountManager.CreateAccount(sender.Address, sender.Balance)
		bc.AccountManager.CreateAccount(recipient.Address, recipient.Balance)
	}

	prevBlock := bc.Chain[len(bc.Chain)-1]
	newBlock := NewBlock(prevBlock.Index+1, transactions, prevBlock.Hash)
	bc.Chain = append(bc.Chain, newBlock)
}

// ValidateTransaction проверяет, что транзакция корректна
func (bc *Blockchain) ValidateTransaction(tx Transaction) bool {
	sender, err := bc.AccountManager.GetAccount(tx.Sender)
	if err != nil {
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
	err := bc.AccountManager.CreateAccount(address, balance)
	if err != nil {
		fmt.Println("Error creating account:", err)
	}
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
