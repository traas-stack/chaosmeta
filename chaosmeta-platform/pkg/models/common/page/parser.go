package page

import (
	"github.com/astaxie/beego"
	"strings"
)

func parsePaginationPathParameter(c *beego.Controller) *PaginationQuery {
	itemsPerPage, err := c.GetInt("page_size", 10)
	if err != nil {
		return NoPagination
	}
	page, err := c.GetInt("page", 1)
	if err != nil {
		return NoPagination
	}
	return NewPaginationQuery(itemsPerPage, int(page-1))
}

func parseFilterPathParameter(c *beego.Controller) *FilterQuery {
	return NewFilterQuery(strings.Split(c.GetString("filterBy"), ","))
}

func parseSortPathParameter(c *beego.Controller) *SortQuery {
	return NewSortQuery(strings.Split(c.GetString("sortBy"), ","))
}

func ParseDataSelectPathParameter(c *beego.Controller) *DataSelectQuery {
	paginationQuery := parsePaginationPathParameter(c)
	sortQuery := parseSortPathParameter(c)
	filterQuery := parseFilterPathParameter(c)
	return NewDataSelectQuery(paginationQuery, sortQuery, filterQuery)
}
