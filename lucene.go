package dragonchain

import "strconv"

func luceneQueryParams(q *QueryOptions) string {
	// Default to offset 0 with a limit of 10
	limit := 10
	if q.Limit > 0 {
		limit = q.Limit
	}
	offset := 0
	if q.Offset > 0 {
		offset = q.Offset
	}
	queryMap := map[string]string{
		"limit":  strconv.Itoa(limit),
		"offset": strconv.Itoa(offset),
	}
	if q.QueryString != "" {
		queryMap["q"] = q.QueryString
	}
	if q.Sort != "" {
		queryMap["sort"] = q.Sort
	}
	var query string
	for k, v := range queryMap {
		query += k + "=" + v + "&"
	}
	// We need to chop off the last '&'
	return query[:len(query)-1]
}
