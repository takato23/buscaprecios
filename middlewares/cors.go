package middlewares

import (
    "os"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

func CORS(router *gin.Engine) {
    corsConfig := cors.Config{
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
        AllowCredentials: true,
    }

    // Si estamos en producción, permitir el dominio de la app y Lovable
    if os.Getenv("GIN_MODE") == "release" {
        webURL := os.Getenv("WEB_URL")
        if webURL != "" {
            corsConfig.AllowOrigins = []string{webURL, "https://lovable.app", "https://preview.lovable.app"}
        } else {
            corsConfig.AllowOrigins = []string{"https://lovable.app", "https://preview.lovable.app"}
        }
    } else {
        // En desarrollo, permitir todo
        corsConfig.AllowOrigins = []string{"*"}
    }

    router.Use(cors.New(corsConfig))
}