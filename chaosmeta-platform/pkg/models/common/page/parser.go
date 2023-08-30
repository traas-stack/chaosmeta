/*
 * Copyright 2022-2023 Chaos Meta Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package page

import (
	beego "github.com/beego/beego/v2/server/web"
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

func ParseDataSelectPathParameterTest() *DataSelectQuery {
	itemsPerPage := 10
	page := 1
	return NewDataSelectQuery(NewPaginationQuery(itemsPerPage, int(page-1)), NewSortQuery(nil), NewFilterQuery(nil))
}
