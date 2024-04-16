package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"oddschecker/payloads"
	"oddschecker/service"
	"strconv"
)

type OddsController struct {
	oddsService *service.OddsService
}

func NewOddsController(oddsService *service.OddsService) *OddsController {
	return &OddsController{
		oddsService: oddsService,
	}
}

func (oc *OddsController) PostOdds(ctx *gin.Context) {
	var payload payloads.OddsPayload
	if err := ctx.BindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format !"})
		return
	}

	status, msg := oc.oddsService.InsertOdds(payload)

	ctx.JSON(status, gin.H{
		"message": msg,
	})
}

func (oc *OddsController) GetOddsByBetID(ctx *gin.Context) {
	betId := ctx.Param("betId")
	castedBetId, err := strconv.ParseInt(betId, 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Bet ID supplied !"})
		return
	}

	odds, err := oc.oddsService.GetOddsByBetID(castedBetId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, odds)
}
