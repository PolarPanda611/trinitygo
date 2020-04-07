package logger

import (
	"time"
)

// SkipperFunc check if method need to be skip
type SkipperFunc func(method string) bool

// Config middleware config
type Config struct {
	ProjectName    bool
	ProjectVersion bool

	// Method displays the request path (bool).
	//
	// Defaults to true.
	Method bool

	// Req show request info
	//
	// Defaults to true.
	Request bool
	// Response show response einfo
	//
	// Defaults to true.
	Response bool
	// Error show error messagee
	//
	// Defaults to true.
	Error bool

	// Latency show latency
	Latency bool

	// Runtime
	// Defaulte to true
	Runtime bool

	// Defaults nil
	LogFunc func(now time.Time, projectName, projectVersion, method string, latency time.Duration, runtimeValue map[string]string)

	// Skippers will according method to see if skip the logger middleware
	Skippers []SkipperFunc
}

// DefaultConfig default log  miiddleware setting
func DefaultConfig() *Config {
	return &Config{
		ProjectName:    true,
		ProjectVersion: true,
		Method:         true,
		Request:        true,
		Response:       true,
		Error:          true,
		Latency:        true,
		Runtime:        true,
		LogFunc:        nil,
		Skippers:       nil,
	}
}
