package services

import (
	"github.com/I1Asyl/ginBerliner/models"
	"github.com/I1Asyl/ginBerliner/pkg/repository"
	"xorm.io/xorm"
)

type Transaction struct {
	*xorm.Session
}

//go:generate mockgen -source=services.go -destination=mocks/services.go

// all authorization services
type Authorization interface {
	AddUser(user models.User) map[string]string
	HashPassword(password string) string
	GenerateToken(user models.AuthorizationForm) (string, error)
	ParseToken(token string) (string, error)
	CheckUserAndPassword(userForm models.AuthorizationForm) (bool, error)
}

// all api services
type Api interface {
	GetFollowing(user models.User) ([]models.User, error)
	GetUserByUsername(username string) (models.User, error)
	GetTeamByTeamName(teamName string) (models.Team, error)
	CreatePost(post models.Post) map[string]string
	GetPostsFromTeams(user models.User) ([]models.Post, error)
	GetPostsFromUsers(user models.User) ([]models.Post, error)
	GetAllPosts(user models.User) ([]models.Post, error)
	GetTeams(user models.User) ([]models.Team, error)
	CreateTeam(team models.Team, user models.User) map[string]string
}

// struct for all needed services
type Services struct {
	Authorization
	Api
}

// returns new Services with all needed authorization and api services
func NewService(repo *repository.Repository) *Services {
	return &Services{Authorization: NewAuthService(*repo.Orm), Api: NewApiService(*repo.Orm)}
}
