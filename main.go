package main

import (
	"fmt"

	"github.com/HHpCpp/AVAF/blockchain"
)

func main() {
	// Создаем новый блокчейн
	bc := blockchain.NewBlockchain()

	// Создаем аккаунт
	address1, err := bc.CreateAccount("strong_password_1", 100.0)
	if err != nil {
		fmt.Println("Error creating account:", err)
		return
	}
	fmt.Println("Created account 1:", address1)

	address2, err := bc.CreateAccount("strong_password_2", 200.0)
	if err != nil {
		fmt.Println("Error creating account:", err)
		return
	}
	fmt.Println("Created account 2:", address2)

	// Получаем баланс
	balance1, err := bc.GetAccountBalance(address1, "strong_password_1")
	if err != nil {
		fmt.Println("Error getting balance:", err)
		return
	}
	fmt.Println("Account 1 balance:", balance1)

	balance2, err := bc.GetAccountBalance(address2, "strong_password_2")
	if err != nil {
		fmt.Println("Error getting balance:", err)
		return
	}
	fmt.Println("Account 2 balance:", balance2)
}
