package blockchain

import (
	"crypto/ecdsa"
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

// CreateAccount создает новый аккаунт
func (bc *Blockchain) CreateAccount(password string, balance float64) (string, *ecdsa.PrivateKey, error) {
	// Создаем аккаунт через AccountManager
	address, privateKey, err := bc.AccountManager.CreateAccount(password, balance)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create account: %w", err)
	}
	return address, privateKey, nil
}

// GetAccountBalance возвращает баланс аккаунта (без пароля)
func (bc *Blockchain) GetAccountBalance(address string) (map[string]float64, error) {
	// Получаем баланс через AccountManager
	balance, err := bc.AccountManager.GetBalance(address)
	if err != nil {
		return nil, fmt.Errorf("failed to get account balance: %w", err)
	}
	return balance, nil
}

// CreateTransaction создает и обрабатывает новую транзакцию
func (bc *Blockchain) CreateTransaction(sender, recipient string, privateKey *ecdsa.PrivateKey, amount float64) (*Transaction, error) {
	// Проверяем, что отправитель и получатель не совпадают
	if sender == recipient {
		return nil, fmt.Errorf("sender and recipient cannot be the same")
	}

	// Получаем баланс отправителя
	senderBalance, err := bc.GetAccountBalance(sender)
	if err != nil {
		return nil, fmt.Errorf("failed to get sender balance: %w", err)
	}

	// Проверяем, что у отправителя достаточно средств
	if senderBalance["AVAF"] < amount {
		return nil, fmt.Errorf("insufficient balance: sender has %.2f, required %.2f", senderBalance["AVAF"], amount)
	}

	// Создаем транзакцию
	tx, err := NewTransaction(sender, recipient, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Подписываем транзакцию с использованием приватного ключа
	if err := tx.Sign(privateKey); err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Проверяем подпись транзакции
	publicKey := &privateKey.PublicKey
	valid, err := tx.Verify(publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to verify transaction: %w", err)
	}
	if !valid {
		return nil, fmt.Errorf("invalid transaction signature")
	}

	// Добавляем транзакцию в новый блок
	bc.AddBlock([]Transaction{*tx})

	return tx, nil
}

// AddBlock добавляет новый блок в блокчейн и обновляет балансы
func (bc *Blockchain) AddBlock(transactions []Transaction) {
	// Проверяем транзакции перед добавлением
	for _, tx := range transactions {
		if !bc.ValidateTransaction(tx) {
			fmt.Println("Invalid transaction:", tx)
			return
		}
	}

	// Обновляем балансы
	for _, tx := range transactions {
		senderBalance, err := bc.GetAccountBalance(tx.Sender)
		if err != nil {
			fmt.Println("Failed to get sender balance:", err)
			return
		}

		recipientBalance, err := bc.GetAccountBalance(tx.Recipient)
		if err != nil {
			fmt.Println("Failed to get recipient balance:", err)
			return
		}

		// Обновляем балансы
		senderBalance["AVAF"] -= tx.Amount
		recipientBalance["AVAF"] += tx.Amount

		// Сохраняем обновленные балансы
		if err := bc.AccountManager.UpdateBalance(tx.Sender, senderBalance); err != nil {
			fmt.Println("Failed to update sender balance:", err)
			return
		}

		if err := bc.AccountManager.UpdateBalance(tx.Recipient, recipientBalance); err != nil {
			fmt.Println("Failed to update recipient balance:", err)
			return
		}
	}

	// Добавляем блок в цепочку
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

	// Проверяем подпись транзакции
	publicKey, err := bc.AccountManager.GetPublicKey(tx.Sender)
	if err != nil {
		fmt.Println("Failed to get public key:", err)
		return false
	}

	valid, err := tx.Verify(publicKey)
	if err != nil || !valid {
		fmt.Println("Invalid transaction signature")
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
