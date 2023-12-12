package database

import "fmt"

// SetName is an example that shows you how to execute insert/update
func (db *appdbimpl) SetUser(id string, username string) error {
	stmt1, _ := db.c.Prepare("INSERT INTO users (userid, username) VALUES (?, ?);")
	stmt2 := "INSERT INTO profiles (userid) VALUES (?);"
	_, err1 := stmt1.Exec(id, username)
	_, err2 := db.c.Exec(stmt2, id)
	if err1 != nil || err2 != nil {
		return fmt.Errorf("Error in Exec\n err1: %v\n err2: %v", err1, err2)
	}
	return nil
}
