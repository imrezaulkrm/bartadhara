
# Bartadhara Project

Bartadhara is a backend project for a news website built using GoLang and MySQL. This project allows users to view news based on their preferred categories and enables admin users to manage news content.

## How to Set Up

Follow these steps to set up the Bartadhara backend project on your machine.

### Step 1: Set Up the Database

1. **Install MySQL** if you haven’t already.
2. **Create the Database**:
   - Open MySQL and create a database for the project:
     ```sql
     CREATE DATABASE news_db;
     ```
3. **Set Up Database Structure**:
   - Import the database structure from the provided `news_db.sql` file:
     ```bash
     mysql -u your_username -p news_db < news_db.sql
     ```
     This will create the required tables without any data.

### Step 2: Set Up GoLang

1. **Install GoLang** if it's not already installed (version 1.18 or above recommended).
2. **Download Project Dependencies**:
   - Navigate to the project folder and download the dependencies:
     ```bash
     go mod tidy
     ```

### Step 3: Configure and Run the Project

1. **Update Database Connection**:
   - Go to the `database/database.go` file, and on **line 49**, update the database connection string with your MySQL username and password:
     ```go
     db, err = sql.Open("mysql", "your_username:your_password@tcp(127.0.0.1:3306)/news_db")
     ```
2. **Run the Project**:
   - Start the Go server by running:
     ```bash
     go run main.go
     ```
   - The server should now be running on `http://localhost:8080`.

## Conclusion

Your Bartadhara backend project should now be up and running. Access the API endpoints to manage and view news content. Once the project is fully developed, a professional README file with additional details and instructions will be added.