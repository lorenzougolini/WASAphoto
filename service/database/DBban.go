package database

import (
	"fmt"
)

// Ban user
func (db *appdbimpl) BanUser(user User, username string) error {

	sqlStmt := "SELECT * FROM bans WHERE userid=? AND ban=?"
	res, err := db.c.Exec(sqlStmt, user.UserID, username)
	if err != nil {
		return fmt.Errorf("error adding the ban: %v", err)
	}
	if rowsAffected, err := res.RowsAffected(); rowsAffected != 0 || err != nil {
		return fmt.Errorf("the user '%s' is already banned", username)
	}

	sqlStmt = "INSERT INTO bans VALUES (?, ?);"
	_, err = db.c.Exec(sqlStmt, user.UserID, username)
	if err != nil {
		return fmt.Errorf("error adding the ban: %v", err)
	}
	return nil
}

// Unban user
func (db *appdbimpl) UnbanUser(id string, username string) error {
	sqlStmt := "DELETE FROM bans WHERE userid=? AND ban=?"
	res, err := db.c.Exec(sqlStmt, id, username)
	if err != nil {
		return fmt.Errorf("error removing the ban: %v", err)
	}
	if affectedRows, _ := res.RowsAffected(); affectedRows == 0 {
		return fmt.Errorf("the user '%s' is not banned at the moment", username)
	}
	return nil
}

// check banned user
func (db *appdbimpl) IsBanned(id string, bannedId string) (bool, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM bans WHERE userid=? AND ban=?", id, bannedId).Scan(&count)
	if err != nil {
		return count > 0, err
	}
	return count > 0, nil
}
