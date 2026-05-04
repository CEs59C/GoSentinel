package message

import (
	"fmt"
	"goSentinel/internal/config"
	"log"
	"net/smtp"
	"time"
)

func SendYandexEmail(conf *config.Conf, body string) error {
	log.Printf("[INFO] [%s] Отправка сообщения от %s к %s\n",
		time.Now().Format("15:04:05"), conf.From, conf.To)

	conf.SmtpHost = "smtp.yandex.ru"
	conf.SmtpPort = "587" // Можно также 465 для SSL
	auth := smtp.PlainAuth("", conf.From, conf.Password, conf.SmtpHost)

	message := []byte("Subject: System Monitor Report\r\n" +
		"To: " + conf.To + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" +
		body)

	err := smtp.SendMail(conf.SmtpHost+":"+conf.SmtpPort, auth, conf.From, []string{conf.To}, message)
	if err != nil {
		return fmt.Errorf("Ошибка отправки:  %w", err)
	}
	log.Println("[INFO] Письмо успешно отправлено!")
	return nil
}
