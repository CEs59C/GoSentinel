package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Mode string

const (
	ModeConsole   Mode = "console"
	ModeEmailText Mode = "email-text"
	ModeEmailHTML Mode = "email-html"
	ModeTest      Mode = "test"
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

type envFileInfo struct {
	path        string
	isEncrypted bool
}

func LoadConfig() (*Config, error) {
	log.Println("[INFO] начата работа слоя Config")

	envPath, err := findEnvFile()
	if err != nil {
		return nil, err
	}

	cfg1, err := buildConfig(envPath)
	if err != nil {
		return nil, err
	}

	err = cfg1.validate()
	if err != nil {
		return nil, err
	}

	cfg2, err := Decrypt(*cfg1)
	if err != nil {
		return nil, err
	}

	return &cfg2, err
}

func findEnvFile() (envFileInfo, error) {
	log.Println("[INFO] поиск .env файла")

	home, err := os.UserHomeDir()
	if err != nil {
		log.Println("[ERROR] домашняя директория не определена")
		home = ""
	}

	paths := []struct {
		path        string
		isEncrypted bool
	}{
		{"build/.env.encrypted", true},          // зашифрованный вариант
		{".env.encrypted", true},                // зашифрованный вариант
		{"build/.env", false},                   // текущая директория
		{".env", false},                         // текущая директория
		{"/etc/sentinel/.env", false},           // системная директория
		{home + "/config/sentinel/.env", false}, // для линукса
	}

	pathsEnv := envFileInfo{}

	for _, p := range paths {
		if _, err := os.Stat(p.path); err == nil {
			if p.isEncrypted {
				log.Println("[INFO] найден шифрованный файл: ", p.path)
				pathsEnv.path = p.path
				pathsEnv.isEncrypted = true
				break
			} else {
				log.Println("[INFO] найден нешифрованный файл:", p.path)
				pathsEnv.path = p.path
				pathsEnv.isEncrypted = false
				break
			}
		}
	}

	if pathsEnv.path == "" {
		return envFileInfo{}, fmt.Errorf("файл не найден")
	}

	return pathsEnv, nil
}

func (c *Config) validate() error {
	if c.From == "" || c.To == "" {
		return fmt.Errorf("email not configured")
	}

	switch c.Mode {
	case ModeConsole, ModeEmailText, ModeEmailHTML, ModeTest:
		return nil
	case "":
		return fmt.Errorf("MODE не задан")
	default:
		return fmt.Errorf("неизвестный MODE: %s", c.Mode)
	}
}

func buildConfig(envPath envFileInfo) (*Config, error) {
	if err := godotenv.Load(envPath.path); err != nil {
		return nil, fmt.Errorf("загрузка %s: %w", envPath.path, err)
	}

	cfg := &Config{
		PathEnv:           envPath.path,
		EncryptedPassword: envPath.isEncrypted,
		From:              os.Getenv("POST_IN"),
		To:                os.Getenv("POST_TO"),
		Password:          os.Getenv("PASSWORD"),
		Mode:              Mode(os.Getenv("MODE")),
		SmtpHost:          os.Getenv("SMTP_HOST"),
		SmtpPort:          os.Getenv("SMTP_PORT"),
	}

	return cfg, nil
}
