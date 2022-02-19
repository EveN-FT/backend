package controllers

import (
	"fmt"
	"net/http"

	"github.com/EveN-FT/backend/db"
	"github.com/gin-gonic/gin"
)

// TicketController ...
type TicketController struct{}

// Redeem ticket
func (ctrl TicketController) Redeem(c *gin.Context) {
	var params struct {
		Address  string `json:"address" binding:"required"`
		TicketID uint64 `json:"ticketId" binding:"required"`
	}

	handle := func(err error) {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error redeeming ticket",
		})
	}

	err := c.ShouldBindJSON(&params)
	if err != nil {
		handle(err)
		return
	}

	// check who ticket owner is
	owns, err := db.CheckTicketAddress(c, params.Address, params.TicketID)
	if err != nil {
		handle(err)
		return
	}
	if !owns {
		handle(fmt.Errorf("You do not own this ticket"))
		return
	}

	// check if ticket is already redeemed by owner
	redeemed, err := db.CheckRedeem(c, params.TicketID)
	if err != nil {
		handle(err)
		return
	}
	if redeemed {
		handle(fmt.Errorf("Ticket already redeemed"))
		return
	}

	// redeem ticket
	err = db.Redeem(c, params.TicketID, params.Address)
	if err != nil {
		handle(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// Transfer ticket
func (ctrl TicketController) Transfer(c *gin.Context) {
	var params struct {
		NewAddress string `json:"newAddress" binding:"required"`
		OldAddress string `json:"oldAddress" binding:"required"`
		TicketID   uint64 `json:"ticketId" binding:"required"`
	}

	handle := func(err error) {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error transfering ticket",
		})
	}

	err := c.ShouldBindJSON(&params)
	if err != nil {
		handle(err)
		return
	}

	// check who ticket owner is
	owns, err := db.CheckTicketAddress(c, params.OldAddress, params.TicketID)
	if err != nil {
		handle(err)
		return
	}
	if !owns {
		handle(fmt.Errorf("You do not own this ticket"))
		return
	}

	// check if ticket is already redeemed by owner
	redeemed, err := db.CheckRedeem(c, params.TicketID)
	if err != nil {
		handle(err)
		return
	}
	if redeemed {
		handle(fmt.Errorf("Ticket already redeemed"))
		return
	}

	// transfer ticket
	err = db.Transfer(c, params.TicketID, params.OldAddress, params.NewAddress)
	if err != nil {
		handle(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// Create ticket redeem
func (ctrl TicketController) CreateRedeem(c *gin.Context) {
	var params struct {
		Address   string   `json:"address" binding:"required"`
		TicketIDs []uint64 `json:"ticketIds" binding:"required"`
	}

	handle := func(err error) {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error creating ticket redeem",
		})
	}

	err := c.ShouldBindJSON(&params)
	if err != nil {
		handle(err)
		return
	}

	// create redeem for ticket
	err = db.CreateRedeemForTicket(c, params.TicketIDs, params.Address)
	if err != nil {
		handle(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
