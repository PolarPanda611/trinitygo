package keyword

var (
	_searchByKey    string  = "SearchBy"
	_pageNumKey     string  = "PageNum"
	_pageSizeKey    string  = "PageSize"
	_orderByKey     string  = "OrderBy"
	_paginationOff  string  = "PaginationOff"
	_defaultKeyword Keyword = Keyword{
		SearchBy:      _searchByKey,
		PageNum:       _pageNumKey,
		PageSize:      _pageSizeKey,
		OrderBy:       _orderByKey,
		PaginationOff: _paginationOff,
	}
)

// Keyword query utils key work definition
type Keyword struct {
	SearchBy      string
	PageNum       string
	PageSize      string
	OrderBy       string
	PaginationOff string
}

// SetKeyword set default  keyword
func SetKeyword(k Keyword) {
	_defaultKeyword = Keyword{
		SearchBy:      k.SearchBy,
		PageNum:       k.PageNum,
		PageSize:      k.PageSize,
		OrderBy:       k.OrderBy,
		PaginationOff: k.PaginationOff,
	}
}

// GetKeyword get keyword list
func GetKeyword() Keyword {
	return _defaultKeyword

}
