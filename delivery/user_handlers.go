package delivery

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/timofef/technopark_subd_sem_project/usecase/interfaces"
	"github.com/valyala/fasthttp"
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

}

func (h * UserHandler) GetUser(ctx *fasthttp.RequestCtx) {

}

func (h * UserHandler) UpdateUser(ctx *fasthttp.RequestCtx) {

}
