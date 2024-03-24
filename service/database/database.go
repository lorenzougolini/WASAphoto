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
	SetUser(userid string, username string) error
	SetName(userid string, username string) error
	CheckIDExistence(userid string) (bool, error)
	GetByUsername(username string) (User, error)
	GetById(userid string) (User, error)
	GetProfile(userid string) (Profile, error)

	FollowUser(followerId string, followedId string) error
	UnfollowUser(followerId string, followedId string) error

	BanUser(passedUsername string, username string) error
	UnbanUser(userid string, username string) error
	IsBannedBy(bannedId string, bannerId string) (bool, error)

	PostPhoto(photo string) error
	GetPhotoById(photoid string) (bool, Photo, error)
	GetPhotoData(photoid string) (PhotoData, error)
	RemovePhoto(photoid string) error

	AddLike(like string) error
	GetLikeByLikeId(likeid string) (bool, Like, error)
	GetLikeByUserId(userid string, photoid string) (bool, error)
	RemoveLike(likeid string) error

	AddComment(comment string) error
	GetCommentByCommentId(id string) (bool, Comment, error)
	RemoveComment(commentid string) error

	GetStream(userid string) (Stream, error)

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
	err5 := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='follows';`).Scan(&tableName)
	err6 := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='bans';`).Scan(&tableName)

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
			authorid TEXT NOT NULL REFERENCES users(userid) ON DELETE CASCADE, 
			picFile TEXT NOT NULL, 
			dateAndTime TEXT NOT NULL,
			description TEXT
			);`
		// FOREIGN KEY (userid) REFERENCES users(userid) ON DELETE CASCADE
		_, err2 = db.Exec(sqlStmt)
		if err2 != nil {
			return nil, fmt.Errorf("error2 creating database structure: %w", err2)
		}
	}

	// likes table
	if errors.Is(err3, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE IF NOT EXISTS likes (
			likeid TEXT NOT NULL PRIMARY KEY,
			photoid TEXT NOT NULL REFERENCES photos(photoid) ON DELETE CASCADE, 
			userid TEXT NOT NULL REFERENCES users(userid) ON DELETE CASCADE,
			dateAndTime TEXT NOT NULL
			);`
		// FOREIGN KEY (photoid) REFERENCES photos(photoid) ON DELETE CASCADE,
		// FOREIGN KEY (userid) REFERENCES users(userid) ON DELETE CASCADE
		_, err3 = db.Exec(sqlStmt)
		if err3 != nil {
			return nil, fmt.Errorf("error3 creating database structure: %w", err3)
		}
	}

	// comments table
	if errors.Is(err4, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE IF NOT EXISTS comments (
			commentid TEXT NOT NULL PRIMARY KEY,
			photoid TEXT NOT NULL REFERENCES photos(photoid) ON DELETE CASCADE, 
			userid TEXT NOT NULL REFERENCES users(userid) ON DELETE CASCADE,
			commentText TEXT NOT NULL,
			dateAndTime TEXT NOT NULL
			);`
		// FOREIGN KEY (photoid) REFERENCES photos(photoid) ON DELETE CASCADE,
		// FOREIGN KEY (userid) REFERENCES users(userid) ON DELETE CASCADE
		_, err4 = db.Exec(sqlStmt)
		if err4 != nil {
			return nil, fmt.Errorf("error4 creating database structure: %w", err4)
		}
	}

	// follows table
	if errors.Is(err5, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE IF NOT EXISTS follows (
			followerid TEXT NOT NULL REFERENCES users(userid) ON DELETE CASCADE,
			followedid TEXT NOT NULL REFERENCES users(userid) ON DELETE CASCADE
			);`
		// FOREIGN KEY (userid) REFERENCES users(userid)
		_, err5 = db.Exec(sqlStmt)
		if err5 != nil {
			return nil, fmt.Errorf("error5 creating database structure: %w", err5)
		}
	}

	// bans table
	if errors.Is(err6, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE IF NOT EXISTS bans (
			bannerid TEXT NOT NULL REFERENCES users(userid) ON DELETE CASCADE,
			bannedid TEXT NOT NULL REFERENCES users(userid) ON DELETE CASCADE
			);`
		// FOREIGN KEY (userid) REFERENCES users(userid)
		_, err6 = db.Exec(sqlStmt)
		if err6 != nil {
			return nil, fmt.Errorf("error4 creating database structure: %w", err6)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
