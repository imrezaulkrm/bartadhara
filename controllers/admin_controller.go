package controllers

import (
	"encoding/json"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"github.com/gorilla/mux"
	"github.com/imrezaulkrm/bartadhara/database"
	"github.com/imrezaulkrm/bartadhara/models"
	"io"
	"os"
	"fmt"
	"log"
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
func (ac *AdminController) InsertAdmin(w http.ResponseWriter, r *http.Request) {
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
}

// UpdateAdmin - Update existing admin
func (ac *AdminController) UpdateAdmin(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	adminID := params["id"]
	var admin models.Admin
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&admin); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update admin in the database
	db := database.GetDB()
	query := "UPDATE admins SET name = ?, username = ?, email = ?, password = ?, picture = ? WHERE id = ?"
	_, err := db.Exec(query, admin.Name, admin.Username, admin.Email, admin.Password, admin.Picture, adminID)
	if err != nil {
		http.Error(w, "Error updating admin", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(admin)
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

    // Optionally, return admin data or a token
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(dbAdmin)
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

	// Save the file to the server
	filePath := fmt.Sprintf("./uploads/%s.jpg", username) // Store the file with the username
	outFile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Create an admin object with hashed password
	admin := models.Admin{
		Name:     name,
		Username: username,
		Email:    email,
		Password: string(hashedPassword), // Store the hashed password
		Picture:  filePath,               // Store the file path in the admin object
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
