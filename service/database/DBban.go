package database

import (
	"database/sql"
	"fmt"
)

// Ban user
func (db *appdbimpl) BanUser(bannerId string, bannedId string) error {

	banned, err := db.IsBannedBy(bannedId, bannerId)
	if banned {
		return fmt.Errorf("the user is already banned")
	}

	if err != nil {
		return fmt.Errorf("error adding the ban: %w", err)

	}

	sqlStmt := "INSERT INTO bans VALUES (?, ?);"
	_, err = db.c.Exec(sqlStmt, bannerId, bannedId)
	if err != nil {
		return fmt.Errorf("error adding the ban: %w", err)
	}
	return nil

}

// Unban user
func (db *appdbimpl) UnbanUser(bannerId string, bannedId string) error {

	if banned, err := db.IsBannedBy(bannedId, bannerId); !banned || err != nil {
		return fmt.Errorf("the user is not banned at the moment")

	}
	sqlStmt := "DELETE FROM bans WHERE bannerid=? AND bannedid=?"
	_, err := db.c.Exec(sqlStmt, bannerId, bannedId)
	if err != nil {
		return fmt.Errorf("error removing the ban: %w", err)
	}

	return nil
}

// check banned user
func (db *appdbimpl) IsBannedBy(bannedId string, bannerId string) (bool, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM bans WHERE bannerid=? AND bannedid=?", bannerId, bannedId).Scan(&count)
	if err == sql.ErrNoRows {
		// No ban found
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
