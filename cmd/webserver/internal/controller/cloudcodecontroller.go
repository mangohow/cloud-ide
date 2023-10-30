package controller

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/mangohow/cloud-ide/cmd/webserver/internal/code"
	"github.com/mangohow/cloud-ide/cmd/webserver/internal/model/reqtype"
	"github.com/mangohow/cloud-ide/cmd/webserver/internal/service"
	"github.com/mangohow/cloud-ide/pkg/logger"
	"github.com/mangohow/cloud-ide/pkg/serialize"
	"github.com/mangohow/cloud-ide/pkg/utils"
	"github.com/sirupsen/logrus"
)

type CloudCodeController struct {
	logger       *logrus.Logger
	spaceService *service.CloudCodeService
}

func NewCloudCodeController() *CloudCodeController {
	return &CloudCodeController{
		logger:       logger.Logger(),
		spaceService: service.NewCloudCodeService(),
	}
}

// CreateSpace 创建一个云空间  method: POST path: /api/space
// Request Param: reqtype.SpaceCreateOption
func (c *CloudCodeController) CreateSpace(ctx *gin.Context) *serialize.Response {
	// 1、用户参数获取和验证
	req := c.creationCheck(ctx)
	if req == nil {
		ctx.Status(http.StatusBadRequest)
		return nil
	}

	// 2、获取用户id，在token验证时已经解析出并放入ctx中了
	idi, _ := ctx.Get("id")
	id := idi.(uint32)

	// 3、调用service处理然后响应结果
	space, err := c.spaceService.CreateWorkspace(req, id)
	switch err {
	case service.ErrNameDuplicate:
		return serialize.Ok(code.SpaceCreateNameDuplicate)
	case service.ErrReachMaxSpaceCount:
		return serialize.Ok(code.SpaceCreateReachMaxCount)
	case service.ErrSpaceCreate:
		return serialize.Ok(code.SpaceCreateFailed)
	case service.ErrReqParamInvalid:
		ctx.Status(http.StatusBadRequest)
		return nil
	}

	if err != nil {
		return serialize.Ok(code.SpaceCreateFailed)
	}

	return serialize.FailWithData(code.SpaceCreateSuccess, space)
}

// creationCheck 用户参数验证
func (c *CloudCodeController) creationCheck(ctx *gin.Context) *reqtype.SpaceCreateOption {
	// 获取用户请求参数
	var req reqtype.SpaceCreateOption
	// 绑定数据
	err := ctx.ShouldBind(&req)
	if err != nil {
		return nil
	}

	c.logger.Debug(req)

	if req.GitRepository != "" {
		matched, err := regexp.MatchString(`^https://\S+.git$`, req.GitRepository)
		if err != nil {
			c.logger.Error("regexp invalid")
			return nil
		}
		if !matched {
			c.logger.Error("git repository invalid")
			return nil
		}
	}

	// 参数验证
	get1, exist1 := ctx.Get("id")
	_, exist2 := ctx.Get("username")
	if !exist1 || !exist2 {
		return nil
	}
	id, ok := get1.(uint32)
	if !ok || id != req.UserId {
		return nil
	}

	return &req
}

// CreateSpaceAndStart 创建一个新的云空间并启动 method: POST path: /api/space_cas
// Request Param: reqtype.SpaceCreateOption
func (c *CloudCodeController) CreateSpaceAndStart(ctx *gin.Context) *serialize.Response {
	req := c.creationCheck(ctx)
	if req == nil {
		ctx.Status(http.StatusBadRequest)
		return nil
	}

	idi, _ := ctx.Get("id")
	id := idi.(uint32)
	uidi, _ := ctx.Get("uid")
	uid := uidi.(string)

	space, err := c.spaceService.CreateAndStartWorkspace(req, id, uid)
	switch err {
	case service.ErrNameDuplicate:
		return serialize.Ok(code.SpaceCreateNameDuplicate)
	case service.ErrReachMaxSpaceCount:
		return serialize.Ok(code.SpaceCreateReachMaxCount)
	case service.ErrSpaceCreate:
		return serialize.Ok(code.SpaceCreateFailed)
	case service.ErrSpaceStart:
		return serialize.Ok(code.SpaceStartFailed)
	case service.ErrOtherSpaceIsRunning:
		return serialize.Ok(code.SpaceOtherSpaceIsRunning)
	case service.ErrReqParamInvalid:
		ctx.Status(http.StatusBadRequest)
		return nil
	case service.ErrSpaceAlreadyExist:
		return serialize.Ok(code.SpaceAlreadyExist)
	case service.ErrResourceExhausted:
		return serialize.Ok(code.ResourceExhausted)
	}

	if err != nil {
		return serialize.Ok(code.SpaceCreateFailed)
	}

	return serialize.FailWithData(code.SpaceStartSuccess, space)
}

// StartSpace 启动一个已存在的云空间 method: POST path: /api/space_start
// request param: space id
func (c *CloudCodeController) StartSpace(ctx *gin.Context) *serialize.Response {
	var req struct {
		Id uint32 `json:"id"`
	}
	err := ctx.ShouldBind(&req)
	if err != nil {
		c.logger.Warnf("bind param error:%v", err)
		ctx.Status(http.StatusBadRequest)
		return nil
	}

	idi, _ := ctx.Get("id")
	id := idi.(uint32)
	uidi, _ := ctx.Get("uid")
	uid := uidi.(string)

	space, err := c.spaceService.StartWorkspace(req.Id, id, uid)
	switch err {
	case service.ErrWorkSpaceNotExist:
		return serialize.Ok(code.SpaceStartNotExist)
	case service.ErrSpaceStart:
		return serialize.Ok(code.SpaceStartFailed)
	case service.ErrOtherSpaceIsRunning:
		return serialize.Ok(code.SpaceOtherSpaceIsRunning)
	case service.ErrSpaceNotFound:
		return serialize.Ok(code.SpaceNotFound)
	}

	if err != nil {
		return serialize.Ok(code.SpaceStartFailed)
	}

	return serialize.FailWithData(code.SpaceStartSuccess, space)
}

// StopSpace 停止正在运行的云空间 method: PUT path: /api/space_stop
// Request Param: sid
func (c *CloudCodeController) StopSpace(ctx *gin.Context) *serialize.Response {
	var req struct {
		Sid string `json:"sid"`
	}
	err := ctx.ShouldBind(&req)
	if err != nil {
		c.logger.Warningf("bind param error:%v", err)
		ctx.Status(http.StatusBadRequest)
		return nil
	}
	uidi, ok := ctx.Get("uid")
	if !ok {
		ctx.Status(http.StatusBadRequest)
		return nil
	}

	uid := uidi.(string)
	err = c.spaceService.StopWorkspace(req.Sid, uid)
	if err != nil {
		if err == service.ErrWorkSpaceIsNotRunning {
			return serialize.Ok(code.SpaceStopIsNotRunning)
		}

		return serialize.Ok(code.SpaceStopFailed)
	}

	return serialize.Ok(code.SpaceStopSuccess)
}

// DeleteSpace 删除已存在的云空间  method: DELETE path: /api/delete
// Request Param: id
func (c *CloudCodeController) DeleteSpace(ctx *gin.Context) *serialize.Response {
	id, err := utils.QueryUint32(ctx, "id")
	if err != nil {
		c.logger.Warningf("get param sid failed:%v", err)
		ctx.Status(http.StatusBadRequest)
		return nil
	}
	c.logger.Debug("id:", id)

	uidi, ok := ctx.Get("uid")
	if !ok {
		ctx.Status(http.StatusBadRequest)
		return nil
	}

	uid := uidi.(string)
	err = c.spaceService.DeleteWorkspace(id, uid)
	if err != nil {
		if err == service.ErrWorkSpaceIsRunning {
			return serialize.Ok(code.SpaceDeleteIsRunning)
		}

		return serialize.Ok(code.SpaceDeleteFailed)
	}

	return serialize.Ok(code.SpaceDeleteSuccess)
}

// ListSpace 获取所有创建的云空间 method: GET path: /api/spaces
// Request param: id uid
func (c *CloudCodeController) ListSpace(ctx *gin.Context) *serialize.Response {
	v1, e1 := ctx.Get("id")
	v2, e2 := ctx.Get("uid")
	if !e1 || !e2 {
		ctx.Status(http.StatusBadRequest)
		return nil
	}
	id := v1.(uint32)
	uid := v2.(string)

	spaces, err := c.spaceService.ListWorkspace(id, uid)
	if err != nil {
		return serialize.Ok(code.QueryFailed)
	}

	return serialize.FailWithData(code.QuerySuccess, spaces)
}

// ModifySpaceName 修改工作空间名称 method: POST path: /api/space_name
func (c *CloudCodeController) ModifySpaceName(ctx *gin.Context) *serialize.Response {
	var req struct {
		Name string `json:"name"` // 新的工作空间的名称
		Id   uint32 `json:"id"`   // 工作空间id
	}
	err := ctx.ShouldBind(&req)
	if err != nil {
		c.logger.Warnf("bind req error:%v", err)
		return serialize.Ok(code.SpaceNameModifyFailed)
	}
	v1, e1 := ctx.Get("id")
	if !e1 {
		ctx.Status(http.StatusBadRequest)
		return nil
	}

	userId := v1.(uint32)
	err = c.spaceService.ModifyName(req.Name, req.Id, userId)
	switch err {
	case service.ErrNameDuplicate:
		return serialize.Ok(code.SpaceCreateNameDuplicate)
	case nil:
		return serialize.Ok(code.SpaceNameModifySuccess)
	default:
		return serialize.Ok(code.SpaceNameModifyFailed)
	}
}
