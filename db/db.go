package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hel2o/go-radius/g"
	"log"
)

var RadiusDb, FireSystemDb *sql.DB

func InitDB() {
	radiusDb, err := sql.Open("mysql", g.Config().GoRadius.RadiusDb)
	if err != nil {
		log.Println(err)
		return
	}
	err = radiusDb.Ping()
	if err != nil {
		log.Println(err)
		return
	}
	RadiusDb = radiusDb
	fireSystemDb, err := sql.Open("mysql", g.Config().GoRadius.FireSystemDb)
	if err != nil {
		log.Println(err)
		return
	}
	err = fireSystemDb.Ping()
	if err != nil {
		log.Println(err)
		return
	}
	FireSystemDb = fireSystemDb
}

//对比radius用户名和密码是否一致
func CheckUserPassword(db *sql.DB, username, password string) bool {
	var count int64
	q := "SELECT COUNT(*) FROM radcheck WHERE username = ? AND value = ?"
	err := db.QueryRow(q, username, password).Scan(&count)
	if err != nil {
		log.Println(err)
		return false
	}
	if count == 1 {
		return true
	}
	return false
}
