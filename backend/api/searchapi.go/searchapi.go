package searchapi

import (
	"github.com/gin-gonic/gin"
)

type Search struct {
	query string
}

func Init(router *gin.Engine, prefix string) {
	router.GET(prefix+"/search", GetSearchResults)
}

func GetSearchResults(c *gin.Context) {

}
