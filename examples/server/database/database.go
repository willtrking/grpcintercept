package database

import (
	"sync"

	"github.com/gocraft/dbr"
	_ "github.com/lib/pq"
)

var connection *dbr.Connection
var err error
var once sync.Once

func GetConnection() (*dbr.Connection, error) {
	once.Do(func() {
		connection, err = dbr.Open("postgres", "My_postgres_string", &EventLog{})
	})

	return connection, err
}
