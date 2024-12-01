package controllers

import (
	"encoding/json"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"github.com/gorilla/mux"
	"github.com/imrezaulkrm/bartadhara-backend/database"
	"github.com/imrezaulkrm/bartadhara-backend/models"
	"io"
	"os"
	"fmt"
	"log"
	"time"
)

// AdminController struct
type AdminController struct{}

// FetchAllAdmins - Get all admins
func (ac *AdminController) FetchAllAdmins(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
	query := "SELECT id, name, username, email, picture FROM admins"
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Error fetching admins", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var admins []models.Admin
	for rows.Next() {
		var admin models.Admin
		if err := rows.Scan(&admin.ID, &admin.Name, &admin.Username, &admin.Email, &admin.Picture); err != nil {
			http.Error(w, "Error scanning admin", http.StatusInternalServerError)
			return
		}
		admins = append(admins, admin)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(admins)
}

// FetchAdminByID - Get admin by ID
func (ac *AdminController) FetchAdminByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	adminID := params["id"]

	db := database.GetDB()
	query := "SELECT id, name, username, email, picture FROM admins WHERE id = ?"
	row := db.QueryRow(query, adminID)

	var admin models.Admin
	if err := row.Scan(&admin.ID, &admin.Name, &admin.Username, &admin.Email, &admin.Picture); err != nil {
		http.Error(w, "Admin not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(admin)
}

// InsertAdmin - Insert new admin
/*func (ac *AdminController) InsertAdmin(w http.ResponseWriter, r *http.Request) {
	var admin models.Admin
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&admin); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert admin into database
	db := database.GetDB()
	query := "INSERT INTO admins (name, username, email, password, picture) VALUES (?, ?, ?, ?, ?)"
	_, err := db.Exec(query, admin.Name, admin.Username, admin.Email, admin.Password, admin.Picture)
	if err != nil {
		http.Error(w, "Error creating admin", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(admin)
}*/

// UpdateAdmin - Update existing admin details
func (ac *AdminController) UpdateAdmin(w http.ResponseWriter, r *http.Request) {
    // Retrieve the admin ID from the URL parameters
    params := mux.Vars(r)
    adminID := params["id"]

    // Parse the form data (allowing for file uploads)
    err := r.ParseMultipartForm(10 << 20) // 10MB limit for file uploads
    if err != nil {
        http.Error(w, "Unable to parse form", http.StatusBadRequest)
        return
    }

    // Retrieve form values for the admin details
    name := r.FormValue("name")
    username := r.FormValue("username")
    email := r.FormValue("email")
    password := r.FormValue("password")

    var pictureURL string

    // Fetch the existing admin data from the database
    existingAdmin, err := database.FetchAdminByID(adminID)
    if err != nil {
        http.Error(w, "Admin not found", http.StatusNotFound)
        return
    }

    // Check if a new picture is uploaded
    if file, _, err := r.FormFile("picture"); err == nil {
        defer file.Close()

        // Create the directory for storing the uploaded image
        currentDate := time.Now().Format("2006-01-02") // Get current date in yyyy-mm-dd format
        pictureFileName := fmt.Sprintf("%s-%s-%s.jpg", name, currentDate, adminID) // Generate the picture file name
        picturePath := fmt.Sprintf("uploads/admins/%s", pictureFileName)

        if err := os.MkdirAll("uploads/admins", os.ModePerm); err != nil {
            http.Error(w, "Unable to create uploads directory", http.StatusInternalServerError)
            return
        }

        // Save the uploaded file to disk
        out, err := os.Create(picturePath)
        if err != nil {
            http.Error(w, "Unable to create file for saving", http.StatusInternalServerError)
            return
        }
        defer out.Close()

        // Copy the content of the uploaded file into the new file
        if _, err = io.Copy(out, file); err != nil {
            http.Error(w, "Error saving the file", http.StatusInternalServerError)
            return
        }

        // Set the picture URL (publicly accessible location)
        pictureURL = fmt.Sprintf("http://localhost:8080/%s", picturePath)
    } else {
        // If no new picture uploaded, retain the existing picture URL
        pictureURL = existingAdmin.Picture
    }

    // Update fields if new values are provided, otherwise keep the old values
    if name == "" {
        name = existingAdmin.Name
    }
    if username == "" {
        username = existingAdmin.Username
    }
    if email == "" {
        email = existingAdmin.Email
    }
    if password == "" {
        password = existingAdmin.Password
    }

    // Hash password if it was changed
    if password != existingAdmin.Password {
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        if err != nil {
            http.Error(w, "Error hashing password", http.StatusInternalServerError)
            return
        }
        password = string(hashedPassword)
    }

    // Update the admin in the database with only the fields that changed
    db := database.GetDB()

    // Prepare the update query with conditionals to only update the fields that changed
    query := "UPDATE admins SET name = ?, username = ?, email = ?, password = ?, picture = ? WHERE id = ?"
    _, err = db.Exec(query, name, username, email, password, pictureURL, adminID)
    if err != nil {
        http.Error(w, "Error updating admin", http.StatusInternalServerError)
        return
    }

    // Create updated admin object for response
    updatedAdmin := models.Admin{
        ID:       existingAdmin.ID,
        Name:     name,
        Username: username,
        Email:    email,
        Password: password,  // Note: You might not want to send password in the response
        Picture:  pictureURL,
    }

    // Send the updated admin details in the response
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(updatedAdmin)
}

// DeleteAdmin - Delete admin by ID
func (ac *AdminController) DeleteAdmin(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	adminID := params["id"]

	db := database.GetDB()
	query := "DELETE FROM admins WHERE id = ?"
	_, err := db.Exec(query, adminID)
	if err != nil {
		http.Error(w, "Error deleting admin", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Login handles admin login
func (ac *AdminController) Login(w http.ResponseWriter, r *http.Request) {
    var admin models.Admin
    err := json.NewDecoder(r.Body).Decode(&admin)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Fetch admin by username or email
    dbAdmin, err := database.FetchAdminByUsernameOrEmail(admin.Username, admin.Email)
    if err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        log.Println("Error fetching admin:", err)
        return
    }

    // Compare the password
    err = bcrypt.CompareHashAndPassword([]byte(dbAdmin.Password), []byte(admin.Password))
    if err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        log.Println("Password comparison failed:", err)
        return
    }

    // Prepare a safe response with necessary fields only
    response := map[string]interface{}{
        "success": true,
        "message": "Login successful",
        "admin": map[string]interface{}{
            "id":       dbAdmin.ID,
            "name":     dbAdmin.Name,
            "username": dbAdmin.Username,
            "email":    dbAdmin.Email,
        },
    }

    // Send response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}


// Register handles admin registration with picture upload
func (ac *AdminController) Register(w http.ResponseWriter, r *http.Request) {
    // Parse multipart form data
    err := r.ParseMultipartForm(10 << 20) // Limit file size to 10MB
    if err != nil {
        http.Error(w, "Error parsing form", http.StatusBadRequest)
        return
    }

    // Extract form values
    name := r.FormValue("name")
    username := r.FormValue("username")
    email := r.FormValue("email")
    password := r.FormValue("password")

    // Handle file upload
    file, _, err := r.FormFile("picture")
    if err != nil {
        http.Error(w, "Error uploading file", http.StatusBadRequest)
        return
    }
    defer file.Close()

    // Create the uploads/admins directory if it doesn't exist
    if err := os.MkdirAll("uploads/admins", os.ModePerm); err != nil {
        http.Error(w, "Unable to create directory for admin pictures", http.StatusInternalServerError)
        return
    }

    // Save the file to the server in the uploads/admins folder
    filePath := fmt.Sprintf("uploads/admins/%s.jpg", username) // Store the file with the username
    outFile, err := os.Create(filePath)
    if err != nil {
        http.Error(w, "Error saving file", http.StatusInternalServerError)
        return
    }
    defer outFile.Close()

    // Copy the uploaded file to the server
    _, err = io.Copy(outFile, file)
    if err != nil {
        http.Error(w, "Error saving file", http.StatusInternalServerError)
        return
    }

    // Create the public URL for the image (assuming you want to access it via HTTP)
    imageURL := fmt.Sprintf("http://localhost:8080/uploads/admins/%s.jpg", username) // URL to access the image

    // Hash the password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Error hashing password", http.StatusInternalServerError)
        return
    }

    // Create an admin object with hashed password and picture URL
    admin := models.Admin{
        Name:     name,
        Username: username,
        Email:    email,
        Password: string(hashedPassword), // Store the hashed password
        Picture:  imageURL,               // Store the image URL in the database
    }

    // Save admin to the database using the database package
    err = database.SaveAdminToDatabase(admin)
    if err != nil {
        http.Error(w, "Error saving admin to database", http.StatusInternalServerError)
        return
    }

    // Send success response
    response := map[string]string{"message": "Admin registered successfully"}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}