package controller

type TableData struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
	Info  interface{} `json:"info"`

	Pagination
}

type Pagination struct {
	PageSize    int `json:"pageSize" url:"pageSize"`
	CurrentPage int `json:"currentPage" url:"currentPage"`
}

func newDefaultTableData(list interface{}, total int64) *TableData {
	return &TableData{
		List:  list,
		Total: total,
		Pagination: Pagination{
			PageSize:    10,
			CurrentPage: 1,
		},
	}
}

type TableDataOption func(d *TableData) *TableData

func TableDataWithInfo(dataInfo interface{}) TableDataOption {
	return func(d *TableData) *TableData {
		d.Info = dataInfo
		return d
	}
}

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
