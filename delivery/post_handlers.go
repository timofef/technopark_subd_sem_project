package delivery

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/timofef/technopark_subd_sem_project/models"
	"github.com/timofef/technopark_subd_sem_project/usecase/interfaces"
	"github.com/valyala/fasthttp"
	"net/http"
)

type PostHandler struct {
	postUsecase interfaces.PostUsecase
}

func NewPostHandler(router *fasthttprouter.Router, usecase interfaces.PostUsecase) {
	handler := &PostHandler{postUsecase: usecase}

	router.GET("/api/post/:id/details", handler.GetPostDetails)
	router.POST("/api/post/:id/details", handler.EditPost)
}

func (h * PostHandler) GetPostDetails(ctx *fasthttp.RequestCtx) {
	id := ctx.UserValue("id").(string)
	related := ctx.QueryArgs().Peek("related")
	postFull, err := h.postUsecase.GetPostDetails(&id, related)

	var response []byte

	switch err {
	case nil:
		ctx.SetStatusCode(http.StatusOK)
		response, _ = postFull.MarshalJSON()
	case models.PostNotExists:
		ctx.SetStatusCode(http.StatusNotFound)
		msg := models.Error{Message: err.Error()}
		response, _ = msg.MarshalJSON()
	}

	ctx.SetContentType("application/json")
	ctx.Write(response)
}

func (h * PostHandler) EditPost(ctx *fasthttp.RequestCtx) {
	id := ctx.UserValue("id").(string)
	update := models.PostUpdate{}
	update.UnmarshalJSON(ctx.PostBody())

	post, err := h.postUsecase.EditPost(&id, &update)

	var response []byte

	switch err {
	case nil:
		ctx.SetStatusCode(http.StatusOK)
		response, _ = post.MarshalJSON()
	case models.PostNotExists:
		ctx.SetStatusCode(http.StatusNotFound)
		msg := models.Error{Message: err.Error()}
		response, _ = msg.MarshalJSON()
	}

	ctx.SetContentType("application/json")
	ctx.Write(response)
}