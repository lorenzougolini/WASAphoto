package database

// GetName is an example that shows you how to query data
func (db *appdbimpl) GetName(username string) (string, string, error) {
	var id string
	err := db.c.QueryRow("SELECT id, username FROM user_table WHERE username=?", username).Scan(&id, &username)
	return id, username, err
}
