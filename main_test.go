package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dberstein/recanati-kvd/controller"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	assert := assert.New(t)
	controller := controller.NewController()
	router := setupRouter(controller)

	// missing key
	val, err := controller.Kv.Get("all")
	assert.Equal(err.Error(), "\tkey not found: \"all\"")
	assert.Equal([]uint8([]byte(nil)), val)

	// create key
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/store?expires=1s", strings.NewReader(`{"key":"all", "value":"123"}`))
	LoggerMiddleware(router).ServeHTTP(w, req)
	assert.Equal(201, w.Code)

	// key exists
	val, err = controller.Kv.Get("all")
	assert.Equal(err, nil)
	assert.Equal("123", string(val))

	assert.True(controller.Kv.Exists("all"))
	assert.False(controller.Kv.Exists("fake"))
}

func TestAddBroken(t *testing.T) {
	assert := assert.New(t)
	controller := controller.NewController()
	router := setupRouter(controller)

	// create key with broken JSON
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/store", strings.NewReader(`{"key":"all" "value":"123"}`))
	LoggerMiddleware(router).ServeHTTP(w, req)
	assert.Equal(400, w.Code)

	// create key with missing key
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/store", strings.NewReader(`{"Xkey":"all", "value":"123"}`))
	LoggerMiddleware(router).ServeHTTP(w, req)
	assert.Equal(400, w.Code)
}

func TestAddPath(t *testing.T) {
	assert := assert.New(t)
	controller := controller.NewController()
	router := setupRouter(controller)

	// missing key
	val, err := controller.Kv.Get("all")
	assert.Equal(err.Error(), "\tkey not found: \"all\"")
	assert.Equal([]uint8([]byte(nil)), val)

	// create key
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/store/all", strings.NewReader(`{"key":"some", "value":"123"}`))
	LoggerMiddleware(router).ServeHTTP(w, req)
	assert.Equal(201, w.Code)

	// key exists
	val, err = controller.Kv.Get("all")
	assert.Equal(err, nil)
	assert.Equal(`{"key":"some", "value":"123"}`, string(val))

	assert.True(controller.Kv.Exists("all"))
	assert.False(controller.Kv.Exists("fake"))
}

func TestGet(t *testing.T) {
	assert := assert.New(t)
	controller := controller.NewController()
	router := setupRouter(controller)

	// create key
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/store", strings.NewReader(`{"key":"all", "value":"123"}`))
	router.ServeHTTP(w, req)
	assert.Equal(201, w.Code)

	// request key
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/store/all", nil)
	LoggerMiddleware(router).ServeHTTP(w, req)
	assert.Equal(200, w.Code)
	assert.Equal("123", w.Body.String())
}

func TestDelete(t *testing.T) {
	assert := assert.New(t)
	controller := controller.NewController()
	router := setupRouter(controller)

	// create key
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/store", strings.NewReader(`{"key":"all", "value":"123"}`))
	LoggerMiddleware(router).ServeHTTP(w, req)
	assert.Equal(201, w.Code)

	assert.True(controller.Kv.Exists("all"))

	// delete key
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/store/all", nil)
	LoggerMiddleware(router).ServeHTTP(w, req)
	assert.Equal(200, w.Code)

	// key does not exists
	_, err := controller.Kv.Get("all")
	assert.Equal(err.Error(), "\tkey not found: \"all\"")

	assert.False(controller.Kv.Exists("all"))
}

func TestList(t *testing.T) {
	assert := assert.New(t)
	controller := controller.NewController()
	router := setupRouter(controller)

	// create keys
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/store", strings.NewReader(`{"key":"all1", "value":"123"}`))
	router.ServeHTTP(w, req)
	assert.Equal(201, w.Code)
	req, _ = http.NewRequest("POST", "/store", strings.NewReader(`{"key":"all2", "value":"321"}`))
	LoggerMiddleware(router).ServeHTTP(w, req)
	assert.Equal(201, w.Code)

	// list keys
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/store-all", nil)
	LoggerMiddleware(router).ServeHTTP(w, req)
	assert.Equal(200, w.Code)
	assert.Equal("{\"all1\":\"0s\",\"all2\":\"0s\"}", w.Body.String())

	assert.True(controller.Kv.Exists("all1"))
	assert.True(controller.Kv.Exists("all2"))
	assert.False(controller.Kv.Exists("all3"))
}

func TestExpire(t *testing.T) {
	assert := assert.New(t)
	controller := controller.NewController()
	router := setupRouter(controller)

	// todo: make background goroutine run
	controller.Kv.Start(500 * time.Millisecond)

	// create keys
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/store?expire=1ns", strings.NewReader(`{"key":"all", "value":"123"}`))
	LoggerMiddleware(router).ServeHTTP(w, req)
	assert.Equal(201, w.Code)

	controller.Kv.Expire()
	controller.Kv.Stop()

	list := controller.Kv.List()
	assert.Equal(1, len(list))
}
