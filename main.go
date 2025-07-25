package main

import (
	"circle-center/globals"
	editor "circle-center/processor"
	"circle-center/reader"
)

func main() {
	r := globals.SetupRouter()

	// register application specific routes
	v1 := r.Group("/v1")
	reader.RegisterRoutes(v1)
	editor.RegisterRoutes(v1)

	// Start the HTTP server on port 8080.
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
