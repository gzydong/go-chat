package ichat

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

type Paginate struct {
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

type PaginateResponse struct {
	Items    interface{} `json:"items"`
	Paginate Paginate    `json:"paginate"`
}
