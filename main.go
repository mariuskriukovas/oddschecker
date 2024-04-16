package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-memdb"
	"oddschecker/controller"
	"oddschecker/model"
	"oddschecker/service"
)

func main() {
	db, err := memdb.NewMemDB(model.CreateSchema())
	if err != nil {
		panic(err)
	}

	userService := service.NewUserService(db)
	betService := service.NewBetService(db)
	oddsService := service.NewOddsService(db, userService, betService)
	oddsController := controller.NewOddsController(oddsService)

	router := gin.Default()

	oddsRoutes := router.Group("/odds")
	{
		oddsRoutes.POST("/", oddsController.PostOdds)
		oddsRoutes.GET("/:betId", oddsController.GetOddsByBetID)
	}

	err = router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
