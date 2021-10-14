package main

import (
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
)

type DBConfig struct {
	Port     string
	Host     string
	Password string
	Database string
	User     string
}

func main() {
	conf := DBConfig{
		Port:     os.Getenv("MYSQL_PORT"),
		Host:     os.Getenv("MYSQL_HOST"),
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Database: os.Getenv("MYSQL_DATABASE"),
	}
	dsn := conf.User + ":" + conf.Password + "@tcp(" + conf.Host + ":" + conf.Port + ")/" + conf.Database + "?parseTime=true&multiStatements=true"
	fmt.Print(dsn)
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		log.Error(err)
		return
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS TEST (id INT PRIMARY KEY, text VARCHAR(100))")
	if err != nil {
		log.Error(err)
		return
	}

	id := time.Now().UnixMicro() % 10000000
	_, err = db.Exec("INSERT INTO TEST (id, text) values (?, 'help!')", id)
	if err != nil {
		log.Error(err)
		return
	}
	rows := []struct {
		Id   int32  `db:"id"`
		Text string `db:"text"`
	}{}
	err = db.Select(&rows, "SELECT * FROM TEST")
	if err != nil {
		log.Error(err)
		return
	}

	for _, row := range rows {
		log.Print(row.Id, row.Text)
	}
	time.Sleep(10)
	fmt.Println("Hello Docker")
}
