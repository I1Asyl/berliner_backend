package handler

// Handler for apis

import (
	"errors"

	"github.com/I1Asyl/ginBerliner/models"
	"github.com/gin-gonic/gin"
)

// method for getting user teams
func (h Handler) getTeams(ctx *gin.Context) {
	res, _ := ctx.Get("user")
	user := res.(models.User)
	ans, err := h.services.Api.GetTeams(user)

	if err != nil {
		ctx.JSON(400, gin.H{})
		return
	}

	ctx.JSON(200, ans)
}

// creating a team for an user
func (h Handler) createTeam(ctx *gin.Context) {
	var team models.Team
	if err := ctx.BindJSON(&team); err != nil {
		ctx.AbortWithError(401, errors.New("input json can not be marshalled to the team model"))
		return
	}
	res, _ := ctx.Get("user")
	user := res.(models.User)
	if invalid := h.services.Api.CreateTeam(team, user); len(invalid) > 0 {
		ctx.Errors = append(ctx.Errors, &gin.Error{Err: errors.New("invalid data")})
		ctx.JSON(422, invalid)
		return
	}

	ctx.JSON(200, gin.H{})
}

// method for posting a post
func (h Handler) createPost(ctx *gin.Context) {
	var post models.Post
	if err := ctx.BindJSON(&post); err != nil {
		ctx.AbortWithError(401, errors.New("input json can not be marshalled to the post model"))
		return
	}

	if invalid := h.services.Api.CreatePost(post); len(invalid) > 0 {
		ctx.Errors = append(ctx.Errors, &gin.Error{Err: errors.New("invalid data")})
		ctx.JSON(422, invalid)
		return
	}
	ctx.JSON(200, gin.H{})
}

// method for reading posts
func (h Handler) getPosts(ctx *gin.Context) {
	res, _ := ctx.Get("user")
	authorType := ctx.DefaultQuery("author", "")
	user := res.(models.User)
	var ans []models.Post
	var err error

	// check if needed post should be written by team, user or all
	if authorType == "team" {
		ans, err = h.services.Api.GetPostsFromTeams(user)
	} else if authorType == "user" {
		ans, err = h.services.Api.GetPostsFromUsers(user)
	} else {
		ans, err = h.services.Api.GetAllPosts(user)
	}
	if err != nil {
		ctx.JSON(400, gin.H{})
		return
	}
	ctx.JSON(200, ans)
}

func (h Handler) getFollowing(ctx *gin.Context) {
	res, _ := ctx.Get("user")
	user := res.(models.User)
	ans, err := h.services.Api.GetFollowing(user)
	if err != nil {
		ctx.JSON(400, gin.H{})
		return
	}
	ctx.JSON(200, ans)
}
