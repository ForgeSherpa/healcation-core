package models

import (
	"time"
)

type History struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         uint      `json:"user_id"`
	Country        string    `json:"country"`
	Town           string    `json:"town"`
	Title          string    `json:"title"`
	StartDate      time.Time `json:"startDate"`
	EndDate        time.Time `json:"endDate"`
	BudgetMin      int       `json:"budgetMin"`
	BudgetMax      int       `json:"budgetMax"`
	Accommodations string    `json:"accommodations" gorm:"type:json"`
	Timelines      string    `json:"timelines" gorm:"type:json"`
	Image          string    `json:"image" gorm:"type:json"`
}
