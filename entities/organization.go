package entities

type Organization struct {
	ID uint `gorm:"primary_key;autoIncrement;"`
	Name  string `gorm:"unique;not null"`
}