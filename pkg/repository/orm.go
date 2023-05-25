package repository

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"xorm.io/xorm"
)

// SetupOrm sets up the database connection
func SetupOrm(dsn string, migrationsUrl string) *xorm.Engine {

	err := migrateAll(dsn, migrationsUrl)

	if err != nil {
		log.Panic(err)
	}
	engine, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		log.Fatalf("Error with creating a database engine: %v", err)
	}
	return engine
}

func migrateAll(dsn string, migrationsUrl string) error {
	dsn += "?multiStatements=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("Error with connecting to the database: %v", err)
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})

	if err != nil {
		return fmt.Errorf("Error with creating a database driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationsUrl,
		"mysql",
		driver)
	if err != nil {
		return fmt.Errorf("error with a database: %v", err)
	}
	err = m.Up()
	if err != nil {
		fmt.Println("No migrations")
	}
	return nil
}
