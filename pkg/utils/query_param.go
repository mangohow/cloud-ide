package utils

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func QueryInt(ctx *gin.Context, key string) (int, error) {
	res := ctx.Query(key)
	if res == "" {
		return 0, errors.New("query failed")
	}

	return strconv.Atoi(res)
}

func ParamUint32(ctx *gin.Context, key string) (uint32, error) {
	id := ctx.Param("id")
	if id == "" {
		return 0, errors.New("query failed")
	}

	if n, err := strconv.Atoi(id); err != nil || n < 0 {
		return 0, err
	} else {
		return uint32(n), nil
	}
}

func QueryUint32(ctx *gin.Context, key string) (uint32, error) {
	res := ctx.Query(key)
	if res == "" {
		return 0, errors.New("query failed")
	}
	if n, err := strconv.Atoi(res); err != nil || n < 0 {
		return 0, err
	} else {
		return uint32(n), nil
	}
}
