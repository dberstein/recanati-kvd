package kv

import (
	"fmt"
	"sync"
	"time"

	"github.com/dberstein/recanati-kvd/log"
)

type Record struct {
	value   []byte
	expires time.Time
}

func (r *Record) IsExpired() bool {
	if r.expires.IsZero() {
		return false
	}

	return time.Now().Sub(r.expires) > 0
}

type KV struct {
	sync.RWMutex
	values map[string]Record
	ticker *time.Ticker
	done   chan bool
}

func NewKV() *KV {
	kv := &KV{
		values: make(map[string]Record),
	}
	return kv
}

// Exists returns `true` for current keys otherwise `false`
func (kv *KV) Exists(key string) bool {
	value, ok := kv.values[key]
	if !ok {
		return false
	}
	return !value.IsExpired()
}

// Add `key`, `value` with `expiry` (expiry `time.Durtion(0)` means no expiry)
func (kv *KV) Add(key string, value []byte, expiry time.Duration) {
	kv.Lock()
	defer kv.Unlock()

	var expires time.Time
	if expiry == 0 {
		expires = time.Time{}
	} else {
		expires = time.Now().Add(expiry)
	}

	kv.values[key] = Record{
		value:   value,
		expires: expires,
	}
	log.Printf("\tadded key: %q (expires: %s)", key, expiry)
}

// Get value of `key` if still not expired
func (kv *KV) Get(key string) ([]byte, error) {
	kv.RLock()
	defer kv.RUnlock()

	value, ok := kv.values[key]
	if !ok {
		return value.value, fmt.Errorf("\tkey not found: %q", key)
	}
	if !value.expires.IsZero() && value.expires.Sub(time.Now()) < 0 {
		kv.delete(key)
		return value.value, fmt.Errorf("\tkey not found: %q", key)
	}
	log.Printf("\taccessed key: %q", key)

	return value.value, nil
}

// Delete removes cache for `key`
func (kv *KV) Delete(key string) {
	kv.Lock()
	defer kv.Unlock()

	kv.delete(key)
}

// Expire deletes expired keys
func (kv *KV) Expire() {
	kv.Lock()
	defer kv.Unlock()

	now := time.Now()
	for k, v := range kv.values {
		if v.expires.IsZero() {
			continue
		}
		if v.expires.After(now) {
			continue
		}
		kv.delete(k)
	}
}

// List returns list of non expired keys, their values and remaining expiry time
func (kv *KV) List() map[string]string {
	kv.RWMutex.Lock()
	defer kv.RWMutex.Unlock()

	// Copy from the original map to the target map
	targetMap := make(map[string]string)
	for key, value := range kv.values {
		var expires time.Duration
		if value.expires.IsZero() {
			expires = 0
		} else {
			expires = value.expires.Sub(time.Now())
		}

		if expires < 0 {
			kv.delete(key)
			continue
		}

		targetMap[key] = fmt.Sprintf("%s", expires)
	}

	return targetMap
}

func (kv *KV) delete(key string) {
	delete(kv.values, key)
	log.Printf("\tdeleted key: %q", key)
}

func (kv *KV) Start(freq time.Duration) {
	// expiry ticker
	kv.ticker = time.NewTicker(freq)
	kv.done = make(chan bool)

	// expire keys in go function and ticker...
	go func() {
		for {
			select {
			case <-kv.done:
				return
			case t := <-kv.ticker.C:
				log.Print("Tick at: ", t)
				kv.Expire()
			}
		}
	}()
}

func (kv *KV) Stop() {
	kv.ticker.Stop()
	kv.done <- true
}
