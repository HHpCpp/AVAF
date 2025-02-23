package main

import (
	"fmt"

	"github.com/HHpCpp/AVAF/accounts"
	"github.com/HHpCpp/AVAF/blockchain"
)

func main() {
	// Создаем AccountManager
	manager := accounts.NewAccountManager("db/accounts/data")

	// Создаем аккаунт
	address, _, err := manager.CreateAccount("strong_password", 100.0)
	if err != nil {
		panic(err)
	}
	fmt.Println("Created account:", address)

	// Получаем приватный ключ
	privateKey, err := manager.GetPrivateKey(address, "strong_password")
	if err != nil {
		panic(err)
	}
	fmt.Println("Private key:", privateKey)
}
func prevmain() {
	// Создаем новый блокчейн
	bc := blockchain.NewBlockchain()

	// Создаем аккаунт
	address1, _, err := bc.CreateAccount("password1", 100.0)
	if err != nil {
		panic(err)
	}
	fmt.Println("Created account 1:", address1)

	address2, privateKey1, err := bc.CreateAccount("password2", 0.0)
	if err != nil {
		panic(err)
	}
	fmt.Println("Created account 2:", address2)

	// Проверяем балансы до транзакции
	balance1, err := bc.GetAccountBalance(address1)
	if err != nil {
		panic(err)
	}
	fmt.Println("Balance of account 1 before transaction:", balance1)

	balance2, err := bc.GetAccountBalance(address2)
	if err != nil {
		panic(err)
	}
	fmt.Println("Balance of account 2 before transaction:", balance2)

	// Создаем транзакцию
	tx, err := bc.CreateTransaction(address1, address2, privateKey1, 50.0)
	if err != nil {
		panic(err)
	}
	fmt.Println("Created transaction:", tx)

	// Проверяем балансы после транзакции
	balance1, err = bc.GetAccountBalance(address1)
	if err != nil {
		panic(err)
	}
	fmt.Println("Balance of account 1 after transaction:", balance1)

	balance2, err = bc.GetAccountBalance(address2)
	if err != nil {
		panic(err)
	}
	fmt.Println("Balance of account 2 after transaction:", balance2)
}
