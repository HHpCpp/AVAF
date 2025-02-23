package accounts

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/HHpCpp/AVAF/crypto"
)

type AccountManager struct {
	mu         sync.RWMutex
	dataDir    string            // Директория для хранения данных аккаунтов
	passwords  map[string]string // Маппинг адресов на пароли
	cipherFile string            // Путь к файлу accountchiper.json
}

type Wallet struct {
	Address string             `json:"address"`
	Crypto  crypto.CryptoJSON  `json:"crypto"`
	Balance map[string]float64 `json:"balances"`
}

func NewAccountManager(dataDir string) *AccountManager {
	cipherFile := filepath.Join(dataDir, "accountchiper.json")
	passwords := make(map[string]string)

	// Загружаем пароли из файла, если он существует
	if _, err := os.Stat(cipherFile); err == nil {
		data, err := os.ReadFile(cipherFile)
		if err != nil {
			panic(fmt.Errorf("failed to read cipher file: %w", err))
		}

		if err := json.Unmarshal(data, &passwords); err != nil {
			panic(fmt.Errorf("failed to unmarshal cipher file: %w", err))
		}
	}

	return &AccountManager{
		dataDir:    dataDir,
		passwords:  passwords,
		cipherFile: cipherFile,
	}
}

// savePasswords сохраняет пароли в файл accountchiper.json
func (am *AccountManager) savePasswords() error {
	am.mu.Lock()
	defer am.mu.Unlock()

	data, err := json.MarshalIndent(am.passwords, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal passwords: %w", err)
	}

	if err := os.WriteFile(am.cipherFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write cipher file: %w", err)
	}

	return nil
}

// CreateAccount создает новый аккаунт и сохраняет пароль
func (am *AccountManager) CreateAccount(password string, balance float64) (string, error) {
	am.mu.Lock()
	defer am.mu.Unlock()

	privateKey, err := GeneratePrivateKey()
	if err != nil {
		return "", fmt.Errorf("failed to generate private key: %w", err)
	}

	address, err := GenerateAddress()
	if err != nil {
		return "", fmt.Errorf("failed to generate address: %w", err)
	}

	cryptoJSON, err := crypto.EncryptData([]byte(privateKey), password)
	if err != nil {
		return "", fmt.Errorf("encryption failed: %w", err)
	}

	wallet := Wallet{
		Address: address,
		Crypto:  *cryptoJSON,
		Balance: map[string]float64{"AVAF": balance},
	}

	filePath := filepath.Join(am.dataDir, address+".json")
	data, err := json.MarshalIndent(wallet, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal error: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0600); err != nil {
		return "", fmt.Errorf("file write error: %w", err)
	}

	// Сохраняем пароль
	am.passwords[address] = password
	if err := am.savePasswords(); err != nil {
		return "", fmt.Errorf("failed to save passwords: %w", err)
	}

	return address, nil
}

// GetAccountBalance возвращает баланс аккаунта
func (am *AccountManager) GetAccountBalance(address, password string) (map[string]float64, error) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	// Проверяем, что пароль совпадает
	storedPassword, exists := am.passwords[address]
	if !exists {
		return nil, fmt.Errorf("account not found")
	}
	if storedPassword != password {
		return nil, fmt.Errorf("incorrect password")
	}

	filePath := filepath.Join(am.dataDir, address+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("file read error: %w", err)
	}

	var wallet Wallet
	if err := json.Unmarshal(data, &wallet); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	_, err = crypto.DecryptData(wallet.Crypto, password)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	return wallet.Balance, nil
}
