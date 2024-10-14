package models

type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"` // পাসওয়ার্ড সাধারণত এনক্রিপ্ট করা উচিত
    Picture  string `json:"picture"`  // পিকচার URL
}
