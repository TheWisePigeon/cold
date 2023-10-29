package repositories

import (
  _ "github.com/mattn/go-sqlite3"
  "github.com/jmoiron/sqlx"
)

var DB *sqlx.DB
var err error

func ConnectToDB() error {
  DB, err = sqlx.Connect("sqlite3", "/home/thewisepigeon/code/cold/.dev.cold.db")
  return err
}
