package main

import (
	"fmt"
	"goSentinel/internal/config"
	"goSentinel/internal/message"
	"goSentinel/internal/report"
	"log"
	"net/http"
)

type app struct {
	cfg *config.Config
}

func New(cfx *config.Config) *app {
	return &app{
		cfg: cfx,
	}
}

func (a *app) Run() error {
	log.Printf("[INFO] стартовало приложение в режиме=%s", a.cfg.Mode)

	rep := report.Collect()

	switch a.cfg.Mode {
	case config.ModeEmailHTML, config.ModeTest:
		htmlBody, err := message.RenderHTMLReport(rep)
		if err != nil {
			return fmt.Errorf("html render failed, fallback to plain text %w", err)
		}

		if a.cfg.Mode == config.ModeTest {
			log.Println("[INFO] режим \"теста-отладки\" с отображением письма на localhost:8080 без отправки")

			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.WriteHeader(http.StatusOK)

				if _, err := w.Write([]byte(htmlBody)); err != nil {
					log.Printf("write response: %v", err)
				}
			})

			log.Println("[INFO] сервер запущен на:8080")
			return http.ListenAndServe(":8080", nil)
		}

		log.Println("[INFO] отправка письма на почту в виде html")
		err = message.SendYandexEmailHTMLView(a.cfg, htmlBody)
		if err != nil {
			return fmt.Errorf("[ERROR] проблемы при отправке отчета %w", err)
		}

	case config.ModeEmailText:
		log.Println("[INFO] отправка письма на почту в виде текста")

		err := message.SendYandexEmail(a.cfg, rep.String())
		if err != nil {
			return fmt.Errorf("[ERROR] проблемы при отправке отчета %w", err)
		}

	case config.ModeConsole:
		log.Println("[INFO] сообщение отправлено в консоль")

		fmt.Println(rep.String())

	default:
		return fmt.Errorf("unknown mode")
	}

	return nil
}
