package models

import "time"

type Timeline struct {
	ID           string         `gorm:"primaryKey" json:"id"`
	Town         string         `json:"town"`
	Country      string         `json:"country"`
	Budget       string         `json:"budget"`
	StartDate    time.Time      `json:"startDate"`
	EndDate      time.Time      `json:"endDate"`
	PlaceVisited []PlaceVisited `gorm:"foreignKey:TimelineID" json:"placeVisited"`
}

type DateData struct {
	Date string         `json:"date"`
	Data []PlaceVisited `gorm:"type:json" json:"data"`
}

type PlaceVisited struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	TimelineID string `gorm:"index" json:"timelineId"` // Foreign Key to Timeline
	Type       string `json:"type"`
	Landmark   string `json:"landmark"`
	RoadName   string `json:"roadName"`
	Town       string `json:"town"`
	Time       string `json:"time"`
	Images     string `json:"images"`
}
