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
	FollowPseudonym(user models.User, pseudonymName string) error
	FollowUser(follower models.User, userName string) error
	DeletePseudonym(pseudonym models.Pseudonym) error
	UpdatePseudonym(pseudonym models.Pseudonym) error
	GetFollowing(user models.User) ([]models.User, error)
	GetUserByUsername(username string) (models.User, error)
	GetPseudonymByPseudonymName(pseudonymName string) (models.Pseudonym, error)
	CreatePost(post models.Post, autthorId int) map[string]string
	GetPostsFromPseudonyms(user models.User) ([]struct {
		models.Pseudonym
		models.PseudonymPost
	}, error)
	GetPostsFromUsers(user models.User) ([]struct {
		models.User
		models.UserPost
	}, error)
	GetNewPostsFromUsers(user models.User) ([]struct {
		models.User
		models.UserPost
	}, error)
	GetNewPostsFromPseudonyms(user models.User) ([]struct {
		models.Pseudonym
		models.PseudonymPost
	}, error)
	//GetAllPosts(user models.User) ([]models.Post, error)
	GetPseudonyms(user models.User) ([]models.Pseudonym, error)
	CreatePseudonym(pseudonym models.Pseudonym, user models.User) map[string]string
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
