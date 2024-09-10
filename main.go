package main

import (
	"log"
	"net/http"

	// "github.com/patrickmn/go-cache"

	"github.com/dberstein/recanati-kvd/controller"
)

func main() {
	router := http.NewServeMux()
	controller := controller.NewController()

	router.HandleFunc("POST /store", controller.Add)
	router.HandleFunc("POST /store/{key}", controller.AddPath)
	router.HandleFunc("GET /store/{key}", controller.Get)
	router.HandleFunc("DELETE /store/{key}", controller.Delete)
	router.HandleFunc("GET /store-all", controller.List)

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
