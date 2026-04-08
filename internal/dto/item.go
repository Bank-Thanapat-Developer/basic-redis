package dto

type ItemCreateRequest struct {
	Name          string `json:"name"`
	Price         int    `json:"price"`
	IsActive      bool   `json:"is_active"`
	RefItemTypeID int    `json:"ref_item_type_id"`
}

type ItemResponse struct {
	ID          int                  `json:"id"`
	Name        string               `json:"name"`
	Price       int                  `json:"price"`
	IsActive    bool                 `json:"is_active"`
	RefItemType *RefItemTypeResponse `json:"ref_item_type,omitempty"`
}

type ItemUpdateRequest struct {
	Name          string `json:"name"`
	Price         int    `json:"price"`
	IsActive      bool   `json:"is_active"`
	RefItemTypeID int    `json:"ref_item_type_id"`
}

type FilterItemRequest struct {
	PaginationRequest
	Name  string `query:"name"`
	Price int    `query:"price"`
}
