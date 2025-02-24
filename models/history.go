package models

import "time"

type History struct {
	ID                   uint                   `gorm:"primaryKey" json:"id"`
	Country              string                 `json:"country"`
	Town                 string                 `json:"town"`
	StartDate            time.Time              `json:"startDate"`
	EndDate              time.Time              `json:"endDate"`
	Image                string                 `json:"image"`
	Description          string                 `json:"description"`
	SelectedAccomodation []SelectedAccomodation `json:"selectedAccomodation" gorm:"foreignKey:HistoryID"`
	SelectedPlaces       []SelectedPlace        `json:"selectedPlaces" gorm:"foreignKey:HistoryID"`
}

type SelectedAccomodation struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Name      string `json:"name"`
	Image     string `json:"image"`
	HistoryID uint   `json:"-"`
}

type SelectedPlace struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	PlaceToVisit string `json:"placeToVisit"`
	Town         string `json:"town"`
	Image        string `json:"image"`
	HistoryID    uint   `json:"-"`
}
