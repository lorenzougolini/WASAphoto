package database

type User struct {
	UserID   string
	Username string
}

// type Profile struct {
// 	photos    []string
// 	followers []string
// 	following []string
// 	banned    []string
// }

type Photo struct {
	PhotoID     string
	UserID      string
	Picture     string
	DateAndTime string
	Description string
	// Likes       []string
	// Comments    []string
}

type Like struct {
	likeID      string
	username    string
	photoID     string
	dateAndTime string
}

type Comment struct {
	commentID   string
	username    string
	photoID     string
	commentText string
	dateAndTime string
}

var Logged = User{}

var UsernameToId = make(map[string]string)
var Users = make(map[string]User)
var Photos = make(map[string]Photo)
var Likes = make(map[string]Like)
var Comments = make(map[string]Comment)

func remove(slice []string, s string) []string {
	for i, el := range slice {
		if el == s {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func checkLogin(id string) bool {
	return id == Logged.UserID
}

func contains(slice []string, s string) bool {
	for _, el := range slice {
		if el == s {
			return true
		}
	}
	return false
}
