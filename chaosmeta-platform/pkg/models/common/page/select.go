package page

type PropertyName string

const (
	NameProperty              = "name"
	CreationTimestampProperty = "creationTimestamp"
	NamespaceProperty         = "namespace"
	LabelProperty             = "label"
	StatusProperty            = "status"
	TypeProperty              = "type"
	PodIpProperty             = "podIP"
)

type DataSelectQuery struct {
	PaginationQuery *PaginationQuery
	SortQuery       *SortQuery
	FilterQuery     *FilterQuery
}

// SortQuery holds options for sort functionality of data select.
type SortQuery struct {
	SortByList []SortBy
}

// SortBy holds the name of the property that should be sorted and whether order should be ascending or descending.
type SortBy struct {
	Property  PropertyName
	Ascending bool
}

// NoSort is as option for no sort.
var NoSort = &SortQuery{
	SortByList: []SortBy{},
}

type FilterQuery struct {
	FilterByList []FilterBy
}

type FilterBy struct {
	Property PropertyName
	Value    ComparableValue
}

var NoFilter = &FilterQuery{
	FilterByList: []FilterBy{},
}

var NoDataSelect = NewDataSelectQuery(NoPagination, NoSort, NoFilter)
var StdMetricsDataSelect = NewDataSelectQuery(NoPagination, NoSort, NoFilter)
var DefaultDataSelect = NewDataSelectQuery(DefaultPagination, NoSort, NoFilter)
var DefaultDataSelectWithMetrics = NewDataSelectQuery(DefaultPagination, NoSort, NoFilter)

func NewDataSelectQuery(paginationQuery *PaginationQuery, sortQuery *SortQuery, filterQuery *FilterQuery) *DataSelectQuery {
	return &DataSelectQuery{
		PaginationQuery: paginationQuery,
		SortQuery:       sortQuery,
		FilterQuery:     filterQuery,
	}
}

func NewSortQuery(sortByListRaw []string) *SortQuery {
	if sortByListRaw == nil || len(sortByListRaw)%2 == 1 {
		return NoSort
	}
	var sortByList []SortBy
	for i := 0; i+1 < len(sortByListRaw); i += 2 {
		var ascending bool
		orderOption := sortByListRaw[i]
		if orderOption == "asc" {
			ascending = true
		} else if orderOption == "desc" {
			ascending = false
		} else {
			return NoSort
		}

		propertyName := sortByListRaw[i+1]
		sortBy := SortBy{
			Property:  PropertyName(propertyName),
			Ascending: ascending,
		}
		// Add to the sort options.
		sortByList = append(sortByList, sortBy)
	}
	return &SortQuery{
		SortByList: sortByList,
	}
}

func NewFilterQuery(filterByListRaw []string) *FilterQuery {
	if filterByListRaw == nil || len(filterByListRaw)%2 == 1 {
		return NoFilter
	}
	var filterByList []FilterBy
	for i := 0; i+1 < len(filterByListRaw); i += 2 {
		propertyName := filterByListRaw[i]
		propertyValue := filterByListRaw[i+1]
		filterBy := FilterBy{
			Property: PropertyName(propertyName),
			Value:    StdComparableString(propertyValue),
		}
		// Add to the filter options.
		filterByList = append(filterByList, filterBy)
	}
	return &FilterQuery{
		FilterByList: filterByList,
	}
}
