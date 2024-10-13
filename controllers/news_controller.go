package controllers

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "news_backend/database" // ডেটাবেস প্যাকেজ
)

// সব নিউজ ফেচ করার জন্য ফাংশন
func GetAllNews(w http.ResponseWriter, r *http.Request) {
    newsList, err := database.FetchAllNews()
    if err != nil {
        http.Error(w, "Unable to fetch news", http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(newsList)
}

// নির্দিষ্ট আইডি দ্বারা নিউজ ফেচ করার জন্য ফাংশন
func GetNewsByID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    // ডেটাবেস থেকে আইডি দ্বারা নিউজ খোঁজা
    news, err := database.FetchNewsByID(id)
    if err != nil {
        http.Error(w, "News not found", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(news)
}

// নতুন নিউজ ইনসার্ট করার জন্য ফাংশন
func CreateNews(w http.ResponseWriter, r *http.Request) {
    var news database.News
    err := json.NewDecoder(r.Body).Decode(&news)
    if err != nil {
        http.Error(w, "Unable to decode JSON", http.StatusBadRequest)
        return
    }

    err = database.InsertNews(news)
    if err != nil {
        http.Error(w, "Unable to insert news", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(news)
}

func UpdateNews(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    var news database.News
    err := json.NewDecoder(r.Body).Decode(&news)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // নিউজ আপডেট করার জন্য ফাংশন কল
    err = database.UpdateNews(id, news)
    if err != nil {
        http.Error(w, "Unable to update news", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

// নিউজ ডিলিট করার জন্য ফাংশন
func DeleteNews(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    err := database.DeleteNews(id)
    if err != nil {
        http.Error(w, "Failed to delete news", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent) // 204 No Content
}
