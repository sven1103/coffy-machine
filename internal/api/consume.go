package api

import (
	"coffy/internal/consume"
	"github.com/gin-gonic/gin"
	"net/http"
)

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
		}
		receipt, err := s.Consume(r.AccountID, r.AccountID, r.Quantity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
		}
		c.JSON(http.StatusCreated, receipt)
	}
}

type ConsumeRequest struct {
	AccountID string `json:"account_id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}
