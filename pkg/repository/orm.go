package repository

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var engine *xorm.Engine

// SetupOrm sets up the database connection
func SetupOrm() *xorm.Engine {
	var err error
	username, password, protocol, address, dbname := os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PROTOCOL"), os.Getenv("DB_ADDRESS"), os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("%v:%v@%v(%v)/%v", username, password, protocol, address, dbname)
	os.Setenv("dsn", dsn)
	engine, err = xorm.NewEngine("mysql", dsn)
	if err != nil {
		log.Fatalf("Error with a database: %v", err)
	}
	return engine
}
