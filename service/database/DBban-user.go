package database

import (
	"fmt"
)

// Ban user
func (db *appdbimpl) BanUser(user map[string]string, username string) error {

	sqlStmt := "SELECT * FROM bans WHERE userid=? AND ban=?"
	res, err := db.c.Exec(sqlStmt, user["id"], username)
	if err != nil {
		return fmt.Errorf("Error adding the ban: %v", err)
	}
	if rowsAffected, err := res.RowsAffected(); rowsAffected != 0 || err != nil {
		return fmt.Errorf("The user '%s' is already banned", username)
	}

	sqlStmt = "INSERT INTO bans VALUES (?, ?);"
	_, err = db.c.Exec(sqlStmt, user["id"], username)
	if err != nil {
		return fmt.Errorf("Error adding the ban: %v", err)
	}
	return nil
}

// Unban user
func (db *appdbimpl) UnbanUser(id string, username string) error {
	sqlStmt := "DELETE FROM bans WHERE userid=? AND ban=?"
	res, err := db.c.Exec(sqlStmt, id, username)
	if err != nil {
		return fmt.Errorf("Error removing the ban: %v", err)
	}
	if affectedRows, _ := res.RowsAffected(); affectedRows == 0 {
		return fmt.Errorf("The user '%s' is not banned at the moment", username)
	}
	return nil
}
