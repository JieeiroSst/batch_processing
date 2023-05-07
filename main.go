package main

import (
	"database/sql"
	"embed"
	"fmt"
	"log"

	"github.com/JIeeiroSst/upload/infrastructure"
	"github.com/JIeeiroSst/upload/repository"
	"github.com/JIeeiroSst/upload/usecase"
	"github.com/JIeeiroSst/upload/utils"
	"github.com/gin-gonic/gin"
	"github.com/pressly/goose/v3"
)

func main() {
	router := gin.Default()

	file := "weather.csv"

	weatherList := utils.CreateWeatherList(file)

	fmt.Println(":========", weatherList)

	infrastructure.LoadEnv()
	database := infrastructure.NewDatabase()
	CreateTable(database.DB)

	repository := repository.NewRepository(database.DB)
	usecase := usecase.NewUsecase(repository)

	usecase.InsertWeather(weatherList)
	router.Run(":8000")
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

func CreateTable(db *sql.DB) {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("mysql"); err != nil {
		log.Println(err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		log.Println(err)
	}
}
