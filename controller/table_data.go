package controller

type TableData struct {
	List        interface{} `json:"list"`
	Total       int64       `json:"total"`
	PageSize    int         `json:"pageSize"`
	CurrentPage int         `json:"currentPage"`
}

func newDefaultTableData(list interface{}, total int64) *TableData {
	return &TableData{
		List:        list,
		Total:       total,
		PageSize:    10,
		CurrentPage: 1,
	}
}

type TableDataOption func(d *TableData) *TableData

func TableDataWithPageSize(pageSize int) TableDataOption {
	return func(d *TableData) *TableData {
		d.PageSize = pageSize
		return d
	}
}

func TableDataWithCurrentPage(currentPage int) TableDataOption {
	return func(d *TableData) *TableData {
		d.CurrentPage = currentPage
		return d
	}
}
