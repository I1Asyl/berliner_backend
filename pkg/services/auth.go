package services

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/I1Asyl/ginBerliner/models"
	"github.com/I1Asyl/ginBerliner/pkg/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// claims for jwt token
type UserClaims struct {
	Username string
	jwt.RegisteredClaims
}

// auth service struct
type AuthService struct {
	//database connection
	repo repository.Repository
}

// NewAuthService returns a new AuthService instance
func NewAuthService(repo repository.Repository) *AuthService {
	return &AuthService{repo: repo}
}

// check if user exists and password is correct
func (a AuthService) CheckUserAndPassword(userForm models.AuthorizationForm) (bool, error) {
	user, err := a.repo.SqlQueries.GetUserByUserame(userForm.Username)
	if err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userForm.Password))
	if err != nil {
		return false, err
	}
	return true, nil
}

// parse jwt token
func (a AuthService) ParseToken(token string) (string, error) {
	claims := UserClaims{}
	ans, err := jwt.ParseWithClaims(token, &claims, func(tok *jwt.Token) (interface{}, error) {
		if _, ok := tok.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", errors.New("Error")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if claims, ok := ans.Claims.(*UserClaims); ok && ans.Valid {
		return claims.Username, nil
	} else {
		return "", err
	}

}

// generate jwt token
func (a AuthService) GenerateToken(user models.AuthorizationForm, issueTime time.Time, expireTime time.Time) (string, error) {
	claims := UserClaims{
		user.Username,

		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(issueTime),
			Issuer:    "test",
			Subject:   "somebody",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET is not set, set your jwt password in .env file")
	}
	signedKey := []byte(os.Getenv("JWT_SECRET"))
	signed, err := token.SignedString(signedKey)
	return signed, err

}

// add user to the database
func (a AuthService) AddUser(user models.User) map[string]string {
	invalid := user.IsValid()
	if len(invalid) == 0 {
		user.Password = a.HashPassword(user.Password)
		if err := a.repo.SqlQueries.AddUser(user); err != nil {
			invalid["common"] = err.Error()
		} else {
			a.repo.SqlQueries.GetUserByUserame(user.Username)
			following := models.Following{UserId: user.Id, FollowerId: user.Id}
			a.repo.SqlQueries.AddFollowing(following)
		}
	}
	return invalid
}

// get User model from username
func (a AuthService) GetUserFromUsername(username string) (models.User, error) {
	user, err := a.repo.SqlQueries.GetUserByUserame(username)

	return user, err
}

// hash password
func (a AuthService) HashPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 11)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	return string(hashed)
}
