package utils

import (
	"database/sql"
	"fmt"
	"os"
)

// DatabaseCreate function with no parameters that returns an error
func DatabaseCreate() error {
	// create Data directory if it does not exist
	if _, err := os.Stat("Data"); os.IsNotExist(err) {
		err := os.Mkdir("Data", 0755)
		if err != nil {
			return fmt.Errorf("Error creating Data directory: %s\n", err)
		}
	}
	// open connection to the SQLite database
	db, err := sql.Open("sqlite3", "./Data/auth.db")
	if err != nil {
		return fmt.Errorf("Error opening database: %s\n", err)
	}
	defer db.Close()
	// create the table if it does not exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (user_id TEXT PRIMARY KEY, token TEXT UNIQUE)`)
	// check if there is an error
	if err != nil {
		return fmt.Errorf("Error creating table: %v", err)
	}

	return nil
}

// DatabaseInsert function with user_id and token as parameters
// that returns an error
func DatabaseInsert(user_id string, token string) error {
	// open connection to the SQLite database
	db, err := sql.Open("sqlite3", "./Data/auth.db")
	if err != nil {
		return fmt.Errorf("Error opening database: %s\n", err)
	}
	defer db.Close()

	// insert the token to the database
	_, err = db.Exec(`INSERT INTO users (user_id, token) VALUES (?, ?)`, user_id, token)
	if err != nil {
		return fmt.Errorf("Error inserting token: %s\n", err)
	}
	return nil
}

// DatabaseDelete function with user_id as parameter that returns an error
func DatabaseDelete(user_id string) error {
	// open connection to the SQLite database
	db, err := sql.Open("sqlite3", "./Data/auth.db")
	if err != nil {
		return fmt.Errorf("Error opening database: %s\n", err)
	}
	defer db.Close()

	// delete the token from the database
	_, err = db.Exec(`DELETE FROM users WHERE user_id = ?`, user_id)
	if err != nil {
		return fmt.Errorf("Error deleting token: %s\n", err)
	}

	return nil
}

// user struct with fields of UserId and Token with json tags
type User struct {
	UserId string `json:"user_id"`
	Token  string `json:"token"`
}

// DatabaseSelectAll function with no parameters that returns a slice of users and an error
func DatabaseSelectAll() ([]User, error) {
	// open connection to the SQLite database
	db, err := sql.Open("sqlite3", "./Data/auth.db")
	if err != nil {
		return nil, fmt.Errorf("Error opening database: %s\n", err)
	}
	defer db.Close()

	// select the token from the database
	rows, err := db.Query(`SELECT user_id, token FROM users`)
	if err != nil {
		return nil, fmt.Errorf("Error selecting token: %s\n", err)
	}
	defer rows.Close()

	// create slice of users
	var users []User
	// iterate over the rows
	for rows.Next() {
		// create a user
		var user User
		// scan the rows into the user
		err = rows.Scan(&user.UserId, &user.Token)
		if err != nil {
			return nil, fmt.Errorf("Error scanning rows: %s\n", err)
		}
		// append the user to the slice of users
		users = append(users, user)
	}
	// check if there is an error
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterating rows: %s\n", err)
	}

	return users, nil
}

// VerifyAdminKey function with key as parameter that returns a boolean
func VerifyAdminKey(key string) bool {
	// return key as os.Args[2]
	return key == os.Args[2]
}

// VerifyToken function with token as parameter that returns a boolean and error
func VerifyToken(token string) (bool, error) {
	// check if token is admin key using VerifyAdminKey function
	if VerifyAdminKey(token) {
		return true, nil
	}
	// open connection to the SQLite database
	db, err := sql.Open("sqlite3", "./Data/auth.db")
	if err != nil {
		return false, fmt.Errorf("Error opening database: %s\n", err)
	}
	defer db.Close()

	// select the token from the database
	rows, err := db.Query(`SELECT token FROM users WHERE token = ?`, token)
	if err != nil {
		return false, fmt.Errorf("Error selecting token: %s\n", err)
	}
	defer rows.Close()

	//check if the token exists
	if rows.Next() {
		return true, nil
	} else {
		return false, nil
	}
}
