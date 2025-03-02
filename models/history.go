package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type History struct {
	ID                   uint                   `gorm:"primaryKey" json:"id"`
	Country              string                 `json:"country"`
	Town                 string                 `json:"town"`
	StartDate            time.Time              `json:"startDate"`
	EndDate              time.Time              `json:"endDate"`
	Image                StringArray            `json:"image" gorm:"type:text"`
	Description          string                 `json:"description"`
	SelectedAccomodation []SelectedAccomodation `json:"selectedAccomodation" gorm:"foreignKey:HistoryID"`
	SelectedPlaces       []SelectedPlace        `json:"selectedPlaces" gorm:"foreignKey:HistoryID"`
}

type SelectedAccomodation struct {
	ID        uint        `gorm:"primaryKey" json:"id"`
	Name      string      `json:"name"`
	Image     StringArray `json:"image" gorm:"type:text"`
	HistoryID uint        `json:"-"`
}

type SelectedPlace struct {
	ID           uint        `gorm:"primaryKey" json:"id"`
	PlaceToVisit string      `json:"placeToVisit"`
	Town         string      `json:"town"`
	Image        StringArray `json:"image" gorm:"type:text"`
	HistoryID    uint        `json:"-"`
}

type StringArray []string

func (s StringArray) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = []string{}
		return nil
	}

	switch v := value.(type) {
	case string:
		return json.Unmarshal([]byte(v), s)
	case []byte:
		return json.Unmarshal(v, s)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
}
