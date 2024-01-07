package models

type Book struct {
	ID        uint    `gorm:"primary key;autoIncrement" json:"id"`
	Title     *string `json:"title"`
	Author    *string `json:"author"`
	Publisher *string `json:"publisher"`
}
