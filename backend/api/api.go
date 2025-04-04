package api

import (
	"github.com/ZakkBob/AskDave/backend/api/crawlerapi"
	"github.com/ZakkBob/AskDave/backend/api/searchapi.go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func Init() {
	router = gin.Default()
	router.Use(cors.Default())
	crawlerapi.Init(router, "/api/v1/crawler")
	searchapi.Init(router, "/api/v1/search")
}

func Run() {
	router.Run() // listen and serve on 0.0.0.0:8080
}
