package models

import (
	"time"
)

type Match struct {
	ID        string    `gorm:"primaryKey; not null; unique" json:"id"`
	Mode      string    `gorm:"not null" json:"mode"`
	CreatedAt time.Time `gorm:"default:current_timestamp"`
	UpdatedAt time.Time `gorm:"default:current_timestamp"`
	Game      Game      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	GameID    int
}
