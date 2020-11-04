package main

import (
	"fmt"

	"github.com/libonomy/node-extract/routes"
)

func main() {
	fmt.Println("AI  training Server Is Running...")
	routes.StartRoutes()
}
