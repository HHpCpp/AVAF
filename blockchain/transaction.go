package blockchain

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
)

// Transaction представляет собой транзакцию между двумя аккаунтами
type Transaction struct {
	Sender    string  `json:"sender"`
	Recipient string  `json:"recipient"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	Signature string  `json:"signature"`
}

// NewTransaction создает новую транзакцию
func NewTransaction(sender, recipient string, amount float64) (*Transaction, error) {
	if sender == recipient {
		return nil, errors.New("sender and recipient cannot be the same")
	}

	if amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	return &Transaction{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
		Currency:  "AVAF",
	}, nil
}

// Sign подписывает транзакцию с использованием приватного ключа
func (t *Transaction) Sign(privateKey *ecdsa.PrivateKey) error {
	// Хешируем данные транзакции
	hash := t.Hash()

	// Подписываем хеш
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Кодируем подпись в hex
	signature := append(r.Bytes(), s.Bytes()...)
	t.Signature = hex.EncodeToString(signature)

	return nil
}

// Verify проверяет подпись транзакции
func (t *Transaction) Verify(publicKey *ecdsa.PublicKey) (bool, error) {
	// Декодируем подпись из hex
	signature, err := hex.DecodeString(t.Signature)
	if err != nil {
		return false, fmt.Errorf("failed to decode signature: %w", err)
	}

	if len(signature) != 64 {
		return false, errors.New("invalid signature length")
	}

	// Разделяем подпись на r и s
	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:])

	// Хешируем данные транзакции
	hash := t.Hash()

	// Проверяем подпись
	return ecdsa.Verify(publicKey, hash[:], r, s), nil
}

// Hash возвращает хеш транзакции
func (t *Transaction) Hash() [32]byte {
	data := fmt.Sprintf("%s-%s-%.18f-%s", t.Sender, t.Recipient, t.Amount, t.Currency)
	return sha256.Sum256([]byte(data))
}
