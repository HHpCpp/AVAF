package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/HHpCpp/AVAF/blockchain"
)

func main() {
	// Создаем новый блокчейн
	bc := blockchain.NewBlockchain()

	// Обрабатываем сигналы завершения программы
	setupCloseHandler(&bc)

	// Создаем новый аккаунт
	account, err := bc.CreateAccount(100.0)
	if err != nil {
		fmt.Println("Error creating account:", err)
		return
	}

	fmt.Printf("Created account:\nAddress: %s\nBalance: %.2f\nPrivate Key: %s\n",
		account.Address, account.Balance, account.PrivateKey)

	// Выполняем транзакцию (пример)
	_, err = bc.CreateAccount(50.0)
	if err != nil {
		fmt.Println("Error creating second account:", err)
		return
	}

	// Выводим все аккаунты
	allAccounts := bc.AccountManager.GetAllAccounts()
	for address, acc := range allAccounts {
		fmt.Printf("Account: %s, Balance: %.2f\n", address, acc.Balance)
	}

	// Ждем завершения программы
	select {}
}

// setupCloseHandler настраивает обработчик сигналов для сохранения данных
func setupCloseHandler(bc *blockchain.Blockchain) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\nSaving data before exit...")
		err := bc.Close()
		if err != nil {
			fmt.Println("Error saving data:", err)
		} else {
			fmt.Println("Data saved successfully.")
		}
		os.Exit(0)
	}()
}
