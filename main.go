package main

import (
    "os"
    "github.com/gin-gonic/gin"
    "ratoneando/config"
    "ratoneando/middlewares"
    "ratoneando/routes"
    "ratoneando/utils/cache"
    "ratoneando/utils/logger"
)

func main() {
    logger.Init()
    config.Init()
    cache.Init() // Inicializa Redis

    // Ajustar el modo de Gin según el entorno
    ginMode := os.Getenv("GIN_MODE")
    if ginMode == "" {
        ginMode = "release" // En producción, usar "release" para mejor rendimiento
    }
    gin.SetMode(ginMode)

    // Obtener el puerto de la variable de entorno (Render lo asigna dinámicamente)
    port := os.Getenv("PORT")
    if port == "" {
        port = "3000" // Usar 3000 si no está definido (para desarrollo local)
    }

    router := gin.Default()

    middlewares.CORS(router)

    // Registrar rutas
    routes.RegisterRoutes(router)

    // Iniciar el servidor en el puerto dinámico de Render
    logger.Log("Starting server on port " + port)
    if err := router.Run(":" + port); err != nil {
        logger.LogFatal("Could not start server: " + err.Error())
    }
}