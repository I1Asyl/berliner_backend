package repository

import (
	"github.com/I1Asyl/ginBerliner/models"
	"github.com/golang-migrate/migrate/v4"
)

type SqlQueries interface {
	FollowPseudonym(user models.User, pseudonym models.Pseudonym) error
	FollowUser(follower models.User, user models.User) error
	GetPseudonymByPseudonymName(pseudonymName string) (models.Pseudonym, error)
	GetUserByUserame(pseudonymName string) (models.User, error)
	GetUserPseudonyms(user models.User) ([]models.Pseudonym, error)
	AddMembership(models.Membership) error
	AddUser(models.User) error
	AddPseudonym(pseudonym models.Pseudonym) error
	AddUserPost(post models.UserPost) error
	AddPseudonymPost(post models.PseudonymPost) error
	AddFollowing(following models.Following) error
	StartTransaction() Transaction
	GetPseudonymPosts(user models.User) ([]struct {
		models.Pseudonym
		models.PseudonymPost
	}, error)
	GetUserPosts(user models.User) ([]struct {
		models.User
		models.UserPost
	}, error)
	GetNewUserPosts(user models.User) ([]struct {
		models.User
		models.UserPost
	}, error)
	GetNewPseudonymPosts(user models.User) ([]struct {
		models.Pseudonym
		models.PseudonymPost
	}, error)
	GetFollowing(user models.User) ([]models.User, error)
	UpdatePseudonym(pseudonym models.Pseudonym) error
	DeletePseudonym(pseudonym models.Pseudonym) error
}

type Repository struct {
	SqlQueries
	Migration *migrate.Migrate
}

//go:generate mockgen -source=repository.go -destination=mocks/repository.go

func NewRepository(dsn string, migrationsPath string) *Repository {
	return &Repository{SqlQueries: NewDatabase(dsn), Migration: NewMigration(dsn, migrationsPath)}
}
