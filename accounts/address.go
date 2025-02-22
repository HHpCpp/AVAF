package accounts

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/akamensky/base58"
)

const (
	addressPrefix = "AVAFu" // Префикс адреса
	addressLength = 42      // Длина адреса
)

// GeneratePrivateKey генерирует приватный ключ (256 символов в hex-формате)
func GeneratePrivateKey() (string, error) {
	bytes := make([]byte, 128) // 128 байт = 256 символов в hex
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate private key: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateAddress генерирует уникальный адрес кошелька
func GenerateAddress() (string, error) {
	// Длина случайной части адреса
	randomPartLength := addressLength - len(addressPrefix)

	// Генерируем случайные байты
	randomBytes := make([]byte, randomPartLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Кодируем байты в Base58
	randomPart := base58.Encode(randomBytes)

	// Обрезаем до нужной длины
	if len(randomPart) > randomPartLength {
		randomPart = randomPart[:randomPartLength]
	}

	// Собираем адрес
	address := addressPrefix + randomPart

	// Проверяем длину адреса
	if len(address) != addressLength {
		return "", errors.New("generated address has incorrect length")
	}

	return address, nil
}

// ValidateAddress проверяет, что адрес соответствует требованиям
func ValidateAddress(address string) bool {
	// Проверяем длину
	if len(address) != addressLength {
		return false
	}

	// Проверяем префикс
	if !strings.HasPrefix(address, addressPrefix) {
		return false
	}

	return true
}
