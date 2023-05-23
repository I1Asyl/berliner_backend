package models

import (
	"net/mail"
	"regexp"
)

// interface for all models
type General interface {
	IsValid() bool
}

func (user User) IsValid() map[string]string {
	validMap := make(map[string]string)

	if !validUsername(user.Username) {
		validMap["username"] = "Invalid username"
	}
	if !validName(user.FirstName) {
		validMap["firstName"] = "Invalid first name"
	}
	if !validName(user.LastName) {
		validMap["lastName"] = "Invalid last name"
	}
	if !validEmail(user.Email) {
		validMap["email"] = "Invalid email"
	}
	if !validPassword(user.Password) {
		validMap["password"] = "Invalid password"
	}

	return validMap
}

func (team Team) IsValid() map[string]string {
	validMap := make(map[string]string)

	if !validName(team.TeamName) {
		validMap["teamName"] = "Invalid team name"
	}
	if team.TeamDescription == "" {
		validMap["teamDescription"] = "Team description can not be empty"
	}
	return validMap
}

func (post Post) IsValid() map[string]string {
	validMap := make(map[string]string)
	if post.Content == "" {
		validMap["content"] = "Invalid content"
		return validMap
	}
	if post.AuthorType == "" {
		validMap["authorType"] = "Invalid author type"
	}

	if post.AuthorType == "user" && !post.UserAuthorId.Valid {
		validMap["userAuthorId"] = "Invalid user author id"
	}
	if post.AuthorType == "team" && !post.TeamAuthorId.Valid {
		validMap["teamAuthorId"] = "Invalid team author id"
	}
	return validMap
}

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func validUsername(username string) bool {
	pattern := "^[a-zA-Z0-9_]{4,25}$"
	ans, _ := regexp.MatchString(pattern, username)
	return ans
}

func validName(name string) bool {
	pattern := "^[A-Z]+[a-z]{1,25}$"
	ans, _ := regexp.MatchString(pattern, name)
	return ans
}

func validPassword(password string) bool {
	patterns := []string{"^[a-zA-Z0-9_@$!%*#?&.]{8,40}$", "[a-z]+", "[A-Z]+", "[\\d]+", "[@$!%*#?&.]+"}
	for _, pattern := range patterns {
		tmp, _ := regexp.MatchString(pattern, password)
		if !tmp {
			return false
		}
	}
	return true
}
