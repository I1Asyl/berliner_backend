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

func (db Database) GetTeamByTeamName(teamName string) (models.Team, error) {
	var team models.Team
	err := db.Get(&team, "SELECT * FROM team WHERE team_name = ?", teamName)
	return team, err
}

func (db Transaction) GetTeamByTeamName(teamName string) (models.Team, error) {
	var team models.Team
	err := db.Get(&team, "SELECT * FROM team WHERE team_name = ?", teamName)
	return team, err
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

func (db Database) GetUserTeams(user models.User) ([]models.Team, error) {
	var teams []models.Team
	err := db.Select(&teams, "SELECT * FROM team WHERE team_leader_id = ?", user.Id)
	return teams, err
}
func (db Transaction) GetUserTeams(user models.User) ([]models.Team, error) {
	var teams []models.Team
	err := db.Select(&teams, "SELECT * FROM team WHERE team_leader_id = ?", user.Id)
	return teams, err
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
	_, err := db.Exec("INSERT INTO membership (team_id, user_id, is_editor) VALUES (?, ?, ?)", membership.TeamId, membership.UserId, membership.IsEditor)
	return err
}
func (db Transaction) AddMembership(membership models.Membership) error {
	_, err := db.Exec("INSERT INTO membership (team_id, user_id, is_editor) VALUES (?, ?, ?)", membership.TeamId, membership.UserId, membership.IsEditor)
	return err
}

func (db Database) AddTeam(team models.Team) error {
	_, err := db.Exec("INSERT INTO team (team_name, team_leader_id, team_description) VALUES (?, ?, ?)", team.TeamName, team.TeamLeaderId, team.TeamDescription)
	return err
}
func (db Transaction) AddTeam(team models.Team) error {
	_, err := db.Exec("INSERT INTO team (team_name, team_leader_id, team_description) VALUES (?, ?, ?)", team.TeamName, team.TeamLeaderId, team.TeamDescription)
	return err
}

func (db Database) AddUserPost(post models.UserPost) error {
	_, err := db.Exec("INSERT INTO user_post (author_type, content, updated_at, created_at, user_id, is_public) VALUES (?, ?, ?, ?, ?, ?);", post.AuthorType, post.Content, post.UpdatedAt, post.CreatedAt, post.UserId, post.IsPublic)
	return err
}

func (db Database) AddTeamPost(post models.TeamPost) error {
	_, err := db.Exec("INSERT INTO team_post (author_type, content, updated_at, created_at, team_id, is_public) VALUES (?, ?, ?, ?, ?, ?);", post.AuthorType, post.Content, post.UpdatedAt, post.CreatedAt, post.TeamId, post.IsPublic)
	fmt.Println(post.TeamId)
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

	err := db.Select(&newTable, fmt.Sprintf("SELECT user_post.*, user.username, user.first_name, user.last_name FROM user_post LEFT JOIN user on user_post.user_id = user.id WHERE user_post.user_id in (%v) OR user_post.user_id = ? ORDER BY updated_at DESC", users), user.Id, user.Id)
	return newTable, err

}

func (db Database) GetTeamPosts(user models.User) ([]struct {
	models.Team
	models.TeamPost
}, error) {
	//var posts []models.UserPost
	teams := "SELECT membership.team_id FROM membership WHERE membership.user_id=?"
	var newTable []struct {
		models.Team
		models.TeamPost
	}
	err := db.Select(&newTable, fmt.Sprintf("SELECT team_post.*, team.team_name FROM team_post LEFT JOIN team on team_post.team_id = team.id WHERE team_post.team_id in (%v) ORDER BY updated_at DESC", teams), user.Id)
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

func (db Database) GetNewTeamPosts(user models.User) ([]struct {
	models.Team
	models.TeamPost
}, error) {
	//var posts []models.UserPost
	users := "SELECT membership.team_id FROM membership WHERE membership.user_id=?"
	var newTable []struct {
		models.Team
		models.TeamPost
	}

	err := db.Select(&newTable, fmt.Sprintf("SELECT team_post.*, team.team_name FROM team_post LEFT JOIN team on team_post.team_id = team.id WHERE team_post.team_id NOT in (%v) AND team_post.is_public = 1 ORDER BY updated_at DESC", users), user.Id)
	return newTable, err
}

func (db Database) FollowTeam(user models.User, team models.Team) error {
	query := "INSERT INTO membership (team_id, user_id, is_editor) VALUES (?, ?, ?)"
	_, err := db.Exec(query, team.Id, user.Id, 0)
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

func (db Database) DeleteTeam(team models.Team) error {
	_, err := db.Exec("DELETE FROM team WHERE team_id = ?", team.Id)
	return err
}

func (db Database) AddFollowing(following models.Following) error {
	_, err := db.Exec("INSERT INTO following (follower_id, user_id) VALUES (?, ?)", following.FollowerId, following.UserId)
	return err
}

func (db Database) UpdateTeam(team models.Team) error {
	if team.TeamName != "" {
		_, err := db.Exec("UPDATE team SET team_name = ? WHERE team_id = ?", team.TeamName, team.Id)
		if err != nil {
			return err
		}
	}
	if team.TeamDescription != "" {
		_, err := db.Exec("UPDATE team SET team_description = ? WHERE team_id = ?", team.TeamName, team.Id)
		if err != nil {
			return err
		}
	}
	return nil
}
