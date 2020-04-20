package queryutil

import (
	"fmt"
	"strings"

	"github.com/PolarPanda611/trinitygo/util"
	"github.com/jinzhu/gorm"
)

var (
	_spiltValue      = "__"
	_filterCondition = []string{"like", "ilike", "in", "notin", "start", "end", "lt", "lte", "gt", "gte", "isnull", "isempty"}
)

// QueryDecoder query decoder
type QueryDecoder interface {
	ConditionSQL() string
	ValueSQL() interface{}
}

// NewDecoder new query decoder
func NewDecoder(queryName string, queryValue string) QueryDecoder {
	decoder := &filterQuery{
		QueryName:  queryName,
		QueryValue: queryValue,
	}
	decoder.decode()
	decoder.decodeFilterAndValue()
	decoder.decodeNestSQL()
	return decoder
}

type filterQuery struct {
	QueryName  string
	QueryValue string

	// processing
	assosiationParam []string
	queryParam       string
	condition        string
	value            string

	//output
	conditionSQL string
	valueSQL     interface{}
}

func (f *filterQuery) decode() {
	params := strings.Split(f.QueryName, _spiltValue)
	paramsLen := len(params)
	if paramsLen == 1 {
		f.queryParam = params[paramsLen-1]
	}
	if paramsLen >= 2 {
		if util.StringInSlice(params[paramsLen-1], _filterCondition) {
			f.condition = params[paramsLen-1]
			f.queryParam = params[paramsLen-2]
			if paramsLen >= 3 {
				f.assosiationParam = params[:paramsLen-2]
			}
		} else {
			f.queryParam = params[paramsLen-1]
			f.assosiationParam = params[:paramsLen-1]
		}
	}

	f.queryParam = gorm.ToColumnName(f.queryParam)
	newAssosiationParam := make([]string, len(f.assosiationParam))
	for i, v := range f.assosiationParam {
		newAssosiationParam[i] = gorm.ToTableName(v)
	}
	f.assosiationParam = newAssosiationParam

}

// decodeFilterAndValue get query sql
func (f *filterQuery) decodeFilterAndValue() {
	switch f.condition {
	case "like", "ilike":
		f.conditionSQL = fmt.Sprintf(" %v %v ? ", f.queryParam, f.condition)
		f.valueSQL = fmt.Sprintf("%v%v%v", "%", f.QueryValue, "%")
		break
	case "in":
		f.conditionSQL = fmt.Sprintf(" %v in (?)  ", f.queryParam)
		f.valueSQL = strings.Split(f.QueryValue, ",")
		break
	case "notin":
		f.conditionSQL = fmt.Sprintf(" %v not in (?)  ", f.queryParam)
		f.valueSQL = strings.Split(f.QueryValue, ",")
		break
	case "start":
		f.conditionSQL = fmt.Sprintf(" %v  >= ? ", f.queryParam)
		f.valueSQL = fmt.Sprintf("%v%v", f.QueryValue, " 00:00:00")
		break
	case "end":
		f.conditionSQL = fmt.Sprintf(" %v  <= ? ", f.queryParam)
		f.valueSQL = fmt.Sprintf("%v%v", f.QueryValue, " 23:59:59")
		break
	case "isnull":
		f.conditionSQL = fmt.Sprintf(" %v is not null ", f.queryParam)
		f.valueSQL = nil
		if f.QueryValue == "true" {
			f.conditionSQL = fmt.Sprintf(" %v is null ", f.queryParam)
		}
		break
	case "lt":
		f.conditionSQL = fmt.Sprintf(" %v  < ? ", f.queryParam)
		f.valueSQL = f.QueryValue
		break
	case "lte":
		f.conditionSQL = fmt.Sprintf(" %v  <= ? ", f.queryParam)
		f.valueSQL = f.QueryValue
		break
	case "gt":
		f.conditionSQL = fmt.Sprintf(" %v  > ? ", f.queryParam)
		f.valueSQL = f.QueryValue
		break
	case "gte":
		f.conditionSQL = fmt.Sprintf(" %v  >= ? ", f.queryParam)
		f.valueSQL = f.QueryValue
		break
	case "isempty":
		f.conditionSQL = fmt.Sprintf(" (COALESCE(\"%v\"::varchar ,'') != '' )  ", f.queryParam)
		f.valueSQL = nil
		if f.QueryValue == "true" {
			f.conditionSQL = fmt.Sprintf(" (COALESCE(\"%v\"::varchar ,'') = '' )  ", f.queryParam)
			f.valueSQL = nil
		}
		break
	default:
		f.conditionSQL = fmt.Sprintf(" %v = ? ", f.queryParam)
		f.valueSQL = f.QueryValue
	}

	// 	break
	return
}

// decodeNestSQL get query
func (f *filterQuery) decodeNestSQL() {
	assosiationParamLen := len(f.assosiationParam)
	for range f.assosiationParam {
		lastIndex := assosiationParamLen - 1
		lastParam := f.assosiationParam[lastIndex]
		f.conditionSQL = fmt.Sprintf(" %v_id in ( select id from %v where %v ) ", lastParam, gorm.DefaultTableNameHandler(nil, lastParam), f.conditionSQL)
		assosiationParamLen = assosiationParamLen - 1

	}
	f.conditionSQL = util.DeleteExtraSpace(f.conditionSQL)
	return
}

// ConditionSQL get condition sql
func (f *filterQuery) ConditionSQL() string {
	return f.conditionSQL
}

// ValueSQL get value sql
func (f *filterQuery) ValueSQL() interface{} {
	return f.valueSQL
}
