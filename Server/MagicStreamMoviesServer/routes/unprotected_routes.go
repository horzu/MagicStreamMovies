package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/horzu/MagicStreamMovies/Server/MagicStreamMoviesServer/controllers"
)

func SetupUnprotectedRoutes(router *gin.Engine) {
	router.POST("/register", controller.RegisterUser())
	router.POST("/login", controller.LoginUser())
	router.GET("/movies", controller.GetMovies())
}
