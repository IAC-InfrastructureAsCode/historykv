package mysql

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	DB                *sql.DB
	LimitHistory      int
	ConsulPrefix      string
	AddUserSTMT       *sql.Stmt
	UpdateUserSTMT    *sql.Stmt
	UpdateTokenSTMT   *sql.Stmt
	DeleteUserSTMT    *sql.Stmt
	AddHistorySTMT    *sql.Stmt
	DeleteHistorySTMT *sql.Stmt
}

func New(dsn string, limit int, consulPrefix string) *DB {
	if dsn == "" {
		log.Fatalln("Please insert MySQL DSN")
	}

	var d DB

	var err error
	d.DB, err = sql.Open("mysql", dsn)

	if err != nil {
		log.Fatalf("Couldn't connect to MySQL DB. DSN: %s\n", dsn)
	}

	d.DB.SetMaxIdleConns(2)
	d.DB.SetMaxOpenConns(4)

	if limit <= 0 {
		limit = 5
	}

	d.LimitHistory = limit
	d.ConsulPrefix = consulPrefix

	return &d
}

func (d *DB) CreateStatement() {
	var err error

	d.AddUserSTMT, err = d.DB.Prepare("INSERT INTO kv_user (user, password, token) VALUES(?,?,?)")

	if err != nil {
		log.Fatalln("Couldn't create Statement: AddUser.", err)
	}

	d.UpdateUserSTMT, err = d.DB.Prepare("UPDATE kv_user SET password = ? WHERE user = ?")

	if err != nil {
		log.Fatalln("Couldn't create Statement: UpdateUser.", err)
	}

	d.DeleteUserSTMT, err = d.DB.Prepare("DELETE FROM kv_user WHERE user = ?")

	if err != nil {
		log.Fatalln("Couldn't create Statement: DeleteUser.", err)
	}

	d.UpdateTokenSTMT, err = d.DB.Prepare("UPDATE kv_user SET token = ? WHERE user = ?")

	if err != nil {
		log.Fatalln("Couldn't create Statement: UpdateToken.", err)
	}

	d.AddHistorySTMT, err = d.DB.Prepare("INSERT INTO kv_history (`key`, value, `by`, time) VALUES(?,?,?,?)")

	if err != nil {
		log.Fatalln("Couldn't create Statement: AddHistory.", err)
	}

	d.DeleteHistorySTMT, err = d.DB.Prepare("DELETE FROM kv_history WHERE id <= ?")

	if err != nil {
		log.Fatalln("Couldn't create Statement: DeleteHistory.", err)
	}
}

func (d *DB) CreateTable() {
	historySTMT, errHistorySTMT := d.DB.Prepare(`
        CREATE TABLE IF NOT EXISTS kv_history(
            id bigint NOT NULL AUTO_INCREMENT,
            ` + "`key`" + ` varchar(512) NOT NULL,
            value BLOB,
            ` + "`by`" + ` varchar(100) NOT NULL,
            time bigint NOT NULL,
            PRIMARY KEY (id),
            INDEX idx_kv_history_key_id (` + "`key`" + `, id DESC)
        )
    `)

	if errHistorySTMT == nil {
		_, errHistorySTMT = historySTMT.Exec()
		if errHistorySTMT != nil {
			log.Fatalln("Error on creating History Table:", errHistorySTMT)
		}
	} else {
		log.Fatalln("Error on creating History Table:", errHistorySTMT)
	}

	userSTMT, errUserSTMT := d.DB.Prepare(`
        CREATE TABLE IF NOT EXISTS kv_user(
            id int NOT NULL AUTO_INCREMENT,
            user varchar(100) NOT NULL UNIQUE,
            password varchar(33) NOT NULL,
            token varchar(255),
            PRIMARY KEY (id)
        )
    `)

	if errUserSTMT == nil {
		_, errUserSTMT = userSTMT.Exec()
		if errUserSTMT != nil {
			log.Fatalln("Error on creating User Table:", errUserSTMT)
		}
	} else {
		log.Fatalln("Error on creating User Table:", errUserSTMT)
	}
}

func (d *DB) AddUser(user string, pass string, token string) bool {
	if user == "" || pass == "" {
		return false
	}

	res, err := d.AddUserSTMT.Exec(user, pass, token)

	if err == nil {
		isAdded, addedErr := res.RowsAffected()
		if addedErr == nil && isAdded >= 1 {
			log.Println("Create User:", user)
			return true
		} else {
			log.Println("Failed to create user:", addedErr)
		}
	} else {
		log.Println("Failed to create user:", err)
	}

	return false
}

func (d *DB) UpdateUser(user string, pass string) bool {
	if user == "" || pass == "" {
		return false
	}

	res, err := d.UpdateUserSTMT.Exec(pass, user)

	if err == nil {
		isUpdated, updateErr := res.RowsAffected()
		if updateErr == nil && isUpdated >= 1 {
			return true
		}
	}

	return false
}

func (d *DB) DeleteUser(user string) bool {
	if user == "" {
		return false
	}

	res, err := d.DeleteUserSTMT.Exec(user)

	if err == nil {
		isDeleted, deleteErr := res.RowsAffected()
		if deleteErr == nil && isDeleted >= 1 {
			log.Println("Delete User:", user)
			return true
		} else {
			log.Println("Failed to delete user:", deleteErr)
		}
	} else {
		log.Println("Failed to delete user:", err)
	}

	return false
}

func (d *DB) GetUser(user string, pass string) string {
	if user == "" || pass == "" {
		return ""
	}

	userRow, userErr := d.DB.Query("SELECT user FROM kv_user WHERE user = ? AND password = ?", user, pass)

	if userErr == nil {
		defer userRow.Close()
		for userRow.Next() {
			var userDB string
			scanErr := userRow.Scan(&userDB)

			if scanErr == nil && userDB == user {
				return userDB
			}
		}
	}

	return ""
}

func (d *DB) IsUserExists(user string) bool {
	if user == "" {
		return false
	}

	userRow, userErr := d.DB.Query("SELECT id FROM kv_user WHERE user = ?", user)

	if userErr == nil {
		defer userRow.Close()
		for userRow.Next() {
			var userID int64
			scanErr := userRow.Scan(&userID)

			if scanErr == nil && userID > 0 {
				return true
			}
		}
	}

	return false
}

func (d *DB) GetToken(user string) string {
	if user == "" {
		return ""
	}

	userRow, userErr := d.DB.Query("SELECT token FROM kv_user WHERE user = ?", user)

	if userErr == nil {
		defer userRow.Close()
		for userRow.Next() {
			var tokenDB string
			scanErr := userRow.Scan(&tokenDB)

			if scanErr == nil {
				return tokenDB
			}
		}
	}

	return ""
}

func (d *DB) UpdateToken(user string, token string) bool {
	if user == "" {
		return false
	}

	res, err := d.UpdateTokenSTMT.Exec(token, user)

	if err == nil {
		isUpdated, updateErr := res.RowsAffected()
		if updateErr == nil && isUpdated >= 1 {
			return true
		}
	}

	return false
}

func (d *DB) AddHistory(key string, value string, by string, time int64) bool {
	if key == "" || by == "" || time <= 0 {
		return false
	}

	if d.ConsulPrefix != "" {
		key = fmt.Sprintf("%s%s", d.ConsulPrefix, key)
	}

	res, err := d.AddHistorySTMT.Exec(key, value, by, time)

	if err == nil {
		isAdded, addedErr := res.RowsAffected()
		if addedErr == nil && isAdded >= 1 {
			log.Println("Update key", key, "by", by)
			return true
		}
	}

	return false
}

func (d *DB) GetHistoryList(key string) []map[string]interface{} {
	if key == "" {
		return nil
	}

	if d.ConsulPrefix != "" {
		key = fmt.Sprintf("%s%s", d.ConsulPrefix, key)
	}

	historyRows, historyErr := d.DB.Query("SELECT id, `by`, time FROM kv_history WHERE `key` = ? ORDER BY id DESC LIMIT ?", key, (d.LimitHistory + 1))

	if historyErr == nil {
		defer historyRows.Close()

		var historyInterface []map[string]interface{}

		historyCount := 0

		for historyRows.Next() {
			var historyID, historyTime int64
			var historyBy string
			scanErr := historyRows.Scan(&historyID, &historyBy, &historyTime)

			if scanErr == nil {
				historyCount++
				if historyCount > d.LimitHistory {
					go d.DeleteHistory(key, historyID)
					break
				}

				historyResult := make(map[string]interface{})
				historyResult["id"] = historyID
				historyResult["user"] = historyBy
				historyResult["time"] = historyTime
				historyInterface = append(historyInterface, historyResult)
			}
		}

		return historyInterface
	}

	return nil
}

func (d *DB) GetHistoryID(id int64) (string, string) {
	if id < 0 {
		return "", ""
	}

	historyRow, historyErr := d.DB.Query("SELECT `key`, value FROM kv_history WHERE id = ?", id)

	if historyErr == nil {
		defer historyRow.Close()
		for historyRow.Next() {
			var historyKey, historyValue string
			scanErr := historyRow.Scan(&historyKey, &historyValue)

			if scanErr == nil {
				if d.ConsulPrefix != "" {
					historyKey = historyKey[len(d.ConsulPrefix):]
				}
				return historyKey, historyValue
			}
		}
	}

	return "", ""
}

func (d *DB) DeleteHistory(key string, lastID int64) int {
	if key == "" || lastID < 0 {
		return 0
	}

	if d.ConsulPrefix != "" {
		key = fmt.Sprintf("%s%s", d.ConsulPrefix, key)
	}

	res, err := d.DeleteHistorySTMT.Exec(key, lastID)

	if err == nil {
		isRemoved, _ := res.RowsAffected()
		if isRemoved >= 1 {
			return int(isRemoved)
		}
	}

	return 0
}

func (d *DB) GetUserList(user string, lastID int64) ([]map[string]interface{}, bool) {
	var dbRows *sql.Rows
	var dbErr error

	perPage := 50
	hasNext := false

	if user == "" {
		if lastID <= 0 {
			dbRows, dbErr = d.DB.Query("SELECT id, user, token FROM kv_user ORDER BY id DESC LIMIT ?", (perPage + 1))
		} else {
			dbRows, dbErr = d.DB.Query("SELECT id, user, token FROM kv_user WHERE id < ? ORDER BY id DESC LIMIT ?", lastID, (perPage + 1))
		}
	} else {
		dbRows, dbErr = d.DB.Query("SELECT id, user, token FROM kv_user WHERE user = ?", user)
	}

	if dbErr == nil {
		defer dbRows.Close()

		var userInterface []map[string]interface{}

		countUser := 0
		for dbRows.Next() {
			countUser++
			if countUser > perPage {
				hasNext = true
				break
			}

			var userID int64
			var userName, userToken string
			scanErr := dbRows.Scan(&userID, &userName, &userToken)

			if scanErr == nil {
				userResult := make(map[string]interface{})
				userResult["id"] = userID
				userResult["user"] = userName
				userResult["token"] = userToken
				userInterface = append(userInterface, userResult)
			}
		}

		return userInterface, hasNext
	}

	return nil, hasNext
}
