package queryutil

import (
	"fmt"
	"net/url"

	"github.com/PolarPanda611/trinitygo/application"
	"github.com/PolarPanda611/trinitygo/util"
	"github.com/jinzhu/gorm"
)

// QueryHandler query handler
// in query
// out gorm scopes
type QueryHandler interface {
	HandleWithPagination() []func(*gorm.DB) *gorm.DB
	Handle() []func(*gorm.DB) *gorm.DB
}

type queryRepositoryImpl struct {
	query               string
	dbBackend           func(db *gorm.DB) *gorm.DB
	pageSize            int
	filterList          []string
	orderByList         map[string]bool
	filterCustomizeFunc map[string]interface{}
}

// QueryConfig query config
type QueryConfig struct {
	DbBackend   func(db *gorm.DB) *gorm.DB
	pageSize    int
	FilterList  []string
	OrderByList map[string]bool
	// SetFilterCustomizeFunc
	// run local
	// option 1 : func(db *gorm.DB, queryValue string) *gorm.DB
	// run remotely
	// option 2 : rpc call --> to do
	//   func(manager interface{}, methodName string) *gorm.DB
	FilterCustomizeFunc map[string]interface{}
}

// DefaultConfig with system setting
func DefaultConfig(tCtx application.Context) *QueryConfig {
	return &QueryConfig{
		DbBackend:           nil,
		pageSize:            tCtx.Application().Conf().GetPageSize(),
		FilterList:          []string{},
		OrderByList:         make(map[string]bool),
		FilterCustomizeFunc: make(map[string]interface{}),
	}

}

// NewWithDefaultConfig query handler with default config
func NewWithDefaultConfig(tCtx application.Context, query string) QueryHandler {
	config := DefaultConfig(tCtx)
	queryHandler := &queryRepositoryImpl{
		pageSize:            config.pageSize,
		filterList:          config.FilterList,
		orderByList:         config.OrderByList,
		filterCustomizeFunc: config.FilterCustomizeFunc,
		query:               query,
	}
	return queryHandler
}

// New query handler with customize handler config
func New(query string, config *QueryConfig) QueryHandler {
	queryHandler := &queryRepositoryImpl{
		pageSize:            config.pageSize,
		filterList:          config.FilterList,
		orderByList:         config.OrderByList,
		filterCustomizeFunc: config.FilterCustomizeFunc,
		query:               query,
	}
	return queryHandler
}
func (q *queryRepositoryImpl) Handle() []func(*gorm.DB) *gorm.DB {
	var queryScope []func(*gorm.DB) *gorm.DB
	//Handle db backend
	if q.dbBackend != nil {
		queryScope = append(queryScope, q.dbBackend)
	}
	//Handle filter
	queryValue, _ := url.ParseQuery(q.query)
	fmt.Println(queryValue)
	for k, v := range queryValue {
		if util.StringInSlice(k, q.filterList) {
			fmt.Println(v)
		}
	}
	return nil
	//Handle query

}

func (q *queryRepositoryImpl) HandleWithPagination() []func(*gorm.DB) *gorm.DB {
	return nil
}
