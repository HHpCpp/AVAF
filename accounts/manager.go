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

// CreateAccount создает новый аккаунт и сохраняет его
func (am *AccountManager) CreateAccount(address string, balance float64) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	if _, exists := am.accounts[address]; exists {
		return fmt.Errorf("account already exists: %s", address)
	}

	am.accounts[address] = NewAccount(address, balance)
	return nil
}

// GetAccount возвращает аккаунт по адресу
func (am *AccountManager) GetAccount(address string) (Account, error) {
	am.mu.RLock()
	defer am.mu.RUnlock()

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
