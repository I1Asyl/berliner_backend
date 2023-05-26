package main

import (
	"fmt"
	"log"
	"os"

	"github.com/I1Asyl/ginBerliner/pkg/handler"
	"github.com/I1Asyl/ginBerliner/pkg/repository"
	"github.com/I1Asyl/ginBerliner/pkg/services"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func main() {

	err := setupConfigs()
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println(os.Getenv("dsn"))
	repository := repository.NewRepository(os.Getenv("dsn"), viper.GetString("db.migrationsUrl"))
	services := services.NewService(repository)
	handler := handler.NewHandler(services)
	err = repository.Migration.Up()
	if err != nil {
		log.Printf("No changes: %v", err)
	}
	router := handler.InitRouter()
	router.Run()
}

func setupConfigs() error {

	if err := godotenv.Load("configs/.env"); err != nil {
		return err
	}
	viper.SetConfigName("config")
	viper.AddConfigPath("configs/")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if os.Getenv("DB_PASSWORD") == "" {
		log.Fatal("DB_PASSWORD is not set, set your database password in .env file")
	}
	username, password, protocol, address, dbname := viper.GetString("db.user"), os.Getenv("DB_PASSWORD"), viper.GetString("db.protocol"), viper.GetString("db.address"), viper.GetString("db.name")

	dsn := fmt.Sprintf("%v:%v@%v(%v)/%v", username, password, protocol, address, dbname)
	os.Setenv("dsn", dsn)

	return nil
}
