package handler

// Handler for apis

import (
	"errors"
	"strconv"

	"github.com/I1Asyl/berliner_backend/models"
	"github.com/gin-gonic/gin"
)

// method for getting user channels
func (h Handler) getChannels(ctx *gin.Context) {
	res, _ := ctx.Get("user")
	user := res.(models.User)
	ans, err := h.services.Api.GetChannels(user)

	if err != nil {
		ctx.JSON(400, gin.H{})
		return
	}

	ctx.JSON(200, ans)
}

// creating a channel for an user
func (h Handler) createChannel(ctx *gin.Context) {
	var channel models.Channel
	if err := ctx.BindJSON(&channel); err != nil {
		ctx.AbortWithError(401, errors.New("input json can not be marshalled to the channel model"))
		return
	}
	res, _ := ctx.Get("user")
	user := res.(models.User)
	if invalid := h.services.Api.CreateChannel(channel, user); len(invalid) > 0 {
		if err, ok := invalid["error"]; ok {
			ctx.AbortWithError(500, errors.New(err))
		}
		ctx.AbortWithError(422, errors.New("invalid data"))

		return
	}

	ctx.JSON(200, gin.H{})
}

// method for posting a post
func (h Handler) createPost(ctx *gin.Context) {
	var post models.Post
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		ctx.AbortWithError(400, err)
	}
	if err := ctx.BindJSON(&post); err != nil {
		ctx.AbortWithError(401, errors.New("input json can not be marshalled to the post model"))
		return
	}
	if invalid := h.services.Api.CreatePost(post, id); len(invalid) > 0 {
		if err, ok := invalid["error"]; ok {
			ctx.AbortWithError(500, errors.New(err))
		}
		ctx.AbortWithError(422, errors.New("invalid data"))
		return
	}
	ctx.JSON(200, gin.H{})
}

func (h Handler) deletePost(ctx *gin.Context) {
	var post models.Post
	if err := ctx.BindJSON(&post); err != nil {
		ctx.AbortWithError(401, err)
		return
	}
	if err := h.services.Api.DeletePost(post); err != nil {
		ctx.AbortWithError(422, err)
		return
	}
	ctx.JSON(200, gin.H{})
}

// method for reading posts
func (h Handler) getPosts(ctx *gin.Context) {
	res, _ := ctx.Get("user")
	authorType := ctx.DefaultQuery("author", "")
	user := res.(models.User)
	var ans interface{}
	var err error

	// check if needed post should be written by channel, user or all
	if authorType == "channel" {
		ans, err = h.services.Api.GetPostsFromChannels(user)
	} else if authorType == "user" {
		ans, err = h.services.Api.GetPostsFromUsers(user)
	} else {
		ctx.AbortWithError(400, errors.New("author type is not specified."))
	}
	if err != nil {
		ctx.AbortWithError(400, err)
		return
	}
	ctx.JSON(200, ans)
}

func (h Handler) getMyChannelPosts(ctx *gin.Context) {
	res, _ := ctx.Get("user")
	user := res.(models.User)
	ans, err := h.services.Api.GetPostsFromMyChannels(user)
	if err != nil {
		ctx.AbortWithError(400, err)
		return
	}
	ctx.JSON(200, ans)
}

func (h Handler) follow(ctx *gin.Context) {
	res, _ := ctx.Get("user")
	followType := ctx.DefaultQuery("follow", "")
	user := res.(models.User)
	var followed models.User
	if err := ctx.BindJSON(&followed); err != nil {
		ctx.AbortWithError(400, err)
	}
	var err error
	if followType == "channel" {
		err = h.services.FollowChannel(user, followed.Username)
	} else if followType == "user" {
		err = h.services.FollowUser(user, followed.Username)
	}
	if err != nil {
		ctx.AbortWithError(400, err)
	} else {
		ctx.JSON(200, "success")
	}
}

func (h Handler) unfollow(ctx *gin.Context) {
	res, _ := ctx.Get("user")
	followType := ctx.DefaultQuery("follow", "")
	user := res.(models.User)
	var followed models.User
	if err := ctx.BindJSON(&followed); err != nil {
		ctx.AbortWithError(400, err)
	}

	var err error
	if followType == "channel" {
		err = h.services.UnfollowChannel(user, followed.Username)
	} else if followType == "user" {
		err = h.services.UnfollowUser(user, followed.Username)
	}
	if err != nil {
		ctx.AbortWithError(422, err)
	} else {
		ctx.JSON(200, "success")
	}
}

func (h Handler) getNewPosts(ctx *gin.Context) {
	res, _ := ctx.Get("user")
	authorType := ctx.DefaultQuery("author", "")
	user := res.(models.User)
	var ans interface{}
	var err error

	// check if needed post should be written by channel, user or all
	if authorType == "channel" {
		ans, err = h.services.Api.GetNewPostsFromChannels(user)
	} else if authorType == "user" {
		ans, err = h.services.Api.GetNewPostsFromUsers(user)
	} else {
		ctx.AbortWithError(400, errors.New("author type is not specified."))
	}
	if err != nil {
		ctx.AbortWithError(400, err)
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

// function for deleting a channel based on its id
func (h Handler) deleteChannel(ctx *gin.Context) {
	var channel models.Channel
	if err := ctx.BindJSON(&channel); err != nil {
		ctx.AbortWithError(401, err)
		return
	}
	if err := h.services.Api.DeleteChannel(channel); err != nil {
		ctx.AbortWithError(422, err)

	}
	ctx.JSON(200, gin.H{})
}

func (h Handler) updateChannel(ctx *gin.Context) {
	var channel models.Channel
	if err := ctx.BindJSON(&channel); err != nil {
		ctx.AbortWithError(401, errors.New("input json can not be marshalled to the channel model"))
		return
	}

	if err := h.services.Api.UpdateChannel(channel); err != nil {
		ctx.AbortWithError(422, err)
	}

	ctx.JSON(200, gin.H{})
}
