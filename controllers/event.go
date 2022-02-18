package controllers

import (
	"fmt"
	"net/http"

	"github.com/EveN-FT/backend/db"
	"github.com/EveN-FT/backend/models"
	"github.com/gin-gonic/gin"
)

// EventController ...
type EventController struct{}

// Create an event
func (ctrl EventController) Create(c *gin.Context) {
	var params struct {
		Address      string `json:"address" binding:"required"`
		OwnerAddress string `json:"ownerAddress" binding:"required"`
	}

	handle := func(err error) {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error creating event",
		})
		return
	}

	err := c.ShouldBindJSON(&params)
	if err != nil {
		handle(err)
		return
	}

	// store set data in db
	id, err := db.CreateEvent(c, &models.Event{
		Address:      params.Address,
		OwnerAddress: params.OwnerAddress,
	})
	if err != nil {
		handle(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"eventId": id,
	})
}

// List events
func (ctrl EventController) ListEvents(c *gin.Context) {
	handle := func(err error) {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error creating event",
		})
		return
	}

	// list events
	events, err := db.ListEvents(c)
	if err != nil {
		handle(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"events": events,
	})
}

// List events by owner
func (ctrl EventController) ListEventsByOwner(c *gin.Context) {
	var params struct {
		OwnerAddress string `json:"ownerAddress" binding:"required"`
	}

	handle := func(err error) {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error creating event",
		})
		return
	}

	err := c.ShouldBindJSON(&params)
	if err != nil {
		handle(err)
		return
	}

	// list events
	events, err := db.ListEventsByOwnerAddress(c, params.OwnerAddress)
	if err != nil {
		handle(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"events": events,
	})
}
