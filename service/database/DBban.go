package database

import (
	"database/sql"
	"errors"
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
func (db *appdbimpl) IsBanned(reqUserId string, banUserId string) (bool, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM bans WHERE userid=? AND bannedid=?", reqUserId, banUserId).Scan(&count)
	if errors.Is(err, sql.ErrNoRows) {
		return count == 0, nil
	}
	return count > 0, err
}
