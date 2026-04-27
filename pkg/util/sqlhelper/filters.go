package sqlhelper

import (
	"fmt"
	"log/slog"
	"strings"

	"strconv"

	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/lib/pq"
	"github.com/ponder2000/rdpms25-cloud-api/pkg/util/parser"
)

type Filter struct {
	ColName  string
	ColVal   interface{}
	Operator string
	IsArray  bool
	IsNum    bool
}

type FilterSlice []*Filter

func (f *Filter) paramValue() (interface{}, error) {
	switch v := f.ColVal.(type) {
	case int, bool, int64, float64:
		return v, nil
	case string:
		if f.isLikeOperator() {
			return "%" + v + "%", nil
		}
		return v, nil
	default:
		return nil, fmt.Errorf("invalid Type for filter")
	}
}

func (f *Filter) isLikeOperator() bool {
	return strings.Contains(f.Operator, "like")
}

func FilterQueryBuilder(table string, filters FilterSlice) []qm.QueryMod {
	queries := make([]qm.QueryMod, 0)
	for _, f := range filters {
		if f.IsArray {
			strVal, ok := f.ColVal.(string)
			if !ok {
				slog.Error("unable to filter array: value is not string", "filter", f)
				continue
			}
			arrayVal := strings.Split(strVal, ",")
			if f.IsNum {
				intArgs := make([]int, 0, len(arrayVal))
				for _, v := range arrayVal {
					if n, err := strconv.Atoi(v); err == nil {
						intArgs = append(intArgs, n)
					}
				}
				queries = append(queries, qm.Where(fmt.Sprintf(`"%s"."%s" = ANY(?)`, table, f.ColName), pq.Array(intArgs)))
			} else {
				queries = append(queries, qm.Where(fmt.Sprintf(`"%s"."%s" = ANY(?)`, table, f.ColName), pq.Array(arrayVal)))
			}
		} else {
			val, e := f.paramValue()
			if e != nil {
				slog.Error("unable to filter", "filter", f)
				continue
			}
			queries = append(queries, qm.Where(fmt.Sprintf(`"%s"."%s" %s ?`, table, f.ColName, f.Operator), val))
		}
	}
	slog.Debug("query formed", "table", table, "query", queries)
	return queries
}

// StringAppendToFilter
// helper append function for string queries
func StringAppendToFilter(filters FilterSlice, colName, operator, colValue string) FilterSlice {
	if len(colValue) > 0 {
		filters = append(filters, &Filter{ColName: colName, ColVal: colValue, Operator: operator})
	}
	return filters
}

// BoolAppendToFilter
// helper append function for bool queries
func BoolAppendToFilter(filters FilterSlice, colName, operator string, colValue string) FilterSlice {
	if v, e := parser.ExtractBool(colValue); e != nil {
		return filters
	} else {
		return append(filters, &Filter{ColName: colName, ColVal: v, Operator: operator})
	}
}

// IntAppendToFilter
// helper append function for string queries
func IntAppendToFilter(filters FilterSlice, colName, operator string, colValue string) FilterSlice {
	if num, e := parser.ExtractInt(colValue, true); e != nil {
		return filters
	} else {
		return append(filters, &Filter{ColName: colName, ColVal: num, Operator: operator})
	}
}

func InStringsAppendToFilter(filters FilterSlice, colName, colValues string) FilterSlice {
	if len(colValues) > 0 {
		filters = append(filters, &Filter{ColName: colName, ColVal: colValues, Operator: "in", IsArray: true})
	}
	return filters
}

func InIntAppendToFilter(filters FilterSlice, colName, colValues string) FilterSlice {
	if len(colValues) > 0 {
		filters = append(filters, &Filter{ColName: colName, ColVal: colValues, Operator: "in", IsArray: true, IsNum: true})
	}
	return filters
}
