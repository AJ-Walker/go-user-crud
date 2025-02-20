package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type UserDTO struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func DBConnAndPing() error {
	log.Print("Inside DBConnAndPing")
	log.Print("lets go database")

	// Capture connectino properties
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "user_crud",
	}

	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return fmt.Errorf("DB Connection error: %v", err)
	}

	// check if db is connected
	if err := db.Ping(); err != nil {
		return fmt.Errorf("DB Connection error: %v", err)
	}
	log.Print("DB Connected! lets go!")
	return nil
}

// Get list of users from DB
func GetUsersDB() ([]UserDTO, error) {
	log.Print("Inside GetUsersDB")

	var users []UserDTO

	rows, err := db.Query("SELECT user_id, name, email FROM users")
	if err != nil {
		return nil, fmt.Errorf("GetUsersDB error: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var user UserDTO

		if err := rows.Scan(&user.Id, &user.Name, &user.Email); err != nil {
			return nil, fmt.Errorf("GetUsersDB error: %v", err)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetUsersDB error: %v", err)
	}

	return users, nil
}

// Add a user to DB
func AddUserDB(user User) error {
	log.Print("Inside AddUserDB")

	_, err := db.Exec("INSERT INTO users (user_id, name, email, password) VALUES (?,?,?,?)", user.Id, user.Name, user.Email, user.Password)
	if err != nil {
		return fmt.Errorf("AddUserDB error: %v", err)
	}
	return nil
}

// Get a single user by their ID from DB
func GetUserByIdDB(id string) (UserDTO, error) {
	log.Print("Inside GetUserByIdDB")

	var user UserDTO
	row := db.QueryRow("SELECT user_id, name, email FROM users WHERE user_id = ?", id)

	if err := row.Scan(&user.Id, &user.Name, &user.Email); err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("GetUserByIdDB error: %v", err)
		}
		return user, fmt.Errorf("GetUserByIdDB error: %v", err)
	}

	return user, nil
}

// Get a single user by their Email from DB
func GetUserByEmailDB(email string) (User, error) {
	log.Print("Inside GetUserByEmailDB")

	var user User
	row := db.QueryRow("SELECT user_id, name, email, password FROM users WHERE email = ?", email)

	if err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("GetUserByEmailDB error: %v", err)
		}
		return user, fmt.Errorf("GetUserByEmailDB error: %v", err)
	}

	return user, nil
}

// Update a user from DB
func UpdateUserDB(id string, user User) error {
	log.Print("Inside UpdateUserDB")

	_, err := db.Exec("UPDATE users SET name=? WHERE user_id=?", user.Name, id)

	if err != nil {
		return fmt.Errorf("UpdateUserDB error: %v", err)
	}

	return nil
}

// Delete a user from DB
func DeleteUserDB(id string) (bool, error) {
	log.Print("Inside DeleteUserDB")

	res, err := db.Exec("DELETE FROM users WHERE user_id = ?", id)
	if err != nil {
		return false, fmt.Errorf("DeleteUserDB error: %v", err)
	}

	res1, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("DeleteUserDB error: %v", err)
	}
	if res1 == 1 {
		return true, nil
	}
	return false, nil
}
