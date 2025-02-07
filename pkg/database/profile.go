package database

import (
	"fmt"
	"os"

	"gorm.io/gorm"
)

type Profile struct {
	gorm.Model
	UserID              uint
	APIActive           bool        `form:"api_active"`
	Language            string      `form:"language"`
	TotalsShow          WorkoutType `form:"totals_show"`
	Timezone            string      `form:"timezone"`
	AutoImportDirectory string      `form:"auto_import_directory"`
	SocialsDisabled     bool        `form:"socials_disabled"`

	PreferredUnits UserPreferredUnits `gorm:"serializer:json"`

	User *User `gorm:"foreignKey:UserID" json:"-"`
}

type UserPreferredUnits struct {
	SpeedRaw     string `form:"speed" json:"speed"`
	DistanceRaw  string `form:"distance" json:"distance"`
	ElevationRaw string `form:"elevation" json:"elevation"`
}

func (u UserPreferredUnits) Tempo() string {
	return "min/" + u.Distance()
}

func (u UserPreferredUnits) Elevation() string {
	switch u.ElevationRaw {
	case "ft":
		return "ft"
	default:
		return "m"
	}
}

func (u UserPreferredUnits) Distance() string {
	switch u.DistanceRaw {
	case "mi":
		return "mi"
	default:
		return "km"
	}
}

func (u UserPreferredUnits) Speed() string {
	switch u.SpeedRaw {
	case "mph":
		return "mph"
	default:
		return "km/h"
	}
}

func (p *Profile) Save(db *gorm.DB) error {
	return db.Save(p).Error
}

func (p *Profile) CanImportFromDirectory() (bool, error) {
	if p == nil {
		return false, nil
	}

	if p.AutoImportDirectory == "" {
		return false, nil
	}

	info, err := os.Stat(p.AutoImportDirectory)
	if err != nil {
		return false, err
	}

	if !info.IsDir() {
		return false, fmt.Errorf("%v is not a directory", p.AutoImportDirectory)
	}

	return true, nil
}
