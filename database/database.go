package database

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql" // MySQL ড্রাইভার, সঠিকভাবে ব্যবহার করা হয়েছে
    "log"
)

var db *sql.DB

// তোমার নিউজ স্ট্রাক্ট
type News struct {
    ID          int    `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
    Image       string `json:"image"`
    Category    string `json:"category"`
    Date        string `json:"date"`
}

// ডেটাবেসের সাথে সংযোগ স্থাপন করা
func ConnectDB() {
    var err error
    db, err = sql.Open("mysql", "root:reza1234@tcp(127.0.0.1:3306)/news_db") // এখানে তোমার ডেটাবেস ক্রিডেনশিয়ালস ব্যবহার করো
    if err != nil {
        log.Fatal(err)
    }
    
    err = db.Ping() // এটি ডেটাবেসে পিং করে দেখবে
    if err != nil {
        log.Fatal("Cannot connect to database:", err)
    }
    log.Println("Connected to database successfully") // সংযোগ সফল হলে লগ করো
}

// সব নিউজ ফেচ করার জন্য ফাংশন
func FetchAllNews() ([]News, error) {
    log.Println("Fetching all news from the database...") // লগিং যুক্ত করা হলো
    var newsList []News
    query := "SELECT id, title, description, image, category, date FROM news"

    rows, err := db.Query(query)
    if err != nil {
        log.Println("Error fetching news:", err) // ত্রুটি লগ করা হচ্ছে
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var news News
        if err := rows.Scan(&news.ID, &news.Title, &news.Description, &news.Image, &news.Category, &news.Date); err != nil {
            log.Println("Error scanning news row:", err) // ত্রুটি লগ করা হচ্ছে
            return nil, err
        }
        newsList = append(newsList, news)
    }
    log.Println("Fetched news successfully:", newsList) // সফলভাবে নিউজ ফেচ করার পর লগ
    return newsList, nil
}

// নির্দিষ্ট আইডি দ্বারা নিউজ ফেচ করার জন্য ফাংশন
func FetchNewsByID(id string) (News, error) {
    log.Printf("Fetching news with ID: %s", id) // লগিং যুক্ত করা হলো
    var news News
    query := "SELECT id, title, description, image, category, date FROM news WHERE id = ?"

    err := db.QueryRow(query, id).Scan(&news.ID, &news.Title, &news.Description, &news.Image, &news.Category, &news.Date)
    if err != nil {
        log.Println("Error fetching news by ID:", err) // ত্রুটি লগ করা হচ্ছে
        return news, err
    }

    log.Println("Fetched news successfully:", news) // সফলভাবে নিউজ ফেচ করার পর লগ
    return news, nil
}

// নতুন নিউজ ইনসার্ট করার জন্য ফাংশন
func InsertNews(news News) error {
    query := "INSERT INTO news (title, description, image, category, date) VALUES (?, ?, ?, ?, ?)"
    _, err := db.Exec(query, news.Title, news.Description, news.Image, news.Category, news.Date)
    return err
}

func UpdateNews(id string, news News) error {
    query := "UPDATE news SET title = ?, description = ?, image = ?, category = ?, date = ? WHERE id = ?"
    _, err := db.Exec(query, news.Title, news.Description, news.Image, news.Category, news.Date, id)
    return err
}

// নিউজ ডিলিট করার জন্য ফাংশন
func DeleteNews(id string) error {
    query := "DELETE FROM news WHERE id = ?"
    _, err := db.Exec(query, id)
    return err
}

