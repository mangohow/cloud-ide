package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mangohow/cloud-ide/pkg/serialize"
)

type Handler func(ctx *gin.Context) *serialize.Response

func Decorate(h Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		r := h(ctx)
		if r != nil {
			ctx.JSON(r.HttpStatus, &r.R)
		}

		serialize.PutResponse(r)
	}
}
