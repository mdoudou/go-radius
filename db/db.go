package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hel2o/go-radius/g"
	"log"
	"strconv"
	"strings"
	"time"
)

type UserPrivilege struct {
	UserName  string
	Privilege int
}

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
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
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

//记录用户认证成功
func AuthSuccess(db *sql.DB, userName, password, result, nasIPAddress, nasIdentifier, framedIPAddress, acctSessionId string) (int64, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	var w = "INSERT INTO `radpostauth` (username, pass, reply, authdate, nasipaddress, clientipaddress, nasidentifier, acctstarttime, acctsessionid) VALUES (?,?,?,now(),?,?,?,now(),?)"
	stmt, err := db.Prepare(w)
	defer stmt.Close()
	HandleErr(err)
	ret, err := stmt.Exec(userName, password, result, nasIPAddress, framedIPAddress, nasIdentifier, acctSessionId)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	rowsAffected, err := ret.RowsAffected()
	HandleErr(err)
	return rowsAffected, err
}

//记录用户登录成功
func Login(db *sql.DB, userName, nasIPAddress, nasIdentifier, framedIPAddress, acctSessionId string) (int64, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	var count int64
	q := "SELECT COUNT(*) FROM radpostauth WHERE acctsessionid = ?"
	err := db.QueryRow(q, acctSessionId).Scan(&count)
	if err != nil {
		log.Println(err)
		return -1, err
	}
	if count == 0 {
		var w = "INSERT INTO `radpostauth` (username, reply, authdate, nasipaddress, clientipaddress, nasidentifier, acctstarttime, acctsessionid) VALUES (?,'Access-Accept',now(),?,?,?,now(),?)"
		stmt, err := db.Prepare(w)
		defer stmt.Close()
		HandleErr(err)
		ret, err := stmt.Exec(userName, nasIPAddress, framedIPAddress, nasIdentifier, acctSessionId)
		if err != nil {
			log.Println(err)
			return 0, err
		}
		rowsAffected, err := ret.RowsAffected()
		HandleErr(err)
		return rowsAffected, err
	}
	return 0, err
}

//记录用户退出登录
func Logout(db *sql.DB, framedIPAddress, acctSessionId string) (int64, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	var w = "UPDATE `radpostauth` SET acctstoptime = now() WHERE acctsessionid = ? AND clientipaddress =?"
	stmt, err := db.Prepare(w)
	defer stmt.Close()
	HandleErr(err)
	ret, err := stmt.Exec(acctSessionId, framedIPAddress)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	rowsAffected, err := ret.RowsAffected()
	HandleErr(err)
	return rowsAffected, err
}

//记录用户登录失败
func AuthFail(db *sql.DB, userName, password, nasIPAddress, nasIdentifier, framedIPAddress, acctSessionId string) (int64, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	var w = "INSERT INTO `radpostauth` (username, pass, reply, authdate, nasipaddress, clientipaddress, nasidentifier, acctsessionid) VALUES (?,?,'Access-Reject',now(),?,?,?,?)"
	stmt, err := db.Prepare(w)
	defer stmt.Close()
	HandleErr(err)
	if acctSessionId == "" {
		acctSessionId = strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	ret, err := stmt.Exec(userName, password, nasIPAddress, framedIPAddress, nasIdentifier, acctSessionId)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	rowsAffected, err := ret.RowsAffected()
	HandleErr(err)
	return rowsAffected, err
}

//读取用户权限
func ReadPrivilege(db *sql.DB, userName, ipAddress string) int {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	var pri int
	var privileges string
	if userName == g.Config().GoRadius.CfgBackName {
		return 3
	}
	q := "SELECT privileges FROM swdb WHERE ipAddress = ?"
	err := db.QueryRow(q, ipAddress).Scan(&privileges)
	if err != nil {
		log.Println(err)
		return 0
	}
	for _, userPri := range decPrivilege(privileges) {
		if userName == userPri.UserName {
			pri = userPri.Privilege
			break
		}
	}
	return pri
}

func decPrivilege(privileges string) []UserPrivilege {
	var userPrivilege []UserPrivilege
	ups := strings.Split(privileges, "|")
	for _, v := range ups {
		up := strings.Split(v, "=")
		p, _ := strconv.Atoi(up[1])
		userPrivilege = append(userPrivilege, UserPrivilege{UserName: up[0], Privilege: p})
	}
	return userPrivilege
}

func HandleErr(err error) {
	if err != nil {
		log.Println(err)
		return
	}
}
