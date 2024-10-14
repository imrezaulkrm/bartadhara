package database

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "log"
    //"fmt"
)

// DB ভেরিয়েবল ডেটাবেস সংযোগের জন্য
var db *sql.DB

// নিউজ স্ট্রাক্ট
type News struct {
    ID          int    `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
    Image       string `json:"image"`
    Category    string `json:"category"`
    Date        string `json:"date"`
}

// এখানে User স্ট্রাক্ট ডিফাইন করুন
type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
    Picture  string `json:"picture"`
}

// ডেটাবেসের সাথে সংযোগ স্থাপন করা
func ConnectDB() {
    var err error
    db, err = sql.Open("mysql", "root:reza1234@tcp(127.0.0.1:3306)/news_db")
    if err != nil {
        log.Fatal(err)
    }
    
    err = db.Ping() // ডেটাবেসে পিং
    if err != nil {
        log.Fatal("Cannot connect to database:", err)
    }
    log.Println("Connected to database successfully")
}

// সব নিউজ ফেচ করার জন্য ফাংশন
func FetchAllNews() ([]News, error) {
    log.Println("Fetching all news from the database...")
    var newsList []News
    query := "SELECT id, title, description, image, category, date FROM news"

    rows, err := db.Query(query)
    if err != nil {
        log.Println("Error fetching news:", err)
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var news News
        if err := rows.Scan(&news.ID, &news.Title, &news.Description, &news.Image, &news.Category, &news.Date); err != nil {
            log.Println("Error scanning news row:", err)
            return nil, err
        }
        newsList = append(newsList, news)
    }
    log.Println("Fetched news successfully:", newsList)
    return newsList, nil
}

// নির্দিষ্ট আইডি দ্বারা নিউজ ফেচ করার জন্য ফাংশন
func FetchNewsByID(id int) (News, error) {
    log.Printf("Fetching news with ID: %d", id)
    var news News
    query := "SELECT id, title, description, image, category, date FROM news WHERE id = ?"

    err := db.QueryRow(query, id).Scan(&news.ID, &news.Title, &news.Description, &news.Image, &news.Category, &news.Date)
    if err != nil {
        log.Println("Error fetching news by ID:", err)
        return news, err
    }

    log.Println("Fetched news successfully:", news)
    return news, nil
}

// নতুন নিউজ ইনসার্ট করার জন্য ফাংশন
func InsertNews(news News) error {
    query := "INSERT INTO news (title, description, image, category, date) VALUES (?, ?, ?, ?, ?)"
    _, err := db.Exec(query, news.Title, news.Description, news.Image, news.Category, news.Date)
    return err
}

// নিউজ আপডেট করার জন্য ফাংশন
func UpdateNews(id int, news News) error {
    query := "UPDATE news SET title = ?, description = ?, image = ?, category = ?, date = ? WHERE id = ?"
    _, err := db.Exec(query, news.Title, news.Description, news.Image, news.Category, news.Date, id)
    return err
}

// নিউজ ডিলিট করার জন্য ফাংশন
func DeleteNews(id int) error {
    query := "DELETE FROM news WHERE id = ?"
    _, err := db.Exec(query, id)
    return err
}

// ------------------------ ইউজারের কাজ শুরু হচ্ছে ------------------------

// নির্দিষ্ট আইডি দ্বারা ব্যবহারকারী ফেচ করার জন্য ফাংশন
// FetchUserByID ফাংশন
func FetchUserByID(id int) (User, error) {
    log.Printf("Fetching user with ID: %d", id)
    var user User
    query := "SELECT id, username, email, password, picture FROM users WHERE id = ?"

    err := db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Picture)
    if err != nil {
        log.Println("Error fetching user by ID:", err)
        return user, err
    }

    log.Println("Fetched user successfully:", user)
    return user, nil
}


// ইউজারের প্রোফাইল পিকচার আপলোড করার জন্য ফাংশন
func UploadUserPicture(userID int, pictureURL string) error {
    query := "UPDATE users SET picture = ? WHERE id = ?"
    _, err := db.Exec(query, pictureURL, userID)
    return err
}

// নতুন ব্যবহারকারী ইনসার্ট করার জন্য ফাংশন
func InsertUser(user User) error {
    query := "INSERT INTO users (username, email, password, picture) VALUES (?, ?, ?, ?)"
    _, err := db.Exec(query, user.Username, user.Email, user.Password, user.Picture)
    return err
}

// সব ব্যবহারকারী ফেচ করার জন্য ফাংশন
func FetchAllUsers() ([]User, error) {
    log.Println("Fetching all users from the database...")
    var userList []User
    query := "SELECT id, username, email, picture FROM users"

    rows, err := db.Query(query)
    if err != nil {
        log.Println("Error fetching users:", err)
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var user User
        if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Picture); err != nil {
            log.Println("Error scanning user row:", err)
            return nil, err
        }
        userList = append(userList, user)
    }
    log.Println("Fetched users successfully:", userList)
    return userList, nil
}

// ব্যবহারকারী ডিলিট করার জন্য ফাংশন
func DeleteUser(id int) error {
    query := "DELETE FROM users WHERE id = ?"
    _, err := db.Exec(query, id)
    return err
}

// UpdateUser - ব্যবহারকারীর তথ্য আপডেট করার জন্য ফাংশন
func UpdateUser(id int, user User) error {
    query := "UPDATE users SET username = ?, email = ?, password = ?, picture = ? WHERE id = ?"
    _, err := db.Exec(query, user.Username, user.Email, user.Password, user.Picture, id)
    return err
}