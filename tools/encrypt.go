package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	file := ".env"
	// генерация ключа
	key := make([]byte, 32)
	io.ReadFull(rand.Reader, key)
	keyB64 := base64.StdEncoding.EncodeToString(key)
	if err := os.WriteFile("./build/.env.key", []byte(keyB64), 0600); err != nil {
		log.Printf("Ошибка сохранения в ./build/.env.key: %v\n", err)
		os.Exit(1)
	}

	dataEnv, err := os.ReadFile(file)
	if err != nil {
		log.Printf("Ошибка: .env файл не найден\n")
		os.Exit(1)
	}

	var encrypted string = ""
	lines := strings.Split(string(dataEnv), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "PASSWORD") {
			parts := strings.SplitN(line, "=", 2)
			log.Println(i, parts)

			if len(parts) == 2 {
				encrypted, _ = encryptPassword(parts[1], key)
				lines[i] = "PASSWORD=" + encrypted
				log.Printf("Пароль зашифрован\n")
			} else {
				log.Println("что-то с паролем в файле", file)
			}
		}
	}

	result := strings.Join(lines, "\n")
	if err := os.WriteFile("./build/.env.encrypted", []byte(result), 0600); err != nil {
		log.Printf("Ошибка сохранения в ./build/.env.encrypted: %v\n", err)
		os.Exit(1)
	}

	printMessage(keyB64)
}

func encryptPassword(password string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(password), nil)
	return "ENC:" + base64.StdEncoding.EncodeToString(ciphertext), nil
}

func printMessage(keyB64 string) {
	const template = `
============================================================

Зашифрованный файл: .env.encrypted
Ключ расшифровки: %s
Создан файл с ключем: ./build/.env.key
Создан файл с зашифрованным паролем: ./build/.env.encrypted

ДЛЯ ЗАПУСКА НА СЕРВЕРЕ:
   export ENCRYPTION_KEY='%[1]s'
   ./sentinel..
============================================================
`
	log.Printf(template, keyB64)
}

//запустить
//надо записать в переменную окружения там где работает мэйк
