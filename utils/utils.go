package utils

import (
    "golang.org/x/crypto/bcrypt"
    "log"
)
// পাসওয়ার্ড হ্যাশ করার ফাংশন
func HashPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        log.Println("Error hashing password:", err)
        return "", err
    }
    return string(hashedPassword), nil
}
// পাসওয়ার্ড কমপেয়ার করার ফাংশন
func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
