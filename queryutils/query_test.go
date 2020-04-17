package queryutil

import (
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
)

func TestQueryHandler(t *testing.T) {
	query := "query=123&Query=1234&code=sasfaf&query__ilike=234&OrderBy=id&SearchBy=124"
	config := &QueryConfig{
		TablePrefix:         "",
		DbBackend:           nil,
		PageSize:            20,
		FilterList:          []string{"query", "query__ilike"},
		OrderByList:         []string{"id"},
		SearchByList:        []string{"code", "name"},
		FilterCustomizeFunc: make(map[string]interface{}),
		IsDebug:             true,
	}
	x := New(config).Handle(query)
	fmt.Println(len(x))
}

func TestQueryHandlerWithPagi(t *testing.T) {
	query := "query=123&Query=1234&code=sasfaf&query__ilike=234&OrderBy=id&SearchBy=124&PageSize=30&PageNum=4&test=wqrqwr"
	config := &QueryConfig{
		TablePrefix:  "",
		DbBackend:    nil,
		PageSize:     20,
		FilterList:   []string{"query", "query__ilike"},
		OrderByList:  []string{"id"},
		SearchByList: []string{"code", "name"},
		FilterCustomizeFunc: map[string]interface{}{
			"test": func(db *gorm.DB, queryValue string) *gorm.DB {
				fmt.Println("Where xxxxx = ?", queryValue)
				return db.Where("xxxxx = ?", queryValue)
			},
		},
		IsDebug: true,
	}
	x := New(config).HandleWithPagination(query)
	fmt.Println(len(x))
}
