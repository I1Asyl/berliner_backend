package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// nullable  int64
type NullInt64 struct {
	sql.NullInt64
}

// method for Marshalling nullable int64
func (ni *NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}

// method for Unmarshalling nullable int64

func (ni *NullInt64) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ni.Int64)
	ni.Valid = (err == nil) && ni.Int64 != -1
	return err
}

// all models and their attributes(collumns) are defined here

type Pseudonym struct {
	Id              int    `json:"id" db:"id"`
	PseudonymLeaderId    int    `json:"pseudonymLeaderId" db:"pseudonym_leader_id"`
	PseudonymName        string `json:"pseudonymName" db:"pseudonym_name"`
	PseudonymDescription string `json:"pseudonymDescription" db:"pseudonym_description"`
}

type User struct {
	Id        int    `json:"id" db:"id"`
	Username  string `json:"username" db:"username"`
	FirstName string `json:"firstName" db:"first_name"`
	LastName  string `json:"lastName" db:"last_name"`
	Password  string `json:"password" db:"password"`
	Email     string `json:"email" db:"email"`
}

type Membership struct {
	Id       int  `json:"id" db:"id"`
	UserId   int  `json:"userId" db:"user_id"`
	PseudonymId   int  `json:"pseudonymId" db:"pseudonym_id"`
	IsEditor bool `json:"isEditor" db:"is_editor"`
}
type Post struct {
	Id         int       `json:"id" db:"id"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updated_at"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	AuthorType string    `json:"authorType" db:"author_type"`
	Content    string    `json:"content" db:"content"`
	IsPublic   bool      `json:"isPublic" db:"is_public"`
}
type UserPost struct {
	UserId int `json:"userId" db:"user_id"`
	Post
}

type PseudonymPost struct {
	Post
	PseudonymId int `json:"pseudonymId" db:"pseudonym_id"`
}

type Following struct {
	Id         int `json:"id" db:"id"`
	UserId     int `json:"userId" db:"user_id"`
	FollowerId int `json:"follwerId" db:"follower_id"`
}

type AuthorizationForm struct {
	Username string
	Password string
}
