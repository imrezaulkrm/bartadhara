package main

import (
    "news_backend/database"
    "news_backend/routes"
    "github.com/gorilla/mux"
    "net/http"
    "log"
)

func main() {
    database.ConnectDB() // ডেটাবেস কানেকশন
    r := mux.NewRouter()

    // রাউট সেটআপ
    routes.NewsRoutes(r)
    routes.UserRoutes(r)

    // সার্ভার শুরু
    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatal(err)
    }
}
