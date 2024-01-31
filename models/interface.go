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

func (channel Channel) IsValid() map[string]string {
	validMap := make(map[string]string)

	if !validChannelName(channel.Name) {
		validMap["name"] = "Invalid channel name"
	}
	if channel.Description == "" {
		validMap["description"] = "Channel description can not be empty"
	}
	return validMap
}

func (post Post) IsValid() map[string]string {
	validMap := make(map[string]string)
	if post.Content == "" {
		validMap["content"] = "Invalid content"
		return validMap
	}
	return validMap
}

func (post UserPost) IsValid() map[string]string {
	validMap := make(map[string]string)

	if post.UserId == 0 {
		validMap["userAuthorId"] = "Invalid user author id"
	}
	return validMap
}

func (post ChannelPost) IsValid() map[string]string {
	validMap := make(map[string]string)

	if post.ChannelId == 0 {
		validMap["channelAuthorId"] = "Invalid channel author id"
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

func validChannelName(name string) bool {
	pattern := "^[a-zA-Z0-9_]{4,25}/[A-Z]+[a-z]{1,25}$"
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
