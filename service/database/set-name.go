package database

import "fmt"

// SetName is an example that shows you how to execute insert/update
func (db *appdbimpl) SetName(id string, username string) {
	statement, err1 := db.c.Prepare("INSERT INTO user_table (id, username) VALUES (?, ?)")
	if err1 != nil {
		fmt.Println("Error in Prepare")
	}
	_, err2 := statement.Exec(id, username)
	if err2 != nil {
		fmt.Println("Error in Exec")
	}
	return
}
