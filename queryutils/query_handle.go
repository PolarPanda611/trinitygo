package queryutil

import (
	"fmt"
	"strings"

	"github.com/PolarPanda611/trinitygo/util"
)

var spiltValue = "__"
var filterCondition = []string{"like", "ilike", "in", "notin", "start", "end", "lt", "lte", "gt", "gte", "isnull", "isempty"}

// FilterQuery for filter query handling
type FilterQuery struct {
	QueryName   string
	QueryValue  string
	TablePrefix string

	// processing
	assosiationParam []string
	queryParam       string
	condition        string
	value            string

	//output
	ConditionSQL string
	ValueSQL     interface{}
}

func (f *FilterQuery) decode() {
	params := strings.Split(f.QueryName, spiltValue)
	paramsLen := len(params)
	if paramsLen == 1 {
		f.queryParam = params[paramsLen-1]
		return
	}
	if paramsLen >= 2 {
		if util.StringInSlice(params[paramsLen-1], filterCondition) {
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
}

// GetFilterConditionSQL get query sql
func (f *FilterQuery) getFilterConditionSQL() {
	switch f.condition {
	case "like":
		f.ConditionSQL = fmt.Sprintf(" %v like ? ", f.queryParam)
		f.ValueSQL = fmt.Sprintf("%v%v%v", "%", f.QueryValue, "%")
		break
	case "ilike":
		f.ConditionSQL = fmt.Sprintf(" %v ilike ? ", f.queryParam)
		f.ValueSQL = fmt.Sprintf("%v%v%v", "%", f.QueryValue, "%")
		break
	case "in":
		f.ConditionSQL = fmt.Sprintf(" %v in (?)  ", f.queryParam)
		f.ValueSQL = strings.Split(f.QueryValue, ",")
		break
	case "notin":
		f.ConditionSQL = fmt.Sprintf(" %v not in (?)  ", f.queryParam)
		f.ValueSQL = strings.Split(f.QueryValue, ",")
		break
	case "start":
		f.ConditionSQL = fmt.Sprintf(" %v  >= ? ", f.queryParam)
		f.ValueSQL = fmt.Sprintf("%v%v", f.QueryValue, " 00:00:00")
		break
	case "end":
		f.ConditionSQL = fmt.Sprintf(" %v  <= ? ", f.queryParam)
		f.ValueSQL = fmt.Sprintf("%v%v", f.QueryValue, " 23:59:59")
		break
	case "isnull":
		f.ConditionSQL = fmt.Sprintf(" %v is not null ", f.queryParam)
		f.ValueSQL = nil
		if f.QueryValue == "true" {
			f.ConditionSQL = fmt.Sprintf(" %v is null ", f.queryParam)
		}
		break
	case "lt":
		f.ConditionSQL = fmt.Sprintf(" %v  < ? ", f.queryParam)
		f.ValueSQL = f.QueryValue
		break
	case "lte":
		f.ConditionSQL = fmt.Sprintf(" %v  <= ? ", f.queryParam)
		f.ValueSQL = f.QueryValue
		break
	case "gt":
		f.ConditionSQL = fmt.Sprintf(" %v  > ? ", f.queryParam)
		f.ValueSQL = f.QueryValue
		break
	case "gte":
		f.ConditionSQL = fmt.Sprintf(" %v  >= ? ", f.queryParam)
		f.ValueSQL = f.QueryValue
		break
	case "isempty":
		f.ConditionSQL = fmt.Sprintf(" (COALESCE(\"%v\"::varchar ,'') != '' )  ", f.queryParam)
		f.ValueSQL = nil
		if f.QueryValue == "true" {
			f.ConditionSQL = fmt.Sprintf(" (COALESCE(\"%v\"::varchar ,'') = '' )  ", f.queryParam)
			f.ValueSQL = nil
		}
		break
	default:
		f.ConditionSQL = fmt.Sprintf(" %v = ? ", f.queryParam)
		f.ValueSQL = f.QueryValue

	}

	// 	break
	return
}

// GetFilterQuerySQL get query
func (f *FilterQuery) GetFilterQuerySQL() {
	f.decode()
	f.getFilterConditionSQL()
	assosiationParamLen := len(f.assosiationParam)
	for range f.assosiationParam {
		lastIndex := assosiationParamLen - 1
		lastParam := f.assosiationParam[lastIndex]
		f.ConditionSQL = fmt.Sprintf(" %v_id in ( select id from %v%v where %v ) ", lastParam, f.TablePrefix, lastParam, f.ConditionSQL)
		assosiationParamLen = assosiationParamLen - 1

	}
	f.ConditionSQL = util.DeleteExtraSpace(f.ConditionSQL)
	return
}
