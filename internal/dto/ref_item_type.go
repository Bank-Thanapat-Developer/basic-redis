package dto

type RefItemTypeCreateRequest struct {
	Name string `json:"name"`
}

type RefItemTypeResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
