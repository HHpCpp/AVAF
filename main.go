package main

import (
	"fmt"

	"github.com/HHpCpp/AVAF/blockchain"
)

func main() {
	// Создаем новый блокчейн
	bc := blockchain.NewBlockchain()

	// Создаем новый аккаунт
	bc.CreateAccount("address1", 100.0)
	bc.CreateAccount("address2", 200.0)

	// Получаем все аккаунты
	accountManager := bc.AccountManager
	allAccounts := accountManager.GetAllAccounts()

	// Выводим все аккаунты
	for address, account := range allAccounts {
		fmt.Printf("Address: %s, Balance: %.2f\n", address, account.Balance)
	}
}
