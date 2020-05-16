package queryutil

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/PolarPanda611/trinitygo/application"
	"github.com/PolarPanda611/trinitygo/keyword"
	"github.com/PolarPanda611/trinitygo/util"
	"github.com/jinzhu/gorm"
)

// KeyWord query utils key work definition
type KeyWord struct {
	SearchBy      string
	PageNum       string
	PageSize      string
	OrderBy       string
	PaginationOff string
}

// QueryHandler query handler
// in query
// out gorm scopes
type QueryHandler interface {
	RemotePreloadlist() []RemotePreloader
	PageSize() int
	PageNum() int
	HandleDBBackend() []func(*gorm.DB) *gorm.DB
	HandleWithPagination(query string) []func(*gorm.DB) *gorm.DB
	Handle(query string) []func(*gorm.DB) *gorm.DB
	HandleRemotePreloader(interface{}) error
}

// RemotePreloader  preload the model from remote resource
type RemotePreloader struct {
	Column      string
	PreloadFunc func(args interface{}, Conditions ...string) (interface{}, error)
	Condition   string
}

type queryRepositoryImpl struct {
	keyword             KeyWord
	query               string
	filterBackend       []func(db *gorm.DB) *gorm.DB
	pageSize            int
	searchByList        []string
	filterList          []string
	orderByList         []string
	preloadList         map[string]func(db *gorm.DB) *gorm.DB
	remotePreloadlist   []RemotePreloader
	filterCustomizeFunc map[string]interface{}
	queryMap            url.Values
	isDebug             bool

	// runtime
	queryScope      []func(*gorm.DB) *gorm.DB
	pageSizeRuntime int
	pageNumRuntime  int
}

// QueryConfig query config
type QueryConfig struct {
	FilterBackend     []func(db *gorm.DB) *gorm.DB
	PageSize          int
	SearchByList      []string
	FilterList        []string
	OrderByList       []string
	PreloadList       map[string]func(db *gorm.DB) *gorm.DB
	RemotePreloadlist []RemotePreloader
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
		FilterBackend:       nil,
		PageSize:            tCtx.Application().Conf().GetPageSize(),
		SearchByList:        []string{},
		FilterList:          []string{},
		OrderByList:         []string{},
		PreloadList:         make(map[string]func(db *gorm.DB) *gorm.DB),
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
		keyword: KeyWord{
			SearchBy:      keyword.GetKeyword().SearchBy,
			PageNum:       keyword.GetKeyword().PageNum,
			PageSize:      keyword.GetKeyword().PageSize,
			OrderBy:       keyword.GetKeyword().OrderBy,
			PaginationOff: keyword.GetKeyword().PaginationOff,
		},
		filterBackend:       config.FilterBackend,
		pageSize:            config.PageSize,
		filterList:          config.FilterList,
		orderByList:         config.OrderByList,
		preloadList:         config.PreloadList,
		remotePreloadlist:   config.RemotePreloadlist,
		searchByList:        config.SearchByList,
		filterCustomizeFunc: config.FilterCustomizeFunc,
		isDebug:             config.IsDebug,
	}
	return queryHandler
}
func (q *queryRepositoryImpl) cleanRuntime() {
	q.queryScope = nil
	q.pageSizeRuntime = 0
	q.pageNumRuntime = 0
}

func (q *queryRepositoryImpl) decodeURL() {
	q.queryMap, _ = url.ParseQuery(q.query)
}
func (q *queryRepositoryImpl) handleDBBackend() {
	if len(q.filterBackend) != 0 {
		for _, backend := range q.filterBackend {
			q.queryScope = append(q.queryScope, backend)
		}
	}
}
func (q *queryRepositoryImpl) handleFilter(k string, v []string) {
	if util.StringInSlice(k, q.filterList) {
		if len(v) != 0 {
			d := NewDecoder(k, v[0])
			if q.isDebug {
				fmt.Printf("where : %v  %v \n ", d.ConditionSQL(), d.ValueSQL())
			}
			if d.ValueSQL() != nil {
				q.queryScope = append(q.queryScope, NewScope(d.ConditionSQL(), d.ValueSQL()))
				return
			}
			q.queryScope = append(q.queryScope, NewScopeWithoutValue(d.ConditionSQL()))
			return
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
	if k == q.keyword.SearchBy {
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
	if k == q.keyword.OrderBy {
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
		IsOff, _ := strconv.ParseBool(paginationOff[0])
		if IsOff {
			return
		}
	}
	var err error
	PageNumFieldInt := 1
	PageSizeFieldInt := q.pageSize
	if len(pageNum) != 0 {
		if PageNumFieldInt, err = strconv.Atoi(pageNum[0]); err != nil || PageNumFieldInt <= 0 {
			PageNumFieldInt = 1
		}
	}
	q.pageNumRuntime = PageNumFieldInt
	if len(pageSize) != 0 {
		if PageSizeFieldInt, err = strconv.Atoi(pageSize[0]); err != nil || PageSizeFieldInt <= 0 {
			PageSizeFieldInt = q.pageSize
		}
	}
	q.pageSizeRuntime = PageSizeFieldInt
	PageNumFieldInt = PageNumFieldInt - 1
	offset := PageNumFieldInt * PageSizeFieldInt
	limit := PageSizeFieldInt
	if q.isDebug {
		fmt.Printf("offset %v  limit %v \n ", offset, limit)
	}
	paginationScope := func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(limit)
	}
	q.queryScope = append(q.queryScope, paginationScope)

}

func (q *queryRepositoryImpl) handlePreload() {
	if len(q.preloadList) > 0 {
		preloadScope := func(db *gorm.DB) *gorm.DB {
			for k, v := range q.preloadList {
				if v == nil {
					db = db.Preload(k)
				} else {
					db = db.Preload(k, v)
				}

			}
			return db
		}
		q.queryScope = append(q.queryScope, preloadScope)
	}
}

func (q *queryRepositoryImpl) Handle(query string) []func(*gorm.DB) *gorm.DB {
	q.query = query
	q.cleanRuntime()
	q.decodeURL()
	q.handleDBBackend()
	for k, v := range q.queryMap {
		q.handleOrderBy(k, v)
		q.handleSearchBy(k, v)
		q.handleCustomizeFilter(k, v)
		q.handleFilter(k, v)
	}
	q.handlePreload()
	newQueryScope := make([]func(*gorm.DB) *gorm.DB, len(q.queryScope))
	copy(newQueryScope, q.queryScope)
	return newQueryScope

}

func (q *queryRepositoryImpl) HandleWithPagination(query string) []func(*gorm.DB) *gorm.DB {
	q.query = query
	q.cleanRuntime()
	q.decodeURL()
	q.handleDBBackend()
	q.handlePagination(q.queryMap[q.keyword.PageNum], q.queryMap[q.keyword.PageSize], q.queryMap[q.keyword.PaginationOff])

	for k, v := range q.queryMap {
		q.handleSearchBy(k, v)
		q.handleCustomizeFilter(k, v)
		q.handleFilter(k, v)
	}
	q.handlePreload()
	newQueryScope := make([]func(*gorm.DB) *gorm.DB, len(q.queryScope))
	copy(newQueryScope, q.queryScope)
	return newQueryScope
}

func (q *queryRepositoryImpl) RemotePreloadlist() []RemotePreloader {
	newRemotePreloader := make([]RemotePreloader, len(q.remotePreloadlist))
	copy(newRemotePreloader, q.remotePreloadlist)
	return newRemotePreloader
}

func (q *queryRepositoryImpl) PageSize() int {
	pageSizeRuntime := q.pageSizeRuntime
	return pageSizeRuntime
}

func (q *queryRepositoryImpl) PageNum() int {
	pageNumRuntime := q.pageNumRuntime
	return pageNumRuntime
}

func (q *queryRepositoryImpl) HandleDBBackend() []func(*gorm.DB) *gorm.DB {
	q.cleanRuntime()
	q.handleDBBackend()
	newQueryScope := make([]func(*gorm.DB) *gorm.DB, len(q.queryScope))
	copy(newQueryScope, q.queryScope)
	return newQueryScope
}

func (q *queryRepositoryImpl) HandleRemotePreloader(obj interface{}) error {
	objType := reflect.TypeOf(obj)
	if objType.Kind() != reflect.Ptr && objType.Kind() != reflect.Slice {
		return errors.New("must be ptr or slice ")
	}
	var objVal reflect.Value
	if objType.Kind() == reflect.Ptr {
		objVal = reflect.Indirect(reflect.ValueOf(obj))
	}
	if objType.Kind() == reflect.Slice {
		objVal = reflect.ValueOf(obj)
	}

	fmt.Println(objType.Kind())
	switch objVal.Kind() {
	case reflect.Struct:
		for _, v := range q.remotePreloadlist {
			structField, _ := objType.Elem().FieldByName(v.Column)
			fmt.Println(structField.Name)
			fmt.Println(structField.Tag.Get("remote_resource"))
			fmt.Println(structField.Tag.Get("remote_condition"))
			targetVal := objVal.FieldByName(v.Column)
			fmt.Println(targetVal.Type()) //PTR
			targetKeyVal := objVal.FieldByName(v.Column + "ID")
			fmt.Println(targetKeyVal.Interface())
			res, err := v.PreloadFunc(targetKeyVal.Interface())
			if err != nil {
				return err
			}
			if !targetVal.CanSet() {
				return errors.New("must be public ")
			}
			targetVal.Set(reflect.ValueOf(res))
		}
		break
	case reflect.Slice:
		break
	default:
		return errors.New("must be ptr of struct or slice ")
	}
	return nil

}

// NewScope create new scope
func NewScope(conditionSQL string, valueSQL interface{}) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(conditionSQL, valueSQL)
	}
}

// NewScopeWithoutValue create new scope
func NewScopeWithoutValue(conditionSQL string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(conditionSQL)
	}
}
