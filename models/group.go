package models

type Group struct {
	// gorm.Model
	ID 	  uint `gorm:"primary_key;autoIncrement;" json:"id"`
	Name  string `gorm:"unique" json:"name"`
}