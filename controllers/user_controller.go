package controllers

import (
    "encoding/json"
    //"errors"
    "net/http"
    "strconv"
    "log"
    "path/filepath"
    "os" // Added for file operations
    "io" 
    "github.com/gorilla/mux"
    "golang.org/x/crypto/bcrypt"

    "github.com/imrezaulkrm/bartadhara/database"
    "github.com/imrezaulkrm/bartadhara/models"
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

// InsertUser handles POST requests to create a new user
func (uc *UserController) InsertUser(w http.ResponseWriter, r *http.Request) {
    // Parse multipart form
    err := r.ParseMultipartForm(10 << 20) // 10 MB limit
    if err != nil {
        http.Error(w, "Unable to parse form", http.StatusBadRequest)
        return
    }

    // Create uploads directory if it doesn't exist
    if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
        http.Error(w, "Unable to create uploads directory", http.StatusInternalServerError)
        return
    }

    // Retrieve user details from form
    name := r.FormValue("name")
    username := r.FormValue("username")
    email := r.FormValue("email")
    password := r.FormValue("password")

    // Check if an image file is uploaded
    var picture string
    if file, _, err := r.FormFile("picture"); err == nil {
        // Process the uploaded file
        defer file.Close()
        picture = filepath.Join("uploads", username+"_profile.jpg") // Example path

        out, err := os.Create(picture)
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
    }

    // Create user model
    user := models.User{
        Name:     name,
        Username: username,
        Email:    email,
        Password: password,
        Picture:  picture, // Picture path included here
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

    // Save user to the database
    if err = database.InsertUser(user); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Save selected categories (if any)
    if len(user.Categories) > 0 {
        if err = database.SaveUserCategories(user.ID, user.Categories); err != nil {
            http.Error(w, "Error saving user categories", http.StatusInternalServerError)
            return
        }
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}


// UpdateUser handles PUT requests to update a user
func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var user models.User
    err = json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Hash the password if provided
    if user.Password != "" {
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
        if err != nil {
            http.Error(w, "Error hashing password", http.StatusInternalServerError)
            return
        }
        user.Password = string(hashedPassword)
    }

    err = database.UpdateUser(id, user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
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
