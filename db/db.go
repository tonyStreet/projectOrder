package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"github.com/tonyStreet/projectOrder/config"
)

var connection *sql.DB
var once sync.Once

func InitDB() (err error) {
	conf := config.GetConfig()
	once.Do(func() {
		dbip := conf.DB.IP
		dbuser := conf.DB.User
		dbpassword := conf.DB.Password
		dbname := conf.DB.Name
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
