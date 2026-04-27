package sqlhelper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/ponder2000/rdpms25-cloud-api/pkg/util/generic"
)

func ComplexFilterAppend(filters []qm.QueryMod, newFilter string, val any) []qm.QueryMod {
	return append(filters, qm.Where(newFilter, val))
}

func ParentChildRelationFilterAppend(filters []qm.QueryMod, rootIds []int, queryColName, descendantViewName, queryType string) []qm.QueryMod {
	iq := InnerQueryParentChildRelation(rootIds, descendantViewName, queryType)
	q := qm.Where(fmt.Sprintf(`%s in (%s)`, queryColName, iq))
	return append(filters, q)
}

func InnerQueryParentChildRelation(rootIds []int, descendantViewName, queryType string) string {
	rootArrayParam := IntArrayToPostgresArray(rootIds)

	var q string
	switch queryType {
	case "self":
		// this will return only root ids
		q = fmt.Sprintf(`(%s)`, strings.Join(generic.Mapper(rootIds, func(i int) string { return strconv.Itoa(i) }), ","))
	case "direct":
		// this will return only root ids and their direct children
		q = fmt.Sprintf(`(select descendant_id from %s where ancestor_id = ANY(%s) and depth <= 1)`, descendantViewName, rootArrayParam)
	case "all":
		// this will return all ids in the table irrespective of which roots are provided
		q = fmt.Sprintf(`(select distinct descendant_id from %s)`, descendantViewName)
	case "ancestor":
		// this will return all ancestors for the given descendant ids
		q = fmt.Sprintf(`(select ancestor_id from %s where descendant_id = ANY(%s))`, descendantViewName, rootArrayParam)
	default:
		// this will return self as well as all childs inside the tree
		q = fmt.Sprintf(`(select descendant_id from %s where ancestor_id = ANY(%s) )`, descendantViewName, rootArrayParam)
	}
	return q
}

func ComplexInStringFilterAppend(filters []qm.QueryMod, newFilter string, rawVal string) []qm.QueryMod {
	if len(rawVal) <= 0 {
		return filters
	}

	arrayVal := strings.Split(rawVal, ",")
	for i := range arrayVal {
		arrayVal[i] = fmt.Sprintf("'%s'", arrayVal[i])
	}

	filters = append(filters, qm.Where(fmt.Sprintf("%s in (%s)", newFilter, strings.Join(arrayVal, ","))))
	return filters
}
