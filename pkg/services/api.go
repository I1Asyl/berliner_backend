package services

import (
	"errors"
	"fmt"

	"github.com/I1Asyl/ginBerliner/models"
	"xorm.io/xorm"
)

// api service struct
type ApiService struct {
	//database connection
	orm xorm.Engine
}

// NewApiService returns a new ApiService instance
func NewApiService(orm xorm.Engine) *ApiService {
	return &ApiService{orm: orm}
}

// gets Team model by its name in the transaction
func (session *Transaction) GetTeamByTeamName(teamName string) (models.Team, error) {
	var team models.Team
	ok, err := session.Where("team_name=?", teamName).Get(&team)
	if err != nil {
		return models.Team{}, err
	}
	if !ok {
		return models.Team{}, errors.New("team does not exist")
	}
	return team, nil
}

// gets User model by username in the transaction
func (session *Transaction) GetTeamByUsername(username string) (models.User, error) {
	var user models.User
	ok, err := session.Where("username=?", username).Get(&user)
	if err != nil {
		return models.User{}, err
	}
	if !ok {
		return models.User{}, errors.New("team does not exist")
	}
	return user, nil
}

// gets Team model by its name from the database
func (a ApiService) GetTeamByTeamName(teamName string) (models.Team, error) {
	var team models.Team
	ok, err := a.orm.Where("team_name=?", teamName).Get(&team)

	if err != nil {
		return models.Team{}, err
	}
	if !ok {
		return models.Team{}, errors.New("team does not exist")
	}
	return team, nil
}

// gets User model by username from the database
func (a ApiService) GetUserByUsername(username string) (models.User, error) {
	var user models.User
	ok, err := a.orm.Where("username=?", username).Get(&user)

	if err != nil {
		return models.User{}, err
	}
	if !ok {
		return models.User{}, errors.New("user does not exist")
	}
	return user, nil
}

// get all teams of the user from the database
func (a ApiService) GetTeams(user models.User) ([]models.Team, error) {

	var teams []models.Team
	a.orm.Where("team_leader_id=?", user.Id).Find(&teams)
	return teams, nil
}

// create a new team in the database for the given user
func (a ApiService) CreateTeam(team models.Team, user models.User) map[string]string {
	session := &Transaction{a.orm.NewSession()}
	defer session.Close()

	session.Begin()
	invalid := team.IsValid()
	team.TeamLeaderId = user.Id
	if len(invalid) == 0 {
		if _, err := session.Insert(team); err != nil {
			invalid["common"] = "Team with this name already exists"
		} else {
			team, _ = session.GetTeamByTeamName(team.TeamName)
			membership := models.Membership{UserId: team.TeamLeaderId, TeamId: team.Id, IsEditor: true}
			_, err := session.Insert(membership)
			if err != nil {
				session.Rollback()
			}
		}
	}
	err := session.Commit()
	if err != nil {
		session.Rollback()
	}

	return invalid
}

// create a new post in the database for the given user or team
func (a ApiService) CreatePost(post models.Post) map[string]string {

	invalid := post.IsValid()

	if len(invalid) == 0 {
		_, err := a.orm.Insert(post)
		if err != nil {
			invalid["common"] = err.Error()
		}
	}

	return invalid
}

// get all user's team posts from the database
func (a ApiService) GetPostsFromTeams(user models.User) ([]models.Post, error) {
	var posts []models.Post
	//var teams []models.Teamconst content = ref("")
	teams := "SELECT membership.team_id FROM membership WHERE membership.user_id=?"
	a.orm.Where(fmt.Sprintf("team_author_id IN (%v)", teams), user.Id).Find(&posts)
	return posts, nil
}

// get all user's following's posts from the database
func (a ApiService) GetPostsFromUsers(user models.User) ([]models.Post, error) {
	var posts []models.Post
	//var teams []models.Teamconst content = ref("")
	users := "SELECT following.user_id FROM following WHERE following.follower_id=?"
	a.orm.Where(fmt.Sprintf("user_author_id IN (%v)", users), user.Id).Find(&posts)
	return posts, nil
}

// get all posts available for the given user from the database
func (a ApiService) GetAllPosts(user models.User) ([]models.Post, error) {
	teamPosts, _ := a.GetPostsFromTeams(user)
	userPosts, _ := a.GetPostsFromUsers(user)
	posts := teamPosts
	posts = append(posts, userPosts...)
	return posts, nil
}

func (a ApiService) GetFollowing(user models.User) ([]models.User, error) {
	var users []models.User
	query := "SELECT following.user_id FROM following WHERE following.follower_id=?"

	a.orm.Where(fmt.Sprintf("id IN (%v)", query), user.Id).Find(&users)
	return users, nil
}
