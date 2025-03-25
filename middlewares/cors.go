package middlewares

import (
    "os"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "time"
)

func CORS(router *gin.Engine) {
    corsConfig := cors.Config{
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour, // Cache de CORS
    }

    // En producción, permitir WEB_URL y Lovable
    if os.Getenv("GIN_MODE") == "release" {
        webURL := os.Getenv("WEB_URL")
        corsConfig.AllowOrigins = []string{"https://lovable.app", "https://preview.lovable.app"}
        if webURL != "" {
            corsConfig.AllowOrigins = append(corsConfig.AllowOrigins, webURL)
        }
    } else {
        // En desarrollo, permitir todo
        corsConfig.AllowOrigins = []string{"*"}
    }

    router.Use(cors.New(corsConfig))
}