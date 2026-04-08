package entities

import (
	"time"
)

type Item struct {
	ID            int        `json:"id"`
	Name          string     `json:"name"`
	Price         int        `json:"price"`
	IsActive      bool       `json:"is_active"`
	RefItemTypeID int        `json:"ref_item_type_id"`
	CreatedAt     time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Item) TableName() string {
	return "items"
}

type ItemWithRefItemType struct {
	Item
	RefItemType RefItemType `json:"ref_item_type"`
}

func (ItemWithRefItemType) TableName() string {
	return "items"
}