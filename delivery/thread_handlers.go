package delivery

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/timofef/technopark_subd_sem_project/models"
	"github.com/timofef/technopark_subd_sem_project/usecase/interfaces"
	"github.com/valyala/fasthttp"
	"net/http"
)

type ThreadHandler struct {
	threadUsecase interfaces.ThreadUsecase
}

func NewThreadHandler(router *fasthttprouter.Router, usecase interfaces.ThreadUsecase) {
	handler := &ThreadHandler{threadUsecase: usecase}

	router.POST("/api/thread/:slug_or_id/create", handler.CreatePosts)
	router.GET("/api/thread/:slug_or_id/details", handler.GetThread)
	router.POST("/api/thread/:slug_or_id/details", handler.UpdateThread)
	router.GET("/api/thread/:slug_or_id/posts", handler.GetPosts)
	router.POST("/api/thread/:slug_or_id/vote", handler.VoteForThread)
}

func (h * ThreadHandler) CreatePosts(ctx *fasthttp.RequestCtx) {
	var posts models.Posts
	posts.UnmarshalJSON(ctx.PostBody())
	slug_or_id := ctx.UserValue("slug_or_id")
	newPosts, err := h.threadUsecase.CreatePosts(slug_or_id, &posts)

	var response []byte

	switch err {
	case nil:
		ctx.SetStatusCode(http.StatusCreated)
		response, _ = newPosts.MarshalJSON()
	case models.ThreadNotExists:
		ctx.SetStatusCode(http.StatusNotFound)
		msg := models.Error{Message: err.Error()}
		response, _ = msg.MarshalJSON()
	case models.ParentNotExists:
		ctx.SetStatusCode(http.StatusConflict)
		msg := models.Error{Message: err.Error()}
		response, _ = msg.MarshalJSON()
	}

	ctx.SetContentType("application/json")
	ctx.Write(response)
}

func (h * ThreadHandler) GetThread(ctx *fasthttp.RequestCtx) {
	slug_or_id := ctx.UserValue("slug_or_id")
	thread, err := h.threadUsecase.GetThread(slug_or_id)

	var response []byte

	switch err {
	case nil:
		ctx.SetStatusCode(http.StatusOK)
		response, _ = thread.MarshalJSON()
	case models.ThreadNotExists:
		ctx.SetStatusCode(http.StatusNotFound)
		msg := models.Error{Message: err.Error()}
		response, _ = msg.MarshalJSON()
	}

	ctx.SetContentType("application/json")
	ctx.Write(response)
}

func (h * ThreadHandler) UpdateThread(ctx *fasthttp.RequestCtx) {
	slug_or_id := ctx.UserValue("slug_or_id")
	var threadUpdate models.ThreadUpdate
	threadUpdate.UnmarshalJSON(ctx.PostBody())
	updatedThread, err := h.threadUsecase.UpdateThread(slug_or_id, &threadUpdate)

	var response []byte

	switch err {
	case nil:
		ctx.SetStatusCode(http.StatusOK)
		response, _ = updatedThread.MarshalJSON()
	case models.ThreadNotExists:
		ctx.SetStatusCode(http.StatusNotFound)
		msg := models.Error{Message: err.Error()}
		response, _ = msg.MarshalJSON()
	}

	ctx.SetContentType("application/json")
	ctx.Write(response)
}

func (h * ThreadHandler) GetPosts(ctx *fasthttp.RequestCtx) {
	slug_or_id := ctx.UserValue("slug_or_id")
	limit := ctx.QueryArgs().Peek("limit")
	since := ctx.QueryArgs().Peek("since")
	sort := ctx.QueryArgs().Peek("sort")
	desc := ctx.QueryArgs().Peek("desc")

	posts, err := h.threadUsecase.GetPosts(slug_or_id, limit, since, sort, desc)

	var response []byte

	switch err {
	case nil:
		ctx.SetStatusCode(http.StatusOK)
		response, _ = posts.MarshalJSON()
	case models.ThreadNotExists:
		ctx.SetStatusCode(http.StatusNotFound)
		msg := models.Error{Message: err.Error()}
		response, _ = msg.MarshalJSON()
	}

	ctx.SetContentType("application/json")
	ctx.Write(response)
}

func (h * ThreadHandler) VoteForThread(ctx *fasthttp.RequestCtx) {
	var voice models.Vote
	voice.UnmarshalJSON(ctx.PostBody())
	slug_or_id := ctx.UserValue("slug_or_id")
	thread, err := h.threadUsecase.VoteForThread(slug_or_id, &voice)

	var response []byte

	//fmt.Println("vote   ", err)
	switch err {
	case nil, models.SameVote:
		ctx.SetStatusCode(http.StatusOK)
		response, _ = thread.MarshalJSON()
	case models.ThreadNotExists, models.UserNotExists:
		ctx.SetStatusCode(http.StatusNotFound)
		msg := models.Error{Message: err.Error()}
		response, _ = msg.MarshalJSON()
	}

	ctx.SetContentType("application/json")
	ctx.Write(response)
}

