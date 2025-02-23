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
	// Указываем директорию для хранения данных аккаунтов
	accountManager := AA.NewAccountManager("db/accounts/data")
	genesisBlock := NewBlock(0, []Transaction{}, "")
	return Blockchain{
		Chain:          []Block{genesisBlock},
		AccountManager: accountManager,
	}
}

// Close сохраняет данные перед завершением работы
func (bc *Blockchain) Close() error {
	// В новой реализации AccountManager не требует явного сохранения,
	// так как данные сохраняются сразу при создании аккаунтов.
	return nil
}

// CreateAccount создает новый аккаунт
func (bc *Blockchain) CreateAccount(password string, balance float64) (string, error) {
	// Создаем аккаунт через AccountManager
	address, err := bc.AccountManager.CreateAccount(password, balance)
	if err != nil {
		return "", fmt.Errorf("failed to create account: %w", err)
	}
	return address, nil
}

// GetAccountBalance возвращает баланс аккаунта
func (bc *Blockchain) GetAccountBalance(address, password string) (map[string]float64, error) {
	// Получаем баланс через AccountManager
	balance, err := bc.AccountManager.GetAccount(address, password)
	if err != nil {
		return nil, fmt.Errorf("failed to get account balance: %w", err)
	}
	return balance, nil
}

// CreateTransaction создает и обрабатывает новую транзакцию
func (bc *Blockchain) CreateTransaction(sender, recipient, password string, amount float64) (*Transaction, error) {
	// Проверяем, что отправитель и получатель не совпадают
	if sender == recipient {
		return nil, fmt.Errorf("sender and recipient cannot be the same")
	}

	// Получаем баланс отправителя
	senderBalance, err := bc.GetAccountBalance(sender, password)
	if err != nil {
		return nil, fmt.Errorf("failed to get sender balance: %w", err)
	}

	// Проверяем, что у отправителя достаточно средств
	if senderBalance["AVAF"] < amount {
		return nil, fmt.Errorf("insufficient balance: sender has %.2f, required %.2f", senderBalance["AVAF"], amount)
	}

	// Получаем баланс получателя
	recipientBalance, err := bc.GetAccountBalance(recipient, password)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipient balance: %w", err)
	}

	// Обновляем балансы
	senderBalance["AVAF"] -= amount
	recipientBalance["AVAF"] += amount

	// TODO: Сохранить обновленные балансы в AccountManager

	// Создаем транзакцию
	tx, err := NewTransaction(sender, recipient, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

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
	// Валидация транзакции (базовая проверка)
	if tx.Sender == "" || tx.Recipient == "" {
		fmt.Println("Invalid transaction: sender or recipient is empty")
		return false
	}

	if tx.Amount <= 0 {
		fmt.Println("Invalid transaction: amount must be greater than 0")
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
