package message

import (
	"fmt"
	"goSentinel/internal/config"
	"log"
	"net/smtp"
	"time"
)

func SendYandexEmail(cfg *config.Config, body string) error {
	log.Printf("[INFO] [%s] отправка сообщения от %s к %s\n",
		time.Now().Format("15:04:05"), cfg.From, cfg.To)

	auth := smtp.PlainAuth("", cfg.From, cfg.Password, cfg.SmtpHost)

	message := []byte("Subject: System Monitor Report\r\n" +
		"To: " + cfg.To + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" +
		body)

	err := smtp.SendMail(cfg.SmtpHost+":"+cfg.SmtpPort, auth, cfg.From, []string{cfg.To}, message)
	if err != nil {
		return fmt.Errorf("ошибка отправки:  %w", err)
	}
	log.Println("[INFO] письмо успешно отправлено!")
	return nil
}

func SendYandexEmailHTMLView(conf *config.Config, body string) error {
	log.Printf("[INFO] отправка сообщения c HTML")

	auth := smtp.PlainAuth("", conf.From, conf.Password, conf.SmtpHost)

	message := []byte(
		"Subject: System Monitor Report\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=UTF-8\r\n" +
			"\r\n" +
			body,
	)

	err := smtp.SendMail(conf.SmtpHost+":"+conf.SmtpPort, auth, conf.From, []string{conf.To}, message)
	if err != nil {
		return fmt.Errorf("ошибка отправки:  %w", err)
	}

	return nil
}
