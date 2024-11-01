package routes

import (
    "github.com/gorilla/mux"
    "github.com/imrezaulkrm/bartadhara/controllers"
)

func NewsRoutes(r *mux.Router) {
    r.HandleFunc("/news", controllers.GetAllNews).Methods("GET")                // সব নিউজ পেতে
    r.HandleFunc("/news/{id}", controllers.GetNewsByID).Methods("GET")          // নির্দিষ্ট আইডি দ্বারা নিউজ পেতে
    r.HandleFunc("/news", controllers.CreateNews).Methods("POST")                // নতুন নিউজ তৈরি করতে
    r.HandleFunc("/news/{id}", controllers.UpdateNews).Methods("PUT")            // নিউজ আপডেট করতে
    r.HandleFunc("/news/{id}", controllers.DeleteNews).Methods("DELETE")         // নিউজ মুছতে
}

