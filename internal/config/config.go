package config

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	EncryptedPassword bool
	Mode              Mode
	PathEnv           string
	Password          string
	From              string
	To                string
	SmtpHost          string
	SmtpPort          string
}
type Mode string

const (
	ModeConsole   Mode = "console"
	ModeEmailText Mode = "email-text"
	ModeEmailHTML Mode = "email-html"
	ModeTest      Mode = "test"
)

func Load() (*Config, error) {
	cfg, err := searchAndLoadEnvFile()
	if err != nil {
		return cfg, fmt.Errorf("не удалось загрузить .env: %w", err)
	}

	cfg.From = os.Getenv("POST_IN")
	cfg.To = os.Getenv("POST_TO")
	cfg.Password = os.Getenv("PASSWORD")
	cfg.Mode = Mode(os.Getenv("MODE"))
	cfg.SmtpHost = os.Getenv("SMTP_HOST")
	cfg.SmtpPort = os.Getenv("SMTP_PORT")

	err = cfg.validate()
	if err != nil {
		return nil, fmt.Errorf("[ERROR] отсутствуют обязательные данные в .ENV: %w", err)
	}

	switch cfg.EncryptedPassword {
	case true:
		cfg.Password, err = decryptPassword(cfg.Password)
		if err != nil {
			return &Config{}, fmt.Errorf("не удалось расшифровать пароль: %w", err)
		}
	case false:
		log.Println("[INFO] Использован не шифрованный пароль")
	default:
		log.Println("[CONFIG] Using plain password (no encryption)")
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.From == "" || c.To == "" {
		return fmt.Errorf("email not configured")
	}

	switch c.Mode {
	case ModeConsole, ModeEmailText, ModeEmailHTML, ModeTest:
		return nil
	default:
		return fmt.Errorf("invalid MODE: %s", c.Mode)
	}
}

func searchAndLoadEnvFile() (*Config, error) {
	log.Println("[INFO] Начат поиск файла .env")
	cfg := &Config{}

	home, err := os.UserHomeDir()
	if err != nil {
		log.Println("[ERROR] Домашняя директория не определена")
		return nil, fmt.Errorf("домашняя директория не определена %w", err)
	}

	paths := []struct {
		path        string
		isEncrypted bool
	}{
		{"build/.env.encrypted", true},          // зашифрованный вариант
		{".env.encrypted", true},                // зашифрованный вариант
		{".env", false},                         // текущая директория
		{"/etc/sentinel/.env", false},           // системная директория
		{home + "/config/sentinel/.env", false}, // для линукса
	}

	for _, p := range paths {
		if err := godotenv.Load(p.path); err == nil {

			if p.isEncrypted {
				log.Println("[INFO] Найден шифрованный файл: ", p.path)
				cfg.EncryptedPassword = true
				cfg.PathEnv = p.path
				break
			} else {
				log.Println("[INFO] Найден нешифрованный файл:", p.path)
				cfg.EncryptedPassword = false
				cfg.PathEnv = p.path
				break
			}
		}
	}

	if cfg.PathEnv == "" {
		return nil, fmt.Errorf(".env файл не найден")
	}

	return cfg, nil
}

func decryptPassword(encrypted string) (string, error) {
	log.Println("[INFO] Начата расшифровка шифрованного пароля")

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
		log.Println("[WARNING] Ключ расшифровки пароля не найден в переменных окружения")
		log.Println(template)
		return "", fmt.Errorf("ENCRYPTION_KEY не установлен")
	} else {
		log.Println("[INFO] Ключ найден")
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
