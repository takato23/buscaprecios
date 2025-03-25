package main

import (
    "os"
    "github.com/gin-gonic/gin"
    "ratoneando/config"
    "ratoneando/middlewares"
    "ratoneando/routes"
    "ratoneando/utils/logger"
)

func main() {
    logger.Init()
    config.Init()
    // cache.Init() // Comentado porque cache.go está vacío y no usamos Redis

    // Ajustar el modo de Gin según el entorno
    ginMode := os.Getenv("GIN_MODE")
    if ginMode == "" {
        ginMode = "debug" // Por defecto en local
    }
    gin.SetMode(ginMode)

    // Leer el puerto desde la variable de entorno PORT
    port := os.Getenv("PORT")
    if port == "" {
        port = "3000" // Valor por defecto para local
    }

    router := gin.Default()

    middlewares.CORS(router)

    // Register routes
    routes.RegisterRoutes(router)

    // Start the server
    logger.Log("Starting server on port " + port)
    if err := router.Run(":" + port); err != nil {
        logger.LogFatal("Could not start server: " + err.Error())
    }
}