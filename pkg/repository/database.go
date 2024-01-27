package repository

import (
	"fmt"
	"log"

	"github.com/I1Asyl/ginBerliner/models"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

type Database struct {
	*sqlx.DB
}

type Transaction struct {
	*sqlx.Tx
}

// SetupOrm sets up the database connection
func NewDatabase(dsn string) Database {
	db, err := sqlx.Open("mysql", dsn+"?parseTime=true")
	if err != nil {
		log.Panic(err)
	}

	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}
	return Database{db}
}

func (db Database) StartTransaction() Transaction {
	tx := db.MustBegin()
	return Transaction{tx}
}

// func (db Database) (tx Transaction) {
// 	tx.Commit()
// }

func (db Database) GetPseudonymByPseudonymName(pseudonymName string) (models.Pseudonym, error) {
	var pseudonym models.Pseudonym
	err := db.Get(&pseudonym, "SELECT * FROM pseudonym WHERE pseudonym_name = ?", pseudonymName)
	return pseudonym, err
}

func (db Transaction) GetPseudonymByPseudonymName(pseudonymName string) (models.Pseudonym, error) {
	var pseudonym models.Pseudonym
	err := db.Get(&pseudonym, "SELECT * FROM pseudonym WHERE pseudonym_name = ?", pseudonymName)
	return pseudonym, err
}

func (db Database) GetUserByUserame(username string) (models.User, error) {
	var user models.User
	err := db.Get(&user, "SELECT * FROM user WHERE username = ?", username)
	return user, err
}
func (db Transaction) GetUserByUserame(username string) (models.User, error) {
	var user models.User
	err := db.Get(&user, "SELECT * FROM user WHERE username = ?", username)
	return user, err
}

func (db Database) GetUserPseudonyms(user models.User) ([]models.Pseudonym, error) {
	var pseudonyms []models.Pseudonym
	err := db.Select(&pseudonyms, "SELECT * FROM pseudonym WHERE pseudonym_leader_id = ?", user.Id)
	return pseudonyms, err
}
func (db Transaction) GetUserPseudonyms(user models.User) ([]models.Pseudonym, error) {
	var pseudonyms []models.Pseudonym
	err := db.Select(&pseudonyms, "SELECT * FROM pseudonym WHERE pseudonym_leader_id = ?", user.Id)
	return pseudonyms, err
}

func (db Database) AddUser(user models.User) error {
	_, err := db.Exec("INSERT INTO user (username, first_name, last_name, email, password) VALUES (?, ?, ?, ?, ?)", user.Username, user.FirstName, user.LastName, user.Email, user.Password)
	return err
}
func (db Transaction) AddUser(user models.User) error {
	_, err := db.Exec("INSERT INTO user (username, first_name, last_name, email, password) VALUES (?, ?, ?, ?, ?)", user.Username, user.FirstName, user.LastName, user.Email, user.Password)
	return err
}

func (db Database) AddMembership(membership models.Membership) error {
	_, err := db.Exec("INSERT INTO membership (pseudonym_id, user_id, is_editor) VALUES (?, ?, ?)", membership.PseudonymId, membership.UserId, membership.IsEditor)
	return err
}
func (db Transaction) AddMembership(membership models.Membership) error {
	_, err := db.Exec("INSERT INTO membership (pseudonym_id, user_id, is_editor) VALUES (?, ?, ?)", membership.PseudonymId, membership.UserId, membership.IsEditor)
	return err
}

func (db Database) AddPseudonym(pseudonym models.Pseudonym) error {
	_, err := db.Exec("INSERT INTO pseudonym (pseudonym_name, pseudonym_leader_id, pseudonym_description) VALUES (?, ?, ?)", pseudonym.PseudonymName, pseudonym.PseudonymLeaderId, pseudonym.PseudonymDescription)
	return err
}
func (db Transaction) AddPseudonym(pseudonym models.Pseudonym) error {
	_, err := db.Exec("INSERT INTO pseudonym (pseudonym_name, pseudonym_leader_id, pseudonym_description) VALUES (?, ?, ?)", pseudonym.PseudonymName, pseudonym.PseudonymLeaderId, pseudonym.PseudonymDescription)
	return err
}

func (db Database) AddUserPost(post models.UserPost) error {
	_, err := db.Exec("INSERT INTO user_post (author_type, content, updated_at, created_at, user_id, is_public) VALUES (?, ?, ?, ?, ?, ?);", post.AuthorType, post.Content, post.UpdatedAt, post.CreatedAt, post.UserId, post.IsPublic)
	return err
}

func (db Database) AddPseudonymPost(post models.PseudonymPost) error {
	_, err := db.Exec("INSERT INTO pseudonym_post (author_type, content, updated_at, created_at, pseudonym_id, is_public) VALUES (?, ?, ?, ?, ?, ?);", post.AuthorType, post.Content, post.UpdatedAt, post.CreatedAt, post.PseudonymId, post.IsPublic)
	fmt.Println(post.PseudonymId)
	return err
}

func (db Database) GetUserPosts(user models.User) ([]struct {
	models.User
	models.UserPost
}, error) {
	//var posts []models.UserPost
	users := "SELECT following.user_id FROM following WHERE following.follower_id=?"
	var newTable []struct {
		models.User
		models.UserPost
	}

	err := db.Select(&newTable, fmt.Sprintf("SELECT user_post.*, user.username, user.first_name, user.last_name FROM user_post LEFT JOIN user on user_post.user_id = user.id WHERE (user_post.user_id in (%v) AND user_post.is_public) OR user_post.user_id = ? ORDER BY updated_at DESC", users), user.Id, user.Id)
	return newTable, err

}

func (db Database) GetPseudonymPosts(user models.User) ([]struct {
	models.Pseudonym
	models.PseudonymPost
}, error) {
	//var posts []models.UserPost
	pseudonyms := "SELECT membership.pseudonym_id FROM membership WHERE membership.user_id=?"
	var newTable []struct {
		models.Pseudonym
		models.PseudonymPost
	}
	err := db.Select(&newTable, fmt.Sprintf("SELECT pseudonym_post.*, pseudonym.pseudonym_name FROM pseudonym_post LEFT JOIN pseudonym on pseudonym_post.pseudonym_id = pseudonym.id WHERE pseudonym_post.pseudonym_id in (%v) AND (pseudonym_post.is_public OR pseudonym.pseudonym_leader_id = ?) ORDER BY updated_at DESC", pseudonyms), user.Id, user.Id)
	return newTable, err

}

func (db Database) GetNewUserPosts(user models.User) ([]struct {
	models.User
	models.UserPost
}, error) {
	//var posts []models.UserPost
	users := "SELECT following.user_id FROM following WHERE following.follower_id=?"
	var newTable []struct {
		models.User
		models.UserPost
	}

	err := db.Select(&newTable, fmt.Sprintf("SELECT user_post.*, user.username, user.first_name, user.last_name FROM user_post LEFT JOIN user on user_post.user_id = user.id WHERE user_post.user_id NOT in (%v) AND NOT user_post.user_id = ? AND user_post.is_public = 1 ORDER BY updated_at DESC", users), user.Id, user.Id)
	return newTable, err
}

func (db Database) GetNewPseudonymPosts(user models.User) ([]struct {
	models.Pseudonym
	models.PseudonymPost
}, error) {
	//var posts []models.UserPost
	users := "SELECT membership.pseudonym_id FROM membership WHERE membership.user_id=?"
	var newTable []struct {
		models.Pseudonym
		models.PseudonymPost
	}

	err := db.Select(&newTable, fmt.Sprintf("SELECT pseudonym_post.*, pseudonym.pseudonym_name FROM pseudonym_post LEFT JOIN pseudonym on pseudonym_post.pseudonym_id = pseudonym.id WHERE pseudonym_post.pseudonym_id NOT in (%v) AND pseudonym_post.is_public = 1 ORDER BY updated_at DESC", users), user.Id)
	return newTable, err
}

func (db Database) FollowPseudonym(user models.User, pseudonym models.Pseudonym) error {
	query := "INSERT INTO membership (pseudonym_id, user_id, is_editor) VALUES (?, ?, ?)"
	_, err := db.Exec(query, pseudonym.Id, user.Id, 0)
	return err
}
func (db Database) FollowUser(follower models.User, user models.User) error {
	query := "INSERT INTO following (user_id, follower_id) VALUES (?, ?)"
	_, err := db.Exec(query, user.Id, follower.Id)
	return err
}

func (db Database) GetFollowing(user models.User) ([]models.User, error) {
	var users []models.User
	err := db.Select(&users, "SELECT * FROM following WHERE follower_id = ?", user.Id)
	return users, err
}

func (db Database) DeletePseudonym(pseudonym models.Pseudonym) error {
	_, err := db.Exec("DELETE FROM pseudonym WHERE pseudonym_id = ?", pseudonym.Id)
	return err
}

func (db Database) AddFollowing(following models.Following) error {
	_, err := db.Exec("INSERT INTO following (follower_id, user_id) VALUES (?, ?)", following.FollowerId, following.UserId)
	return err
}

func (db Database) UpdatePseudonym(pseudonym models.Pseudonym) error {
	if pseudonym.PseudonymName != "" {
		_, err := db.Exec("UPDATE pseudonym SET pseudonym_name = ? WHERE pseudonym_id = ?", pseudonym.PseudonymName, pseudonym.Id)
		if err != nil {
			return err
		}
	}
	if pseudonym.PseudonymDescription != "" {
		_, err := db.Exec("UPDATE pseudonym SET pseudonym_description = ? WHERE pseudonym_id = ?", pseudonym.PseudonymName, pseudonym.Id)
		if err != nil {
			return err
		}
	}
	return nil
}
