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

// CreateSpace 创建一个云空间  method: POST path: /api/workspace
// Request Param: reqtype.SpaceCreateOption
func (c *CloudCodeController) CreateSpace(ctx *gin.Context) *serialize.Response {
	// 1、用户参数获取和验证
	req := c.creationCheck(ctx)
	if req == nil {
		return serialize.Error(http.StatusBadRequest)
	}

	// 2、获取用户id，在token验证时已经解析出并放入ctx中了
	userId := utils.MustGet[uint32](ctx, "id")

	// 3、调用service处理然后响应结果
	space, err := c.spaceService.CreateWorkspace(req, userId)
	switch err {
	case service.ErrNameDuplicate:
		return serialize.Fail(code.SpaceCreateNameDuplicate)
	case service.ErrReachMaxSpaceCount:
		return serialize.Fail(code.SpaceCreateReachMaxCount)
	case service.ErrSpaceCreate:
		return serialize.Fail(code.SpaceCreateFailed)
	case service.ErrReqParamInvalid:
		return serialize.Error(http.StatusBadRequest)
	}

	if err != nil {
		return serialize.Fail(code.SpaceCreateFailed)
	}

	return serialize.OkData(space)
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
		return serialize.Error(http.StatusBadRequest)
	}

	userId := utils.MustGet[uint32](ctx, "id")
	uid := utils.MustGet[string](ctx, "uid")

	space, err := c.spaceService.CreateAndStartWorkspace(req, userId, uid)
	switch err {
	case service.ErrNameDuplicate:
		return serialize.Fail(code.SpaceCreateNameDuplicate)
	case service.ErrReachMaxSpaceCount:
		return serialize.Fail(code.SpaceCreateReachMaxCount)
	case service.ErrSpaceCreate:
		return serialize.Fail(code.SpaceCreateFailed)
	case service.ErrSpaceStart:
		return serialize.Fail(code.SpaceStartFailed)
	case service.ErrOtherSpaceIsRunning:
		return serialize.Fail(code.SpaceOtherSpaceIsRunning)
	case service.ErrReqParamInvalid:
		return serialize.Error(http.StatusBadRequest)
	case service.ErrSpaceAlreadyExist:
		return serialize.Fail(code.SpaceAlreadyExist)
	case service.ErrResourceExhausted:
		return serialize.Fail(code.ResourceExhausted)
	}

	if err != nil {
		return serialize.Fail(code.SpaceCreateFailed)
	}

	return serialize.OkData(space)
}

// StartSpace 启动一个已存在的云空间 method: POST path: /api/workspace/start
// request param: space id
func (c *CloudCodeController) StartSpace(ctx *gin.Context) *serialize.Response {
	var req reqtype.SpaceId
	err := ctx.ShouldBind(&req)
	if err != nil {
		c.logger.Warnf("bind param error:%v", err)
		return serialize.Error(http.StatusBadRequest)
	}

	userId := utils.MustGet[uint32](ctx, "id")
	uid := utils.MustGet[string](ctx, "uid")

	space, err := c.spaceService.StartWorkspace(req.Id, userId, uid)
	switch err {
	case service.ErrWorkSpaceNotExist:
		return serialize.Fail(code.SpaceStartNotExist)
	case service.ErrSpaceStart:
		return serialize.Fail(code.SpaceStartFailed)
	case service.ErrOtherSpaceIsRunning:
		return serialize.Fail(code.SpaceOtherSpaceIsRunning)
	case service.ErrSpaceNotFound:
		return serialize.Fail(code.SpaceNotFound)
	}

	if err != nil {
		return serialize.Fail(code.SpaceStartFailed)
	}

	return serialize.OkData(space)
}

// StopSpace 停止正在运行的云空间 method: PUT path: /api/workspace/stop
// Request Param: sid
func (c *CloudCodeController) StopSpace(ctx *gin.Context) *serialize.Response {
	var req reqtype.SpaceId
	err := ctx.ShouldBind(&req)
	if err != nil {
		c.logger.Warnf("bind param error:%v", err)
		return serialize.Error(http.StatusBadRequest)
	}

	uid := utils.MustGet[string](ctx, "uid")
	userId := utils.MustGet[uint32](ctx, "id")

	err = c.spaceService.StopWorkspace(req.Id, userId, uid)
	if err != nil {
		if err == service.ErrWorkSpaceIsNotRunning {
			return serialize.Ok()
		}

		return serialize.Fail(code.SpaceStopFailed)
	}

	return serialize.Ok()
}

// DeleteSpace 删除已存在的云空间  method: DELETE path: /api/workspace
// Request Param: id
func (c *CloudCodeController) DeleteSpace(ctx *gin.Context) *serialize.Response {
	var req reqtype.SpaceId
	err := ctx.ShouldBind(&req)
	if err != nil {
		c.logger.Warnf("bind param error:%v", err)
		return serialize.Error(http.StatusBadRequest)
	}
	c.logger.Debug("space id:", req.Id)

	// 获取用户id和用户uid
	userId := utils.MustGet[uint32](ctx, "id")
	uid := utils.MustGet[string](ctx, "uid")

	err = c.spaceService.DeleteWorkspace(req.Id, userId, uid)
	if err != nil {
		if err == service.ErrWorkSpaceIsRunning {
			return serialize.Fail(code.SpaceDeleteIsRunning)
		}

		return serialize.Fail(code.SpaceDeleteFailed)
	}

	return serialize.Ok()
}

// ListSpace 获取所有创建的云空间 method: GET path: /api/workspace/list
// Request param: id uid
func (c *CloudCodeController) ListSpace(ctx *gin.Context) *serialize.Response {
	userId := utils.MustGet[uint32](ctx, "id")
	uid := utils.MustGet[string](ctx, "uid")

	spaces, err := c.spaceService.ListWorkspace(userId, uid)
	if err != nil {
		return serialize.Fail(code.QueryFailed)
	}

	return serialize.OkData(spaces)
}

// ModifySpaceName 修改工作空间名称 method: POST path: /api/workspace/name
func (c *CloudCodeController) ModifySpaceName(ctx *gin.Context) *serialize.Response {
	var req struct {
		Name string `json:"name"` // 新的工作空间的名称
		Id   uint32 `json:"id"`   // 工作空间id
	}
	err := ctx.ShouldBind(&req)
	if err != nil {
		c.logger.Warnf("bind req error:%v", err)
		return serialize.Fail(code.SpaceNameModifyFailed)
	}

	userId := utils.MustGet[uint32](ctx, "id")

	err = c.spaceService.ModifyName(req.Name, req.Id, userId)
	switch err {
	case service.ErrNameDuplicate:
		return serialize.Fail(code.SpaceCreateNameDuplicate)
	case nil:
		return serialize.Ok()
	default:
		return serialize.Fail(code.SpaceNameModifyFailed)
	}
}
