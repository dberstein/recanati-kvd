package main

import (
	"flag"
	"net"
	"net/http"
	"time"

	// "github.com/patrickmn/go-cache"

	"github.com/dberstein/recanati-kvd/controller"
	"github.com/dberstein/recanati-kvd/log"
	"github.com/dberstein/recanati-kvd/rw"
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

			rw := rw.New(w)
			next.ServeHTTP(rw, r)

			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				panic(err)
			}

			log.Print(r.Method, r.URL, ip, rw.StatusCode, time.Now().Sub(start), r.Header["User-Agent"])
		},
	)
}

func main() {
	freqString := flag.String("f", "5m", "cleanup frequency")
	addrString := flag.String("l", ":8080", "listen address")

	flag.Parse()

	freqDuration, err := time.ParseDuration(*freqString)
	if err != nil {
		panic(err)
	}

	log.Printf("Listening address %q\n", *addrString)
	log.Printf("Cleanup frequency %q\n", freqDuration)

	controller := controller.NewController()
	router := setupRouter(controller)

	controller.Kv.Start(freqDuration)
	chain := LoggerMiddleware(router)
	if err := http.ListenAndServe(*addrString, chain); err != nil {
		controller.Kv.Stop()
		log.Fatal(err)
	}
}
