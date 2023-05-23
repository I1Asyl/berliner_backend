package main

import (
	"log"

	"github.com/I1Asyl/ginBerliner/pkg/handler"
	"github.com/I1Asyl/ginBerliner/pkg/repository"
	"github.com/I1Asyl/ginBerliner/pkg/services"
	"github.com/joho/godotenv"
)

func init() {
	setupEnv()
}

func main() {
	repository := repository.NewRepository()
	services := services.NewService(*repository)
	handler := handler.NewHandler(services)

	router := handler.InitRouter()
	router.Run()
}

func setupEnv() {

	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

}
