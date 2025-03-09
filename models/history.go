package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type History struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Country   string    `json:"country"`
	Town      string    `json:"town"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	Image     string    `json:"image" gorm:"type:json"`
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
