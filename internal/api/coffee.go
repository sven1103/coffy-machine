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
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		list, err := allToCoffeeInfo(bev)
		if err != nil {
			log.Println("conversion to coffee info failed:", err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusOK, list)
	}
}

// PatchCoffeeDetails updates the information about a coffee
//
//	@Summary		changes the details for a coffee
//	@Schemes		http
//	@Description	Changes the coffee details information.
//	@ID				change-details-coffee
//	@Tags			coffees
//	@Param			request	body	DetailsUpdateRequest	true	"coffee details update request"
//	@Param			id		path	string					true	"coffee ID"
//	@Produce		json
//	@Success		200	{object}	CoffeeInfo
//	@Failure		404	{ object }	map[string]string
//	@Router			/coffees/{id}/info [patch]
func PatchCoffeeDetails(service *product.Service) func(*gin.Context) {
	if service == nil {
		return func(c *gin.Context) {
			c.JSON(http.StatusInternalServerError, gin.H{})
		}
	}
	return func(c *gin.Context) {
		coffeeID := c.Param("id")
		coffee, err := service.Find(coffeeID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{})
		}
		details := product.Details{}
		if err := c.ShouldBindJSON(&details); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
		}
		err = coffee.UpdateDetails(details)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
		}
		info, err := toCoffeeInfo(coffee)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
		}
		c.JSON(http.StatusOK, info)
	}
}

// PatchCoffeePrice updates the price for a coffee
//
//	@Summary		changes the price for a coffee
//	@Schemes		http
//	@Description	Changes the coffee price.
//	@ID				change-price-coffee
//	@Tags			coffees
//	@Param			request	body	PriceUpdateRequest	true	"coffee price change request"
//	@Param			id		path	string				true	"coffee ID"
//	@Produce		json
//	@Success		200	{object}	CoffeeInfo
//	@Failure		404	{ object }	map[string]string
//	@Router			/coffees/{id}/price [patch]
func PatchCoffeePrice(service *product.Service) func(*gin.Context) {
	if service == nil {
		return func(c *gin.Context) {
			c.JSON(http.StatusInternalServerError, gin.H{})
		}
	}
	return func(c *gin.Context) {
		coffeeID := c.Param("id")
		coffee, err := service.Find(coffeeID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{})
			return
		}
		request := PriceUpdateRequest{}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}
		err = coffee.ChangePrice(request.Price, request.Reason)
		if err != nil {
			if err.Error() == "invalid price" {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		info, err := toCoffeeInfo(coffee)
		c.JSON(http.StatusOK, info)
	}
}

// CreateCoffee creates a new coffee in coffy with an initial price.
//
//	@Summary		create new coffee
//	@Schemes		http
//	@Description	Creates a new coffee in coffy.
//	@ID				create-new-coffee
//	@Tags			coffees
//	@Param			request	body	CreateCoffeeRequest	true	"coffee creation request"
//	@Produce		json
//	@Success		201	{object}	CoffeeInfo
//	@Router			/coffees [post]
func CreateCoffee(service *product.Service) func(*gin.Context) {
	if service == nil {
		return func(c *gin.Context) {
			c.JSON(http.StatusInternalServerError, gin.H{})
		}
	}
	return func(c *gin.Context) {
		var json CreateCoffeeRequest
		if err := c.ShouldBindJSON(&json); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		b, err := service.Create(json.Name, json.Price, json.CuppingScore, &json.Details)
		if err != nil {
			log.Println(err)
			switch {
			case errors.Is(err, product.InvalidPropertyError):
				c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create coffee: " + err.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{})
			}
			return
		}

		info := CoffeeInfo{ID: b.AggregateID, Name: b.Type, Price: b.Price(), CuppingScore: b.CoffeeValue().Value, Details: toCoffeeDetails(b.Details())}

		c.JSON(http.StatusCreated, info)
	}
}

type CreateCoffeeRequest struct {
	Name         string                `form:"name" json:"name" binding:"required"`
	Price        float64               `form:"price" json:"price" binding:"required"`
	CuppingScore *int                  `form:"cupping_score" json:"cupping_score,omitempty"`
	Details      product.CoffeeDetails `form:"info" json:"info"`
}

type PriceUpdateRequest struct {
	Price  float64 `form:"price" json:"price" binding:"required"`
	Reason string  `form:"reason" json:"reason" binding:"required"`
}

type DetailsUpdateRequest struct {
	Details product.CoffeeDetails `form:"info" json:"info"`
}

type CoffeeInfo struct {
	ID           string                `json:"id"`
	Name         string                `json:"name"`
	Price        float64               `json:"price"`
	CuppingScore int                   `json:"cupping_score"`
	Details      product.CoffeeDetails `json:"info"`
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
	d := toCoffeeDetails(b.Details())
	return CoffeeInfo{ID: b.AggregateID, Name: b.Type, Price: b.Price(), CuppingScore: b.CoffeeValue().Value, Details: d}, nil
}

func toCoffeeDetails(d product.Details) product.CoffeeDetails {
	return product.CoffeeDetails{Origin: d.Origin, Description: d.Description, Misc: d.Misc}
}
