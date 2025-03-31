package api

import (
	"github.com/ZakkBob/AskDave/backend/api/crawlerapi"
	"github.com/ZakkBob/AskDave/backend/api/searchapi.go"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func Init() {
	router = gin.Default()
	crawlerapi.Init(router, "/api/crawler")
	searchapi.Init(router, "/api/search")
}

func Run() {
	router.Run() // listen and serve on 0.0.0.0:8080
}
