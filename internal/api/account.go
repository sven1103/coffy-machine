package api

import (
	"coffy/internal/account"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

// GetAccounts returns all available account IDs.
//
//	@Summary		requests existing accounts
//	@Schemes		http
//	@ID				get-accounts
//	@Description	Request a list of all accounts.
//	@Tags			accounts
//	@Produce		json
//	@Success		200	{array}	AccountAlias
//	@Router			/accounts [get]
func GetAccounts(service *account.Accounting) func(*gin.Context) {
	if service == nil {
		return func(c *gin.Context) {
			log.Println(errors.New("service is nil"))
			c.JSON(http.StatusInternalServerError, gin.H{})
		}
	}
	return func(c *gin.Context) {
		accounts, err := service.ListAll()
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		alias := make([]AccountAlias, 0)
		for _, a := range accounts {
			entry, err := convertAccount(&a)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{})
			}
			alias = append(alias, entry)
		}
		c.JSON(http.StatusOK, alias)
	}
}

// GetAccountById returns the account associated with the provided ID.
//
//	@Summary		access account info by ID
//	@Schemes		http
//	@Description	Request account by ID.
//	@ID				request-account-by-id
//	@Tags			accounts
//	@Param			id	path	string	true	"account ID"
//	@Produce		json
//	@Success		200	{object}	AccountAlias
//	@Router			/accounts/{id} [get]
func GetAccountById(service *account.Accounting) func(*gin.Context) {
	if service == nil {
		return func(c *gin.Context) {
			log.Println(errors.New("service is nil"))
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

// CreateAccount creates a new account with a unique ID.
//
//	@Summary		create an account
//	@Schemes		http
//	@Description	Creates a new account in coffy.
//	@ID				create-new-account
//	@Tags			accounts
//	@Param			id	body	AccountCreationRequest	true	"account creation request"
//	@Produce		json
//	@Success		200	{object}	AccountAlias
//	@Router			/accounts [post]
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
	return AccountAlias{
		ID:            a.ID(),
		Owner:         a.Owner(),
		Balance:       a.Balance(),
		ConsumedTotal: a.ConsumedTotal()}, nil
}

type AccountCreatedResponse struct {
	ID    string `json:"id"`
	Owner string `json:"owner"`
}
