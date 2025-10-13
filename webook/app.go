package main

import (
	"Book_Exp/webook/internal/events"

	"github.com/gin-gonic/gin"
)

type App struct {
	server   *gin.Engine
	consumer []events.Consumer
}
