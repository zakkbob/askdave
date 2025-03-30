package main

import (
	"github.com/ZakkBob/AskDave/backend/crawlerapi"
	"github.com/ZakkBob/AskDave/backend/orm"
)

func main() {
	orm.Connect("")
	crawlerapi.Init()
	defer orm.Close()
}
