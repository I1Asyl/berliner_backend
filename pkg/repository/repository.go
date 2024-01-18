package repository

import (
	"github.com/I1Asyl/ginBerliner/models"
	"github.com/golang-migrate/migrate/v4"
)

type SqlQueries interface {
	GetTeamByTeamName(teamName string) (models.Team, error)
	GetUserByUserame(teamName string) (models.User, error)
	GetUserTeams(user models.User) ([]models.Team, error)
	AddMembership(models.Membership) error
	AddUser(models.User) error
	AddTeam(team models.Team) error
	AddUserPost(post models.UserPost) error
	AddTeamPost(post models.TeamPost) error
	AddFollowing(following models.Following) error
	StartTransaction() Transaction
	GetTeamPosts(user models.User) ([]struct {
		models.Team
		models.TeamPost
	}, error)
	GetUserPosts(user models.User) ([]struct {
		models.User
		models.UserPost
	}, error)
	GetNewUserPosts(user models.User) ([]struct {
		models.User
		models.UserPost
	}, error)
	GetNewTeamPosts(user models.User) ([]struct {
		models.Team
		models.TeamPost
	}, error)
	GetFollowing(user models.User) ([]models.User, error)
	UpdateTeam(team models.Team) error
	DeleteTeam(team models.Team) error
}

type Repository struct {
	SqlQueries
	Migration *migrate.Migrate
}

//go:generate mockgen -source=repository.go -destination=mocks/repository.go

func NewRepository(dsn string, migrationsPath string) *Repository {
	return &Repository{SqlQueries: NewDatabase(dsn), Migration: NewMigration(dsn, migrationsPath)}
}
