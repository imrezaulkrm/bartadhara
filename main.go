package main

import (
    "log"
    "net/http"
    "news_backend/database"
    "news_backend/routes"
    "github.com/gorilla/mux"
)

// CORS Middleware
func enableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*") // সব ডোমেইন থেকে অ্যাক্সেসের অনুমতি
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        // যদি OPTIONS মেথড হয় (preflight request), তাহলে তা হ্যান্ডেল করা
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}

func main() {
    database.ConnectDB() // ডেটাবেস কানেকশন
    r := mux.NewRouter()

    // রাউট সেটআপ
    routes.NewsRoutes(r)
    routes.UserRoutes(r)

    // CORS middleware যুক্ত করে সার্ভার চালু করা
    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", enableCORS(r)); err != nil {
        log.Fatal(err)
    }
}

