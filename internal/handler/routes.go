// Code generated by goctl. DO NOT EDIT.
// goctl 1.8.1

package handler

import (
	"net/http"

	"shortener/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/:shortUrl",
				Handler: ShowHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/convert",
				Handler: ConvertHandler(serverCtx),
			},
		},
	)
}
