package sqlhelper

import (
	"database/sql"
	"fmt"

	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/shivamkkas/rdpms25-go-rpc-service/pkg/util/generic"
)

func IntArrayToPostgresArray(intSlice []int) string {
	idsStrSlice := generic.Mapper(intSlice, func(elm int) string { return strconv.Itoa(elm) })
	return fmt.Sprintf(`ARRAY[%s]`, strings.Join(idsStrSlice, ","))
}

func StringArrayFormat(commaSeperaetdVal string) string {
	vals := strings.Split(commaSeperaetdVal, ",")
	for i := range vals {
		vals[i] = fmt.Sprintf(`'%s'`, vals[i])
	}
	return fmt.Sprintf(`(%s)`, strings.Join(vals, ","))
}
func IntSliceFormat(commaSeparatedVal string) ([]int, error) {

	vals := strings.Split(commaSeparatedVal, ",")
	result := make([]int, 0, len(vals))

	for _, v := range vals {
		trimmed := strings.TrimSpace(v)
		n, err := strconv.Atoi(trimmed)
		if err != nil {
			return nil, fmt.Errorf("invalid integer found: %s", trimmed)
		}
		result = append(result, n)
	}

	return result, nil
}

func IntArrayFormat(commaSeperaetdVal string) string {
	return fmt.Sprintf(`(%s)`, commaSeperaetdVal)
}

func JsonContainQuery(columnName, commaSeperaetdVal string) string {
	queries := make([]string, 0)

	for _, kv := range strings.Split(commaSeperaetdVal, ",") {
		tmp := strings.Split(kv, ":")
		if len(tmp) == 2 {
			key, value := tmp[0], tmp[1]
			value = "%" + value + "%"
			queries = append(queries, fmt.Sprintf(`"%s"->>'%s' ilike '%s'`, columnName, key, value))
		} else {
			slog.Warn("Invalid json contain filter")
		}
	}
	return strings.Join(queries, " and ")
}

func SelectAlias(colName, alias string) string {
	return fmt.Sprintf(`%s as %s`, colName, alias)
}

func DateTimeString(t time.Time) string {
	return fmt.Sprintf(`%d-%d-%d %d:%d:%d`, t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

func LikePrefixSuffixAdd(arg string) string {
	return "%" + arg + "%"
}

func StringToNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}
