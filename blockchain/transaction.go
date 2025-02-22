package blockchain

// Transaction представляет собой транзакцию между двумя аккаунтами
type Transaction struct {
	Sender    string  // Адрес отправителя
	Recipient string  // Адрес получателя
	Amount    float64 // Сумма перевода
}

// NewTransaction создает новую транзакцию
func NewTransaction(sender, recipient string, amount float64) Transaction {
	return Transaction{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
	}
}
