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
	rep := report.Collect()

	switch a.cfg.Mode {
	case config.ModeEmailHTML, config.ModeTest:
		htmlBody, err := message.RenderHTMLReport(rep)
		if err != nil {
			return fmt.Errorf("html render failed, fallback to plain text %w", err)
		}

		// режим "теста-отладки" с отображением письма на localhost:8080 без отправки
		if a.cfg.Mode == config.ModeTest {
			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(htmlBody))
			})

			log.Println("Server started on :8080")
			http.ListenAndServe(":8080", nil)
			return nil
		}

		// отправка письма на почту в виде html
		err = message.SendYandexEmailHTMLView(a.cfg, htmlBody)
		if err != nil {
			return fmt.Errorf("[ERROR] проблемы при отправке отчета %w", err)
		}

	case config.ModeEmailText: // отправка письма на почту в виде текста
		err := message.SendYandexEmail(a.cfg, rep.String())
		if err != nil {
			return fmt.Errorf("[ERROR] проблемы при отправке отчета %w", err)
		}

	case config.ModeConsole: // печать текста в консоль
		log.Println("[INFO] Сообщение отправлено в консоль")
		fmt.Println(rep.String())

	default:
		return fmt.Errorf("unknown mode")
	}

	return nil
}
