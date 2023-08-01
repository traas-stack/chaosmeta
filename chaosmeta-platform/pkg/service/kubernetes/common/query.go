package common

import (
	"chaosmeta-platform/util"
	"strconv"

	"github.com/emicklei/go-restful/v3"
	"k8s.io/apimachinery/pkg/labels"
)

type Field string
type Value string

const (
	FieldName                = "name"
	FieldNames               = "names"
	FieldUID                 = "uid"
	FieldCreationTimeStamp   = "creationTimestamp"
	FieldCreateTime          = "createTime"
	FieldLastUpdateTimestamp = "lastUpdateTimestamp"
	FieldUpdateTime          = "updateTime"
	FieldLabel               = "label"
	FieldAnnotation          = "annotation"
	FieldNamespace           = "namespace"
	FieldStatus              = "status"
	FieldOwnerReference      = "ownerReference"
	FieldOwnerKind           = "ownerKind"

	FieldType = "type"
)

var SortableFields = []Field{
	FieldCreationTimeStamp,
	FieldCreateTime,
	FieldUpdateTime,
	FieldLastUpdateTimestamp,
	FieldName,
}

// Field contains all the query field that can be compared
var ComparableFields = []Field{
	FieldName,
	FieldUID,
	FieldLabel,
	FieldAnnotation,
	FieldNamespace,
	FieldStatus,
	FieldOwnerReference,
	FieldOwnerKind,
}

const (
	ParameterName          = "name"
	ParameterLabelSelector = "labelSelector"
	ParameterFieldSelector = "fieldSelector"
	ParameterPage          = "page"
	ParameterLimit         = "limit"
	ParameterOrderBy       = "sortBy"
	ParameterAscending     = "ascending"
)

// Query represents api search terms
type Query struct {
	Pagination *Pagination

	// sort result in which field, default to FieldCreationTimeStamp
	SortBy Field

	// sort result in ascending or descending order, default to descending
	Ascending bool

	//
	Filters map[Field]Value

	LabelSelector string
}

type Pagination struct {
	// items per page
	Limit int

	// offset
	Offset int
}

var NoPagination = newPagination(-1, 0)

// make sure that pagination is valid
func newPagination(limit int, offset int) *Pagination {
	return &Pagination{
		Limit:  limit,
		Offset: offset,
	}
}

func (q *Query) Selector() labels.Selector {
	if selector, err := labels.Parse(q.LabelSelector); err != nil {
		return labels.Everything()
	} else {
		return selector
	}
}

func (p *Pagination) GetValidPagination(total int) (startIndex, endIndex int) {

	// no pagination
	if p.Limit == NoPagination.Limit {
		return 0, total
	}

	// out of range
	if p.Limit < 0 || p.Offset < 0 || p.Offset > total {
		return 0, 0
	}

	startIndex = p.Offset
	endIndex = startIndex + p.Limit

	if endIndex > total {
		endIndex = total
	}

	return startIndex, endIndex
}

func New() *Query {
	return &Query{
		Pagination: NoPagination,
		SortBy:     "",
		Ascending:  false,
		Filters:    map[Field]Value{},
	}
}

type Filter struct {
	Field Field
	Value Value
}

func ParseQueryParameter(request *restful.Request) *Query {
	query := New()

	limit, err := strconv.Atoi(request.QueryParameter(ParameterLimit))
	// equivalent to undefined, use the default value
	if err != nil {
		limit = -1
	}
	page, err := strconv.Atoi(request.QueryParameter(ParameterPage))
	// equivalent to undefined, use the default value
	if err != nil {
		page = 1
	}

	query.Pagination = newPagination(limit, (page-1)*limit)

	query.SortBy = Field(defaultString(request.QueryParameter(ParameterOrderBy), FieldCreationTimeStamp))

	ascending, err := strconv.ParseBool(defaultString(request.QueryParameter(ParameterAscending), "false"))
	if err != nil {
		query.Ascending = false
	} else {
		query.Ascending = ascending
	}

	query.LabelSelector = request.QueryParameter(ParameterLabelSelector)

	for key, values := range request.Request.URL.Query() {
		if !util.HasString([]string{ParameterPage, ParameterLimit, ParameterOrderBy, ParameterAscending, ParameterLabelSelector}, key) {
			// support multiple query condition
			for _, value := range values {
				query.Filters[Field(key)] = Value(value)
			}
		}
	}

	return query
}

func defaultString(value, defaultValue string) string {
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
