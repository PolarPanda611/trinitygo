package sd

import (
	"errors"
	"sync/atomic"

	"github.com/PolarPanda611/trinitygo/httputil"
)

// Balancer balancer
type Balancer interface {
	Client() (*httputil.ServiceClient, error)
}

// NewRoundRobin returns a load balancer that returns services in sequence.
func NewRoundRobin(s []httputil.ServiceClient) Balancer {
	return &roundRobin{
		s: s,
		c: 0,
	}
}

type roundRobin struct {
	s []httputil.ServiceClient
	c uint64
}

func (rr *roundRobin) Client() (*httputil.ServiceClient, error) {
	if len(rr.s) <= 0 {
		return nil, errors.New("not found service client ")
	}
	old := atomic.AddUint64(&rr.c, 1) - 1
	idx := old % uint64(len(rr.s))
	return &rr.s[idx], nil
}
