package api

import (
	"coffy/internal/equipment"
	"coffy/internal/product"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type MachineAlias struct {
	MachineID string `json:"id"`
	Brand     string `json:"brand"`
	Model     string `json:"model"`
	CoffeeID  string `json:"coffee_id"`
}

// GetMachines lists all available machines in coffy.
//
//	@Summary		list all machines
//	@Schemes		http
//	@Description	Lists all available machines in coffy.
//	@ID				list-machines
//	@Tags			machines
//	@Produce		json
//	@Success		200	{array}	MachineAlias
//	@Router			/machines [get]
func GetMachines(service *equipment.Service) func(*gin.Context) {
	if service == nil {
		return func(c *gin.Context) {
			log.Println(errors.New("equipment service is nil"))
			c.JSON(http.StatusInternalServerError, gin.H{})
		}
	}
	return func(c *gin.Context) {
		machines, err := service.ListAll()
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{})
		}
		alias := make([]MachineAlias, 0)
		for _, m := range machines {
			alias = append(alias, toAlias(m))
		}
		c.JSON(http.StatusOK, alias)
	}
}

// CreateMachine creates a new coffee machine entry with a unique ID.
//
//	@Summary		creates a machine
//	@Schemes		http
//	@Description	Creates a new machine in coffy.
//	@ID				create-new-machine
//	@Tags			machines
//	@Param			id	body	MachineCreationRequest	true	"machine creation request"
//	@Produce		json
//	@Success		201	{object}	MachineAlias
//	@Router			/machines [post]
func CreateMachine(service *equipment.Service) func(*gin.Context) {
	if service == nil {
		return func(c *gin.Context) {
			log.Println(errors.New("equipment service is nil"))
			c.JSON(http.StatusInternalServerError, gin.H{})
		}
	}
	return func(c *gin.Context) {
		var request MachineCreationRequest

		if err := c.BindJSON(&request); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}

		m, err := service.Create(request.Brand, request.Model)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{})
		}
		alias := toAlias(*m)
		c.JSON(http.StatusCreated, alias)
	}
}

// PatchMachines creates a new coffee machine entry with a unique ID.
//
//	@Summary		changes the loaded coffee
//	@Schemes		http
//	@Description	Changes the coffee currently loaded in the machine.
//	@ID				change-machine-coffee
//	@Tags			machines
//	@Param			request	body	LoadCoffeeRequest	true	"machine loading request"
//	@Param			id		path	string				true	"machine ID"
//	@Produce		json
//	@Success		200	{object}	MachineAlias
//	@Failure		404	{ object }	map[string]string
//	@Router			/machines/{id} [patch]
func PatchMachines(service *equipment.Service, coffeeService *product.Service) func(*gin.Context) {
	if service == nil || coffeeService == nil {
		return func(c *gin.Context) {
			log.Println(errors.New("machine service or coffee service is nil"))
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
	}
	return func(c *gin.Context) {
		machineId := c.Param("id")
		var request = LoadCoffeeRequest{}
		if err := c.BindJSON(&request); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}
		coffee, err := coffeeService.Find(request.CoffeeID)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusNotFound, "unknown resource: coffee not found")
			return
		}
		machine, err := service.LoadCoffee(machineId, coffee.AggregateID)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		c.JSON(http.StatusOK, toAlias(*machine))
	}
}

func toAlias(m equipment.Machine) MachineAlias {
	coffeeId, err := m.Coffee()
	if err != nil {
		coffeeId = ""
	}
	return MachineAlias{
		MachineID: m.AggregateID,
		Brand:     m.Brand,
		Model:     m.Model,
		CoffeeID:  coffeeId,
	}
}

type MachineCreationRequest struct {
	Model string `json:"model"`
	Brand string `json:"brand"`
}

type LoadCoffeeRequest struct {
	CoffeeID string `json:"coffee_id"`
}
