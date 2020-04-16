package queryutil

import "testing"

func TestQueryHandler(t *testing.T) {
	query := "query=123&query=1234"
	config := &QueryConfig{
		DbBackend:           nil,
		pageSize:            20,
		FilterList:          []string{},
		OrderByList:         make(map[string]bool),
		FilterCustomizeFunc: make(map[string]interface{}),
	}
	handler := New(query, config)

	handler.Handle()
}
