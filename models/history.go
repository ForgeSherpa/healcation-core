package models

import (
	"time"
)

type History struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `json:"user_id"`
	Country   string    `json:"country"`
	Town      string    `json:"town"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	Budget    string    `json:"budget"`
	Timelines string    `json:"timelines" gorm:"type:json"`
	Image     string    `json:"image"`
}
