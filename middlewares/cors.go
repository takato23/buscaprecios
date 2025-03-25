package middlewares

import (
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

func CORS(router *gin.Engine) {
    corsConfig := cors.Config{
        AllowOrigins:     []string{"*"}, // Permitir todos los orígenes (para testing)
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
        AllowCredentials: true,
        AllowWildcard:    true,
    }

    router.Use(cors.New(corsConfig))
}