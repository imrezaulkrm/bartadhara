package database

import (
    "database/sql"
    "log"
    "errors"
    _ "github.com/go-sql-driver/mysql" // MySQL ড্রাইভার
    "news_backend/models"               // এখানে models প্যাকেজটি আমদানি করুন
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
    Name     string `json:"name"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
    Picture  string `json:"picture"`
}

// ConnectDB ডেটাবেসের সাথে সংযোগ স্থাপন করে
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

// UpdateUserCategories updates the categories for a specific user
func UpdateUserCategories(userID int, categories []string) error {
    // Start a transaction
    tx, err := db.Begin() // Use `db` instead of `DB`
    if err != nil {
        return err
    }

    // Prepare the query to delete existing categories
    deleteQuery := "DELETE FROM user_categories WHERE user_id = ?"
    if _, err := tx.Exec(deleteQuery, userID); err != nil {
        tx.Rollback()
        return err
    }

    // Prepare the query to insert new categories
    insertQuery := "INSERT INTO user_categories (user_id, category) VALUES (?, ?)"
    for _, category := range categories {
        if _, err := tx.Exec(insertQuery, userID, category); err != nil {
            tx.Rollback()
            return err
        }
    }

    // Commit the transaction
    return tx.Commit()
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

/// ------------------------ ইউজারের কাজ শুরু হচ্ছে ------------------------

// FetchAllUsers ফাংশন সব ব্যবহারকারী ফেচ করে
func FetchAllUsers() ([]models.User, error) {
    log.Println("Fetching all users from the database...")
    var userList []models.User
    query := "SELECT id, name, username, email, password, picture FROM users" // picture যুক্ত করুন

    rows, err := db.Query(query)
    if err != nil {
        log.Println("Error fetching users:", err)
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var user models.User
        if err := rows.Scan(&user.ID, &user.Name, &user.Username, &user.Email, &user.Password, &user.Picture); err != nil { // picture যুক্ত করুন
            log.Println("Error scanning user row:", err)
            return nil, err
        }
        userList = append(userList, user)
    }
    log.Println("Fetched users successfully:", userList)
    return userList, nil
}


// FetchUserByID ফাংশন ব্যবহারকারী আইডি দ্বারা ব্যবহারকারী ফেচ করে
// FetchUserByID ফাংশন ব্যবহারকারী আইডি দ্বারা ব্যবহারকারী ফেচ করে
func FetchUserByID(id int) (models.User, error) {
    log.Printf("Fetching user with ID: %d", id)
    var user models.User
    query := "SELECT id, name, username, email, password, picture FROM users WHERE id = ?" // picture যুক্ত করুন

    err := db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Username, &user.Email, &user.Password, &user.Picture) // picture যুক্ত করুন
    if err != nil {
        log.Println("Error fetching user by ID:", err)
        return user, err
    }

    log.Println("Fetched user successfully:", user)
    return user, nil
}


// FetchUserByUsername retrieves a user by username from the database
// FetchUserByUsername retrieves a user by username from the database
func FetchUserByUsername(username string) (*models.User, error) {
    var user models.User
    query := "SELECT id, name, username, email, password, picture FROM users WHERE username = ?" // picture যুক্ত করুন
    
    err := db.QueryRow(query, username).Scan(&user.ID, &user.Name, &user.Username, &user.Email, &user.Password, &user.Picture) // picture যুক্ত করুন
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("user not found")
        }
        return nil, err
    }
    
    return &user, nil
}



// FetchUserByEmail ফাংশন ব্যবহারকারীর ইমেইল দ্বারা ব্যবহারকারী সন্ধান করে
// FetchUserByEmail ফাংশন ব্যবহারকারীর ইমেইল দ্বারা ব্যবহারকারী সন্ধান করে
func FetchUserByEmail(email string) (*models.User, error) {
    var user models.User
    err := db.QueryRow("SELECT id, name, username, email, password, picture FROM users WHERE email = ?", email).Scan(&user.ID, &user.Name, &user.Username, &user.Email, &user.Password, &user.Picture) // picture যুক্ত করুন
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil // কোন ব্যবহারকারী পাওয়া যায়নি
        }
        return nil, err // অন্য একটি ত্রুটি ঘটেছে
    }
    return &user, nil
}


// InsertUser ফাংশন নতুন ব্যবহারকারী যুক্ত করে
func InsertUser(user models.User) error {
    _, err := db.Exec("INSERT INTO users (name, username, email, password, picture) VALUES (?, ?, ?, ?, ?)", user.Name, user.Username, user.Email, user.Password, user.Picture) // picture যুক্ত করুন
    return err
}

// UpdateUser ফাংশন বিদ্যমান ব্যবহারকারী আপডেট করে
// UpdateUser ফাংশন বিদ্যমান ব্যবহারকারী আপডেট করে
func UpdateUser(id int, user models.User) error {
    _, err := db.Exec("UPDATE users SET name = ?, username = ?, email = ?, password = ?, picture = ? WHERE id = ?", user.Name, user.Username, user.Email, user.Password, user.Picture, id) // picture যুক্ত করুন
    return err
}


// DeleteUser ফাংশন ব্যবহারকারী ডিলিট করে
func DeleteUser(id int) error {
    query := "DELETE FROM users WHERE id = ?"
    _, err := db.Exec(query, id)
    return err
}

// SaveUserCategories ফাংশন ব্যবহারকারীর ক্যাটেগরিগুলি সংরক্ষণ করে
func SaveUserCategories(userID int, categories []string) error {
    // প্রথমে ক্যাটেগরি মুছে ফেলুন
    _, err := db.Exec("DELETE FROM user_categories WHERE user_id = ?", userID)
    if err != nil {
        return err
    }

    // নতুন ক্যাটেগরি যোগ করুন
    for _, category := range categories {
        _, err = db.Exec("INSERT INTO user_categories (user_id, category) VALUES (?, ?)", userID, category)
        if err != nil {
            return err
        }
    }
    return nil
}

// FetchUserCategories ফাংশন নির্দিষ্ট ব্যবহারকারীর জন্য ক্যাটেগরি সংগ্রহ করে
func FetchUserCategories(userID int) ([]string, error) {
    rows, err := db.Query("SELECT category FROM user_categories WHERE user_id = ?", userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var categories []string
    for rows.Next() {
        var category string
        if err := rows.Scan(&category); err != nil {
            return nil, err
        }
        categories = append(categories, category)
    }
    return categories, nil
}
// FetchUserByUsernameOrEmail retrieves a user by username or email from the database
// FetchUserByUsernameOrEmail retrieves a user by username or email from the database
func FetchUserByUsernameOrEmail(username, email string) (*models.User, error) {
    var user models.User
    query := "SELECT id, name, username, email, password, picture FROM users WHERE username = ? OR email = ?" // picture যুক্ত করুন
    
    err := db.QueryRow(query, username, email).Scan(&user.ID, &user.Name, &user.Username, &user.Email, &user.Password, &user.Picture) // picture যুক্ত করুন
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("user not found")
        }
        return nil, err
    }
    
    return &user, nil
}
