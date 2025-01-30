package main

import (
	"coffy/internal/account"
	"coffy/internal/api"
	"coffy/internal/cmd"
	"coffy/internal/coffy"
	"coffy/internal/product"
	"coffy/internal/storage"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

var version string = "1.0.0"

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

	// init the event repo
	repo, err := storage.CreateEventRepository("test.db")
	if err != nil {
		log.Fatal(err)
	}

	// create app services first
	accService := account.NewAccounting(&repo)
	beverageService := product.NewService(&repo)

	router := gin.Default()

	// beverages API
	router.GET("/beverages", api.GetBeverages(beverageService))
	router.POST("/beverages", api.CreateBeverage(beverageService))

	// accounts API
	const pathAccounts = "/accounts"
	router.GET(pathAccounts, api.GetAccounts(accService))
	router.GET(pathAccounts+"/:id", api.GetAccountById(accService))
	router.POST(pathAccounts, api.CreateAccount(accService))

	// run the server
	err = router.Run(fmt.Sprintf(":%d", config.Server.Port))
	if err != nil {
		log.Fatal(err)
	}
}

func logStartup() {
	log.Printf("Starting Coffy server (version: %s) ...", version)
}
