package middlewares

import (
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "os"
)

func CORS(router *gin.Engine) {
    corsConfig := cors.Config{
        AllowOrigins: []string{
            "https://preview--cocina-organizada.lovable.app",
            "https://cocina-organizada.lovable.app", // En caso de que la versión final sea esta
        }, // Asegúrate de agregar la URL correcta
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
        AllowCredentials: true,
    }

    router.Use(cors.New(corsConfig))
}