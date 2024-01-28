package services

import (
	"time"

	"github.com/I1Asyl/ginBerliner/models"
	"github.com/I1Asyl/ginBerliner/pkg/repository"
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

// gets Pseudonym model by its name in the transaction

// gets Pseudonym model by its name from the database
func (a ApiService) GetPseudonymByPseudonymName(pseudonymName string) (models.Pseudonym, error) {
	var pseudonym models.Pseudonym
	pseudonym, err := a.repo.SqlQueries.GetPseudonymByPseudonymName(pseudonymName)

	return pseudonym, err
}

// gets User model by username from the database
func (a ApiService) GetUserByUsername(username string) (models.User, error) {
	var user models.User
	user, err := a.repo.SqlQueries.GetUserByUserame(username)

	return user, err
}

// get all pseudonyms of the user from the database
func (a ApiService) GetPseudonyms(user models.User) ([]models.Pseudonym, error) {

	var pseudonyms []models.Pseudonym
	pseudonyms, err := a.repo.SqlQueries.GetUserPseudonyms(user)

	return pseudonyms, err
}

// create a new pseudonym in the database for the given user
func (a ApiService) CreatePseudonym(pseudonym models.Pseudonym, user models.User) map[string]string {

	invalid := pseudonym.IsValid()
	pseudonym.PseudonymLeaderId = user.Id
	tx := a.repo.SqlQueries.StartTransaction()

	if len(invalid) == 0 {
		if err := tx.AddPseudonym(pseudonym); err != nil {
			invalid["error"] = err.Error()
		} else {
			pseudonym, _ = tx.GetPseudonymByPseudonymName(pseudonym.PseudonymName)
			membership := models.Membership{UserId: pseudonym.PseudonymLeaderId, PseudonymId: pseudonym.Id, IsEditor: true}
			tx.AddMembership(membership)

		}
	}
	err := tx.Commit()
	if err != nil {
		tx.Rollback()
	}

	return invalid
}

// create a new post in the database for the given user or pseudonym
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
			post := models.PseudonymPost{PseudonymId: authorId, Post: post}
			err := a.repo.SqlQueries.AddPseudonymPost(post)
			if err != nil {
				invalid["error"] = err.Error()
			}
		}

	}
	return invalid
}

// func (a ApiService) DeletePost(authorType string, authorId int) map[string]string {
// 	post.CreatedAt = time.Now()
// 	post.UpdatedAt = time.Now()
// 	invalid := post.IsValid()
// 	if len(invalid) == 0 {
// 		if post.AuthorType == "user" {
// 			post := models.UserPost{UserId: authorId, Post: post}
// 			err := a.repo.SqlQueries.AddUserPost(post)
// 			if err != nil {
// 				invalid["error"] = err.Error()
// 			}

// 		} else {
// 			post := models.PseudonymPost{PseudonymId: authorId, Post: post}
// 			err := a.repo.SqlQueries.AddPseudonymPost(post)
// 			if err != nil {
// 				invalid["error"] = err.Error()
// 			}
// 		}

// 	}
// 	return invalid
// }


func (a ApiService) GetPostsFromPseudonyms(user models.User) ([]struct {
	models.Pseudonym
	models.PseudonymPost
}, error) {
	posts, err := a.repo.SqlQueries.GetPseudonymPosts(user)
	return posts, err
}

func (a ApiService) GetNewPostsFromPseudonyms(user models.User) ([]struct {
	models.Pseudonym
	models.PseudonymPost
}, error) {
	posts, err := a.repo.SqlQueries.GetNewPseudonymPosts(user)
	return posts, err
}

func (a ApiService) FollowPseudonym(user models.User, pseudonymName string) error {
	pseudonym, err := a.GetPseudonymByPseudonymName(pseudonymName)
	if err != nil {
		return err
	}
	return a.repo.FollowPseudonym(user, pseudonym)
}

func (a ApiService) FollowUser(follower models.User, userName string) error {
	user, err := a.GetUserByUsername(userName)
	if err != nil {
		return err
	}
	return a.repo.FollowUser(follower, user)
}

func (a ApiService) UnfollowPseudonym(user models.User, pseudonymName string) error {
	pseudonym, err := a.GetPseudonymByPseudonymName(pseudonymName)
	if err != nil {
		return err
	}
	return a.repo.UnfollowPseudonym(user, pseudonym)
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
// 	pseudonymPosts, _ := a.GetPostsFromPseudonyms(user)
// 	userPosts, _ := a.GetPostsFromUsers(user)
// 	posts := pseudonymPosts
// 	posts = append(posts, userPosts...)
// 	return posts, nil
// }

func (a ApiService) GetFollowing(user models.User) ([]models.User, error) {
	users, err := a.repo.SqlQueries.GetFollowing(user)
	return users, err
}

func (a ApiService) DeletePseudonym(pseudonym models.Pseudonym) error {
	err := a.repo.SqlQueries.DeletePseudonym(pseudonym)
	return err
}

func (a ApiService) UpdatePseudonym(pseudonym models.Pseudonym) error {
	err := a.repo.SqlQueries.UpdatePseudonym(pseudonym)
	return err
}
