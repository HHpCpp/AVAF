package accounts

// Account представляет собой аккаунт пользователя
type Account struct {
	Address    string  // Уникальный идентификатор аккаунта
	Balance    float64 // Баланс аккаунта
	PrivateKey string  // Приватный ключ для восстановления
}

// NewAccount создает новый аккаунт
func NewAccount(address string, balance float64, privateKey string) Account {
	return Account{
		Address:    address,
		Balance:    balance,
		PrivateKey: privateKey,
	}
}
