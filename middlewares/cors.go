package middlewares

import (
    "os"
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

    // En producción, permitir solo la URL de Lovable
    if os.Getenv("GIN_MODE") == "release" {
        webURL := os.Getenv("WEB_URL")
        if webURL != "" {
            corsConfig.AllowOrigins = []string{webURL}
        }
    }

    router.Use(cors.New(corsConfig))
}