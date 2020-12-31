package delivery

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/timofef/technopark_subd_sem_project/models"
	"github.com/timofef/technopark_subd_sem_project/usecase/interfaces"
	"github.com/valyala/fasthttp"
	"net/http"
)

type ForumHandler struct {
	forumUsecase interfaces.ForumUsecase
}

func NewForumHandler(router *fasthttprouter.Router, usecase interfaces.ForumUsecase) {
	handler := &ForumHandler{forumUsecase: usecase}

	// Not using router.POST("/api/forum/create", handler.CreateForum)
	// because fasthttp will have conflict in routing

	router.POST("/api/forum/:slug", handler.CreateForum)
	router.POST("/api/forum/:slug/create", handler.CreateThread)
	router.GET("/api/forum/:slug/details", handler.GetForumDetails)
	router.GET("/api/forum/:slug/threads", handler.GetForumThreads)
	router.GET("/api/forum/:slug/users", handler.GetForumUsers)
}

func (h * ForumHandler) CreateForum(ctx *fasthttp.RequestCtx) {
	var forum models.Forum
	forum.UnmarshalJSON(ctx.PostBody())

	newForum, err := h.forumUsecase.CreateForum(&forum)

	var response []byte

	switch err {
	case nil:
		ctx.SetStatusCode(http.StatusCreated)
		forum.Slug = newForum.Slug
		forum.Title = newForum.Title
		forum.User = newForum.User
		response, _ = forum.MarshalJSON()
	case models.UserNotExists:
		ctx.SetStatusCode(http.StatusNotFound)
		msg := models.Error{Message: err.Error()}
		response, _ = msg.MarshalJSON()
	case models.ForumExists:
		ctx.SetStatusCode(http.StatusConflict)
		var existingForum models.Forum
		existingForum.Slug = newForum.Slug
		existingForum.Title = newForum.Title
		existingForum.User = newForum.User
		response, _ = existingForum.MarshalJSON()
	}

	ctx.SetContentType("application/json")
	ctx.Write(response)
}

func (h * ForumHandler) CreateThread(ctx *fasthttp.RequestCtx) {
	var thread models.Thread
	thread.UnmarshalJSON(ctx.PostBody())
	slug := ctx.UserValue("slug")
	thread.Forum = slug.(string)

	newThread, err := h.forumUsecase.CreateThread(&thread)

	var response []byte

	switch err {
	case nil:
		ctx.SetStatusCode(http.StatusCreated)
		response, _ = newThread.MarshalJSON()
	case models.UserNotExists, models.ForumNotExists:
		ctx.SetStatusCode(http.StatusNotFound)
		msg := models.Error{Message: err.Error()}
		response, _ = msg.MarshalJSON()
	case models.ThreadExists:
		ctx.SetStatusCode(http.StatusConflict)
		response, _ = newThread.MarshalJSON()
	}

	ctx.SetContentType("application/json")
	ctx.Write(response)
}

func (h * ForumHandler) GetForumDetails(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	forum, err := h.forumUsecase.GetForumDetails(&slug)

	var response []byte

	switch err {
	case nil:
		ctx.SetStatusCode(http.StatusOK)
		response, _ = forum.MarshalJSON()
	case models.ForumNotExists:
		ctx.SetStatusCode(http.StatusNotFound)
		msg := models.Error{Message: err.Error()}
		response, _ = msg.MarshalJSON()
	}

	ctx.SetContentType("application/json")
	ctx.Write(response)
}

func (h * ForumHandler) GetForumThreads(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	since := ctx.QueryArgs().Peek("since")
	desc := ctx.QueryArgs().Peek("desc")
	limit := ctx.QueryArgs().Peek("limit")

	threads, err := h.forumUsecase.GetForumThreads(&slug, since, desc, limit)

	var response []byte

	switch err {
	case nil:
		ctx.SetStatusCode(http.StatusOK)
		response, _ = threads.MarshalJSON()
	case models.ForumNotExists:
		ctx.SetStatusCode(http.StatusNotFound)
		msg := models.Error{Message: err.Error()}
		response, _ = msg.MarshalJSON()
	}

	ctx.SetContentType("application/json")
	ctx.Write(response)
}

func (h * ForumHandler) GetForumUsers(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	since := ctx.QueryArgs().Peek("since")
	desc := ctx.QueryArgs().Peek("desc")
	limit := ctx.QueryArgs().Peek("limit")

	users, err := h.forumUsecase.GetForumUsers(&slug, since, desc, limit)

	var response []byte

	switch err {
	case nil:
		ctx.SetStatusCode(http.StatusOK)
		response, _ = users.MarshalJSON()
	case models.ForumNotExists:
		ctx.SetStatusCode(http.StatusNotFound)
		msg := models.Error{Message: err.Error()}
		response, _ = msg.MarshalJSON()
	}

	ctx.SetContentType("application/json")
	ctx.Write(response)
}