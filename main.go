package main

import (
	"fmt"

	"github.com/HHpCpp/AVAF/blockchain"
)

func main() {
	// Создаем новый блокчейн
	bc := blockchain.NewBlockchain()

	// Создаем новый аккаунт
	account, err := bc.CreateAccount(100.0)
	if err != nil {
		fmt.Println("Error creating account:", err)
		return
	}

	fmt.Printf("Created account: %+v\n", account)

	// Получаем все аккаунты
	allAccounts := bc.AccountManager.GetAllAccounts()
	for address, acc := range allAccounts {
		fmt.Printf("Address: %s, Balance: %.2f\n", address, acc.Balance)
	}
}
