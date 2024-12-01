package controllers

import (
    "github.com/gorilla/mux"
    "encoding/json"
    "net/http"
    "github.com/imrezaulkrm/bartadhara-backend/database"
    "github.com/imrezaulkrm/bartadhara-backend/models"
    "strconv"
    "fmt"
    "os"
    "io"
    "time"
)

// সমস্ত নিউজ ফেচ করার জন্য ফাংশন
func GetAllNews(w http.ResponseWriter, r *http.Request) {
    news, err := database.FetchAllNews()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(news)
}

// নির্দিষ্ট নিউজ ফেচ করার জন্য ফাংশন
func GetNewsByID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    news, err := database.FetchNewsByID(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(news)
}


func CreateNews(w http.ResponseWriter, r *http.Request) {
    // Parse form data
    err := r.ParseMultipartForm(10 << 20) // 10 MB limit
    if err != nil {
        http.Error(w, "Unable to parse form", http.StatusBadRequest)
        return
    }

    // Retrieve news details from form
    title := r.FormValue("title")
    description := r.FormValue("description")
    category := r.FormValue("category")
    date := r.FormValue("date")

    // Handle file upload for the news image
    var imageURL string
    if file, _, err := r.FormFile("image"); err == nil {
        defer file.Close()

        // Generate a unique image name using title and date
        // Format: "title-YYYY-MM-DD.jpg"
        imageName := fmt.Sprintf("%s-%s.jpg", title, date)

        // File path with the generated image name
        imagePath := fmt.Sprintf("uploads/news/%s", imageName)
        
        // Ensure the directory exists
        if err := os.MkdirAll("uploads/news", os.ModePerm); err != nil {
            http.Error(w, "Unable to create uploads/news directory", http.StatusInternalServerError)
            return
        }

        // Create the file for saving
        out, err := os.Create(imagePath)
        if err != nil {
            http.Error(w, "Unable to create the file for writing", http.StatusInternalServerError)
            return
        }
        defer out.Close()

        // Save the uploaded image file
        if _, err = io.Copy(out, file); err != nil {
            http.Error(w, "Error saving the file", http.StatusInternalServerError)
            return
        }

        // Set image URL for accessing the file
        imageURL = fmt.Sprintf("http://localhost:8080/%s", imagePath)
    }

    // Parse and validate date
    parsedDate, err := time.Parse("2006-01-02", date)
    if err != nil {
        http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
        return
    }

    // Create news model
    news := database.News{ // Use database.News instead of models.News
        Title:       title,
        Description: description,
        Image:       imageURL,
        Category:    category,
        Date:        parsedDate.Format("2006-01-02"),
    }

    // Insert news into the database
    if err = database.InsertNews(news); err != nil {
        http.Error(w, "Failed to insert news", http.StatusInternalServerError)
        return
    }

    // Send success response
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(news)
}


func UpdateNews(w http.ResponseWriter, r *http.Request) {
    // Get the news ID from URL
    newsID := mux.Vars(r)["id"]
    
    // Parse form data
    err := r.ParseMultipartForm(10 << 20) // 10 MB limit
    if err != nil {
        http.Error(w, "Unable to parse form", http.StatusBadRequest)
        return
    }

    // Retrieve updated news details from form
    title := r.FormValue("title")
    description := r.FormValue("description")
    category := r.FormValue("category")
    date := r.FormValue("date")

    var imageURL string
    // Check if a new image is provided
    if file, _, err := r.FormFile("image"); err == nil {
        defer file.Close()

        // File path with title and date as identifier
        imagePath := fmt.Sprintf("uploads/news/%s_%s.jpg", title, date)
        if err := os.MkdirAll("uploads/news", os.ModePerm); err != nil {
            http.Error(w, "Unable to create uploads/news directory", http.StatusInternalServerError)
            return
        }

        out, err := os.Create(imagePath)
        if err != nil {
            http.Error(w, "Unable to create the file for writing", http.StatusInternalServerError)
            return
        }
        defer out.Close()
        if _, err = io.Copy(out, file); err != nil {
            http.Error(w, "Error saving the file", http.StatusInternalServerError)
            return
        }

        imageURL = fmt.Sprintf("http://localhost:8080/%s", imagePath)
    } else {
        // If no new image is provided, fetch the existing image URL from the database
        existingNews, err := database.GetNewsByID(newsID)
        if err != nil {
            http.Error(w, "News not found", http.StatusNotFound)
            return
        }
        imageURL = existingNews.Image // Use the old image URL if no new image is uploaded
    }

    // Parse and validate date
    parsedDate, err := time.Parse("2006-01-02", date)
    if err != nil {
        http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
        return
    }

    // Update news in the database
    updatedNews := models.News{
        Title:       title,
        Description: description,
        Image:       imageURL,
        Category:    category,
        Date:        parsedDate.Format("2006-01-02"),
    }

    // Call the database function to update the news
    if err := database.UpdateNews(newsID, updatedNews); err != nil {
        http.Error(w, "Failed to update news", http.StatusInternalServerError)
        return
    }

    // Send success response
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(updatedNews)
}


// নিউজ ডিলিট করার জন্য ফাংশন
func DeleteNews(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    if err := database.DeleteNews(id); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}
