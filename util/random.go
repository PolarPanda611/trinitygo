package util

import (
	"math/rand"
)

// RandomGenerater random generater interface{}
type RandomGenerater interface {
	RandomString() string
}

// RandomGenerater  random generater
type randomGeneraterImpl struct {
	r   *rand.Rand
	len int
}

// NewRandomGenerater new generater
func NewRandomGenerater(r *rand.Rand, len int) RandomGenerater {
	return &randomGeneraterImpl{
		r:   r,
		len: len,
	}
}

// RandomString generate random string
func (r *randomGeneraterImpl) RandomString() string {
	bytes := make([]byte, r.len)
	for i := 0; i < r.len; i++ {
		b := r.r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}
