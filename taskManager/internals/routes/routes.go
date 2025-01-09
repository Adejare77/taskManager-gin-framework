package routes

import (
	"github.com/Adejare77/go/taskManager/internals/controllers"
	"github.com/gin-gonic/gin"
)

var PublicRoutes = func(route *gin.RouterGroup) {
	route.POST("/login", controllers.Login)
	route.POST("/register", controllers.Register)
}

var ProtectedRoutes = func(route *gin.RouterGroup) {
	route.GET("/task", controllers.GetTasks)
	route.GET("/task/:taskID", controllers.GetTasksByID)
	route.POST("/task", controllers.PostTask)
	route.PUT("/task/:id", controllers.UpdateTask)
	route.DELETE("/task/:id", controllers.DeleteTask)

	route.DELETE("/user", controllers.DeleteUser)
}
