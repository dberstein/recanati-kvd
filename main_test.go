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
