package utils

import (
	"math/rand"
	"time"
)

func RandI64(end int) int {
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)
	return rng.Intn(end)
}
