package models

type AvailableLater struct {
	User         User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID       int
	GameMode     GameMode `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	GameModeID   int
	SeekingRole  string
	OfferingRole string
}
