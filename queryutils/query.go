package queryutil

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/PolarPanda611/trinitygo/application"
	"github.com/PolarPanda611/trinitygo/util"
	"github.com/jinzhu/gorm"
)

var (
	_searchByKey   string = "SearchBy"
	_pageNumKey    string = "PageNum"
	_pageSizeKey   string = "PageSize"
	_orderByKey    string = "OrderBy"
	_paginitionOff string = "PaginationOff"
)

// QueryHandler query handler
// in query
// out gorm scopes
type QueryHandler interface {
	HandleWithPagination(query string) []func(*gorm.DB) *gorm.DB
	Handle(query string) []func(*gorm.DB) *gorm.DB
}

type queryRepositoryImpl struct {
	query               string
	tablePrefix         string
	dbBackend           func(db *gorm.DB) *gorm.DB
	pageSize            int
	searchByList        []string
	filterList          []string
	orderByList         []string
	filterCustomizeFunc map[string]interface{}

	queryMap   url.Values
	queryScope []func(*gorm.DB) *gorm.DB

	isDebug bool
}

// QueryConfig query config
type QueryConfig struct {
	TablePrefix  string
	DbBackend    func(db *gorm.DB) *gorm.DB
	PageSize     int
	SearchByList []string
	FilterList   []string
	OrderByList  []string
	// SetFilterCustomizeFunc
	// run local
	// option 1 : func(db *gorm.DB, queryValue string) *gorm.DB
	// run remotely
	// option 2 : rpc call --> to do
	//   func(manager interface{}, methodName string) *gorm.DB
	FilterCustomizeFunc map[string]interface{}
	IsDebug             bool
}

// DefaultConfig with system setting
func DefaultConfig(tCtx application.Context) *QueryConfig {
	return &QueryConfig{
		TablePrefix:         tCtx.Application().Conf().GetTablePrefix(),
		DbBackend:           nil,
		PageSize:            tCtx.Application().Conf().GetPageSize(),
		SearchByList:        []string{},
		FilterList:          []string{},
		OrderByList:         []string{},
		FilterCustomizeFunc: make(map[string]interface{}),
		IsDebug:             false,
	}

}

// NewWithDefaultConfig query handler with default config
func NewWithDefaultConfig(tCtx application.Context, query string) QueryHandler {
	config := DefaultConfig(tCtx)
	return New(config)
}

// New query handler with customize handler config
func New(config *QueryConfig) QueryHandler {
	queryHandler := &queryRepositoryImpl{
		tablePrefix:         config.TablePrefix,
		pageSize:            config.PageSize,
		filterList:          config.FilterList,
		orderByList:         config.OrderByList,
		searchByList:        config.SearchByList,
		filterCustomizeFunc: config.FilterCustomizeFunc,
		isDebug:             config.IsDebug,
	}
	return queryHandler
}

func (q *queryRepositoryImpl) decodeURL() {
	q.queryMap, _ = url.ParseQuery(q.query)
}
func (q *queryRepositoryImpl) handleDBBackend() {
	if q.dbBackend != nil {
		q.queryScope = append(q.queryScope, q.dbBackend)
	}
}
func (q *queryRepositoryImpl) handleFilter(k string, v []string) {
	if util.StringInSlice(k, q.filterList) {
		if len(v) != 0 {
			d := NewDecoder(k, v[0], q.tablePrefix)
			if q.isDebug {
				fmt.Printf("where : %v  %v \n ", d.ConditionSQL(), d.ValueSQL())
			}
			q.queryScope = append(q.queryScope, newScope(d.ConditionSQL(), d.ValueSQL()))
		}

	}
}

func (q *queryRepositoryImpl) handleCustomizeFilter(k string, v []string) {
	if customizeFunc, ok := q.filterCustomizeFunc[k]; ok {
		queryValue := strings.Join(v, ",")
		if f, ok := customizeFunc.(func(db *gorm.DB, queryValue string) *gorm.DB); ok {
			queryFunc := func(db *gorm.DB) *gorm.DB {
				return f(db, queryValue)
			}
			q.queryScope = append(q.queryScope, queryFunc)
		}
	}
}
func (q *queryRepositoryImpl) handleSearchBy(k string, v []string) {
	if k == _searchByKey {
		if len(v) != 0 {
			queryValue := v[0]
			var queryString string
			ilikeFormatter := " \"%v\" ilike ? "
			orJoinFormatter := " %v or %v "
			var querySlice []string
			var queryValueSlice []interface{}
			for _, searchByColumn := range q.searchByList {
				querySlice = append(querySlice, fmt.Sprintf(ilikeFormatter, gorm.ToColumnName(searchByColumn)))
			}
			for i, query := range querySlice {
				queryValueSlice = append(queryValueSlice, "%"+queryValue+"%")
				if i == 0 {
					queryString = query
					continue
				}
				queryString = fmt.Sprintf(orJoinFormatter, queryString, query)
			}
			if q.isDebug {
				fmt.Printf("where : %v  %v \n ", queryString, queryValueSlice)
			}

			searchScope := func(db *gorm.DB) *gorm.DB {
				return db.Where(queryString, queryValueSlice...)
			}
			q.queryScope = append(q.queryScope, searchScope)
		}
	}

}

func (q *queryRepositoryImpl) handleOrderBy(k string, v []string) {
	if k == _orderByKey {
		ordercondition := ""
		queryValue := strings.Join(v, ",")
		for _, orderField := range strings.Split(queryValue, ",") {
			if len(strings.Split(orderField, "-")) == 2 {
				if util.StringInSlice(strings.Split(orderField, "-")[1], q.orderByList) {
					ordercondition += strings.Split(orderField, "-")[1] + " desc ,"
				}
			} else if len(strings.Split(orderField, "-")) == 1 {
				if util.StringInSlice(orderField, q.orderByList) {
					ordercondition += orderField + " asc ,"
				}
			}
		}
		if q.isDebug {
			fmt.Printf("order by  :   %v \n ", strings.TrimSuffix(ordercondition, ","))
		}
		orderByScope := func(db *gorm.DB) *gorm.DB {
			return db.Order(strings.TrimSuffix(ordercondition, ","))
		}
		q.queryScope = append(q.queryScope, orderByScope)
	}

}

func (q *queryRepositoryImpl) handlePagination(pageNum []string, pageSize []string, paginationOff []string) {
	if len(paginationOff) != 0 {
		return
	}
	var err error
	PageNumFieldInt := 1
	PageSizeFieldInt := q.pageSize
	if len(pageNum) != 0 {
		if PageNumFieldInt, err = strconv.Atoi(pageNum[0]); err != nil || PageNumFieldInt <= 0 {
			PageNumFieldInt = 1
		}
	}
	if len(pageSize) != 0 {
		if PageSizeFieldInt, err = strconv.Atoi(pageSize[0]); err != nil || PageSizeFieldInt <= 0 {
			PageSizeFieldInt = q.pageSize
		}
	}
	PageNumFieldInt = PageNumFieldInt - 1
	offset := PageNumFieldInt * PageSizeFieldInt
	limit := PageSizeFieldInt
	if q.isDebug {
		fmt.Printf("offset %v  limit %v \n ", offset, limit)
	}
	pagiScope := func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(limit)
	}
	q.queryScope = append(q.queryScope, pagiScope)

}
func (q *queryRepositoryImpl) Handle(query string) []func(*gorm.DB) *gorm.DB {
	q.query = query
	q.decodeURL()
	q.handleDBBackend()
	for k, v := range q.queryMap {
		q.handleOrderBy(k, v)
		q.handleSearchBy(k, v)
		q.handleCustomizeFilter(k, v)
		q.handleFilter(k, v)
	}
	newQueryScope := make([]func(*gorm.DB) *gorm.DB, len(q.queryScope))
	copy(newQueryScope, q.queryScope)
	return newQueryScope

}

func (q *queryRepositoryImpl) HandleWithPagination(query string) []func(*gorm.DB) *gorm.DB {
	q.query = query
	q.decodeURL()
	q.handleDBBackend()
	q.handlePagination(q.queryMap[_pageNumKey], q.queryMap[_pageSizeKey], q.queryMap[_paginitionOff])

	for k, v := range q.queryMap {
		q.handleSearchBy(k, v)
		q.handleCustomizeFilter(k, v)
		q.handleFilter(k, v)
	}
	newQueryScope := make([]func(*gorm.DB) *gorm.DB, len(q.queryScope))
	copy(newQueryScope, q.queryScope)
	return newQueryScope
}

func newScope(conditionSQL string, valueSQL interface{}) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(conditionSQL, valueSQL)
	}
}
