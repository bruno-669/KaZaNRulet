package postgres

import (
	"database/sql"
	_ "embed"
	"log/slog"
	"strconv"
	"strings"
)

//go:embed requests/create_table.sql
var createTableSQL string

//go:embed requests/drop_table.sql
var dropTableSQL string

func executeSQLFromEmbed(db *sql.DB, script string) error {
	commands := strings.Split(script, ";")
	slog.Info("INFO: DataBase:  queries retrieved:" + strconv.Itoa(len(commands)))
	for i, cmd := range commands {
		slog.Info("comand " + strconv.Itoa(i) + " :\"" + cmd + "\"")
	}
	for _, cmd := range commands {
		cmd = strings.TrimSpace(cmd)
		if cmd == "" {
			continue
		}
		if _, err := db.Exec(cmd); err != nil {
			return err
		}
	}

	return nil
}

func CreateTablesIfNotExists(db *sql.DB) error {
	slog.Info("INFO: DataBase: Attempt create tables")
	return executeSQLFromEmbed(db, createTableSQL)
}

func DropTablesIfExists(db *sql.DB) error {
	slog.Info("INFO: DataBase: Attempt drop tables")
	return executeSQLFromEmbed(db, dropTableSQL)
}
