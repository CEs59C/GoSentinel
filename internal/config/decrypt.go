package config

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"
)

func Decrypt(cfg Config) (Config, error) {
	if !cfg.EncryptedPassword {
		return cfg, nil
	}

	newCfg := cfg
	decrypted, err := decryptPassword(cfg.Password)
	if err != nil {
		return Config{}, fmt.Errorf("decrypt password: %w", err)
	}

	newCfg.Password = decrypted
	newCfg.EncryptedPassword = false

	return newCfg, nil
}

func decryptPassword(encrypted string) (string, error) {
	log.Println("[INFO] начата расшифровка шифрованного пароля")

	if !strings.HasPrefix(encrypted, "ENC:") {
		return "", fmt.Errorf("invalid encrypted format")
	}
	encrypted = strings.TrimPrefix(encrypted, "ENC:")

	keyB64 := os.Getenv("ENCRYPTION_KEY")
	if keyB64 == "" {
		const template = `
			============================================================
			ДЛЯ ЗАПУСКА:
			   export ENCRYPTION_KEY='Ключ расшифровки'
			============================================================
			`
		log.Println("[WARNING] ключ расшифровки пароля не найден в переменных окружения")
		log.Println(template)
		return "", fmt.Errorf("ENCRYPTION_KEY не установлен")
	} else {
		log.Println("[INFO] ключ найден")
	}

	// Декодируем ключ
	key, err := base64.StdEncoding.DecodeString(keyB64)
	if err != nil {
		return "", fmt.Errorf("неверный формат ключа: %w", err)
	}

	// Декодируем зашифрованные данные
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", fmt.Errorf("неверный формат зашифрованных данных: %w", err)
	}

	// Расшифровываем
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("ошибка создания шифра: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("ошибка создания GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("данные повреждены (слишком короткие)")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("расшифровка не удалась (возможно неверный ключ): %w", err)
	}

	return string(plaintext), nil
}
