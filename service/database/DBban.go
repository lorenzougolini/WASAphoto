package database

import (
	"encoding/json"
	"fmt"
)

// Ban user
func (db *appdbimpl) BanUser(loggedUser string, username string) error {

	logged_user := User{}
	_ = json.Unmarshal([]byte(loggedUser), &logged_user)
	banned_user, err := db.GetByUsername(username)
	if err != nil {
		return fmt.Errorf("error adding the ban: %w", err)
	}

	banned, err := db.IsBanned(logged_user.UserID, banned_user.UserID)
	if banned {
		return fmt.Errorf("the user '%s' is already banned", username)

	} else if err != nil {
		return fmt.Errorf("error adding the ban: %w", err)

	} else {
		sqlStmt := "INSERT INTO bans VALUES (?, ?);"
		_, err = db.c.Exec(sqlStmt, logged_user.UserID, banned_user.UserID)
		if err != nil {
			return fmt.Errorf("error adding the ban: %w", err)
		}
		return nil
	}
}

// Unban user
func (db *appdbimpl) UnbanUser(loggedId string, bannedId string) error {
	if count, err := db.IsBanned(loggedId, bannedId); !count || err != nil {
		return fmt.Errorf("the user '%s' is not banned at the moment", bannedId)

	} else {
		sqlStmt := "DELETE FROM bans WHERE userid=? AND bannedid=?"
		_, err := db.c.Exec(sqlStmt, loggedId, bannedId)
		if err != nil {
			return fmt.Errorf("error removing the ban: %w", err)
		}
	}
	return nil
}

// check banned user
func (db *appdbimpl) IsBanned(id string, bannedId string) (bool, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM bans WHERE userid=? AND bannedid=?", id, bannedId).Scan(&count)
	if err != nil {
		return count > 0, err
	}
	return count > 0, nil
}
