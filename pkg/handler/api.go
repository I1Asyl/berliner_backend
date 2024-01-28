package handler

// Handler for apis

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/I1Asyl/ginBerliner/models"
	"github.com/gin-gonic/gin"
)

// method for getting user pseudonyms
func (h Handler) getPseudonyms(ctx *gin.Context) {
	res, _ := ctx.Get("user")
	user := res.(models.User)
	ans, err := h.services.Api.GetPseudonyms(user)

	if err != nil {
		ctx.JSON(400, gin.H{})
		return
	}

	ctx.JSON(200, ans)
}

// creating a pseudonym for an user
func (h Handler) createPseudonym(ctx *gin.Context) {
	var pseudonym models.Pseudonym
	if err := ctx.BindJSON(&pseudonym); err != nil {
		ctx.AbortWithError(401, errors.New("input json can not be marshalled to the pseudonym model"))
		return
	}
	res, _ := ctx.Get("user")
	user := res.(models.User)
	if invalid := h.services.Api.CreatePseudonym(pseudonym, user); len(invalid) > 0 {
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
		ctx.AbortWithError(400, errors.New("invalid id"))
	}
	if err := ctx.BindJSON(&post); err != nil {
		fmt.Println(post, err)
		var js interface{}
		err = ctx.BindJSON(&js)
		fmt.Println(js, err)
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

// func (h Handler) deletePost(ctx *gin.Context) {
// 	var post models.Post
// 	id, err := strconv.Atoi(ctx.Query("id"))
// 	if err != nil {
// 		ctx.AbortWithError(400, errors.New("invalid id"))
// 	}

// 	if err := h.services.Api.DeletePost(id); err != nil {
// 		ctx.AbortWithError(422, err)
// 		return
// 	}
// 	ctx.JSON(200, gin.H{"success"})
// }

// method for reading posts
func (h Handler) getPosts(ctx *gin.Context) {
	res, _ := ctx.Get("user")
	authorType := ctx.DefaultQuery("author", "")
	user := res.(models.User)
	var ans interface{}
	var err error

	// check if needed post should be written by pseudonym, user or all
	if authorType == "pseudonym" {
		ans, err = h.services.Api.GetPostsFromPseudonyms(user)
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

func (h Handler) follow(ctx *gin.Context) {
	res, _ := ctx.Get("user")
	followType := ctx.DefaultQuery("follow", "")
	user := res.(models.User)
	var followed models.User
	if err := ctx.BindJSON(&followed); err != nil {
		ctx.AbortWithError(400, err)
	}
	var err error
	if followType == "pseudonym" {
		err = h.services.FollowPseudonym(user, followed.Username)
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
	if followType == "pseudonym" {
		err = h.services.UnfollowPseudonym(user, followed.Username)
	} else if followType == "user" {
		err = h.services.UnfollowUser(user, followed.Username)
	}
	if err != nil {
		ctx.AbortWithError(400, err)
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

	// check if needed post should be written by pseudonym, user or all
	if authorType == "pseudonym" {
		ans, err = h.services.Api.GetNewPostsFromPseudonyms(user)
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

// function for deleting a pseudonym based on its id
func (h Handler) deletePseudonym(ctx *gin.Context) {
	var pseudonym models.Pseudonym
	if err := ctx.BindJSON(&pseudonym); err != nil {
		ctx.AbortWithError(401, errors.New("input json can not be marshalled to the pseudonym model"))
		return
	}
	if err := h.services.Api.DeletePseudonym(pseudonym); err != nil {
		ctx.AbortWithError(422, errors.New("invalid data"))

	}
	ctx.JSON(200, gin.H{})
}

func (h Handler) updatePseudonym(ctx *gin.Context) {
	var pseudonym models.Pseudonym
	if err := ctx.BindJSON(&pseudonym); err != nil {
		ctx.AbortWithError(401, errors.New("input json can not be marshalled to the pseudonym model"))
		return
	}

	if err := h.services.Api.UpdatePseudonym(pseudonym); err != nil {
		ctx.AbortWithError(422, errors.New("invalid data"))
	}

	ctx.JSON(200, gin.H{})
}
