package AVAF

import (
	"encoding/json"
	"fmt"

	"github.com/HHpCpp/AVAF/blockchain"
)

func main() {
	// Создаем новый блокчейн
	bc := blockchain.NewBlockchain()

	// Создаем аккаунты
	bc.CreateAccount("Alice", 1000)
	bc.CreateAccount("Bob", 500)

	// Выводим начальные балансы
	fmt.Println("Initial balances:")
	fmt.Println("Alice:", bc.Accounts["Alice"].Balance)
	fmt.Println("Bob:", bc.Accounts["Bob"].Balance)

	// Создаем транзакции
	transactions := []blockchain.Transaction{
		blockchain.NewTransaction("Alice", "Bob", 200),
		blockchain.NewTransaction("Bob", "Alice", 50),
	}

	// Добавляем блок с транзакциями
	bc.AddBlock(transactions)

	// Выводим блокчейн
	for _, block := range bc.Chain {
		blockJSON, _ := json.MarshalIndent(block, "", "  ")
		fmt.Println(string(blockJSON))
	}

	// Выводим итоговые балансы
	fmt.Println("Final balances:")
	fmt.Println("Alice:", bc.Accounts["Alice"].Balance)
	fmt.Println("Bob:", bc.Accounts["Bob"].Balance)

	// Проверяем валидность блокчейна
	fmt.Println("Is blockchain valid?", bc.IsValid())
}
