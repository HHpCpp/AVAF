package tests

import (
	"fmt"

	"github.com/HHpCpp/AVAF/adb"
)

func PrintAllDBEntries(db *adb.LevelDB) {
	iter := db.NewIterator()
	defer iter.Release()

	fmt.Println("All entries in LevelDB:")
	for iter.Next() {
		key := string(iter.Key())
		value := string(iter.Value())
		fmt.Printf("Key: %s, Value: %s\n", key, value)
	}

	if err := iter.Error(); err != nil {
		fmt.Printf("Iterator error: %v\n", err)
	}
}
