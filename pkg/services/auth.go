package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/I1Asyl/ginBerliner/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"xorm.io/xorm"
)

// claims for jwt token
type UserClaims struct {
	Username string
	jwt.RegisteredClaims
}

// auth service struct
type AuthService struct {
	//database connection
	orm xorm.Engine
}

// NewAuthService returns a new AuthService instance
func NewAuthService(orm xorm.Engine) *AuthService {
	return &AuthService{orm: orm}
}

// check if user exists and password is correct
func (a AuthService) CheckUserAndPassword(userForm models.AuthorizationForm) (bool, error) {
	var user models.User
	a.orm.Where("username = ?", userForm.Username).Get(&user)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userForm.Password))
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
func (a AuthService) GenerateToken(user models.AuthorizationForm) (string, error) {
	claims := UserClaims{
		user.Username,

		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 365 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "test",
			Subject:   "somebody",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedKey := []byte(os.Getenv("JWT_SECRET"))
	signed, err := token.SignedString(signedKey)
	return signed, err

}

// add user to the database
func (a AuthService) AddUser(user models.User) map[string]string {
	invalid := user.IsValid()
	if len(invalid) == 0 {
		user.Password = a.HashPassword(user.Password)
		if _, err := a.orm.Insert(user); err != nil {
			invalid["common"] = err.Error()
		} else {
			a.orm.Where("username=?", user.Username).Get(&user)
			following := models.Following{UserId: user.Id, FollowerId: user.Id}
			a.orm.Insert(following)
		}
	}
	return invalid
}

// get User model from username
func (a AuthService) GetUserFromUsername(username string) (models.User, error) {
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

// hash password
func (a AuthService) HashPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 11)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	return string(hashed)
}
