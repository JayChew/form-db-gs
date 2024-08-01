package core

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var DBConn *sqlx.DB

func InitalizeDB(host string, port string, user string, password string, name string) *sqlx.DB {
	mysqlHost := host
	mysqlPort := port
	mysqlUsername := user
	mysqlPassword := password
	mysqlTable := name

	var connection string
	connection = mysqlUsername +
		":" +
		mysqlPassword +
		"@tcp(" +
		mysqlHost +
		":" +
		mysqlPort +
		")/" +
		mysqlTable + "?parseTime=true"

	conn, err := sqlx.Open("mysql", connection)
	if err != nil {
		panic(err)
	}

	conn = conn.Unsafe()
	// See "Important settings" section.
	conn.SetConnMaxLifetime(3 * time.Second)
	conn.SetMaxOpenConns(100)
	conn.SetMaxIdleConns(100)
	DBConn = conn
	return DBConn
}