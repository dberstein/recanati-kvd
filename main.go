package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

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

	ticker := time.NewTicker(5 * time.Minute)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
				controller.Kv.Expire()
			}
		}
	}()

	if err := http.ListenAndServe(":8080", router); err != nil {
		ticker.Stop()
		done <- true
		log.Fatal(err)
	}
}
