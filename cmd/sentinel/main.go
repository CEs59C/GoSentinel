package main

import (
	"goSentinel/internal/config"
	"log"
)

func main() {
	log.SetPrefix("[sentinel] ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("[ERROR] проблемы при парсинге файлов:", err)
	}

	app := New(cfg)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
