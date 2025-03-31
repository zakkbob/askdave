package api

import (
	"github.com/ZakkBob/AskDave/backend/api/crawlerapi"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func Init() {
	router = gin.Default()
	crawlerapi.Init(router, "/api/crawler")
}

func Run() {
	router.Run() // listen and serve on 0.0.0.0:8080
}
