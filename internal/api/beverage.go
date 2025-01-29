package api

import (
	"coffy/internal/product"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetBeverages(service *product.Service) func(*gin.Context) {
	if service == nil {
		return func(c *gin.Context) {
			c.JSON(http.StatusInternalServerError, gin.H{})
		}
	}
	return func(c *gin.Context) {
		ids, err := service.ListAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusOK, ids)
	}
}

func CreateBeverage(service *product.Service) func(*gin.Context) {
	if service == nil {
		return func(c *gin.Context) {
			c.JSON(http.StatusInternalServerError, gin.H{})
		}
	}
	return func(c *gin.Context) {
		var json CreateBeverageRequest
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		b, err := service.Create(json.Name, json.Price)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		c.JSON(http.StatusCreated, Beverage{b.AggregateID, b.BeverageType, b.Price()})
	}
}

type CreateBeverageRequest struct {
	Name  string  `form:"name" json:"name" binding:"required"`
	Price float64 `form:"price" json:"price" binding:"required"`
}

type Beverage struct {
	ID   string  `json:"id"`
	Name string  `json:"name"`
	Cost float64 `json:"price"`
}
