package delivery

import (
	"encoding/json"
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

	router.POST("/api/forum/create", handler.CreateForum)
	//router.POST("/api/forum/:forum_slug/create", handler.CreateThread)
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
		response, _ = json.Marshal(err)
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

}

func (h * ForumHandler) GetForumDetails(ctx *fasthttp.RequestCtx) {

}

func (h * ForumHandler) GetForumThreads(ctx *fasthttp.RequestCtx) {

}

func (h * ForumHandler) GetForumUsers(ctx *fasthttp.RequestCtx) {

}