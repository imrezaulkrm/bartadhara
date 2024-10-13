package main

import (
    "log"
    "net/http"
    "news_backend/database"
    "news_backend/routes"
    "github.com/gorilla/mux" // মাক্স প্যাকেজটি ইমপোর্ট করা
)

func main() {
    // ডেটাবেস কানেকশন স্থাপন করা
    database.ConnectDB()
    
    r := mux.NewRouter()
    routes.NewsRoutes(r)

    log.Fatal(http.ListenAndServe(":8000", r))
}

