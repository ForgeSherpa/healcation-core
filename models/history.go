package models

import (
	"time"
)

type History struct {
	ID        string    `json:"id"`
	Town      string    `json:"town"`
	Country   string    `json:"country"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	Image     string    `json:"image"`
	UserID    string    `json:"user_id"`
	User      User      `json:"user" gorm:"foreignKey:UserID;references:ID"`
}
