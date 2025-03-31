package crawlerapi

import (
	"fmt"
	"net/http"

	"github.com/ZakkBob/AskDave/backend/orm"
	"github.com/ZakkBob/AskDave/gocommon/tasks"
	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine, prefix string) {
	router.GET(prefix+"/tasks", GetTasks)
	router.POST(prefix+"/results", PostResults)
}

func GetTasks(c *gin.Context) {
	t, _ := orm.NextTasks(10)

	c.IndentedJSON(http.StatusOK, t)
}

func PostResults(c *gin.Context) {
	var results tasks.Results

	if err := c.BindJSON(&results); err != nil {
		return
	}

	err := orm.SaveResults(&results)

	if err != nil {
		fmt.Println(err)
		return
	}

	// fmt.Printf("%v", results)
}
