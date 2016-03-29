package database

import (
  "sync"

	_ "github.com/lib/pq"
  "github.com/gocraft/dbr"
)


var connection *dbr.Connection
var err error
var once sync.Once

func GetConnection() (*dbr.Connection, error) {
  once.Do(func() {
    connection, err  = dbr.Open("postgres","My_postgres_string",&EventLog{})
  })

  return connection, err
}
