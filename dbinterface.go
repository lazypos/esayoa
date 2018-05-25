package esayoa

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

const (
	connString = "postgres://%s:%s@%s/%s?sslmode=%s" //用户-密码-IP-库名
)

type DBpgopt struct {
	DB *sql.DB
}

var GDBpgopt *DBpgopt = &DBpgopt{}

func (this *DBpgopt) Init() error {
	return this.Connect("postgres", "514ddwddw", "127.0.0.1", "esayoa", "disable")
}

func (this *DBpgopt) Connect(user, pass, host, dbname, sslmode string) error {
	this.Close()
	var err error = nil
	this.DB, err = sql.Open("postgres", fmt.Sprintf(connString, user, pass, host, dbname, sslmode))
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("连接数据库成功")
	return nil
}

func (this *DBpgopt) Close() {
	if this.DB != nil {
		this.DB.Close()
		this.DB = nil
		log.Println("断开数据库连接")
	}
}

func (this *DBpgopt) Query(sql string) (*sql.Rows, error) {
	rows, err := this.DB.Query(sql)
	return rows, err
}

func (this *DBpgopt) Execute(sql string) (sql.Result, error) {
	ret, err := this.DB.Exec(sql)
	return ret, err
}

//查询某个值
func (this *DBpgopt) QueryVal(sql string) (string, error) {
	rows, err := this.DB.Query(sql)
	if err != nil {
		return "", err
	}
	if !rows.Next() {
		return "", nil
	}
	defer rows.Close()

	rst := ""
	err = rows.Scan(&rst)
	return rst, err
}
