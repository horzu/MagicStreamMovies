package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/horzu/MagicStreamMovies/Server/MagicStreamMoviesServer/controllers"
	middleware "github.com/horzu/MagicStreamMovies/Server/MagicStreamMoviesServer/middleware"
)

func SetupProtectedRoutes(router *gin.Engine) {
	router.Use(middleware.AuthMiddleware())

	router.GET("/movie/:imdb_id", controller.GetMovie())
	router.POST("/addmovie", controller.AddMovie())
	router.GET("/recommendedmovies", controller.GetRecommendedMovies())
	router.PATCH("/updatereview/:imdb_id", controller.AdminReviewUpdate())
}
