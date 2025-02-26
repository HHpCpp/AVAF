package accounts

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/HHpCpp/AVAF/adb"
	"github.com/HHpCpp/AVAF/crypto"
)

var (
	publicKeyCache = make(map[string]*ecdsa.PublicKey)
	cacheMutex     sync.RWMutex
)

type AccountManager struct {
	mu sync.RWMutex
	db *adb.LevelDB
}

type Wallet struct {
	Address   string             `json:"address"`
	Crypto    crypto.CryptoJSON  `json:"crypto"`
	Balance   map[string]float64 `json:"balances"` // Изменено на float64
	PublicKey string             `json:"publicKey"`
}

func NewAccountManager(db *adb.LevelDB) *AccountManager {
	return &AccountManager{db: db}
}

func (am *AccountManager) SaveAccount(wallet Wallet) error {
	data, err := json.Marshal(wallet)
	if err != nil {
		return fmt.Errorf("failed to marshal wallet: %w", err)
	}

	key := "account_" + wallet.Address
	fmt.Printf("Saving account with key: %s\n", key) // Отладочный вывод
	return am.db.Save(key, data)
}

func (am *AccountManager) LoadAccount(address string) (Wallet, error) {
	data, err := am.db.Load("account_" + address)
	if err != nil {
		return Wallet{}, fmt.Errorf("failed to load account: %w", err)
	}

	var wallet Wallet
	if err := json.Unmarshal(data, &wallet); err != nil {
		return Wallet{}, fmt.Errorf("failed to unmarshal wallet: %w", err)
	}

	return wallet, nil
}

func (am *AccountManager) GetPrivateKey(address, password string) (*ecdsa.PrivateKey, error) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	wallet, err := am.LoadAccount(address)
	if err != nil {
		return nil, err
	}

	privateKeyBytes, err := crypto.DecryptData(wallet.Crypto, password)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt private key: %w", err)
	}

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

func (am *AccountManager) CreateAccount(password string, balance float64) (string, *ecdsa.PrivateKey, error) {
	am.mu.Lock()
	defer am.mu.Unlock()

	privateKey, address, err := GenerateKeyPair()
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate key pair: %w", err)
	}

	privateKeyHex := hex.EncodeToString(privateKey.D.Bytes())

	cryptoJSON, err := crypto.EncryptData([]byte(privateKeyHex), password)
	if err != nil {
		return "", nil, fmt.Errorf("failed to encrypt private key: %w", err)
	}

	publicKeyHex := hex.EncodeToString(append(
		privateKey.PublicKey.X.Bytes(),
		privateKey.PublicKey.Y.Bytes()...,
	))

	wallet := Wallet{
		Address:   address,
		Crypto:    *cryptoJSON,
		Balance:   map[string]float64{"AVAF": balance},
		PublicKey: publicKeyHex,
	}

	if err := am.SaveAccount(wallet); err != nil {
		return "", nil, fmt.Errorf("failed to save account: %w", err)
	}

	cacheMutex.Lock()
	publicKeyCache[address] = &privateKey.PublicKey
	cacheMutex.Unlock()

	return address, privateKey, nil
}
func (am *AccountManager) GetAllAccounts() ([]string, error) {
	var Addresses []string
	// Создаем итератор для LevelDB
	iter := am.db.NewIterator()
	defer iter.Release()

	// Проходим по всем ключам
	for iter.Next() {
		key := string(iter.Key())
		value := iter.Value()

		// Проверяем, что ключ начинается с "account_"
		if strings.HasPrefix(key, "account_") {
			var account Account
			if err := json.Unmarshal(value, &account); err != nil {
				return nil, fmt.Errorf("failed to unmarshal account: %w", err)
			}
			Addresses = append(Addresses, account.Address)
		}
	}

	if err := iter.Error(); err != nil {
		return nil, fmt.Errorf("iterator error: %w", err)
	}

	return Addresses, nil
}
func (am *AccountManager) GetBalance(address string) (map[string]float64, error) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	wallet, err := am.LoadAccount(address)
	if err != nil {
		return nil, err
	}

	return wallet.Balance, nil // Возвращаем баланс напрямую
}

func (am *AccountManager) UpdateBalance(address string, balances map[string]float64) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	wallet, err := am.LoadAccount(address)
	if err != nil {
		return err
	}

	// Обновляем баланс для каждой валюты
	for currency, balance := range balances {
		wallet.Balance[currency] = balance // Просто обновляем значение
	}

	return am.SaveAccount(wallet)
}

func (am *AccountManager) GetPublicKey(address string) (*ecdsa.PublicKey, error) {
	cacheMutex.RLock()
	if pubKey, ok := publicKeyCache[address]; ok {
		cacheMutex.RUnlock()
		return pubKey, nil
	}
	cacheMutex.RUnlock()

	wallet, err := am.LoadAccount(address)
	if err != nil {
		return nil, err
	}

	pubKeyBytes, err := hex.DecodeString(wallet.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %w", err)
	}

	pubKey := new(ecdsa.PublicKey)
	pubKey.Curve = elliptic.P256()
	pubKey.X = new(big.Int).SetBytes(pubKeyBytes[:32])
	pubKey.Y = new(big.Int).SetBytes(pubKeyBytes[32:])

	cacheMutex.Lock()
	publicKeyCache[address] = pubKey
	cacheMutex.Unlock()

	return pubKey, nil
}
