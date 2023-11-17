package main

import (
	"admin-api/db"
	"admin-api/router"
	"admin-api/validation"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file!!")
	}

	println(">>>>>>>>> Before Init", db.GetDB())
	err := db.InitDB()
	if err != nil {
		log.Fatal("Error connecting database: ", err.Error())
	}

	println(">>>>>>>>> After Init", db.GetDB())

	err = db.Migrate()
	if err != nil {
		log.Fatal("Error migration: ", err.Error())
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Validator = validation.NewCustomValidator()

	v1 := e.Group("/v1")
	router.VersionOne(v1)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("APP_PORT"))))

}
