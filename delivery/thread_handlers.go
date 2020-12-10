package delivery

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/timofef/technopark_subd_sem_project/usecase/interfaces"
	"github.com/valyala/fasthttp"
)

type TreadHandler struct {
	threadUsecase interfaces.ThreadUsecase
}

func NewTreadHandler(router *fasthttprouter.Router, usecase interfaces.ThreadUsecase) {
	handler := &TreadHandler{threadUsecase: usecase}

	router.POST("/api/thread/:slug_or_id/create", handler.CreatePost)
	router.GET("/api/thread/:slug_or_id/details", handler.GetThread)
	router.POST("/api/thread/:slug_or_id/details", handler.UpdateThread)
	router.GET("/api/thread/:slug_or_id/posts", handler.GetPosts)
	router.POST("/api/thread/:slug_or_id/vote", handler.VoteForThread)
}

func (h * TreadHandler) CreatePost(ctx *fasthttp.RequestCtx) {

}

func (h * TreadHandler) GetThread(ctx *fasthttp.RequestCtx) {

}

func (h * TreadHandler) UpdateThread(ctx *fasthttp.RequestCtx) {

}

func (h * TreadHandler) GetPosts(ctx *fasthttp.RequestCtx) {

}

func (h * TreadHandler) VoteForThread(ctx *fasthttp.RequestCtx) {

}

