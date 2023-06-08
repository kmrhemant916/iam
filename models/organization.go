package models

type Organization struct {
	ID uint `gorm:"primary_key;autoIncrement;" json:"id"`
	Name  string `gorm:"unique" json:"name"`
}