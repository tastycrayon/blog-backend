package main

import (
	"log"

	"github.com/pocketbase/pocketbase"
)

func main() {

	var foo string = "bar"

	println([]byte(foo))
	app := pocketbase.New()

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
