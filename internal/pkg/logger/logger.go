package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

func setupFileLogger() (*os.File, error) {

	appDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	logFileName := filepath.Join(appDir, "../tmp/logs", "app.log")
	os.Remove(logFileName)
	logDir := filepath.Join(appDir, "../tmp/logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}
	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func InitLogger(consoleWrite bool) (cleanup func(), err error) {
	logFile, err := setupFileLogger()
	if err != nil {
		fmt.Printf("ERROR setupFileLogger: %v\n", err)
	}

	var logWriter io.Writer
	if consoleWrite && logFile != nil {
		logWriter = io.MultiWriter(os.Stdout, logFile)
	} else if logFile != nil {
		logWriter = logFile
	} else {
		logWriter = os.Stdout
	}

	logHandler := slog.NewTextHandler(logWriter, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})

	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	cleanup = func() {
		slog.Info("Exit program")
		logFile.Close()
	}
	return cleanup, nil
}
