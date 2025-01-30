package api

import (
	"coffy/internal/account"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func GetAccounts(service *account.Accounting) func(*gin.Context) {
	if service == nil {
		return func(c *gin.Context) {
			c.JSON(http.StatusServiceUnavailable, gin.H{})
		}
	}
	return func(c *gin.Context) {
		ids, err := service.ListAll()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{})
			return
		}

		c.JSON(http.StatusOK, ids)
	}
}

func GetAccountById(service *account.Accounting) func(*gin.Context) {
	if service == nil {
		return func(c *gin.Context) {
			c.JSON(http.StatusServiceUnavailable, gin.H{})
		}
	}
	return func(c *gin.Context) {
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
	}
}

func CreateAccount(service *account.Accounting) func(*gin.Context) {
	if service == nil {
		return func(c *gin.Context) {
			c.JSON(http.StatusServiceUnavailable, gin.H{})
		}
	}
	return func(c *gin.Context) {
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
	}
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
