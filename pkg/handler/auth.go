package handler

import (
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
		ctx.JSON(401, gin.H{})
		return
	}

	//check if user data is valid
	if invalid := h.services.Authorization.AddUser(user); len(invalid) > 0 {
		ctx.AbortWithStatusJSON(422, invalid)
		return
	}
	ctx.JSON(200, gin.H{})
}

// login method for logging in user
func (h *Handler) login(ctx *gin.Context) {
	var user models.AuthorizationForm
	//check if user is valid json type
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(401, gin.H{})
		return
	}

	//check if user data is valid
	exist, err := h.services.Authorization.CheckUserAndPassword(user)
	if !exist || err != nil {
		ctx.JSON(401, gin.H{
			"error": "username or password is incorrect",
		})
		return
	}
	// generate token
	token, err := h.services.Authorization.GenerateToken(user)
	if err != nil {
		ctx.JSON(401, gin.H{})
		return
	}
	ctx.JSON(200, gin.H{
		"token": token,
	})

}
