package avafdb

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDB struct {
	db *leveldb.DB
}

// NewLevelDB создает новое подключение к LevelDB
func NewLevelDB(path string) (*LevelDB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open LevelDB: %w", err)
	}
	return &LevelDB{db: db}, nil
}

// Close закрывает подключение к LevelDB
func (l *LevelDB) Close() error {
	return l.db.Close()
}

// Save сохраняет данные по ключу
func (l *LevelDB) Save(key string, value []byte) error {
	return l.db.Put([]byte(key), value, nil)
}

// Load загружает данные по ключу
func (l *LevelDB) Load(key string) ([]byte, error) {
	return l.db.Get([]byte(key), nil)
}

// Delete удаляет данные по ключу
func (l *LevelDB) Delete(key string) error {
	return l.db.Delete([]byte(key), nil)
}
