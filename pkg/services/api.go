package services

import (
	"time"

	"github.com/I1Asyl/berliner_backend/models"
	"github.com/I1Asyl/berliner_backend/pkg/repository"
)

// api service struct
type ApiService struct {
	//database connection
	repo repository.Repository
}

// NewApiService returns a new ApiService instance
func NewApiService(repo repository.Repository) *ApiService {
	return &ApiService{repo: repo}
}

// gets Channel model by its name in the transaction

// gets Channel model by its name from the database
func (a ApiService) GetChannelByName(name string) (models.Channel, error) {
	var channel models.Channel
	channel, err := a.repo.SqlQueries.GetChannelByName(name)

	return channel, err
}

// gets User model by username from the database
func (a ApiService) GetUserByUsername(username string) (models.User, error) {
	var user models.User
	user, err := a.repo.SqlQueries.GetUserByUserame(username)

	return user, err
}

// get all channels of the user from the database
func (a ApiService) GetChannels(user models.User) ([]models.Channel, error) {

	var channels []models.Channel
	channels, err := a.repo.SqlQueries.GetUserChannels(user)

	return channels, err
}

// create a new channel in the database for the given user
func (a ApiService) CreateChannel(channel models.Channel, user models.User) map[string]string {

	invalid := channel.IsValid()
	channel.LeaderId = user.Id
	tx := a.repo.SqlQueries.StartTransaction()

	if len(invalid) == 0 {
		if err := tx.AddChannel(channel); err != nil {
			invalid["error"] = err.Error()
		} else {
			channel, _ = tx.GetChannelByName(channel.Name)
			membership := models.Membership{UserId: channel.LeaderId, ChannelId: channel.Id, IsEditor: true}
			tx.AddMembership(membership)

		}
	}
	err := tx.Commit()
	if err != nil {
		tx.Rollback()
	}

	return invalid
}

// create a new post in the database for the given user or channel
func (a ApiService) CreatePost(post models.Post, authorId int) map[string]string {
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	invalid := post.IsValid()
	if len(invalid) == 0 {
		if post.AuthorType == "user" {
			post := models.UserPost{UserId: authorId, Post: post}
			err := a.repo.SqlQueries.AddUserPost(post)
			if err != nil {
				invalid["error"] = err.Error()
			}

		} else {
			post := models.ChannelPost{ChannelId: authorId, Post: post}
			err := a.repo.SqlQueries.AddChannelPost(post)
			if err != nil {
				invalid["error"] = err.Error()
			}
		}

	}
	return invalid
}

func (a ApiService) DeletePost(post models.Post) error {
	var err error

	if post.AuthorType == "user" {
		userPost := models.UserPost{Post: post}
		err = a.repo.SqlQueries.DeleteUserPost(userPost)
	} else {
		channelPost := models.ChannelPost{Post: post}
		err = a.repo.SqlQueries.DeleteChannelPost(channelPost)
	}
	return err
}

func (a ApiService) GetPostsFromChannels(user models.User) ([]struct {
	models.Channel
	models.ChannelPost
}, error) {
	posts, err := a.repo.SqlQueries.GetChannelPosts(user)
	return posts, err
}

func (a ApiService) GetPostsFromMyChannels(user models.User) ([]struct {
	models.Channel
	models.ChannelPost
}, error) {
	posts, err := a.repo.SqlQueries.GetMyChannelPosts(user)
	return posts, err
}

func (a ApiService) GetNewPostsFromChannels(user models.User) ([]struct {
	models.Channel
	models.ChannelPost
}, error) {
	posts, err := a.repo.SqlQueries.GetNewChannelPosts(user)
	return posts, err
}

func (a ApiService) FollowChannel(user models.User, name string) error {
	channel, err := a.GetChannelByName(name)
	if err != nil {
		return err
	}
	return a.repo.FollowChannel(user, channel)
}

func (a ApiService) FollowUser(follower models.User, userName string) error {
	user, err := a.GetUserByUsername(userName)
	if err != nil {
		return err
	}
	return a.repo.FollowUser(follower, user)
}

func (a ApiService) UnfollowChannel(user models.User, name string) error {
	channel, err := a.GetChannelByName(name)

	if err != nil {
		return err
	}

	return a.repo.UnfollowChannel(user, channel)
}

func (a ApiService) UnfollowUser(follower models.User, userName string) error {
	user, err := a.GetUserByUsername(userName)
	if err != nil {
		return err
	}
	return a.repo.UnfollowUser(follower, user)
}

// get all user's following's posts from the database
func (a ApiService) GetPostsFromUsers(user models.User) ([]struct {
	models.User
	models.UserPost
}, error) {
	posts, err := a.repo.SqlQueries.GetUserPosts(user)
	return posts, err
}

func (a ApiService) GetNewPostsFromUsers(user models.User) ([]struct {
	models.User
	models.UserPost
}, error) {
	posts, err := a.repo.SqlQueries.GetNewUserPosts(user)
	return posts, err
}

// get all posts available for the given user from the database
// func (a ApiService) GetAllPosts(user models.User) ([]models.Post, error) {
// 	channelPosts, _ := a.GetPostsFromChannels(user)
// 	userPosts, _ := a.GetPostsFromUsers(user)
// 	posts := channelPosts
// 	posts = append(posts, userPosts...)
// 	return posts, nil
// }

func (a ApiService) GetFollowing(user models.User) ([]models.User, error) {
	users, err := a.repo.SqlQueries.GetFollowing(user)
	return users, err
}

func (a ApiService) DeleteChannel(channel models.Channel) error {
	err := a.repo.SqlQueries.DeleteChannel(channel)
	return err
}

func (a ApiService) UpdateChannel(channel models.Channel) error {
	err := a.repo.SqlQueries.UpdateChannel(channel)
	return err
}
