package main

import (
	"coffy/internal/account"
	"coffy/internal/api"
	"coffy/internal/product"
	"coffy/internal/storage"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func main() {
	repo, err := storage.CreateEventRepository("test.db")
	if err != nil {
		log.Fatal(err)
	}
	service := account.NewAccounting(&repo)
	beverageService := product.NewService(&repo)

	router := gin.Default()

	setupRoutes(router, service)

	router.GET("/beverages", api.GetBeverages(beverageService))

	router.POST("/beverages", api.CreateBeverage(beverageService))

	router.Run(":8088")
}

func setupRoutes(router *gin.Engine, service *account.Accounting) {

	router.GET("/accounts", func(c *gin.Context) {
		result, err := service.ListAll()
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{})
		}
		c.IndentedJSON(http.StatusOK, result)
	})

	router.GET("/accounts/:id", func(c *gin.Context) {
		id := c.Param("id")
		result, err := service.Find(id)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		alias, err := convertAccount(result)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		if alias.ID == "" {
			c.JSON(http.StatusNotFound, gin.H{})
			return
		}
		c.IndentedJSON(http.StatusOK, alias)
	})

	router.POST("/accounts", func(c *gin.Context) {
		var request AccountCreationRequest

		if err := c.BindJSON(&request); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}

		if strings.TrimSpace(request.Owner) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "owner required"})
			return
		}

		acc, err := service.Create(request.Owner)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		alias, err := convertAccount(acc)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusCreated, alias)
	})
}

type AccountAlias struct {
	ID            string  `json:"id"`
	Owner         string  `json:"owner"`
	Balance       float64 `json:"balance"`
	ConsumedTotal int     `json:"consumedTotal"`
}

type AccountCreationRequest struct {
	Owner string `json:"owner"`
}

func convertAccount(a *account.Account) (AccountAlias, error) {
	return AccountAlias{ID: a.ID(), Owner: a.Owner(), Balance: a.Balance(), ConsumedTotal: a.ConsumedTotal()}, nil
}

type AccountCreatedResponse struct {
	ID    string `json:"id"`
	Owner string `json:"owner"`
}
