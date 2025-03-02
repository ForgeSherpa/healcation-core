package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type History struct {
	AutoID    uint        `gorm:"primaryKey;autoIncrement" json:"-"`
	ID        string      `json:"id" gorm:"uniqueIndex"`
	Country   string      `json:"country"`
	Town      string      `json:"town"`
	StartDate time.Time   `json:"startDate"`
	EndDate   time.Time   `json:"endDate"`
	Image     StringArray `json:"image" gorm:"type:json"`
}

func (h *History) BeforeCreate(tx *gorm.DB) (err error) {
	var lastID int64
	tx.Model(&History{}).Select("coalesce(max(auto_id),0)").Scan(&lastID)
	h.AutoID = uint(lastID + 1)
	h.ID = fmt.Sprintf("%d", h.AutoID)
	return
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
