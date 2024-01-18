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

// gets Team model by its name in the transaction

// gets Team model by its name from the database
func (a ApiService) GetTeamByTeamName(teamName string) (models.Team, error) {
	var team models.Team
	team, err := a.repo.SqlQueries.GetTeamByTeamName(teamName)

	return team, err
}

// gets User model by username from the database
func (a ApiService) GetUserByUsername(username string) (models.User, error) {
	var user models.User
	user, err := a.repo.SqlQueries.GetUserByUserame(username)

	return user, err
}

// get all teams of the user from the database
func (a ApiService) GetTeams(user models.User) ([]models.Team, error) {

	var teams []models.Team
	teams, err := a.repo.SqlQueries.GetUserTeams(user)

	return teams, err
}

// create a new team in the database for the given user
func (a ApiService) CreateTeam(team models.Team, user models.User) map[string]string {

	invalid := team.IsValid()
	team.TeamLeaderId = user.Id
	tx := a.repo.SqlQueries.StartTransaction()

	if len(invalid) == 0 {
		if err := tx.AddTeam(team); err != nil {
			invalid["error"] = err.Error()
		} else {
			team, _ = tx.GetTeamByTeamName(team.TeamName)
			membership := models.Membership{UserId: team.TeamLeaderId, TeamId: team.Id, IsEditor: true}
			tx.AddMembership(membership)

		}
	}
	err := tx.Commit()
	if err != nil {
		tx.Rollback()
	}

	return invalid
}

// create a new post in the database for the given user or team
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
			post := models.TeamPost{TeamId: authorId, Post: post}
			err := a.repo.SqlQueries.AddTeamPost(post)
			if err != nil {
				invalid["error"] = err.Error()
			}
		}

	}
	return invalid
}

func (a ApiService) GetPostsFromTeams(user models.User) ([]struct {
	models.Team
	models.TeamPost
}, error) {
	posts, err := a.repo.SqlQueries.GetTeamPosts(user)
	return posts, err
}

func (a ApiService) GetNewPostsFromTeams(user models.User) ([]struct {
	models.Team
	models.TeamPost
}, error) {
	posts, err := a.repo.SqlQueries.GetNewTeamPosts(user)
	return posts, err
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
// 	teamPosts, _ := a.GetPostsFromTeams(user)
// 	userPosts, _ := a.GetPostsFromUsers(user)
// 	posts := teamPosts
// 	posts = append(posts, userPosts...)
// 	return posts, nil
// }

func (a ApiService) GetFollowing(user models.User) ([]models.User, error) {
	users, err := a.repo.SqlQueries.GetFollowing(user)
	return users, err
}

func (a ApiService) DeleteTeam(team models.Team) error {
	err := a.repo.SqlQueries.DeleteTeam(team)
	return err
}

func (a ApiService) UpdateTeam(team models.Team) error {
	err := a.repo.SqlQueries.UpdateTeam(team)
	return err
}
