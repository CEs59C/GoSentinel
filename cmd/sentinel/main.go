package main

import (
	"fmt"
	"goSentinel/internal/config"
	"goSentinel/internal/message"
	"goSentinel/internal/report"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.ParseEnv()
	if err != nil {
		log.Fatal("[ERROR] проблемы при парсинге файлов:", err)
	}

	rep := report.Collect()

	htmlBody, err := message.RenderHTMLReport(rep)
	if err != nil {
		log.Fatal("[WARN] HTML render failed, fallback to plain text", err)
	}

	cfg.EnableHTTPPreview = false
	if cfg.EnableHTTPPreview == true {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(htmlBody))
		})

		log.Println("Server started on :8080")
		http.ListenAndServe(":8080", nil)
	}

	if cfg.IsHTMLView {
		err = message.SendYandexEmailHTMLView(cfg, htmlBody)
		if err != nil {
			log.Println("[ERROR] проблемы при отправке отчета", err)
			return
		}
	} else {
		if cfg.IsSendMail {
			err = message.SendYandexEmail(cfg, rep.String())
			if err != nil {
				log.Println("[ERROR] проблемы при отправке отчета", err)
				return
			}
		} else {
			fmt.Println(rep.String())
		}
	}

	log.Println("[INFO] Сообщение отправлено")
}
