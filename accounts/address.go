package accounts

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/sha3"
)

// GenerateKeyPair генерирует пару ключей (приватный и публичный) и адрес
func GenerateKeyPair() (*ecdsa.PrivateKey, string, error) {
	// Генерация приватного ключа
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate private key: %w", err)
	}

	// Генерация адреса из публичного ключа
	pubKey := append(
		privateKey.PublicKey.X.Bytes(),
		privateKey.PublicKey.Y.Bytes()...,
	)

	// Хешируем публичный ключ с помощью Keccak-256
	hash := sha3.NewLegacyKeccak256()
	hash.Write(pubKey)
	address := "AVAFu" + hex.EncodeToString(hash.Sum(nil)[12:])

	return privateKey, address, nil
}
