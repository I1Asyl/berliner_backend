package repository

import (
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var engine *xorm.Engine

// SetupOrm sets up the database connection
func SetupOrm() *xorm.Engine {
	var err error
	dsn := os.Getenv("dsn")
	engine, err = xorm.NewEngine("mysql", dsn)
	if err != nil {
		log.Fatalf("Error with a database: %v", err)
	}
	return engine
}
