package page

var NoPagination = NewPaginationQuery(-1, -1)
var EmptyPagination = NewPaginationQuery(0, 0)
var DefaultPagination = NewPaginationQuery(10, 0)

type PaginationQuery struct {
	// How many items per page should be returned
	ItemsPerPage int
	// Number of page that should be returned when pagination is applied to the list
	Page int
}

func NewPaginationQuery(itemsPerPage, page int) *PaginationQuery {
	return &PaginationQuery{itemsPerPage, page}
}

func (p *PaginationQuery) IsValidPagination() bool {
	return p.ItemsPerPage >= 0 && p.Page >= 0
}

func (p *PaginationQuery) IsPageAvailable(itemsCount, startingIndex int) bool {
	return itemsCount > startingIndex && p.ItemsPerPage > 0
}

func (p *PaginationQuery) GetPaginationSettings(itemsCount int) (startIndex int, endIndex int) {
	startIndex = p.ItemsPerPage * p.Page
	endIndex = startIndex + p.ItemsPerPage

	if endIndex > itemsCount {
		endIndex = itemsCount
	}

	return startIndex, endIndex
}
