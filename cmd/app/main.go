package main

import (
	"log"
	"project/internal/app"
)

func main() {
	log.Println("Application start up!")
	a := app.New()
	log.Println("Application created")
	a.StartServer()
	log.Println("Application terminated!")
}
