// AVAFu59efebb01cd950f880169d1fae5527f9f328d763 YaSosyDildaki_pass  AVAFu376ac1c128faa35c21477c75a04ad1673131f934

package main

import (
	"fmt"
	"log"

	"github.com/HHpCpp/AVAF/accounts"
	"github.com/HHpCpp/AVAF/adb"
	"github.com/HHpCpp/AVAF/pos"
)

func main() {
	db, err := adb.NewLevelDB("db/LevelDB")
	if err != nil {
		log.Fatalf("Failed to open LevelDB: %v", err)
	}
	defer db.Close()

	// Создаем AccountManager
	accountManager := accounts.NewAccountManager(db)

	// Создаем StakingWallet
	stakingWallet := pos.NewStakingWallet(db)

	// Пример создания аккаунтов
	address, _, err := accountManager.CreateAccount("password2", 2000.0)
	if err != nil {
		log.Fatalf("Failed to generate key pair: %v", err)
	}

	prv, _ := accountManager.GetPrivateKey(address, "password2")
	err = stakingWallet.StakeTokens(address, 1000.0, prv)
	if err != nil {
		log.Fatalf("Failed to load account 1: %v", err)
	}
	bal, _ := accountManager.GetBalance(address)
	fmt.Println(bal)
}

/* db, err := avafdb.NewLevelDB("db/leveldb")
if err != nil {
	fmt.Println("Failed to create LevelDB:", err)
	return
}
defer db.Close()

// Создаем блокчейн
bc, err := blockchain.NewBlockchain(db)
Address, _, _ := bc.AccountManager.CreateAccount("pass1", 100000)
Address1, prv, _ := bc.AccountManager.CreateAccount("pass2", 0)
bcdl11, _ := bc.AccountManager.GetBalance(Address)
bcdl22, _ := bc.AccountManager.GetBalance(Address1)
fmt.Println(bcdl11, bcdl22)

bc.CreateTransaction(Address, Address1, prv, 10, "test")
// bc.CreateTransaction(Address, Address1, prv, 10, "test")
// bc.CreateTransaction(Address, Address1, prv, 10, "test")
// bcdl1, _ := bc.AccountManager.GetBalance(Address)
// bcdl2, _ := bc.AccountManager.GetBalance(Address1)
// fmt.Println(tx, bcdl1, bcdl2)

// GetBlockByHash("c364c3b9e3ac3371e4d72d4c7341121e8e04d3121e29abc008558a3cbc7a2e0d") NOT WORKING
gta, _ := bc.AccountManager.GetAllAccounts()
cua, _ := bc.GetAllBlocks()
// AVAFu62d2b688ef6679219e007cee17933cc587857c68 AVAFue89cf58aff73c4d22bc23ce331ad452801ab03e4
// ht, _ := bc.AccountManager.GetPrivateKey("AVAFu62d2b688ef6679219e007cee17933cc587857c68", "pass2")
fmt.Println(gta)
fmt.Println(cua)

bcdl1, _ := bc.AccountManager.GetBalance(Address)
bcdl2, _ := bc.AccountManager.GetBalance(Address1)
fmt.Println(bcdl1, bcdl2)
// fmt.Println("PrivateKey", ht)
/* )
fmt.Println(bcdl11, bcdl22)

if err != nil {
	fmt.Println("Failed to create blockchain:", err)
	return
}

// Создаем транзакцию
tx, _ := bc.CreateTransaction(Address, Address1, prv, 10, "test")

bcdl1, _ := bc.AccountManager.GetBalance(Address)
bcdl2, _ := bc.AccountManager.GetBalance(Address1)
fmt.Println(tx, bcdl1, bcdl2)*/
