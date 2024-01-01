package main

import (
	"github.com/syndtr/goleveldb/leveldb"
	meals "meal-planning"
	"net"
	"net/http"
)

func main() {
	db, err := leveldb.OpenFile("level.db", nil)
	if err != nil {
		panic(err)
	}

	handler, err := meals.NewMealPlanningHandler(db)
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
