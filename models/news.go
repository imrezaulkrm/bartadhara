package models

type News struct {
    ID          int    `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
    Image       string `json:"image"`
    Category    string `json:"category"`
    Date        string `json:"date"`
}
