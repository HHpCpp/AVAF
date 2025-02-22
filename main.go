package main

import (
	"fmt"

	"github.com/HHpCpp/AVAF/blockchain"
)

func main() {
	// Создаем новый блокчейн
	bc := blockchain.NewBlockchain()

	// Создаем два аккаунта
	account1, err := bc.CreateAccount(100.0)
	if err != nil {
		fmt.Println("Error creating account 1:", err)
		return
	}

	account2, err := bc.CreateAccount(50.0)
	if err != nil {
		fmt.Println("Error creating account 2:", err)
		return
	}

	fmt.Printf("Account 1: %s, Balance: %.2f\n", account1.Address, account1.Balance)
	fmt.Printf("Account 2: %s, Balance: %.2f\n", account2.Address, account2.Balance)

	// Выполняем транзакцию
	tx, err := bc.CreateTransaction(account1.Address, account2.Address, 30.0)
	if err != nil {
		fmt.Println("Error creating transaction:", err)
		return
	}

	fmt.Println("Transaction successful:", tx)

	// Проверяем балансы после транзакции
	updatedAccount1, _ := bc.AccountManager.GetAccount(account1.Address)
	updatedAccount2, _ := bc.AccountManager.GetAccount(account2.Address)

	fmt.Printf("Account 1: %s, Balance: %.2f\n", updatedAccount1.Address, updatedAccount1.Balance)
	fmt.Printf("Account 2: %s, Balance: %.2f\n", updatedAccount2.Address, updatedAccount2.Balance)
}
