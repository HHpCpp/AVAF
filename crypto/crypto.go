package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/scrypt"
)

// CipherParams содержит параметры шифрования
type CipherParams struct {
	IV string `json:"iv"`
}

// KDFParams содержит параметры KDF-функции
type KDFParams struct {
	DKLen int    `json:"dklen"`
	N     int    `json:"n"`
	P     int    `json:"p"`
	R     int    `json:"r"`
	Salt  string `json:"salt"`
}

// CryptoJSON хранит зашифрованные данные и параметры
type CryptoJSON struct {
	Cipher       string       `json:"cipher"`
	CipherCode   string       `json:"ciphercode"` // Переименовано из CipherText
	CipherParams CipherParams `json:"cipherparams"`
	KDF          string       `json:"kdf"`
	KDFParams    KDFParams    `json:"kdfparams"`
	MAC          string       `json:"mac"`
}

// EncryptData шифрует данные с использованием пароля
func EncryptData(data []byte, password string) (*CryptoJSON, error) {
	// Генерируем соль
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	// Генерируем ключ с помощью scrypt
	derivedKey, err := scrypt.Key(
		[]byte(password),
		salt,
		262144, // N
		8,      // R
		1,      // P
		32,     // DKLen
	)
	if err != nil {
		return nil, fmt.Errorf("scrypt error: %w", err)
	}

	// Генерируем IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("iv generation error: %w", err)
	}

	// Шифруем данные
	block, err := aes.NewCipher(derivedKey[:16])
	if err != nil {
		return nil, fmt.Errorf("aes error: %w", err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(data))
	mode.CryptBlocks(ciphertext, data)

	// Вычисляем MAC
	mac := sha256.Sum256(append(derivedKey[16:], ciphertext...))

	return &CryptoJSON{
		Cipher:     "aes-128-cbc",
		CipherCode: hex.EncodeToString(ciphertext),
		CipherParams: CipherParams{
			IV: hex.EncodeToString(iv),
		},
		KDF: "scrypt",
		KDFParams: KDFParams{
			DKLen: 32,
			N:     262144,
			P:     1,
			R:     8,
			Salt:  hex.EncodeToString(salt),
		},
		MAC: hex.EncodeToString(mac[:]),
	}, nil
}

// DecryptData расшифровывает данные с использованием пароля
func DecryptData(cryptoJSON CryptoJSON, password string) ([]byte, error) {
	// Декодируем соль
	salt, err := hex.DecodeString(cryptoJSON.KDFParams.Salt)
	if err != nil {
		return nil, fmt.Errorf("salt decode error: %w", err)
	}

	// Генерируем ключ
	derivedKey, err := scrypt.Key(
		[]byte(password),
		salt,
		cryptoJSON.KDFParams.N,
		cryptoJSON.KDFParams.R,
		cryptoJSON.KDFParams.P,
		cryptoJSON.KDFParams.DKLen,
	)
	if err != nil {
		return nil, fmt.Errorf("scrypt error: %w", err)
	}

	// Декодируем ciphercode
	ciphertext, err := hex.DecodeString(cryptoJSON.CipherCode)
	if err != nil {
		return nil, fmt.Errorf("ciphertext decode error: %w", err)
	}

	// Проверяем MAC
	mac := sha256.Sum256(append(derivedKey[16:], ciphertext...))
	expectedMAC, err := hex.DecodeString(cryptoJSON.MAC)
	if err != nil {
		return nil, fmt.Errorf("mac decode error: %w", err)
	}

	if !bytes.Equal(mac[:], expectedMAC) {
		return nil, errors.New("mac mismatch")
	}

	// Декодируем IV
	iv, err := hex.DecodeString(cryptoJSON.CipherParams.IV)
	if err != nil {
		return nil, fmt.Errorf("iv decode error: %w", err)
	}

	// Расшифровываем данные
	block, err := aes.NewCipher(derivedKey[:16])
	if err != nil {
		return nil, fmt.Errorf("aes error: %w", err)
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	return plaintext, nil
}
