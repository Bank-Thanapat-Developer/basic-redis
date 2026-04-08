package dto

type PaginationRequest struct {
	Page     int `query:"page"`
	PageSize int `query:"page_size"`
}

type PaginationResponse struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int   `json:"total_items"`
	TotalPages int   `json:"total_pages"`
	Items      []any `json:"items"`
}
