/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetUser(username string) (User, error)
	SetUser(id string, username string) error
	CheckID(id string) (int, error)
	FollowUser(user map[string]string, username string) error
	// UnfollowUser(id string, username string) error

	BanUser(user map[string]string, username string) error
	UnbanUser(id string, username string) error

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err1 := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	err2 := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='photos';`).Scan(&tableName)
	err3 := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='likes';`).Scan(&tableName)
	err4 := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='comments';`).Scan(&tableName)
	err5 := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='profiles';`).Scan(&tableName)

	// users table
	if errors.Is(err1, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE IF NOT EXISTS users (
			userid TEXT NOT NULL PRIMARY KEY, 
			username TEXT NOT NULL UNIQUE
			);`
		_, err1 = db.Exec(sqlStmt)
		if err1 != nil {
			return nil, fmt.Errorf("error1 creating database structure: %w", err1)
		}
	}

	// photos table
	if errors.Is(err2, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE IF NOT EXISTS photos (
			photoid TEXT NOT NULL PRIMARY KEY, 
			userid TEXT NOT NULL, 
			picture TEXT NOT NULL, 
			dateAndTime TEXT NOT NULL,
			FOREING KEY userid REFERENCES users(userid)
			);`
		_, err2 = db.Exec(sqlStmt)
		if err2 != nil {
			return nil, fmt.Errorf("error2 creating database structure: %w", err2)
		}
	}

	// likes table
	if errors.Is(err3, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE IF NOT EXISTS likes (
			likeid TEXT NOT NULL PRIMARY KEY,
			photoid TEXT NOT NULL, 
			username TEXT NOT NULL,
			dateAndTime TEXT NOT NULL,
			FOREIGN KEY (photoid) REFERENCES photos(photoid)
			);`
		_, err3 = db.Exec(sqlStmt)
		if err3 != nil {
			return nil, fmt.Errorf("error3 creating database structure: %w", err3)
		}
	}

	// comments table
	if errors.Is(err4, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE IF NOT EXISTS comments (
			commentid TEXT NOT NULL PRIMARY KEY,
			photoid TEXT NOT NULL, 
			username TEXT NOT NULL,
			commentText TEXT NOT NULL,
			dateAndTime TEXT NOT NULL,
			FOREIGN KEY (photoid) REFERENCES photos(photoid)
			);`
		_, err4 = db.Exec(sqlStmt)
		if err4 != nil {
			return nil, fmt.Errorf("error4 creating database structure: %w", err4)
		}
	}

	// profiles table
	if errors.Is(err5, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE IF NOT EXISTS profiles (
			userid TEXT NOT NULL,
			photos JSON DEFAULT('[]'),
			followers JSON DEFAULT('[]'),
			following JSON DEFAULT('[]'),
			banned JSON DEFAULT('[]'),
			FOREIGN KEY (userid) REFERENCES users(userid)
			);`
		_, err4 = db.Exec(sqlStmt)
		if err4 != nil {
			return nil, fmt.Errorf("error4 creating database structure: %w", err4)
		}
	}
	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
