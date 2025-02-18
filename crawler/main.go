package main

import (
	"ZakkBob/AskDave/crawler/daveapiclient"
)

func main() {
	c := daveapiclient.Create("http://localhost:3000")

	c.Run()
}
