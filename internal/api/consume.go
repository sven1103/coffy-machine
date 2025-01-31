package api

import (
	"coffy/internal/consume"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Consume applies an actual consume request to a user's account
//
//	@Summary		consume a coffee
//	@Schemes		http
//	@Description	Informs coffy about a user consumed a coffee.
//	@ID				consume-a-coffee
//	@Tags			consume
//	@Param			request	body	ConsumeRequest	true	"consume coffee request"
//	@Produce		json
//	@Success		200	{object}	consume.Receipt
//	@Router			/consume [post]
func Consume(s *consume.Service) func(c *gin.Context) {
	if s == nil {
		return func(c *gin.Context) {
			c.JSON(http.StatusInternalServerError, gin.H{})
		}
	}
	return func(c *gin.Context) {
		r := &ConsumeRequest{}
		if err := c.ShouldBind(r); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		receipt, err := s.Consume(r.AccountID, r.ProductID, r.Quantity)
		if err != nil {
			log.Println(err)
			switch {
			case errors.Is(err, consume.ErrorProductNotFound):
				c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			case errors.Is(err, consume.ErrorAccountNotFound):
				c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{})
			}
			return
		}
		c.JSON(http.StatusCreated, receipt)
	}
}

type ConsumeRequest struct {
	AccountID string `json:"account_id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}
