package handler

import (
	"errors"
	"time"
	"fmt"
	"github.com/I1Asyl/ginBerliner/models"
	"github.com/gin-gonic/gin"
)

type UserRepository interface {
	AddUser()
}

// signUp method for creating new user
// recieves an User model

func (h *Handler) signUp(ctx *gin.Context) {
	var user models.User
	//check if user is valid json type
	if err := ctx.BindJSON(&user); err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	//check if user data is valid
	if invalid := h.services.Authorization.AddUser(user); len(invalid) > 0 {
		if err, ok := invalid["error"]; ok {
			ctx.AbortWithError(500, errors.New(err))
		}
		ctx.AbortWithStatusJSON(422, invalid)
		fmt.Println(invalid);
		return
	}
	ctx.JSON(200, gin.H{})
}

// login method for logging in user
func (h *Handler) login(ctx *gin.Context) {
	var user models.AuthorizationForm
	//check if user is valid json type
	if err := ctx.BindJSON(&user); err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	//check if user data is valid
	exist, err := h.services.Authorization.CheckUserAndPassword(user)
	if !exist || err != nil {
		ctx.AbortWithError(401, errors.New("username or password is incorrect"))
		return
	}
	// generate token
	token, err := h.services.Authorization.GenerateToken(user, time.Now(), time.Now().Add(time.Hour*24))
	if err != nil {
		ctx.AbortWithError(500, errors.New(err.Error()))
		return
	}
	ctx.JSON(200, gin.H{
		"token": token,
	})

}
