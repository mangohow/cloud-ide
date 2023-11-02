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
		authGroup.POST("/login", router.HandlerAdapter(userController.Login))
		authGroup.GET("/username/check", router.HandlerAdapter(userController.CheckUsernameAvailable))
		authGroup.POST("/register", router.HandlerAdapter(userController.Register))
		authGroup.GET("/emailCode", router.HandlerAdapter(userController.GetEmailValidateCode))
	}

	apiGroup := engine.Group("/api", middleware.Auth())
	tmplController := controller.NewSpaceTmplController()
	{
		apiGroup.GET("/template/list", router.HandlerAdapter(tmplController.SpaceTmpls))
		apiGroup.GET("/spec/list", router.HandlerAdapter(tmplController.SpaceSpecs))
	}

	spaceController := controller.NewCloudCodeController()
	{
		apiGroup.GET("/workspace/list", router.HandlerAdapter(spaceController.ListSpace))
		apiGroup.DELETE("/workspace", router.HandlerAdapter(spaceController.DeleteSpace))
		apiGroup.POST("/workspace", router.HandlerAdapter(spaceController.CreateSpace))
		apiGroup.POST("/workspace/cas", router.HandlerAdapter(spaceController.CreateSpaceAndStart))
		apiGroup.PUT("/workspace/start", router.HandlerAdapter(spaceController.StartSpace))
		apiGroup.PUT("/workspace/stop", router.HandlerAdapter(spaceController.StopSpace))
		apiGroup.PUT("/workspace/name", router.HandlerAdapter(spaceController.ModifySpaceName))
	}
}
