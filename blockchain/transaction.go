package blockchain

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"
)

type Transaction struct {
	Hash       string  `json:"hash"`
	Type       string  `json:"type"`      // transfer/smartcontract
	Sender     string  `json:"from"`      // Адрес отправителя
	Recipient  string  `json:"to"`        // Адрес получателя
	ValueType  string  `json:"valueType"` // AVAF
	Value      float64 `json:"value"`     // Количество
	Afuel      float64 `json:"afuel"`     // Единицы вычислительной работы
	AfuelPrice float64 `json:"afuelPrice"`
	Data       string  `json:"data"` // Сообщение
	Signature  string  `json:"signature"`
	Timestamp  string  `json:"timestamp"`
}

func Ntr(sender, recipient string, amount float64, data string) (*Transaction, error) {
	if sender == recipient {
		return nil, errors.New("sender and recipient cannot be the same")
	}

	if amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	// Стандартные значения комиссии
	afuel := 1000.0
	afuelPrice := 0.0001

	tx := &Transaction{
		Type:       "transfer",
		Sender:     sender,
		Recipient:  recipient,
		ValueType:  "AVAF",
		Value:      amount,
		Afuel:      afuel,
		AfuelPrice: afuelPrice,
		Data:       data,
		Timestamp:  time.Now().UTC().Format(time.RFC3339), // Добавляем текущее время
	}

	// Генерируем хеш транзакции
	hash := tx.Hashdo()
	tx.Hash = hex.EncodeToString(hash[:])
	return tx, nil
}

func (t *Transaction) Sign(privateKey *ecdsa.PrivateKey) error {
	hash := t.Hashdo()
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	signature := append(r.Bytes(), s.Bytes()...)
	t.Signature = hex.EncodeToString(signature)
	return nil
}

func (t *Transaction) Verify(publicKey *ecdsa.PublicKey) (bool, error) {
	signature, err := hex.DecodeString(t.Signature)
	if err != nil {
		return false, fmt.Errorf("failed to decode signature: %w", err)
	}

	if len(signature) != 64 {
		return false, errors.New("invalid signature length")
	}

	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:])
	hash := t.Hashdo()

	return ecdsa.Verify(publicKey, hash[:], r, s), nil
}

func (t *Transaction) Hashdo() [32]byte {
	data := fmt.Sprintf(
		"%s-%s-%s-%s-%.18f-%.18f-%.18f-%s-%s",
		t.Type,
		t.Sender,
		t.Recipient,
		t.ValueType,
		t.Value,
		t.Afuel,
		t.AfuelPrice,
		t.Data,
		t.Timestamp,
	)
	return sha256.Sum256([]byte(data))
}
