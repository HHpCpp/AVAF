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
	accountManager := AA.NewAccountManager("accounts.json") // Указываем имя файла для сохранения
	genesisBlock := NewBlock(0, []Transaction{}, "")
	return Blockchain{
		Chain:          []Block{genesisBlock},
		AccountManager: accountManager,
	}
}

// Close сохраняет данные перед завершением работы
func (bc *Blockchain) Close() error {
	return bc.AccountManager.Close()
}

// CreateAccount создает новый аккаунт
func (bc *Blockchain) CreateAccount(balance float64) (AA.Account, error) {
	return bc.AccountManager.CreateAccount(balance)
}

// CreateTransaction создает и обрабатывает новую транзакцию
func (bc *Blockchain) CreateTransaction(sender, recipient string, amount float64) (*Transaction, error) {
	// Создаем транзакцию
	tx, err := NewTransaction(sender, recipient, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Проверяем валидность транзакции
	if err := tx.Validate(); err != nil {
		return nil, fmt.Errorf("invalid transaction: %w", err)
	}

	// Проверяем, что у отправителя достаточно средств
	senderAccount, err := bc.AccountManager.GetAccount(sender)
	if err != nil {
		return nil, fmt.Errorf("sender account not found: %w", err)
	}

	if senderAccount.Balance < amount {
		return nil, fmt.Errorf("insufficient balance: sender has %.2f, required %.2f", senderAccount.Balance, amount)
	}

	// Получаем аккаунт получателя
	recipientAccount, err := bc.AccountManager.GetAccount(recipient)
	if err != nil {
		return nil, fmt.Errorf("recipient account not found: %w", err)
	}

	// Обновляем балансы
	senderAccount.Balance -= amount
	recipientAccount.Balance += amount

	// Сохраняем обновленные аккаунты в AccountManager
	bc.AccountManager.UpdateAccount(senderAccount)
	bc.AccountManager.UpdateAccount(recipientAccount)

	// Добавляем транзакцию в новый блок
	bc.AddBlock([]Transaction{*tx})

	return tx, nil
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
