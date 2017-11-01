package cache

import (
	"net/http"
)

// Key address a value in a Cache.
//
// Caches takes Keys instead of strings to encourage compile time safety.
type Key interface {
	String() string
}

// Version is used to test if the cache is recent.
//
// Caches takes Versions instead of strings to encourage compile time safety.
type Version interface {
	String() string
}

type Cache interface {

	// Put sets a value in the cache.
	Put(key Key, version Version, value []byte) error

	// GetBytes does a cache lookup and, if necessary, fill.
	//
	// If the cache has the key, and the version matches, the cached value will be returned.
	// Otherwise, filler will be called to fill the cache.
	// If it succeeds, its result will be put in the cache and returned.
	// Otherwise, an error will be returned.
	GetBytes(key Key, version Version, filler func() (Version, []byte, error)) ([]byte, error)

	// GetHandler does a cache lookup and, if necessary, fill.
	//
	// If the cache has the key, and the version matches, an http.Handler serving the value will be returned.
	// Otherwise, filler will be called to fill the cache.
	// If it succeeds, its result will be put in the cache and returned.
	// Otherwise, an error will be returned.
	//
	// The behavior of the http.Handler if called more than once is undefined.
	// Specific implementations may document their own behavior.
	GetHandler(key Key, version Version, filler func() (Version, []byte, error), contentType string) (http.Handler, error)

	// Close terminates open connections.
	//
	// The behavior of Close after the first call is undefined.
	// Specific implementations may document their own behavior.
	Close()
}
