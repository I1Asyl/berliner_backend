package services

import (
	"time"

	"github.com/I1Asyl/berliner_backend/models"
	"github.com/I1Asyl/berliner_backend/pkg/repository"
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
	FollowChannel(user models.User, name string) error
	FollowUser(follower models.User, userName string) error
	UnfollowChannel(user models.User, name string) error
	UnfollowUser(follower models.User, userName string) error
	DeleteChannel(channel models.Channel) error
	UpdateChannel(channel models.Channel) error
	GetFollowing(user models.User) ([]models.User, error)
	GetUserByUsername(username string) (models.User, error)
	GetChannelByName(name string) (models.Channel, error)
	CreatePost(post models.Post, autthorId int) map[string]string
	DeletePost(post models.Post) error
	GetPostsFromMyChannels(user models.User) ([]struct {
		models.Channel
		models.ChannelPost
	}, error)
	GetPostsFromChannels(user models.User) ([]struct {
		models.Channel
		models.ChannelPost
	}, error)
	GetPostsFromUsers(user models.User) ([]struct {
		models.User
		models.UserPost
	}, error)
	GetNewPostsFromUsers(user models.User) ([]struct {
		models.User
		models.UserPost
	}, error)
	GetNewPostsFromChannels(user models.User) ([]struct {
		models.Channel
		models.ChannelPost
	}, error)
	//GetAllPosts(user models.User) ([]models.Post, error)
	GetChannels(user models.User) ([]models.Channel, error)
	CreateChannel(channel models.Channel, user models.User) map[string]string
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
