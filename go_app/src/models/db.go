package models

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func init() {
	var err error
	driver_name := "mysql"
	if driver_name == "" {
		log.Fatal("Invalid driver name")
	}
	dsn := "instachat:instachat@tcp(instachat_development_database:3306)/instachat_development?charset=utf8mb4&parseTime=True&loc=Local"
	if dsn == "" {
		log.Fatal("Invalid DSN")
	}
	DB, err = sqlx.Connect(driver_name, dsn)
	if err != nil {
		log.Fatal(err)
	}
}
