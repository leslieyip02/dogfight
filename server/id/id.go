package id

import (
	"sync"

	"github.com/sqids/sqids-go"
)

var (
	counter    uint64 = 0
	encoder, _        = sqids.New(sqids.Options{MinLength: 6})
	mu         sync.Mutex
)

func NewShortId() (string, error) {
	mu.Lock()
	defer mu.Unlock()

	encoded, err := encoder.Encode([]uint64{counter})
	if err != nil {
		return "", err
	}

	counter++
	return encoded, nil
}
