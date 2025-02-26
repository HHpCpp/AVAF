package blockchain

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	AA "github.com/HHpCpp/AVAF/accounts"
	avafdb "github.com/HHpCpp/AVAF/adb"
	pos "github.com/HHpCpp/AVAF/pos"
)

type Blockchain struct {
	Chain          []Block
	AccountManager *AA.AccountManager
	db             *avafdb.LevelDB // LevelDB для хранения данных
	StakingWallet  *pos.StakingWallet
}

func (bc *Blockchain) NewTransaction(Address string, Address1 string, prv *ecdsa.PrivateKey, i int) {
	panic("unimplemented")
}

func GetLastBlockIndex(db *avafdb.LevelDB) (int, error) {
	// Создаем итератор для LevelDB
	iter := db.NewIterator()
	defer iter.Release()

	lastIndex := -1

	// Проходим по всем ключам
	for iter.Next() {
		key := string(iter.Key())

		// Проверяем, что ключ начинается с "block_"
		if strings.HasPrefix(key, "block_") {
			// Извлекаем индекс из ключа
			var index int
			_, err := fmt.Sscanf(key, "block_%d", &index)
			if err != nil {
				return -1, fmt.Errorf("failed to parse block index: %w", err)
			}

			// Обновляем последний индекс
			if index > lastIndex {
				lastIndex = index
			}
		}
	}

	if err := iter.Error(); err != nil {
		return -1, fmt.Errorf("iterator error: %w", err)
	}

	return lastIndex, nil
}

func NewBlockchain(db *avafdb.LevelDB) (*Blockchain, error) {
	// Создаем AccountManager
	accountManager := AA.NewAccountManager(db)

	// Создаем StakingWallet
	stakingWallet := pos.NewStakingWallet(db)
	/* if err := stakingWallet.LoadStakes(); err != nil {
		return nil, fmt.Errorf("failed to load stakes: %w", err)
	} */

	// Создаем генезис-блок
	genesisBlock := NewBlock(0, []Transaction{}, "")
	chain := []Block{genesisBlock}

	// Сохраняем генезис-блок в LevelDB
	if err := saveBlock(db, genesisBlock); err != nil {
		return nil, fmt.Errorf("failed to save genesis block: %w", err)
	}

	return &Blockchain{
		Chain:          chain,
		AccountManager: accountManager,
		StakingWallet:  stakingWallet,
		db:             db,
	}, nil
}

func LoadAllBlocks(db *avafdb.LevelDB) ([]Block, error) {
	var blocks []Block

	// Создаем итератор для LevelDB
	iter := db.NewIterator()
	defer iter.Release()

	// Проходим по всем ключам
	for iter.Next() {
		key := string(iter.Key())
		value := iter.Value()

		// Проверяем, что ключ начинается с "block_"
		if strings.HasPrefix(key, "block_") {
			var block Block
			if err := json.Unmarshal(value, &block); err != nil {
				return nil, fmt.Errorf("failed to unmarshal block: %w", err)
			}
			blocks = append(blocks, block)
		}
	}

	if err := iter.Error(); err != nil {
		return nil, fmt.Errorf("iterator error: %w", err)
	}

	// Сортируем блоки по индексу
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].Index < blocks[j].Index
	})

	return blocks, nil
}

// saveBlock сохраняет блок в LevelDB
func saveBlock(db *avafdb.LevelDB, block Block) error {
	data, err := json.Marshal(block)
	if err != nil {
		return fmt.Errorf("failed to marshal block: %w", err)
	}
	key := fmt.Sprintf("block_%d", block.Index) // Уникальный ключ для каждого блока
	return db.Save(key, data)
}

func NewTransaction(sender string, recipient string, amount float64, data string) (*Transaction, error) {
	if sender == recipient {
		return nil, errors.New("sender and recipient cannot be the same")
	}

	if amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	// Стандартные значения комиссии
	afuel := 1000.0
	afuelPrice := 0.0001

	tx := &Transaction{
		Type:       "transfer",
		Sender:     sender,
		Recipient:  recipient,
		ValueType:  "AVAF",
		Value:      amount,
		Afuel:      afuel,
		AfuelPrice: afuelPrice,
		Data:       data,
		Timestamp:  time.Now().Format(time.RFC3339), // Добавляем текущее время
	}

	// Генерируем хеш транзакции
	hash := tx.Hashdo()
	tx.Hash = hex.EncodeToString(hash[:])
	return tx, nil
}
func (bc *Blockchain) CreateTransaction(sender string, recipient string, privateKey *ecdsa.PrivateKey, amount float64, data string) (*Transaction, error) {
	// Проверяем, что отправитель и получатель не совпадают
	if sender == recipient {
		return nil, fmt.Errorf("sender and recipient cannot be the same")
	}

	// Получаем баланс отправителя
	sb, err := bc.AccountManager.GetBalance(sender)
	if err != nil {
		return nil, fmt.Errorf("failed to get sender balance: %w", err)
	}

	// Рассчитываем комиссию (1000 Afuel * 0.0001 AVAF)
	commission := 1000.0 * 0.0001

	// Проверяем, что у отправителя достаточно средств (сумма + комиссия)
	if sb["AVAF"] < amount+commission {
		return nil, fmt.Errorf("insufficient balance: sender has %.2f, required %.2f", sb["AVAF"], amount+commission)
	}

	// Создаем транзакцию
	tx, err := NewTransaction(sender, recipient, amount, data)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Подписываем транзакцию с использованием приватного ключа
	if err := tx.Sign(privateKey); err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Проверяем подпись транзакции
	publicKey := &privateKey.PublicKey
	valid, err := tx.Verify(publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to verify transaction: %w", err)
	}
	if !valid {
		return nil, fmt.Errorf("invalid transaction signature")
	}

	// Добавляем транзакцию в новый блок
	bc.AddBlock([]Transaction{*tx})

	return tx, nil
}

func (bc *Blockchain) GetBlockByHash(hash string) (*Block, error) {
	// Формируем ключ для блока
	key := fmt.Sprintf("block_%s", hash)

	// Загружаем данные из LevelDB
	data, err := bc.db.Load(key)
	if err != nil {
		return nil, fmt.Errorf("failed to load block: %w", err)
	}

	// Десериализуем данные в структуру Block
	var block Block
	if err := json.Unmarshal(data, &block); err != nil {
		return nil, fmt.Errorf("failed to unmarshal block: %w", err)
	}

	return &block, nil
}

func (bc *Blockchain) GetBlockByIndex(index int) (*Block, error) {
	// Формируем ключ для блока
	key := fmt.Sprintf("block_%d", index)

	// Загружаем данные из LevelDB
	data, err := bc.db.Load(key)
	if err != nil {
		return nil, fmt.Errorf("failed to load block: %w", err)
	}

	// Десериализуем данные в структуру Block
	var block Block
	if err := json.Unmarshal(data, &block); err != nil {
		return nil, fmt.Errorf("failed to unmarshal block: %w", err)
	}

	return &block, nil
}
func (bc *Blockchain) GetCurrentBlock() (*Block, error) {
	// Проверяем, что блокчейн не пустой
	if len(bc.Chain) == 0 {
		return nil, fmt.Errorf("blockchain is empty")
	}

	// Возвращаем последний блок
	lastBlock := bc.Chain[len(bc.Chain)-1]
	return &lastBlock, nil
}

func (bc *Blockchain) GetAllBlocks() ([]Block, error) {
	var blocks []Block

	// Создаем итератор для LevelDB
	iter := bc.db.NewIterator()
	defer iter.Release()

	// Проходим по всем ключам
	for iter.Next() {
		key := string(iter.Key())
		value := iter.Value()

		// Проверяем, что ключ начинается с "block_"
		if strings.HasPrefix(key, "block_") {
			var block Block
			if err := json.Unmarshal(value, &block); err != nil {
				return nil, fmt.Errorf("failed to unmarshal block: %w", err)
			}
			blocks = append(blocks, block)
		}
	}

	if err := iter.Error(); err != nil {
		return nil, fmt.Errorf("iterator error: %w", err)
	}

	// Проверяем, что блокчейн не пустой
	if len(blocks) == 0 {
		return nil, fmt.Errorf("blockchain is empty")
	}

	return blocks, nil
}

func (bc *Blockchain) GetAccountBalance(sender string) (any, any) {
	panic("unimplemented")
}

// loadBlock загружает блок из LevelDB
func LoadBlock(db *avafdb.LevelDB, index int) (Block, error) {
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
	// Получаем последний блок
	prevBlock := bc.Chain[len(bc.Chain)-1]

	// Создаем новый блок с индексом на 1 больше, чем у предыдущего
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
	// Retrieve the sender's public key
	publicKey, err := bc.AccountManager.GetPublicKey(tx.Sender)
	if err != nil {
		return false
	}

	// Verify the transaction's signature
	valid, err := tx.Verify(publicKey)
	if err != nil || !valid {
		return false
	}

	// Verify the transaction's hash
	computedHash := tx.Hashdo()
	computedHashStr := hex.EncodeToString(computedHash[:])
	if computedHashStr != tx.Hash {
		return computedHashStr == tx.Hash
	}

	return true
}

// saveTransaction сохраняет транзакцию в LevelDB
func SaveTransaction(db *avafdb.LevelDB, tx Transaction) error {
	data, err := json.Marshal(tx)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %w", err)
	}
	key := fmt.Sprintf("tx_%s", tx.Hash)
	return db.Save(key, data)
}

// loadTransaction загружает транзакцию из LevelDB
func LoadTransaction(db *avafdb.LevelDB, hash string) (Transaction, error) {
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
