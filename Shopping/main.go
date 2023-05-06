package main

import (
	"Shopping/Routes"
)

func main() {
	router := Routes.SetupRouter()
	router.Run()
}
