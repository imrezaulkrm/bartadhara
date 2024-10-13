package models

type News struct {
    ID          uint   `json:"id" gorm:"primaryKey"`
    Title       string `json:"title"`
    Description string `json:"description"`
    Image       string `json:"image"`
    Category    string `json:"category"`
    Date        string `json:"date"`
}
