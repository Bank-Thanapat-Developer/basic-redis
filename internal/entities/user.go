package entities

import "time"

type User struct {
	ID        int        `json:"id"`
	Username  string     `json:"username" gorm:"uniqueIndex;not null"`
	Password  string     `json:"-"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}
