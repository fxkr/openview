package safe

import (
	"encoding/json"
)

type Key struct {
	raw []byte
}

// NewKey constructs a Key from components.
//
// It is guaranteed that the components are safely delimited.
func NewKey(components ...interface{}) Key {
	bytes, err := json.Marshal(components)
	if err != nil {
		panic("Failed to convert key to JSON")
	}
	return Key{bytes}
}

// Converts Key to a string.
//
// There are no guarantees about the format of the string.
// In particular, it is NOT guaranteed to be a safe filename!
func (k Key) String() string {
	return string(k.raw)
}
