package main

import (
	"fmt"
	"goSentinel/internal/config"
	"goSentinel/internal/message"
	"goSentinel/internal/report"
	"log"
)

func main() {
	cfg, err := config.ParseEnv()
	if err != nil {
		log.Fatal("[ERROR] проблемы при парсинге файлов:", err)
	}

	rep := report.Report()

	// опция отправки на почту или перчати в терминал
	// true - почта, false - печать в терминал
	//config.IsSendMail = true

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
