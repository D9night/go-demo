package wk2

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

)

type NotFoundError struct {
	Query string
	Err error
}

func main() {
	list, err := GetUserList("du")
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Printf("cannot find user like %s", "du")
	}
	var res []*User
	list.Scan(res)
	fmt.Printf("查询到的数目： %d", len(res))

}


func GetUserList(n string) (*sql.Rows, error) {
	Db, err := ConnectToDb()
	defer Db.Close()
	if err != nil {
		return nil, err
	}
	return Db.Query("select id, name from user where name regexp ? ", n)
}

// db连接自行处理错误, 不适合给用户看到
func ConnectToDb() (*sql.DB, error) {
	dns := "user:password@tcp(127.0.0.1:3306)/timeout=5s"
	Db, err := sql.Open("mysql", dns)
	if err != nil {
		log.Printf("连接数据库错误：%s\n", err.Error())
		return Db, err
	}

	Db.SetMaxIdleConns(1)
	Db.SetMaxOpenConns(1)
	Db.SetConnMaxIdleTime(time.Second * 60)
	Db.SetConnMaxLifetime(time.Second * 60)

	err = Db.Ping()
	if err != nil {
		return Db, err
	}
	return Db, err

}
