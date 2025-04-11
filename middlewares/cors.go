package middlewares

import (
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

func CORS(router *gin.Engine) {
    corsConfig := cors.Config{
    // AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174"}, // Comentado temporalmente
    AllowAllOrigins:  true, // <-- Añadir esto temporalmente
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
    AllowCredentials: true,
    // AllowWildcard:    true, // Esto podría ser redundante con AllowAllOrigins
}

    router.Use(cors.New(corsConfig))
}
