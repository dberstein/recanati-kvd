package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
	req, _ := http.NewRequest("POST", "/store", strings.NewReader(`{"key":"all", "value":"123"}`))
	router.ServeHTTP(w, req)
	assert.Equal(201, w.Code)

	// key exists
	val, err = controller.Kv.Get("all")
	assert.Equal(err, nil)
	assert.Equal("123", string(val))
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
	router.ServeHTTP(w, req)
	assert.Equal(200, w.Code)
	body, _ := w.Body.ReadBytes('\n')
	assert.Equal([]uint8([]byte{0x31, 0x32, 0x33}), body)
}

func TestDelete(t *testing.T) {
	assert := assert.New(t)
	controller := controller.NewController()
	router := setupRouter(controller)

	// delete key
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/store/all", nil)
	router.ServeHTTP(w, req)
	assert.Equal(200, w.Code)

	// key does not exists
	_, err := controller.Kv.Get("all")
	assert.Equal(err.Error(), "\tkey not found: \"all\"")
}

// func TestList(t *testing.T) {
// 	assert := assert.New(t)
// 	controller := controller.NewController()
// 	router := setupRouter(controller)

// }
