package dto

type Paginate struct {
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

type PaginateResponse struct {
	Rows     interface{} `json:"rows"`
	Paginate Paginate    `json:"paginate"`
}
