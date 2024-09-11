package main

import (
	"net"
	"net/http"
	"time"

	// "github.com/patrickmn/go-cache"

	"github.com/dberstein/recanati-kvd/controller"
	"github.com/dberstein/recanati-kvd/log"
)

func setupRouter(controller *controller.Controller) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("POST /store", controller.Add)
	router.HandleFunc("POST /store/{key}", controller.AddPath)
	router.HandleFunc("GET /store/{key}", controller.Get)
	router.HandleFunc("DELETE /store/{key}", controller.Delete)
	router.HandleFunc("GET /store-all", controller.List)

	return router
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r)

			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				panic(err)
			}

			log.Print(r.Method, r.URL, ip, time.Now().Sub(start))
		},
	)
}

func main() {
	controller := controller.NewController()
	router := setupRouter(controller)

	ticker := time.NewTicker(5 * time.Second)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				log.Print("Tick at", t)
				controller.Kv.Expire()
			}
		}
	}()

	chain := LoggerMiddleware(router)
	if err := http.ListenAndServe(":8080", chain); err != nil {
		ticker.Stop()
		done <- true
		log.Fatal(err)
	}
}
