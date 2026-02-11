package db

import (
	"bytes"
	"database/sql"
	"encoding/gob"

	_ "embed"

	"github.com/codeshelldev/gotl/pkg/logger"
	"github.com/codeshelldev/secured-signal-api/internals/config"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

//go:embed schema.sql
var schema string

func Init() {
	var err error

	db, err = sql.Open("sqlite3", config.ENV.DB_PATH)

	if err != nil {
		logger.Fatal("Error opening database: ", err.Error())
		return
	}

	db.SetMaxOpenConns(1)

	err = db.Ping()

	if err != nil {
		logger.Fatal("Error opening database connection: ", err.Error())
		return
	}

	_, err = db.Exec(schema)

	if err != nil {
		logger.Fatal("Could not apply database schema: ", err.Error())
		return
	}

	logger.Debug("Successfully opened database")
}

func Close() {
	ShutdownRequestDB()

	db.Close()
}

func Serialize(value any) []byte {
	var valueBytes bytes.Buffer

	enc := gob.NewEncoder(&valueBytes)
	enc.Encode(value)

	return valueBytes.Bytes()
}

func Deserialize[T any](valueBytes []byte) T {
	var out T

	buf := bytes.NewBuffer(valueBytes)
	dec := gob.NewDecoder(buf)

	dec.Decode(&out)

	return out
}