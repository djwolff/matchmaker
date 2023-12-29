package models

import (
	"time"
)

type Game struct {
	ID        string    `gorm:"primaryKey; not null; unique" json:"id"`
	Title     string    `gorm:"not null" json:"title"`
	CreatedAt time.Time `gorm:"default:current_timestamp"`
	UpdatedAt time.Time `gorm:"default:current_timestamp"`
}
