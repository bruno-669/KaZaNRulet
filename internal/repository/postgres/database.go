package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib" // Импортируем драйвер pgx для database/sql
)

var (
	err_db_empty_name = errors.New("ERROR: DataBase: Empty dataBase name")
)

type Config struct {
	host     string
	port     string
	user     string
	password string
	dbName   string
}

func CreateConfig(host, port, user, password, dbName string) (*Config, error) {
	return &Config{
		host:     host,
		port:     port,
		password: password,
		dbName:   dbName,
	}, nil

}

func (c Config) DSN() string {
	base := fmt.Sprintf("postgres://%s:%s@%s:%s/", c.user, c.password, c.host, c.port)
	if c.dbName != "" {
		base += c.dbName
	}
	return base + "?sslmode=disable"
}

func ConnectServer(config Config) (*sql.DB, error) {
	serverConfig := config
	serverConfig.dbName = ""
	dsn := serverConfig.DSN()

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	slog.Info("INFO: DataBase: Connected to PostgreSQL server", "host", config.host, "port", config.port)
	return db, nil
}

func EnsureDatabase(serverDB *sql.DB, dbName string) error {
	var exists bool
	err := serverDB.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", dbName).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		_, err = serverDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			return err
		}
		slog.Info("INFO: DataBase: Database created", "name", dbName)
	} else {
		slog.Info("INFO: DataBase: Database already exists", "name", dbName)
	}
	return nil
}

func ConnectDataBase(config Config) (*sql.DB, error) {
	if config.dbName == "" {
		return nil, err_db_empty_name
	}
	dsn := config.DSN()

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	slog.Info("INFO: DataBase: Connected to database", "name", config.dbName)
	return db, nil
}

func InitConnect(config Config) (*sql.DB, error) {
	var err error
	defer func() {
		if err != nil {
			slog.Error("ERROR: DataBase: " + err.Error())
		}
	}()
	serverDB, err := ConnectServer(config)
	if err != nil {
		return nil, err
	}
	defer serverDB.Close()
	if err = EnsureDatabase(serverDB, config.dbName); err != nil {
		return nil, err
	}
	db, err := ConnectDataBase(config)
	if err != nil {
		return nil, err
	}
	return db, nil
}
