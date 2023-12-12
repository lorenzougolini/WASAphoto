package database

import (
	"fmt"
)

// GetName is an example that shows you how to query data
func (db *appdbimpl) FollowUser(user map[string]string, username string) error {

	sqlStmt := "UPDATE profiles SET following = JSON_INSERT(following, '$[#]', ?) WHERE userid=?;"
	_, err := db.c.Exec(sqlStmt, username, user["id"])
	if err != nil {
		return fmt.Errorf("Error adding the follow: %v", err)
	}

	followed, err := db.GetUser(username)
	if err != nil {
		return fmt.Errorf("Error adding the follow: %v", err)
	}
	sqlStmt = "UPDATE profiles SET followers = JSON_INSERT(followers, '$[#]', ?) WHERE userid=?;"
	_, err = db.c.Exec(sqlStmt, user["username"], followed.UserID)
	if err != nil {
		return fmt.Errorf("Error adding the follow: %v", err)
	}

	return nil
}
