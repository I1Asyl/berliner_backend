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

type Team struct {
	Id              int    `json:"id"`
	TeamLeaderId    int    `json:"teamLeaderId"`
	TeamName        string `json:"teamName"`
	TeamDescription string `json:"teamDescription"`
}

type User struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
	Email     string `json:"email"`
}

type Membership struct {
	Id       int  `json:"id"`
	UserId   int  `json:"userId"`
	TeamId   int  `json:"teamId"`
	IsEditor bool `json:"isEditor"`
}

type Post struct {
	Id           int       `json:"id"`
	AuthorType   string    `json:"authorType"`
	UserAuthorId NullInt64 `json:"userAuthorId"`
	TeamAuthorId NullInt64 `json:"teamAuthorId"`
	UpdatedAt    time.Time `json:"updatedAt" xorm:"updated"`
	CreatedAt    time.Time `json:"createdAt" xorm:"created"`
	Content      string    `json:"content"`
}

type Following struct {
	Id         int `json:"id"`
	UserId     int `json:"userId"`
	FollowerId int `json:"follwerId"`
}

type AuthorizationForm struct {
	Username string
	Password string
}
