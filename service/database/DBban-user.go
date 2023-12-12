package database

import (
	"fmt"
)

// GetName is an example that shows you how to query data
func (db *appdbimpl) BanUser(user map[string]string, username string) error {
	sqlStmt := "UPDATE profiles SET banned = JSON_INSERT(banned, '$[#]', ?) WHERE userid=?;"
	_, err := db.c.Exec(sqlStmt, username, user["id"])
	if err != nil {
		return fmt.Errorf("Error adding the ban: %v", err)
	}
	return nil
}

// NOT WORKING
func (db *appdbimpl) UnbanUser(id string, username string) error {
	sqlStmt := "UPDATE profiles SET banned = JSON_REMOVE(banned, JSON_UNQUOTE(JSON_SEARCH(banned, 'one', ?))) WHERE userid = ?"
	res, err := db.c.Exec(sqlStmt, username, id)
	if err != nil {
		return fmt.Errorf("Error removing the ban: %v", err)
	}
	if affectedRows, _ := res.RowsAffected(); affectedRows == 0 {
		return fmt.Errorf("The user '%s' is not banned at the moment", username)
	}
	return nil
}
