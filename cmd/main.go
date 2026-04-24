package main

import (
	"flag"
	"log/slog"
	"medvisitlog/internal/pkg/logger"
)

var (
	Host     = "localhost"
	Port     = "5433"
	User     = "ilalaguzin"
	Password = "12345"
	DBName   = "medvisitlog_db"
)

func main() {
	// логгер
	debug := flag.Bool("debug", false, "enable debug mod")
	flag.Parse()
	cleanup, err := logger.InitLogger(*debug)
	if err != nil {
		panic("ERROR logger" + err.Error())
	}
	defer cleanup()

	slog.Info("Start program")

	// реализация

}
