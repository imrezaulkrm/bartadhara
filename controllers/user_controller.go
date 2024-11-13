package controllers

import (
    "log"
    "time"
    //"path/filepath"  // Add this import
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "strconv"
    "github.com/gorilla/mux"
    "golang.org/x/crypto/bcrypt"
    "github.com/imrezaulkrm/bartadhara/models"
    "github.com/imrezaulkrm/bartadhara/database"
)
// UserController struct
type UserController struct{}

// FetchAllUsers handles GET requests to fetch all users
func (uc *UserController) FetchAllUsers(w http.ResponseWriter, r *http.Request) {
    users, err := database.FetchAllUsers()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

// FetchUserByID handles GET requests to fetch a user by ID
func (uc *UserController) FetchUserByID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil || id <= 0 {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    user, err := database.FetchUserByID(id)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

// InsertUser handles user registration, picture upload, and data insertion
func (uc *UserController) InsertUser(w http.ResponseWriter, r *http.Request) {
    // Parse multipart form
    err := r.ParseMultipartForm(10 << 20) // 10 MB limit
    if err != nil {
        http.Error(w, "Unable to parse form", http.StatusBadRequest)
        return
    }

    // Create uploads/users directory if it doesn't exist
    if err := os.MkdirAll("uploads/users", os.ModePerm); err != nil {
        http.Error(w, "Unable to create uploads/users directory", http.StatusInternalServerError)
        return
    }

    // Retrieve user details from form
    name := r.FormValue("name")
    username := r.FormValue("username")
    email := r.FormValue("email")
    password := r.FormValue("password")

    // Handle file upload for the profile picture (added here)
    var picturePath string
    var imageURL string
    if file, _, err := r.FormFile("picture"); err == nil {
        defer file.Close()

        // File path with username
        picturePath = fmt.Sprintf("uploads/users/%s.jpg", username)

        out, err := os.Create(picturePath)
        if err != nil {
            http.Error(w, "Unable to create the file for writing", http.StatusInternalServerError)
            return
        }
        defer out.Close()

        // Copy the uploaded file to the server
        if _, err = io.Copy(out, file); err != nil {
            http.Error(w, "Error saving the file", http.StatusInternalServerError)
            return
        }

        // Create the public URL for the image (assuming you want to access it via HTTP)
        imageURL = fmt.Sprintf("http://localhost:8080/uploads/users/%s.jpg", username)
    }

    // Create user model with the basic information, including picture URL
    user := models.User{
        Name:     name,
        Username: username,
        Email:    email,
        Password: password,
        Picture:  imageURL, // Store the image URL in the database
    }

    // Validate user input
    if err := user.Validate(); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Check if username or email already exists
    if existingUser, _ := database.FetchUserByUsername(user.Username); existingUser != nil {
        http.Error(w, "Username already exists", http.StatusConflict)
        return
    }

    if existingUserByEmail, _ := database.FetchUserByEmail(user.Email); existingUserByEmail != nil {
        http.Error(w, "Email already exists", http.StatusConflict)
        return
    }

    // Hash the password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Error hashing password", http.StatusInternalServerError)
        return
    }
    user.Password = string(hashedPassword)

    // Save the user to the database to generate an ID
    if err = database.InsertUser(user); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Send response back to the user
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

// UpdateUser updates user details
func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
    userIDStr := mux.Vars(r)["id"]
    
    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        http.Error(w, "Invalid user ID format", http.StatusBadRequest)
        return
    }

    // Parse form data
    err = r.ParseMultipartForm(10 << 20)
    if err != nil {
        http.Error(w, "Unable to parse form", http.StatusBadRequest)
        return
    }

    // Retrieve updated details
    name := r.FormValue("name")
    username := r.FormValue("username")
    email := r.FormValue("email")
    password := r.FormValue("password")

    var pictureURL string

    // Fetch existing user data
    existingUser, err := database.FetchUserByID(userID)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    // Check if a new picture is uploaded
    if file, _, err := r.FormFile("picture"); err == nil {
        defer file.Close()

        // Handle picture upload
        currentDate := time.Now().Format("2006-01-02") // Get current date in yyyy-mm-dd format
        pictureFileName := fmt.Sprintf("%s-%s-%d.jpg", name, currentDate, userID) // Generate picture file name
        picturePath := fmt.Sprintf("uploads/users/%s", pictureFileName)

        if err := os.MkdirAll("uploads/users", os.ModePerm); err != nil {
            http.Error(w, "Unable to create uploads directory", http.StatusInternalServerError)
            return
        }

        out, err := os.Create(picturePath)
        if err != nil {
            http.Error(w, "Unable to create file for saving", http.StatusInternalServerError)
            return
        }
        defer out.Close()
        if _, err = io.Copy(out, file); err != nil {
            http.Error(w, "Error saving the file", http.StatusInternalServerError)
            return
        }

        pictureURL = fmt.Sprintf("http://localhost:8080/%s", picturePath)
    } else {
        // If no new picture uploaded, retain existing picture
        pictureURL = existingUser.Picture
    }

    // Update only the fields provided
    if name == "" {
        name = existingUser.Name
    }
    if username == "" {
        username = existingUser.Username
    }
    if email == "" {
        email = existingUser.Email
    }
    if password == "" {
        password = existingUser.Password
    }

    // Hash password if changed
    if password != existingUser.Password {
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        if err != nil {
            http.Error(w, "Error hashing password", http.StatusInternalServerError)
            return
        }
        password = string(hashedPassword)
    }

    // Update user data
    updatedUser := models.User{
        Name:     name,
        Username: username,
        Email:    email,
        Password: password,
        Picture:  pictureURL,
    }

    if err := database.UpdateUser(userID, updatedUser); err != nil {
        http.Error(w, "Failed to update user", http.StatusInternalServerError)
        return
    }

    // Send the updated user in the response
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(updatedUser)
}

// DeleteUser handles DELETE requests to remove a user
func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    err = database.DeleteUser(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

// Login handles user login
func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
    var user models.User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Fetch user by username or email
    dbUser, err := database.FetchUserByUsernameOrEmail(user.Username, user.Email)
    if err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        log.Println("Error fetching user:", err)
        return
    }

    // Compare the password
    err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
    if err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        log.Println("Password comparison failed:", err)
        return
    }

    // Optionally, you can return user data or a token
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(dbUser)
}

// FetchUserCategories ফাংশন ব্যবহারকারীর ক্যাটেগরি ফেচ করে
func (uc *UserController) FetchUserCategories(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    categories, err := database.FetchUserCategories(userID)
    if err != nil {
        http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{"categories": categories})
}

// UpdateUserCategories ফাংশন ব্যবহারকারীর ক্যাটেগরি আপডেট করে
func (uc *UserController) UpdateUserCategories(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    var input struct {
        Categories []string `json:"categories"`
    }

    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    if err := database.UpdateUserCategories(userID, input.Categories); err != nil {
        http.Error(w, "Failed to update categories", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{"message": "Categories updated successfully"})
}
