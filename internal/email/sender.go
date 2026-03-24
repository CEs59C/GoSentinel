package email

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func SendYandexEmail(body string) error {
	// 1. Загружаем .env (ищем в нескольких местах)
	if err := loadEnvFile(); err != nil {
		return fmt.Errorf("не удалось загрузить .env: %w", err)
	}

	from := os.Getenv("POST_IN")
	to := os.Getenv("POST_TO")
	encryptedPassword := os.Getenv("PASSWORD")

	if from == "" || to == "" || encryptedPassword == "" {
		return fmt.Errorf("POST_IN, POST_TO или PASSWORD не установлены")
	}

	password, err := decryptPassword(encryptedPassword)
	if err != nil {
		return fmt.Errorf("не удалось расшифровать пароль: %w", err)
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
		return fmt.Errorf("Ошибка отправки:", err)
	}
	log.Println("Письмо успешно отправлено!")
	return nil
}

func decryptPassword(encrypted string) (string, error) {
	if !strings.HasPrefix(encrypted, "ENC:") {
		return encrypted, nil
	}
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

func loadEnvFile() error {
	paths := []string{
		"build/.env.encrypted", // зашифрованный вариант
		".env.encrypted",       // зашифрованный вариант
		".env",                 // текущая директория
		"/etc/sentinel/.env",   // системная директория
		os.Getenv("HOME") + "/.config/sentinel/.env",
	}
	for _, path := range paths {
		if err := godotenv.Load(path); err == nil {
			log.Println("Загружен файл:", path)
			return nil
		}
	}
	return fmt.Errorf(".env файл не найден")
}
