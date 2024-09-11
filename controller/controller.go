package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dberstein/recanati-kvd/kv"
)

type Controller struct {
	Kv kv.KV
}

func NewController() *Controller {
	return &Controller{
		Kv: *kv.NewKV(),
	}
}

// Add adds JSON payload with `key` and `value`
func (c *Controller) Add(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Key   string `json:"key"`
		Value []byte `json:"value"`
	}{}

	body, err := io.ReadAll(r.Body) // from body
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	key := payload.Key
	if key == "" {
		http.Error(w, "missing key", http.StatusBadRequest)
	}

	var expiry time.Duration
	if r.URL.Query().Has("expires") {
		expiryString := r.URL.Query().Get("expires") // from query string
		expiry, err = time.ParseDuration(expiryString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	c.Kv.Add(key, payload.Value, expiry)
}

// AddPath adds `key` from path with `value` being request body (ie. not JSON)
func (c *Controller) AddPath(w http.ResponseWriter, r *http.Request) {
	value, err := io.ReadAll(r.Body) // from body
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key := r.PathValue("key")
	if key == "" {
		http.Error(w, "missing key", http.StatusBadRequest)
	}

	var expiry time.Duration
	if r.URL.Query().Has("expires") {
		expiryString := r.URL.Query().Get("expires") // from query string
		expiry, err = time.ParseDuration(expiryString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	c.Kv.Add(key, value, expiry)
}

// Get retrieves `key` from path
func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	value, err := c.Kv.Get(key)
	if err != nil {
		http.Error(w, fmt.Sprintf("key not found: %q", key), http.StatusNotFound)
		return
	}
	w.Write([]byte(value))
}

// Delete deletes `key` from path
func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	c.Kv.Delete(key)
}

// List lists active items
func (c *Controller) List(w http.ResponseWriter, r *http.Request) {
	list := c.Kv.List()

	res, err := json.Marshal(list)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Write(res)
}
