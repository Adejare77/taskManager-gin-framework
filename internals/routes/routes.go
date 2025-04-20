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
	route.GET("/tasks", controllers.GetTasks)
	route.POST("/tasks", controllers.PostTask)
	route.GET("/tasks/:task_id", controllers.GetTasksByID)
	route.PATCH("/tasks/:task_id", controllers.UpdateTask)
	route.DELETE("/tasks/:task_id", controllers.DeleteTask)

	route.GET("/user/logout", controllers.Logout)
	route.DELETE("/user", controllers.DeleteUser)
}
