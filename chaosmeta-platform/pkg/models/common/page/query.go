package page

import (
	"sort"
)

type DataCell interface {
	GetProperty(PropertyName) ComparableValue
}

type ComparableValue interface {
	Compare(ComparableValue) int
	Contains(ComparableValue) bool
}

type DataSelector struct {
	GenericDataList []DataCell
	DataSelectQuery *DataSelectQuery
}

func (s DataSelector) Len() int { return len(s.GenericDataList) }

func (s DataSelector) Swap(i, j int) {
	s.GenericDataList[i], s.GenericDataList[j] = s.GenericDataList[j], s.GenericDataList[i]
}

func (s DataSelector) Less(i, j int) bool {
	for _, sortBy := range s.DataSelectQuery.SortQuery.SortByList {
		a := s.GenericDataList[i].GetProperty(sortBy.Property)
		b := s.GenericDataList[j].GetProperty(sortBy.Property)
		if a == nil || b == nil {
			break
		}
		cmp := a.Compare(b)
		if cmp == 0 {
			continue
		} else {
			return (cmp == -1 && sortBy.Ascending) || (cmp == 1 && !sortBy.Ascending)
		}
	}
	return false
}

func (s *DataSelector) Sort() *DataSelector {
	sort.Sort(*s)
	return s
}

func (s *DataSelector) Filter() *DataSelector {
	var filteredList []DataCell

	for _, c := range s.GenericDataList {
		matches := true
		for _, filterBy := range s.DataSelectQuery.FilterQuery.FilterByList {
			v := c.GetProperty(filterBy.Property)
			if v == nil || !v.Contains(filterBy.Value) {
				matches = false
				break
			}
		}
		if matches {
			filteredList = append(filteredList, c)
		}
	}

	s.GenericDataList = filteredList
	return s
}

func (s *DataSelector) Paginate() *DataSelector {
	pQuery := s.DataSelectQuery.PaginationQuery
	dataList := s.GenericDataList
	startIndex, endIndex := pQuery.GetPaginationSettings(len(dataList))

	if !pQuery.IsValidPagination() {
		return s
	}
	if !pQuery.IsPageAvailable(len(s.GenericDataList), startIndex) {
		s.GenericDataList = []DataCell{}
		return s
	}
	s.GenericDataList = dataList[startIndex:endIndex]
	return s
}

func GenericDataSelect(dataList []DataCell, dsQuery *DataSelectQuery) []DataCell {
	SelectableData := DataSelector{
		GenericDataList: dataList,
		DataSelectQuery: dsQuery,
	}
	return SelectableData.Sort().Paginate().GenericDataList
}

func GenericDataSelectWithFilter(dataList []DataCell, dsQuery *DataSelectQuery) ([]DataCell, int) {
	SelectableData := DataSelector{
		GenericDataList: dataList,
		DataSelectQuery: dsQuery,
	}
	filtered := SelectableData.Filter()
	filteredTotal := len(filtered.GenericDataList)
	processed := filtered.Sort().Paginate()
	return processed.GenericDataList, filteredTotal
}
