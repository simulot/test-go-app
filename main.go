package main

import (
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func main() {
	// Initialize web application storage and state
	app.Route("/", &page{})
	app.RunWhenOnBrowser()

	http.Handle("/",
		&app.Handler{
			Name:        "Hello",
			Description: "An Hello World! example",
		})
	// Starting here, the server side
	if err := http.ListenAndServe(":8001", nil); err != nil {
		log.Fatal(err)
	}
}
