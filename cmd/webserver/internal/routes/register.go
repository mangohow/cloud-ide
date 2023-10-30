package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mangohow/cloud-ide/cmd/webserver/internal/controller"
	"github.com/mangohow/cloud-ide/cmd/webserver/internal/middleware"
	"github.com/mangohow/cloud-ide/pkg/router"
)

func Register(engine *gin.Engine) {
	authGroup := engine.Group("/auth")
	userController := controller.NewUserController()
	{
		authGroup.POST("/login", router.Decorate(userController.Login))
		authGroup.GET("/username/check", router.Decorate(userController.CheckUsernameAvailable))
		authGroup.POST("/register", router.Decorate(userController.Register))
		authGroup.GET("/emailCode", router.Decorate(userController.GetEmailValidateCode))
	}

	apiGroup := engine.Group("/api", middleware.Auth())
	tmplController := controller.NewSpaceTmplController()
	{
		apiGroup.GET("/template/list", router.Decorate(tmplController.SpaceTmpls))
		apiGroup.GET("/spec/list", router.Decorate(tmplController.SpaceSpecs))
	}

	spaceController := controller.NewCloudCodeController()
	{
		apiGroup.GET("/workspace/list", router.Decorate(spaceController.ListSpace))
		apiGroup.DELETE("/workspace", router.Decorate(spaceController.DeleteSpace))
		apiGroup.POST("/workspace", router.Decorate(spaceController.CreateSpace))
		apiGroup.POST("/workspace/cas", router.Decorate(spaceController.CreateSpaceAndStart))
		apiGroup.PUT("/workspace/start", router.Decorate(spaceController.StartSpace))
		apiGroup.PUT("/workspace/stop", router.Decorate(spaceController.StopSpace))
		apiGroup.PUT("/workspace/name", router.Decorate(spaceController.ModifySpaceName))
	}
}
