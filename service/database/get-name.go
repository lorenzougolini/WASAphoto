package database

// GetName is an example that shows you how to query data
func (db *appdbimpl) GetName(id string) (string, string, error) {
	var username string
	err := db.c.QueryRow("SELECT id, username FROM user_table WHERE id=?", id).Scan(&id, &username)
	return id, username, err
}
