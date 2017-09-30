package cache

import (
	"net/http"
)

// Keys address values in a Cache.
//
// Caches takes Keys instead of strings to encourage compile time safety.
type Key interface {
	String() string
}

type Cache interface {

	// Put sets a value in the cache.
	Put(key Key, value []byte) error

	// GetBytes does a cache lookup and, if necessary, fill.
	//
	// If the cache has the key, the cached value will be returned.
	// Otherwise, filler will be called to fill the cache.
	// If it succeeds, its result will be put in the cache and returned.
	// Otherwise, an error will be returned.
	GetBytes(key Key, filler func() ([]byte, error)) ([]byte, error)

	// GetHandler does a cache lookup and, if necessary, fill.
	//
	// If the cache has the key, an http.Handler serving the value will be returned.
	// Otherwise, filler will be called to fill the cache.
	// If it succeeds, its result will be put in the cache and returned.
	// Otherwise, an error will be returned.
	//
	// The behavior of the http.Handler if called more than once is undefined.
	// Specific implementations may document their own behavior.
	GetHandler(key Key, filler func() ([]byte, error), contentType string) (http.Handler, error)

	// Close terminates open connections.
	//
	// The behavior of Close after the first call is undefined.
	// Specific implementations may document their own behavior.
	Close()
}
