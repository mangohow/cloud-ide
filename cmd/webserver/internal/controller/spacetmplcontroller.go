package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mangohow/cloud-ide/cmd/webserver/internal/code"
	"github.com/mangohow/cloud-ide/cmd/webserver/internal/service"
	"github.com/mangohow/cloud-ide/pkg/logger"
	"github.com/mangohow/cloud-ide/pkg/serialize"
	"github.com/sirupsen/logrus"
)

type SpaceTmplController struct {
	logger  *logrus.Logger
	service *service.SpaceTmplService
}

func NewSpaceTmplController() *SpaceTmplController {
	return &SpaceTmplController{
		logger:  logger.Logger(),
		service: service.NewSpaceTmplService(),
	}
}

// SpaceTmpls 获取所有模板 method: GET path:/api/template/list
func (s *SpaceTmplController) SpaceTmpls(ctx *gin.Context) *serialize.Response {
	tmpls, kinds, err := s.service.GetAllUsingTmpl()
	if err != nil {
		s.logger.Warnf("get tmpls err:%v", err)
		return serialize.Fail(code.QueryFailed)
	}

	return serialize.OkData(gin.H{
		"tmpls": tmpls,
		"kinds": kinds,
	})
}

// SpaceSpecs 获取空间规格 method: GET path:/api/spec/list
func (s *SpaceTmplController) SpaceSpecs(ctx *gin.Context) *serialize.Response {
	specs, err := s.service.GetAllSpec()
	if err != nil {
		s.logger.Warnf("get specs error:%v", err)
		return serialize.Fail(code.QueryFailed)
	}

	return serialize.OkData(specs)
}
