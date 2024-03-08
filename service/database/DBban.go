package database

import (
	"database/sql"
	"fmt"
)

// Ban user
func (db *appdbimpl) BanUser(reqUserId string, banUserId string) error {

	banned, err := db.IsBanned(reqUserId, banUserId)
	if banned {
		return fmt.Errorf("the user is already banned")
	}

	if err != nil {
		return fmt.Errorf("error adding the ban: %w", err)

	}

	sqlStmt := "INSERT INTO bans VALUES (?, ?);"
	_, err = db.c.Exec(sqlStmt, reqUserId, banUserId)
	if err != nil {
		return fmt.Errorf("error adding the ban: %w", err)
	}
	return nil

}

// Unban user
func (db *appdbimpl) UnbanUser(reqUserId string, banUserId string) error {

	if banned, err := db.IsBanned(reqUserId, banUserId); !banned || err != nil {
		return fmt.Errorf("the user is not banned at the moment")

	}
	sqlStmt := "DELETE FROM bans WHERE userid=? AND bannedid=?"
	_, err := db.c.Exec(sqlStmt, reqUserId, banUserId)
	if err != nil {
		return fmt.Errorf("error removing the ban: %w", err)
	}

	return nil
}

// check banned user
func (db *appdbimpl) IsBanned(reqUserID, banUserID string) (bool, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM bans WHERE userid=? AND bannedid=?", reqUserID, banUserID).Scan(&count)
	if err == sql.ErrNoRows {
		// No ban found
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("error checking ban existence: %w", err)
	}

	return count > 0, nil
}
