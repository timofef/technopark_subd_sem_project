package delivery

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/timofef/technopark_subd_sem_project/usecase/interfaces"
	"github.com/valyala/fasthttp"
	"net/http"
)

type ServiceHandler struct {
	serviceUsecase interfaces.ServiceUsecase
}

func NewServiceHandler(router *fasthttprouter.Router, usecase interfaces.ServiceUsecase) {
	handler := &ServiceHandler{serviceUsecase: usecase}

	router.POST("/api/service/clear", handler.ClearBase)
	router.GET("/api/service/status", handler.GetStatus)
}

func (h *ServiceHandler) GetStatus(ctx *fasthttp.RequestCtx) {
	info, _ := h.serviceUsecase.GetStatus()

	var response []byte

	ctx.SetStatusCode(http.StatusOK)
	response, _ = info.MarshalJSON()

	ctx.SetContentType("application/json")
	ctx.Write(response)
}

func (h *ServiceHandler) ClearBase(ctx *fasthttp.RequestCtx) {
	_ = h.serviceUsecase.Clear()

	var response []byte

	ctx.SetStatusCode(http.StatusOK)

	ctx.SetContentType("application/json")
	ctx.Write(response)
}
