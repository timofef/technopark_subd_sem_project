package delivery

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/timofef/technopark_subd_sem_project/usecase/interfaces"
	"github.com/valyala/fasthttp"
)

type ForumHandler struct {
	forumUsecase interfaces.ForumUsecase
}

func NewForumHandler(router *fasthttprouter.Router, usecase interfaces.ForumUsecase) {
	handler := &ForumHandler{forumUsecase: usecase}

	router.POST("/api/forum/create", handler.CreateForum)
	router.POST("/api/forum/create", handler.CreateThread)
	router.GET("/api/forum/:slug/details", handler.GetForumDetails)
	router.GET("/api/forum/:slug/threads", handler.GetForumThreads)
	router.GET("/api/forum/:slug/users", handler.GetForumUsers)
}

func (h * ForumHandler) CreateForum(ctx *fasthttp.RequestCtx) {

}

func (h * ForumHandler) CreateThread(ctx *fasthttp.RequestCtx) {

}

func (h * ForumHandler) GetForumDetails(ctx *fasthttp.RequestCtx) {

}

func (h * ForumHandler) GetForumThreads(ctx *fasthttp.RequestCtx) {

}

func (h * ForumHandler) GetForumUsers(ctx *fasthttp.RequestCtx) {

}