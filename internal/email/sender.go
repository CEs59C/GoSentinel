package email

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type EnvMode int

const (
	EnvUnknown EnvMode = iota
	EnvPlain
	EnvEncrypted
)

func SendYandexEmail(body string) error {
	// 1. Загружаем .env (ищем в нескольких местах)
	mode, err := loadEnvFile()
	if err != nil {
		return fmt.Errorf("не удалось загрузить .env: %w", err)
	}

	from := os.Getenv("POST_IN")
	to := os.Getenv("POST_TO")
	encryptedPassword := os.Getenv("PASSWORD")

	if from == "" || to == "" || encryptedPassword == "" {
		return fmt.Errorf("POST_IN, POST_TO или PASSWORD не установлены")
	}

	var password string

	switch mode {
	case EnvEncrypted:
		password, err = decryptPassword(encryptedPassword)
		if err != nil {
			return fmt.Errorf("не удалось расшифровать пароль: %w", err)
		}
	case EnvPlain:
		password = encryptedPassword
		log.Println("Использован не шифрованный пароль")
	default:
		log.Println("[CONFIG] Using plain password (no encryption)")
	}

	log.Printf("[%s] Отправка email от %s к %s\n",
		time.Now().Format("15:04:05"), from, to)

	to1 := []string{to}
	smtpHost := "smtp.yandex.ru"
	smtpPort := "587" // Можно также 465 для SSL

	message := []byte("Subject: System Monitor Report\r\n" +
		"To: " + to1[0] + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" +
		body)

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to1, message)
	if err != nil {
		return fmt.Errorf("Ошибка отправки:  %w", err)
	}
	log.Println("Письмо успешно отправлено!")
	return nil
}

func decryptPassword(encrypted string) (string, error) {
	encrypted = encrypted[4:]
	keyB64 := os.Getenv("ENCRYPTION_KEY")
	if keyB64 == "" {
		return "", fmt.Errorf("ENCRYPTION_KEY не установлен")
	} else {
		log.Printf("Ключ найден\n")
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

func loadEnvFile() (EnvMode, error) {
	paths := []struct {
		path string
		mode EnvMode
	}{
		{"build/.env.encrypted", EnvEncrypted}, // зашифрованный вариант
		{".env.encrypted", EnvEncrypted},       // зашифрованный вариант
		{".env", EnvPlain},                     // текущая директория
		{"/etc/sentinel/.env", EnvPlain},       // системная директория
		{os.Getenv("HOME") + "/.config/sentinel/.env", EnvPlain},
	}
	for _, p := range paths {
		if err := godotenv.Load(p.path); err == nil {
			switch p.mode {
			case EnvEncrypted:
				log.Println("Загружен шифрованный файл:", p)
			case EnvPlain:
				log.Println("Загружен нешифрованный файл:", p)
			}
			return p.mode, nil
		}
	}
	return EnvUnknown, fmt.Errorf(".env файл не найден")
}
