package delivery

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/timofef/technopark_subd_sem_project/models"
	"github.com/timofef/technopark_subd_sem_project/usecase/interfaces"
	"github.com/valyala/fasthttp"
	"net/http"
)

type UserHandler struct {
	userUsecase interfaces.UserUsecase
}

func NewUserHandler(router *fasthttprouter.Router, usecase interfaces.UserUsecase) {
	handler := &UserHandler{userUsecase: usecase}

	router.POST("/api/user/:nickname/create", handler.CreateUser)
	router.GET("/api/user/:nickname/profile", handler.GetUser)
	router.POST("/api/user/:nickname/profile", handler.UpdateUser)
}

func (h * UserHandler) CreateUser(ctx *fasthttp.RequestCtx) {
	var user models.User
	user.UnmarshalJSON(ctx.PostBody())
	nickname := ctx.UserValue("nickname").(string)

	users, err := h.userUsecase.CreateUser(&user, &nickname)

	var response []byte

	switch err {
	case nil:
		ctx.SetStatusCode(http.StatusCreated)
		user.Nickname = ctx.UserValue("nickname").(string)
		response, _ = user.MarshalJSON()
	case models.UserExists:
		ctx.SetStatusCode(http.StatusConflict)
		response, _ = users.MarshalJSON()
	}

	ctx.SetContentType("application/json")
	ctx.Write(response)
}

func (h * UserHandler) GetUser(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)

	user, err := h.userUsecase.GetUser(&nickname)

	var response []byte
	switch err {
	case nil:
		ctx.SetStatusCode(http.StatusOK)
		response, _ = user.MarshalJSON()
	case models.UserNotExists:
		ctx.SetStatusCode(http.StatusNotFound)
		msg := models.Error{Message: err.Error()}
		response, _ = msg.MarshalJSON()
	}

	ctx.SetContentType("application/json")
	ctx.Write(response)
}

func (h * UserHandler) UpdateUser(ctx *fasthttp.RequestCtx) {
	var info models.UserUpdate
	info.UnmarshalJSON(ctx.PostBody())
	nickname := ctx.UserValue("nickname").(string)

	newProfile, err := h.userUsecase.UpdateUser(&info, &nickname)
	var response []byte

	switch err {
	case nil:
		ctx.SetStatusCode(http.StatusOK)
		response, _ = newProfile.MarshalJSON()
	case models.UserConflict:
		ctx.SetStatusCode(http.StatusConflict)
		msg := models.Error{Message: err.Error()}
		response, _ = msg.MarshalJSON()
	case models.UserNotExists:
		ctx.SetStatusCode(http.StatusNotFound)
		msg := models.Error{Message: err.Error()}
		response, _ = msg.MarshalJSON()
	}

	ctx.SetContentType("application/json")
	ctx.Write(response)
}
