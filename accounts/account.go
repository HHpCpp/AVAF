package accounts

// Account представляет собой аккаунт пользователя
type Account struct {
	Address    string             `json:"address"`    // Уникальный идентификатор аккаунта
	Balance    map[string]float64 `json:"balances"`   // Баланс аккаунта (валюта -> float64)
	PrivateKey string             `json:"privateKey"` // Приватный ключ для восстановления
}

// NewAccount создает новый аккаунт
func NewAccount(address string, balance float64, privateKey string) Account {
	// Инициализируем баланс как map[string]float64
	balances := make(map[string]float64)
	balances["AVAF"] = balance // Пример для валюты "AVAF"

	return Account{
		Address:    address,
		Balance:    balances,
		PrivateKey: privateKey,
	}
}
