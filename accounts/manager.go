package accounts

import (
	"fmt"
	"sync"
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
func (am *AccountManager) CreateAccount(balance float64) (Account, error) {
	am.mu.Lock()
	defer am.mu.Unlock()

	// Генерируем уникальный адрес
	address, err := GenerateAddress()
	if err != nil {
		return Account{}, fmt.Errorf("failed to generate address: %w", err)
	}

	// Проверяем, что адрес уникален
	if _, exists := am.accounts[address]; exists {
		return Account{}, fmt.Errorf("address collision: %s", address)
	}

	// Создаем аккаунт
	account := NewAccount(address, balance)
	am.accounts[address] = account

	return account, nil
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
