package api

import (
	"coffy/internal/product"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// GetCoffees returns all available coffees in coffy.
//
//	@Summary		get all coffees
//	@Schemes		http
//	@Description	Lists all available coffees in coffy.
//	@ID				get-all-coffees
//	@Tags			coffees
//	@Produce		json
//	@Success		200	{array}	CoffeeInfo
//	@Router			/coffees [get]
func GetCoffees(service *product.Service) func(*gin.Context) {
	if service == nil {
		return func(c *gin.Context) {
			c.JSON(http.StatusInternalServerError, gin.H{})
		}
	}
	return func(c *gin.Context) {
		bev, err := service.ListAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		list, err := allToCoffeeInfo(bev)
		if err != nil {
			log.Println("conversion to beverage info failed:", err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusOK, list)
	}
}

// CreateBeverage creates a new coffee in coffy with an initial price.
//
//	@Summary		create new coffee
//	@Schemes		http
//	@Description	Creates a new coffee in coffy.
//	@ID				create-new-coffee
//	@Tags			coffees
//	@Param			request	body	CreateBeverageRequest	true	"coffee creation request"
//	@Produce		json
//	@Success		200	{object}	CoffeeInfo
//	@Router			/coffees [post]
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
		c.JSON(http.StatusCreated, CoffeeInfo{b.AggregateID, b.Type, b.Price()})
	}
}

type CreateBeverageRequest struct {
	Name  string  `form:"name" json:"name" binding:"required"`
	Price float64 `form:"price" json:"price" binding:"required"`
}

type CoffeeInfo struct {
	ID   string  `json:"id"`
	Name string  `json:"name"`
	Cost float64 `json:"price"`
}

func allToCoffeeInfo(list []product.Coffee) ([]CoffeeInfo, error) {
	if list == nil {
		return nil, errors.New("beverage conversion failed, list was nil")
	}
	r := make([]CoffeeInfo, 0)
	for _, v := range list {
		info, err := toCoffeeInfo(&v)
		if err != nil {
			return nil, errors.Join(errors.New("could not convert to beverage info"), err)
		}
		r = append(r, info)
	}
	return r, nil
}

func toCoffeeInfo(b *product.Coffee) (CoffeeInfo, error) {
	if b == nil {
		return CoffeeInfo{}, errors.New("beverage is nil")
	}
	return CoffeeInfo{ID: b.AggregateID, Name: b.Type, Cost: b.Price()}, nil
}
