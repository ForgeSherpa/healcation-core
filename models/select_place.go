package models

type SelectPlace struct {
	ID                         string                       `gorm:"primaryKey;type:varchar(191)" json:"id"`
	City                       string                       `gorm:"type:varchar(191)" json:"city"`
	Country                    string                       `gorm:"type:varchar(191)" json:"country"`
	Description                string                       `gorm:"type:text" json:"description"`
	PlaceRecommendation        []PlaceRecommendation        `gorm:"foreignKey:SelectPlaceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"placeRecommendation"`
	AccomodationRecommendation []AccomodationRecommendation `gorm:"foreignKey:SelectPlaceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"accomodationRecommendation"`
	SelectedPlace              []Time                       `gorm:"many2many:select_place_times;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"selectedPlace"`
	SelectedAccomodation       []string                     `gorm:"type:json" json:"selectedAccomodation"`
}

type PlaceRecommendation struct {
	ID            string   `gorm:"primaryKey;type:varchar(191)" json:"id"`
	PlaceToVisit  string   `gorm:"type:varchar(191)" json:"placeToVisit"`
	Town          string   `gorm:"type:varchar(191)" json:"town"`
	Image         []string `gorm:"type:json" json:"image"`
	SelectPlaceID string   `gorm:"type:varchar(191)" json:"-"`
}

type AccomodationRecommendation struct {
	ID            string   `gorm:"primaryKey;type:varchar(191)" json:"id"`
	Name          string   `gorm:"type:varchar(191)" json:"name"`
	Town          string   `gorm:"type:varchar(191)" json:"town"`
	Image         []string `gorm:"type:json" json:"image"`
	SelectPlaceID string   `gorm:"type:varchar(191)" json:"-"`
}

type Time struct {
	ID        string          `gorm:"primaryKey;type:varchar(191)" json:"id"`
	TimeOfDay string          `gorm:"type:varchar(191)" json:"timeOfDay"`
	Places    []SelectedPlace `gorm:"many2many:time_selected_places;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"places"`
}

type SelectedPlace struct {
	ID           string `gorm:"primaryKey;type:varchar(191)" json:"id"`
	PlaceToVisit string `gorm:"type:varchar(191)" json:"placeToVisit"`
	Town         string `gorm:"type:varchar(191)" json:"town"`
}
