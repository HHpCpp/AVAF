package blockchain

import (
	"errors"
	"fmt"
)

// Transaction представляет собой транзакцию между двумя аккаунтами
type Transaction struct {
	Sender    string  // Адрес отправителя
	Recipient string  // Адрес получателя
	Amount    float64 // Сумма перевода
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
	}, nil
}

// Validate проверяет, что транзакция корректна
func (t *Transaction) Validate() error {
	if t.Sender == "" || t.Recipient == "" {
		return errors.New("sender and recipient addresses cannot be empty")
	}

	if t.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	return nil
}

// String возвращает строковое представление транзакции
func (t *Transaction) String() string {
	return fmt.Sprintf("Transaction(Sender: %s, Recipient: %s, Amount: %.2f)", t.Sender, t.Recipient, t.Amount)
}
