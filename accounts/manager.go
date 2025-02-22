package accounts

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// AccountManager управляет аккаунтами
type AccountManager struct {
	accounts map[string]Account
	mu       sync.RWMutex
	filename string // Полный путь к файлу для сохранения данных
}

// NewAccountManager создает новый менеджер аккаунтов
func NewAccountManager(filename string) *AccountManager {
	// Создаем полный путь к файлу
	fullPath := filepath.Join("accounts", "data", filename)

	// Создаем директорию, если она не существует
	err := os.MkdirAll(filepath.Dir(fullPath), 0755)
	if err != nil {
		panic(fmt.Errorf("failed to create data directory: %w", err))
	}

	am := &AccountManager{
		accounts: make(map[string]Account),
		filename: fullPath,
	}
	am.loadFromFile() // Загружаем данные при создании
	return am
}

// SaveToFile сохраняет данные аккаунтов в файл
func (am *AccountManager) SaveToFile() error {
	am.mu.RLock()
	defer am.mu.RUnlock()

	// Сериализуем данные в JSON
	data, err := json.MarshalIndent(am.accounts, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal accounts: %w", err)
	}

	// Записываем данные в файл
	err = os.WriteFile(am.filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

// loadFromFile загружает данные аккаунтов из файла
func (am *AccountManager) loadFromFile() error {
	am.mu.Lock()
	defer am.mu.Unlock()

	// Проверяем, существует ли файл
	if _, err := os.Stat(am.filename); os.IsNotExist(err) {
		return nil // Файл не существует, пропускаем загрузку
	}

	// Читаем данные из файла
	data, err := os.ReadFile(am.filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Десериализуем данные
	err = json.Unmarshal(data, &am.accounts)
	if err != nil {
		return fmt.Errorf("failed to unmarshal accounts: %w", err)
	}

	return nil
}

// Close сохраняет данные перед завершением работы
func (am *AccountManager) Close() error {
	return am.SaveToFile()
}

// CreateAccount создает новый аккаунт с уникальным адресом и приватным ключом
func (am *AccountManager) CreateAccount(balance float64) (Account, error) {
	am.mu.Lock()
	defer am.mu.Unlock()

	// Генерируем уникальный адрес
	address, err := GenerateAddress()
	if err != nil {
		return Account{}, fmt.Errorf("failed to generate address: %w", err)
	}

	// Генерируем приватный ключ
	privateKey, err := GeneratePrivateKey()
	if err != nil {
		return Account{}, fmt.Errorf("failed to generate private key: %w", err)
	}

	// Проверяем, что адрес уникален
	if _, exists := am.accounts[address]; exists {
		return Account{}, fmt.Errorf("address collision: %s", address)
	}

	// Создаем аккаунт
	account := NewAccount(address, balance, privateKey)
	am.accounts[address] = account

	return account, nil
}

// GetAccountByPrivateKey возвращает аккаунт по приватному ключу
func (am *AccountManager) GetAccountByPrivateKey(privateKey string) (Account, error) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	for _, account := range am.accounts {
		if account.PrivateKey == privateKey {
			return account, nil
		}
	}

	return Account{}, fmt.Errorf("account not found for the given private key")
}

// UpdateAccount обновляет существующий аккаунт
func (am *AccountManager) UpdateAccount(account Account) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	// Проверяем, что аккаунт существует
	if _, exists := am.accounts[account.Address]; !exists {
		return fmt.Errorf("account not found: %s", account.Address)
	}

	// Обновляем аккаунт
	am.accounts[account.Address] = account
	return nil
}

// GetAccount возвращает аккаунт по адресу
func (am *AccountManager) GetAccount(address string) (Account, error) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	// Проверяем валидность адреса
	if !ValidateAddress(address) {
		return Account{}, fmt.Errorf("invalid address: %s", address)
	}

	account, exists := am.accounts[address]
	if !exists {
		return Account{}, fmt.Errorf("account not found: %s", address)
	}

	return account, nil
}

// GetAllAccounts возвращает все аккаунты
func (am *AccountManager) GetAllAccounts() map[string]Account {
	am.mu.RLock()
	defer am.mu.RUnlock()

	return am.accounts
}
