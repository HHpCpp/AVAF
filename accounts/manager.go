package accounts

import (
	"fmt"
	"sync"

	AA "github.com/HHpCpp/AVAF/accounts"
)

// AccountManager управляет аккаунтами
type AccountManager struct {
	accounts map[string]Account
	mu       sync.RWMutex
}

// NewAccountManager создает новый менеджер аккаунтов
func NewAccountManager() *AccountManager {
	return &AccountManager{
		accounts: make(map[string]Account),
	}
}

// CreateAccount создает новый аккаунт с уникальным адресом
func (am *AccountManager) CreateAccount(balance float64) (AA.Account, error) {
	am.mu.Lock()
	defer am.mu.Unlock()

	// Генерируем уникальный адрес
	address, err := GenerateAddress()
	if err != nil {
		return AA.Account{}, fmt.Errorf("failed to generate address: %w", err)
	}

	// Проверяем, что адрес уникален
	if _, exists := am.accounts[address]; exists {
		return AA.Account{}, fmt.Errorf("address collision: %s", address)
	}

	// Создаем аккаунт
	account := AA.NewAccount(address, balance)
	am.accounts[address] = account

	return account, nil
}

// UpdateAccount обновляет существующий аккаунт
func (am *AccountManager) UpdateAccount(account AA.Account) error {
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
