package routes

import (
	"final-project-golang/controllers"
	"final-project-golang/database"
	"final-project-golang/middlewares"

	"github.com/gin-gonic/gin"
)

func StartApp() *gin.Engine {
	router := gin.Default()

	db := database.ConnectDB()
	userController := controllers.NewUserController(db)
	photoController := controllers.NewPhotoController(db)
	commentController := controllers.NewCommentController(db)
	socialController := controllers.NewSocialController(db)

	userGroup := router.Group("/users")
	{
		userGroup.POST("/register", userController.Register)
		userGroup.POST("/login", userController.Login)
		userGroup.PUT("/", middlewares.Auth(), userController.Update)
		userGroup.DELETE("/", middlewares.Auth(), userController.Delete)
	}

	photoGroup := router.Group("/photos")
	{
		photoGroup.POST("/", middlewares.Auth(), photoController.Create)
		photoGroup.GET("/", middlewares.Auth(), photoController.Get)
		photoGroup.PUT("/:photoId", middlewares.Auth(), photoController.Update)
		photoGroup.DELETE("/:photoId", middlewares.Auth(), photoController.Delete)
	}

	commentGroup := router.Group("/comments")
	{
		commentGroup.POST("/", middlewares.Auth(), commentController.Create)
		commentGroup.GET("/", middlewares.Auth(), commentController.Get)
		commentGroup.PUT("/:commentId", middlewares.Auth(), commentController.Update)
		commentGroup.DELETE("/:commentId", middlewares.Auth(), commentController.Delete)
	}

	socialGroup := router.Group("/socialmedias")
	{
		socialGroup.POST("/", middlewares.Auth(), socialController.Create)
		socialGroup.GET("/", middlewares.Auth(), socialController.Get)
		socialGroup.PUT("/:socialMediaId", middlewares.Auth(), socialController.Update)
		socialGroup.DELETE("/:socialMediaId", middlewares.Auth(), socialController.Delete)
	}

	return router

}
