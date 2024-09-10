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
	if time.Now().Sub(r.expires) > 0 {
		return true
	}
	return false
}

type ListRecord struct {
	Value   []byte `json:"value"`
	Expires string `json:"expires"`
}

type KV struct {
	sync.Mutex
	values map[string]Record
}

func NewKV() *KV {
	kv := &KV{
		values: make(map[string]Record),
	}
	return kv
}

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
	log.Print("added key: ", key, ":", expiry)
}

func (kv *KV) Get(key string) ([]byte, error) {
	kv.Lock()
	defer kv.Unlock()

	value, ok := kv.values[key]
	if !ok {
		return value.value, fmt.Errorf("key not found: %q", key)
	}
	if !value.expires.IsZero() && value.expires.Sub(time.Now()) < 0 {
		delete(kv.values, key)
		log.Print("deleted key: ", key)
		return value.value, fmt.Errorf("key not found: %q", key)
	}

	return value.value, nil
}

func (kv *KV) Delete(key string) {
	kv.Lock()
	defer kv.Unlock()

	delete(kv.values, key)
	log.Print("deleted key: ", key)
}

func (kv *KV) Expire() {
	now := time.Now()
	for k, v := range kv.values {
		if v.expires.IsZero() {
			continue
		}
		if v.expires.Before(now) {
			continue
		}
		kv.Delete(k)
	}
}

func (kv *KV) List() map[string]ListRecord {
	kv.Mutex.Lock()
	defer kv.Mutex.Unlock()

	// Copy from the original map to the target map
	targetMap := make(map[string]ListRecord)
	for key, value := range kv.values {
		var expires time.Duration
		if value.expires.IsZero() {
			expires = 0
		} else {
			expires = value.expires.Sub(time.Now())
		}

		if expires < 0 {
			delete(kv.values, key)
			log.Print("deleted key: ", key)
			continue
		}

		targetMap[key] = ListRecord{
			Value:   value.value,
			Expires: fmt.Sprintf("%s", expires),
		}
	}

	return targetMap
}
