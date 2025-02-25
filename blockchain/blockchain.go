package blockchain

import (
	"encoding/json"
	"fmt"

	AA "github.com/HHpCpp/AVAF/accounts"
	"github.com/HHpCpp/AVAF/avafdb"
)

type Blockchain struct {
	Chain          []Block
	AccountManager *AA.AccountManager
	db             *avafdb.LevelDB // LevelDB для хранения данных
}

func NewBlockchain(dbPath string) (*Blockchain, error) {
	// Открываем LevelDB
	db, err := avafdb.NewLevelDB(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create LevelDB: %w", err)
	}

	// Создаем AccountManager
	accountManager := AA.NewAccountManager(db)

	// Загружаем генезис-блок
	genesisBlock := NewBlock(0, []Transaction{}, "")
	chain := []Block{genesisBlock}

	// Сохраняем генезис-блок в LevelDB
	if err := saveBlock(db, genesisBlock); err != nil {
		return nil, fmt.Errorf("failed to save genesis block: %w", err)
	}

	return &Blockchain{
		Chain:          chain,
		AccountManager: accountManager,
		db:             db,
	}, nil
}

// saveBlock сохраняет блок в LevelDB
func saveBlock(db *avafdb.LevelDB, block Block) error {
	data, err := json.Marshal(block)
	if err != nil {
		return fmt.Errorf("failed to marshal block: %w", err)
	}
	key := fmt.Sprintf("block_%d", block.Index)
	return db.Save(key, data)
}

// loadBlock загружает блок из LevelDB
func loadBlock(db *avafdb.LevelDB, index int) (Block, error) {
	key := fmt.Sprintf("block_%d", index)
	data, err := db.Load(key)
	if err != nil {
		return Block{}, fmt.Errorf("failed to load block: %w", err)
	}

	var block Block
	if err := json.Unmarshal(data, &block); err != nil {
		return Block{}, fmt.Errorf("failed to unmarshal block: %w", err)
	}

	return block, nil
}

// AddBlock добавляет новый блок в блокчейн
func (bc *Blockchain) AddBlock(transactions []Transaction) error {
	for _, tx := range transactions {
		if !bc.ValidateTransaction(tx) {
			return fmt.Errorf("invalid transaction: %s", tx.Hash)
		}

		// Сохраняем транзакцию в LevelDB
		if err := saveTransaction(bc.db, tx); err != nil {
			return fmt.Errorf("failed to save transaction: %w", err)
		}
	}

	// Обновление балансов
	for _, tx := range transactions {
		commission := tx.Afuel * tx.AfuelPrice

		// Отправитель
		senderBalance, _ := bc.AccountManager.GetBalance(tx.Sender)
		senderBalance["AVAF"] -= tx.Value + commission
		bc.AccountManager.UpdateBalance(tx.Sender, senderBalance)

		// Получатель
		recipientBalance, _ := bc.AccountManager.GetBalance(tx.Recipient)
		recipientBalance["AVAF"] += tx.Value
		bc.AccountManager.UpdateBalance(tx.Recipient, recipientBalance)
	}

	// Создание нового блока
	prevBlock := bc.Chain[len(bc.Chain)-1]
	newBlock := NewBlock(prevBlock.Index+1, transactions, prevBlock.Hash)

	// Сохраняем блок в LevelDB
	if err := saveBlock(bc.db, newBlock); err != nil {
		return fmt.Errorf("failed to save block: %w", err)
	}

	// Добавляем блок в цепочку
	bc.Chain = append(bc.Chain, newBlock)
	return nil
}

func (bc *Blockchain) ValidateTransaction(tx Transaction) bool {
	panic("unimplemented")
}

// saveTransaction сохраняет транзакцию в LevelDB
func saveTransaction(db *avafdb.LevelDB, tx Transaction) error {
	data, err := json.Marshal(tx)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %w", err)
	}
	key := fmt.Sprintf("tx_%s", tx.Hash)
	return db.Save(key, data)
}

// loadTransaction загружает транзакцию из LevelDB
func loadTransaction(db *avafdb.LevelDB, hash string) (Transaction, error) {
	key := fmt.Sprintf("tx_%s", hash)
	data, err := db.Load(key)
	if err != nil {
		return Transaction{}, fmt.Errorf("failed to load transaction: %w", err)
	}

	var tx Transaction
	if err := json.Unmarshal(data, &tx); err != nil {
		return Transaction{}, fmt.Errorf("failed to unmarshal transaction: %w", err)
	}

	return tx, nil
}
