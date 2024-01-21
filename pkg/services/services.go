package services

import (
	"time"

	"github.com/I1Asyl/ginBerliner/models"
	"github.com/I1Asyl/ginBerliner/pkg/repository"
)

//go:generate mockgen -source=services.go -destination=mocks/services.go

// all authorization services
type Authorization interface {
	AddUser(user models.User) map[string]string
	HashPassword(password string) string
	GenerateToken(user models.AuthorizationForm, issueTime time.Time, expireTime time.Time) (string, error)
	ParseToken(token string) (string, error)
	CheckUserAndPassword(userForm models.AuthorizationForm) (bool, error)
}

// all api services
type Api interface {
	FollowTeam(user models.User, teamName string) error
	FollowUser(follower models.User, userName string) error
	DeleteTeam(team models.Team) error
	UpdateTeam(team models.Team) error
	GetFollowing(user models.User) ([]models.User, error)
	GetUserByUsername(username string) (models.User, error)
	GetTeamByTeamName(teamName string) (models.Team, error)
	CreatePost(post models.Post, autthorId int) map[string]string
	GetPostsFromTeams(user models.User) ([]struct {
		models.Team
		models.TeamPost
	}, error)
	GetPostsFromUsers(user models.User) ([]struct {
		models.User
		models.UserPost
	}, error)
	GetNewPostsFromUsers(user models.User) ([]struct {
		models.User
		models.UserPost
	}, error)
	GetNewPostsFromTeams(user models.User) ([]struct {
		models.Team
		models.TeamPost
	}, error)
	//GetAllPosts(user models.User) ([]models.Post, error)
	GetTeams(user models.User) ([]models.Team, error)
	CreateTeam(team models.Team, user models.User) map[string]string
}

// func clearAllData() {
// 	query := "DELETE FROM users"
// }

// struct for all needed services
type Services struct {
	Authorization
	Api
}

// returns new Services with all needed authorization and api services
func NewService(repo *repository.Repository) *Services {
	return &Services{Authorization: NewAuthService(*repo), Api: NewApiService(*repo)}
}
