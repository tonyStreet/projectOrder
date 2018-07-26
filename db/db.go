package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"sync"
)

var connection *sql.DB
var once sync.Once

func InitDB() (err error) {
	once.Do(func() {
		dbip := "192.168.0.2:3306"
		dbuser := "order-app"
		dbpassword := "5sEjLqbLxs"
		dbname := "logistics"
		var conStr string
		conStr = dbuser + ":" + dbpassword + "@" + "tcp(" + dbip + ")/" + dbname + "?charset=utf8"
		connection, err = sql.Open("mysql", conStr)
		connection.SetMaxOpenConns(10)
	})
	return err
}

func GetDataSource() (conn *sql.DB, err error) {
	err = connection.Ping()
	return connection, err
}
