package models

// AdministrativeBoundary maps to the administrative_boundaries table
type AdministrativeBoundary struct {
	OBJECTID     int     `gorm:"primaryKey;column:objectid"`
	ADMIN_ID3    string  `gorm:"column:admin_id3"`
	NAME1        string  `gorm:"column:name1"`
	NAME_ENG1    string  `gorm:"column:name_eng1"`
	NAME2        string  `gorm:"column:name2"`
	NAME_ENG2    string  `gorm:"column:name_eng2"`
	NAME3        string  `gorm:"column:name3"`
	NAME_ENG3    string  `gorm:"column:name_eng3"`
	ADMIN_ID1    string  `gorm:"column:admin_id1"`
	ADMIN_ID2    string  `gorm:"column:admin_id2"`
	Type         int     `gorm:"column:type"`
	Version      string  `gorm:"column:version"`
	POP_YEAR     int     `gorm:"column:pop_year"`
	POPULATION   float64 `gorm:"column:population"`
	MALE         float64 `gorm:"column:male"`
	FEMALE       float64 `gorm:"column:female"`
	HOUSE        float64 `gorm:"column:house"`
	Shape_Area   float64 `gorm:"column:shape__area"`
	Shape_Length float64 `gorm:"column:shape__length"`
	SearchVector string  `gorm:"-"`
}

func (AdministrativeBoundary) TableName() string {
	return "administrative_boundaries"
}