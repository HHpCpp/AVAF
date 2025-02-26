package pos

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"time"

	crg "crypto/rand"

	"github.com/HHpCpp/AVAF/accounts"
	"github.com/HHpCpp/AVAF/adb"
	ye "github.com/HHpCpp/AVAF/crypto"
)

type StakingWallet struct {
	Address     string                       // Адрес кошелька для стейкинга
	db          *adb.LevelDB                 // LevelDB для хранения данных
	privateKeys map[string]*ecdsa.PrivateKey // Хранение приватных ключей для подписи
}

func NewStakingWallet(db *adb.LevelDB) *StakingWallet {
	return &StakingWallet{
		Address:     "AVAFuNETWORKaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		db:          db,
		privateKeys: make(map[string]*ecdsa.PrivateKey),
	}
}

func (sw *StakingWallet) GetAccount(address string) (*accounts.Account, error) {
	key := "account_" + address
	fmt.Printf("Loading account with key: %s\n", key) // Отладочный вывод

	data, err := sw.db.Load(key)
	if err != nil {
		return nil, fmt.Errorf("failed to load account: %w", err)
	}

	var account accounts.Account
	if err := json.Unmarshal(data, &account); err != nil {
		return nil, fmt.Errorf("failed to unmarshal account: %w", err)
	}

	fmt.Printf("Account loaded: Address=%s, Balance=%v\n", account.Address, account.Balance) // Отладочный вывод
	return &account, nil
}

func (sw *StakingWallet) StakeTokens(accountAddress string, amount float64, privateKey *ecdsa.PrivateKey) error {
	// Проверяем, что приватный ключ не равен nil
	if privateKey == nil {
		return fmt.Errorf("private key is required for staking")
	}

	// Проверяем, что приватный ключ действителен
	if !privateKey.Curve.IsOnCurve(privateKey.PublicKey.X, privateKey.PublicKey.Y) {
		return fmt.Errorf("invalid private key: public key is not on the curve")
	}

	// Получаем аккаунт
	fmt.Printf("Attempting to load account: %s\n", accountAddress) // Отладочный вывод
	account, err := sw.GetAccount(accountAddress)
	if err != nil {
		return fmt.Errorf("failed to get account: %w", err)
	}
	fmt.Printf("Account loaded successfully: %+v\n", account) // Отладочный вывод

	// Извлекаем баланс для валюты "AVAF"
	balance, ok := account.Balance["AVAF"]
	if !ok {
		return fmt.Errorf("currency 'AVAF' not found in account balance")
	}
	fmt.Printf("Current balance: %f\n", balance) // Отладочный вывод

	// Проверяем, что у аккаунта достаточно токенов
	if amount > balance {
		return fmt.Errorf("insufficient balance")
	}

	// Проверяем, что приватный ключ соответствует адресу аккаунта
	expectedAddress := ye.PubkeyToAddress(privateKey.PublicKey)
	if expectedAddress != accountAddress {
		return fmt.Errorf("private key does not match the account address")
	}

	// Подписываем транзакцию стейкинга
	stakeTx := &StakeTransaction{
		AccountAddress: accountAddress,
		Amount:         amount,
		Timestamp:      time.Now().UTC().Format(time.RFC3339),
	}

	if err := stakeTx.Sign(privateKey); err != nil {
		return fmt.Errorf("failed to sign stake transaction: %w", err)
	}

	// Проверяем подпись
	publicKey := &privateKey.PublicKey
	valid, err := stakeTx.Verify(publicKey)
	if err != nil {
		return fmt.Errorf("failed to verify stake transaction: %w", err)
	}
	if !valid {
		return fmt.Errorf("invalid stake transaction signature")
	}

	// Вычитаем сумму из баланса
	account.Balance["AVAF"] = balance - amount
	fmt.Printf("New balance after staking: %f\n", account.Balance["AVAF"]) // Отладочный вывод

	// Сохраняем обновленный аккаунт
	if err := sw.saveAccount(account); err != nil {
		return fmt.Errorf("failed to save account: %w", err)
	}
	fmt.Println("Account saved successfully") // Отладочный вывод

	// Сохраняем стейк в LevelDB
	if err := sw.SaveStake(accountAddress, amount); err != nil {
		return fmt.Errorf("failed to save stake: %w", err)
	}
	fmt.Println("Stake saved successfully") // Отладочный вывод

	return nil
}

func (sw *StakingWallet) saveAccount(account *accounts.Account) error {
	data, err := json.Marshal(account)
	if err != nil {
		return fmt.Errorf("failed to marshal account: %w", err)
	}

	return sw.db.Save("account_"+account.Address, data)
}

func (sw *StakingWallet) SaveStake(address string, stake float64) error {
	stakeBytes := []byte(fmt.Sprintf("%f", stake))
	return sw.db.Save("stake_"+address, stakeBytes)
}

func (sw *StakingWallet) AllValidators() (map[string]float64, error) {
	validators := make(map[string]float64)

	// Создаем итератор для LevelDB
	iter := sw.db.NewIterator()
	defer iter.Release()

	// Проходим по всем ключам
	for iter.Next() {
		key := string(iter.Key())
		value := iter.Value()

		// Проверяем, что ключ начинается с "stake_"
		if strings.HasPrefix(key, "stake_") {
			address := strings.TrimPrefix(key, "stake_")
			stake, err := strconv.ParseFloat(string(value), 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse stake for address %s: %w", address, err)
			}

			validators[address] = stake
		}
	}

	if err := iter.Error(); err != nil {
		return nil, fmt.Errorf("iterator error: %w", err)
	}

	if len(validators) == 0 {
		return nil, fmt.Errorf("no validators available")
	}

	return validators, nil
}

func (sw *StakingWallet) SelectValidator() (string, error) {
	// Получаем всех валидаторов
	validators, err := sw.AllValidators()
	if err != nil {
		return "", fmt.Errorf("failed to get validators: %w", err)
	}

	// Считаем общий стейк
	totalStake := 0.0
	for _, stake := range validators {
		totalStake += stake
	}

	if totalStake == 0 {
		return "", fmt.Errorf("no validators available")
	}

	// Генерируем случайное число
	rand.Seed(time.Now().UnixNano())
	r := rand.Float64() * totalStake

	// Выбираем валидатора
	for address, stake := range validators {
		r -= stake
		if r <= 0 {
			return address, nil
		}
	}

	return "", fmt.Errorf("failed to select validator")
}

// StakeTransaction представляет транзакцию стейкинга
type StakeTransaction struct {
	AccountAddress string  `json:"accountAddress"`
	Amount         float64 `json:"amount"`
	Timestamp      string  `json:"timestamp"`
	Signature      string  `json:"signature"`
}

// Sign подписывает транзакцию стейкинга
func (st *StakeTransaction) Sign(privateKey *ecdsa.PrivateKey) error {
	hash := st.Hash()
	r, s, err := ecdsa.Sign(crg.Reader, privateKey, hash[:])
	if err != nil {
		return fmt.Errorf("failed to sign stake transaction: %w", err)
	}

	signature := append(r.Bytes(), s.Bytes()...)
	st.Signature = hex.EncodeToString(signature)
	return nil
}

// Verify проверяет подпись транзакции
func (st *StakeTransaction) Verify(publicKey *ecdsa.PublicKey) (bool, error) {
	signature, err := hex.DecodeString(st.Signature)
	if err != nil {
		return false, fmt.Errorf("failed to decode signature: %w", err)
	}

	if len(signature) != 64 {
		return false, errors.New("invalid signature length")
	}

	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:])
	hash := st.Hash()

	return ecdsa.Verify(publicKey, hash[:], r, s), nil
}

// Hash возвращает хеш транзакции
func (st *StakeTransaction) Hash() [32]byte {
	data := fmt.Sprintf(
		"%s-%.18f-%s",
		st.AccountAddress,
		st.Amount,
		st.Timestamp,
	)
	return sha256.Sum256([]byte(data))
}
