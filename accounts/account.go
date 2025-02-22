package accounts

// Account представляет собой аккаунт пользователя
type Account struct {
	Address string  // Уникальный идентификатор аккаунта
	Balance float64 // Баланс аккаунта
}

// NewAccount создает новый аккаунт
func NewAccount(address string, balance float64) Account {
	return Account{
		Address: address,
		Balance: balance,
	}
}
