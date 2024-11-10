package models

import (
    "errors"
    "regexp"
)
// Admin স্ট্রাকচার
type Admin struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
    Picture  string `json:"picture"`
}
// Validate ফাংশন এডমিনের ইনপুট যাচাই করে
func (a *Admin) Validate() error {
    // নামের ভ্যালিডেশন
    if len(a.Name) == 0 {
        return errors.New("নাম অবশ্যই প্রদান করতে হবে")
    }

    // ব্যবহারকারীর নামের ভ্যালিডেশন
    if len(a.Username) < 3 {
        return errors.New("ব্যবহারকারীর নাম কমপক্ষে ৩ অক্ষরের হতে হবে")
    }

    // ইমেইলের ভ্যালিডেশন
    if !isValidAdminEmail(a.Email) {
        return errors.New("ইমেইল ঠিক নয়")
    }

    // পাসওয়ার্ডের ভ্যালিডেশন
    if len(a.Password) < 6 {
        return errors.New("পাসওয়ার্ড কমপক্ষে ৬ অক্ষরের হতে হবে")
    }

    return nil
}
// isValidEmail ফাংশন ইমেইল ঠিক কিনা যাচাই করে
func isValidAdminEmail(email string) bool {
    regex := `^[a-zA-Z0-9._%+-]+@[a-zAZ0-9.-]+\.[a-zA-Z]{2,}$`
    re := regexp.MustCompile(regex)
    return re.MatchString(email)
}
