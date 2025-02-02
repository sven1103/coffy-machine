package main

import (
	"coffy/docs"
	"coffy/internal/account"
	"coffy/internal/api"
	"coffy/internal/cmd"
	"coffy/internal/coffy"
	"coffy/internal/consume"
	"coffy/internal/product"
	"coffy/internal/storage"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"os"
)

// Holds the current version value of the application.
var version string = "1.0.0"

// @title			coffy-server API
// @version		1.0
// @description	This is a description of the coffy-server API capabilities.
// @host			localhost:8080
// @BasePath		/api/v1
func main() {
	logStartup()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Something went wrong:", err)
			os.Exit(1)
		}
		os.Exit(0)
	}()
	cmd.Execute(startCoffy)
}

func startCoffy(config *coffy.Config) {
	log.Println("Received app configuration")
	docs.SwaggerInfo.BasePath = "/api/v1"

	// init the event repo
	repo, err := storage.CreateEventRepository(config.Database.Path)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database in use:", config.Database.Path)

	// create app services first
	accService := account.NewAccounting(&repo)
	beverageService := product.NewService(&repo)
	consumeService := consume.NewService(accService, beverageService)

	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		// accounts API
		const pathAccounts = "/accounts"
		v1.GET(pathAccounts, api.GetAccounts(accService))
		v1.GET(pathAccounts+"/:id", api.GetAccountById(accService))
		v1.POST(pathAccounts, api.CreateAccount(accService))

		// beverages API
		v1.GET("/coffees", api.GetCoffees(beverageService))
		v1.POST("/coffees", api.CreateBeverage(beverageService))

		// consume API
		v1.POST("/consume", api.Consume(consumeService))
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// run the server
	err = router.Run(fmt.Sprintf(":%d", config.Server.Port))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Coffy Machine is running and listening on port", config.Server.Port)
}

func logStartup() {
	log.Printf("Starting Coffy server (version: %s) ...", version)
}
