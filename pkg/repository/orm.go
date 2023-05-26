package repository

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"xorm.io/xorm"
)

// SetupOrm sets up the database connection
func NewOrm(dsn string) *xorm.Engine {
	engine, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		log.Fatalf("Error with creating a database engine: %v", err)
	}
	return engine
}

// func NewMigrate(dsn string) *migrate.Migrate {

func NewMigration(dsn string, migrationPath string) *migrate.Migrate {
	dsn += "?multiStatements=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error with connecting to the database: %v", err)
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})

	if err != nil {
		log.Fatalf("Error with creating a database driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		"mysql",
		driver)
	if err != nil {
		log.Fatalf("error with a database: %v", err)
	}
	return m
}
