package main

import (
	meals "meal-planning"
	"net"
	"net/http"
)

func main() {
	handler, err := meals.NewMealPlanningHandler()
	if err != nil {
		panic(err)
	}

	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		panic(err)
	}

	err = http.Serve(listener, handler)
	if err != nil {
		panic(err)
	}
}
