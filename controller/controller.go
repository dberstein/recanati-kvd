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

func (c *Controller) add(w http.ResponseWriter, r *http.Request, key string, value []byte) {
	var expiry time.Duration
	if !r.URL.Query().Has("expires") {
		c.Kv.Add(key, value, expiry)
	} else {
		expiryString := r.URL.Query().Get("expires") // from query string
		expiry, err := time.ParseDuration(expiryString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if expiry < 0 {
			http.Error(w, fmt.Sprintf("expiration must be positive: %s\n", expiry), http.StatusBadRequest)
			return
		}

		c.Kv.Add(key, value, expiry)
	}
	w.WriteHeader(http.StatusCreated)
}

// Add adds JSON payload with `key` and `value`
func (c *Controller) Add(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}{}

	body, err := io.ReadAll(r.Body) // from body
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key := payload.Key
	if key == "" {
		http.Error(w, "missing key", http.StatusBadRequest)
	}

	c.add(w, r, key, []byte(payload.Value))
}

// AddPath adds `key` from path with `value` being request body (ie. not JSON)
func (c *Controller) AddPath(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if key == "" {
		http.Error(w, "missing key", http.StatusBadRequest)
	}

	value, err := io.ReadAll(r.Body) // from body
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c.add(w, r, key, value)
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
	w.WriteHeader(http.StatusOK)
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
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
