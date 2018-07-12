package db

import (
	"log"

	historyDatabaseMySQL "github.com/tokopedia/historykv/src/db/mysql"
	historyDatabaseSQLite "github.com/tokopedia/historykv/src/db/sqlite"
	historyUtil "github.com/tokopedia/historykv/src/util"
)

type DB interface {
	AddUser(user string, pass string, token string) bool
	AddHistory(key string, value string, by string, time int64) bool
	UpdateToken(user string, token string) bool
	UpdateUser(user string, pass string) bool
	DeleteUser(user string) bool
	IsUserExists(user string) bool
	GetToken(user string) string
	GetUser(user string, password string) string
	GetHistoryList(key string) []map[string]interface{}
	GetHistoryID(id int64) (string, string)
	DeleteHistory(key string, lastID int64) int
	GetUserList(user string, lastID int64) ([]map[string]interface{}, bool)
	CreateTable()
	CreateStatement()
}

func New(dbType string, dbPath string, dbHistoryLimit int, consulPrefix string) DB {
	var d DB

	if dbType == "mysql" {
		log.Println("> Database : Using MySQL. Addr:", dbPath)
		d = historyDatabaseMySQL.New(dbPath, dbHistoryLimit, consulPrefix)
	} else {
		log.Println("> Database : Using SQLite. Path:", dbPath)
		d = historyDatabaseSQLite.New(dbPath, dbHistoryLimit, consulPrefix)
	}

	log.Println("--------------------------------------------------")
	log.Println("Preparing Database...")
	d.CreateTable()
	log.Println("OK")

	log.Println("--------------------------------------------------")
	log.Println("Preparing Statement...")
	d.CreateStatement()
	log.Println("OK")

	// Create Admin
	if !d.IsUserExists("admin") {
		if d.AddUser("admin", historyUtil.CreateHashPassword("admin", "admin"), "") {
			log.Println("--------------------------------------------------")
			log.Println("Admin created..")
			log.Println("   User: admin")
			log.Println("   Pass: admin")
			log.Println("It is recommended to change admin password later.")
			log.Println("--------------------------------------------------")
		}
	}

	return d
}
