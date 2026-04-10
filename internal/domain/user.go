package domain

import "time"

type User struct {
	ID       uint      `json:"id" gorm:"primaryKey;autoIncrement;not null"`
	Name     string    `json:"name" gorm:"not null"`
	Email    string    `json:"email" gorm:"uniqueIndext;not null"`
	Password string    `json:"-" gorm:"not null"`
	CreateAt time.Time `json:"created_at"`
	UpdateAt time.Time `json:"updated_at"`
}
