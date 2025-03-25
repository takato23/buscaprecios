package middlewares

import (
    "os"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

func CORS(router *gin.Engine) {
    corsConfig := cors.Config{
        AllowOrigins:     []string{"*"}, // Por defecto, permitir todos los orígenes
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
        AllowCredentials: true,
    }

    // Si estamos en producción, restringir los orígenes a WEB_URL (si está definida)
    if os.Getenv("GIN_MODE") == "release" {
        webURL := os.Getenv("WEB_URL")
        if webURL != "" {
            corsConfig.AllowOrigins = []string{webURL}
        }
    }

    router.Use(cors.New(corsConfig))
}