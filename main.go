package main

import (
    "log"
    "net/http"
    "github.com/imrezaulkrm/bartadhara-backend/database"
    "github.com/imrezaulkrm/bartadhara-backend/routes"
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
    // ডেটাবেস কানেকশন
    database.ConnectDB()

    // রাউট সেটআপ
    r := mux.NewRouter()

    r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads/"))))

    // News, User এবং Admin রাউটস সেটআপ
    routes.NewsRoutes(r)
    routes.UserRoutes(r)
    routes.AdminRoutes(r) // Admin রাউট যোগ করা

    // CORS middleware যুক্ত করে সার্ভার চালু করা
    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", enableCORS(r)); err != nil {
        log.Fatal(err)
    }
}
