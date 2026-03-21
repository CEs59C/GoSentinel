package email

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func MainSendPost() {
	sendYandexEmail("hello world")
}

func sendYandexEmail(body string) {
	err := godotenv.Load(".env")
	from := os.Getenv("POST_IN")
	poTo := os.Getenv("POST_TO")
	password := os.Getenv("PASSWORD")

	to := []string{poTo}
	smtpHost := "smtp.yandex.ru"
	smtpPort := "587" // Можно также 465 для SSL

	message := []byte("Subject: System Monitor Report\r\n" +
		"To: " + to[0] + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" +
		body)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println("Ошибка отправки:", err)
		return
	}
	fmt.Println("Письмо успешно отправлено!")
}
