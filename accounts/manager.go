package accounts

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"sync"

	"github.com/HHpCpp/AVAF/crypto"
)

var (
	publicKeyCache = make(map[string]*ecdsa.PublicKey)
	cacheMutex     sync.RWMutex
)

type AccountManager struct {
	mu      sync.RWMutex
	dataDir string
}

type Wallet struct {
	Address   string            `json:"address"`
	Crypto    crypto.CryptoJSON `json:"crypto"`    // Зашифрованные данные (приватный ключ)
	Balance   map[string]string `json:"balances"`  // Баланс в hex-коде (открытый)
	PublicKey string            `json:"publicKey"` // Публичный ключ в hex-формате
}

func NewAccountManager(dataDir string) *AccountManager {
	return &AccountManager{
		dataDir: dataDir,
	}
}

// GetPrivateKey возвращает приватный ключ, расшифрованный с использованием пароля
func (am *AccountManager) GetPrivateKey(address, password string) (*ecdsa.PrivateKey, error) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	// Читаем файл кошелька
	filePath := filepath.Join(am.dataDir, address+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("file read error: %w", err)
	}

	// Декодируем JSON
	var wallet Wallet
	if err := json.Unmarshal(data, &wallet); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	// Расшифровываем приватный ключ
	privateKeyBytes, err := crypto.DecryptData(wallet.Crypto, password)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt private key: %w", err)
	}

	// Преобразуем приватный ключ из hex в структуру ecdsa.PrivateKey
	privateKeyHex := string(privateKeyBytes)
	privateKeyBytes, err = hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key hex: %w", err)
	}

	privateKey := new(ecdsa.PrivateKey)
	privateKey.Curve = elliptic.P256()
	privateKey.D = new(big.Int).SetBytes(privateKeyBytes)
	privateKey.PublicKey.X, privateKey.PublicKey.Y = privateKey.Curve.ScalarBaseMult(privateKey.D.Bytes())

	return privateKey, nil
}

// CreateAccount создает новый аккаунт
func (am *AccountManager) CreateAccount(password string, balance float64) (string, *ecdsa.PrivateKey, error) {
	am.mu.Lock()
	defer am.mu.Unlock()

	// Генерация пары ключей и адреса
	privateKey, address, err := GenerateKeyPair()
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate key pair: %w", err)
	}

	// Преобразуем приватный ключ в hex
	privateKeyHex := hex.EncodeToString(privateKey.D.Bytes())

	// Шифруем приватный ключ
	cryptoJSON, err := crypto.EncryptData([]byte(privateKeyHex), password)
	if err != nil {
		return "", nil, fmt.Errorf("failed to encrypt private key: %w", err)
	}

	// Преобразуем баланс в hex-код
	balanceHex := hex.EncodeToString([]byte(fmt.Sprintf("%.18f", balance)))

	// Сохраняем публичный ключ в hex-формате
	publicKeyHex := hex.EncodeToString(append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...))

	wallet := Wallet{
		Address:   address,
		Crypto:    *cryptoJSON,
		Balance:   map[string]string{"AVAF": balanceHex},
		PublicKey: publicKeyHex,
	}

	// Сохраняем кошелек в файл
	filePath := filepath.Join(am.dataDir, address+".json")
	data, err := json.MarshalIndent(wallet, "", "  ")
	if err != nil {
		return "", nil, fmt.Errorf("marshal error: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0600); err != nil {
		return "", nil, fmt.Errorf("file write error: %w", err)
	}

	// Сохраняем публичный ключ в кэш
	cacheMutex.Lock()
	publicKeyCache[address] = &privateKey.PublicKey
	cacheMutex.Unlock()

	return address, privateKey, nil
}

// GetBalance возвращает баланс без необходимости ввода пароля
func (am *AccountManager) GetBalance(address string) (map[string]float64, error) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	filePath := filepath.Join(am.dataDir, address+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("file read error: %w", err)
	}

	var wallet Wallet
	if err := json.Unmarshal(data, &wallet); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	// Декодируем баланс из hex-кода
	balances := make(map[string]float64)
	for currency, balanceHex := range wallet.Balance {
		balanceBytes, err := hex.DecodeString(balanceHex)
		if err != nil {
			return nil, fmt.Errorf("failed to decode balance: %w", err)
		}
		var balance float64
		_, err = fmt.Sscanf(string(balanceBytes), "%f", &balance)
		if err != nil {
			return nil, fmt.Errorf("failed to parse balance: %w", err)
		}
		balances[currency] = balance
	}

	return balances, nil
}

// UpdateBalance обновляет баланс аккаунта
func (am *AccountManager) UpdateBalance(address string, balances map[string]float64) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	filePath := filepath.Join(am.dataDir, address+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("file read error: %w", err)
	}

	var wallet Wallet
	if err := json.Unmarshal(data, &wallet); err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}

	// Обновляем баланс
	wallet.Balance = make(map[string]string)
	for currency, balance := range balances {
		balanceHex := hex.EncodeToString([]byte(fmt.Sprintf("%.18f", balance)))
		wallet.Balance[currency] = balanceHex
	}

	// Сохраняем обновленный кошелек
	data, err = json.MarshalIndent(wallet, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0600); err != nil {
		return fmt.Errorf("file write error: %w", err)
	}

	return nil
}

// GetPublicKey возвращает публичный ключ из кэша
func (am *AccountManager) GetPublicKey(address string) (*ecdsa.PublicKey, error) {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	if pubKey, ok := publicKeyCache[address]; ok {
		return pubKey, nil
	}
	return nil, fmt.Errorf("public key not found for address: %s", address)
}
