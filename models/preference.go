package models

type Preferences struct {
	ID                 string       `gorm:"primaryKey;type:varchar(191)" json:"id"`
	StartDate          string       `gorm:"type:date" json:"startDate"`
	EndDate            string       `gorm:"type:date" json:"endDate"`
	SelectedPreference []Preference `gorm:"many2many:preferences_preferences;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"selectedPreference"`
}

type Preference struct {
	ID    string `gorm:"primaryKey;type:varchar(191)" json:"id"`
	Type  string `gorm:"type:varchar(191)" json:"type"`
	Image string `gorm:"type:varchar(191)" json:"image"`
}

type PreferenceLink struct {
	PreferencesID string `gorm:"primaryKey;type:varchar(191)" json:"preferencesId"`
	PreferenceID  string `gorm:"primaryKey;type:varchar(191)" json:"preferenceId"`
	Type          string `gorm:"type:varchar(191)" json:"type"`
	Image         string `gorm:"type:varchar(191)" json:"image"`
}
