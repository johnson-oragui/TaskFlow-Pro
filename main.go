package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/johnson-oragui/TaskFlow-Pro/api/config"
	"github.com/johnson-oragui/TaskFlow-Pro/api/database"
	"github.com/johnson-oragui/TaskFlow-Pro/api/middleware"
	"github.com/johnson-oragui/TaskFlow-Pro/api/routers"
	"github.com/johnson-oragui/TaskFlow-Pro/api/utils"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, continuing...")

	}

}

func main() {
	// ---------- LOAD CONFIG -------------------
	config := config.Load()

	// ---------- DATABASE ---------------------
	db, err := database.ConnectDatabase(config.DBURL)
	if err != nil {
		log.Println(err.Error())
	}

	// --------- REDIS INIT ---------------
	database.InitRedis(config.RedisUrl)

	// ------------- INIT GIN ---------------------
	gin.SetMode(gin.ReleaseMode)
	ginEngine := gin.New()

	// ---------- CORS Protection -------------
	ginEngine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "PUT", "POST", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "X-Refresh-Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	ginEngine.SetTrustedProxies([]string{config.AppUrl})

	// -----------Panic Recovery Middleware ------------------
	ginEngine.Use(gin.Recovery())

	// ------------ ROUTER MIDDLEWARE ---------------------
	ginEngine.Use(middleware.RouteLoggerMiddleware())

	// ------------- PATH MIDDLEWARE ---------------
	ginEngine.Use(middleware.PathMiddleware())

	// -------------- RATE LIMITER --------------------
	ginEngine.Use(middleware.RateLimiterMiddleware(database.RedisClient, 5, time.Minute, "global"))

	// ------------ ROUTES -------------------

	routers.DefaultRouters(ginEngine)
	routers.RouterV1(ginEngine)

	// -------------- GRACEFUL SHUTDOWN -------------------
	utils.StartGinServer(ginEngine, config.Port, func() {
		// ---------- CLOSE REDIS -----------------
		database.CloseRedis()

		// ------------ CLOSE DATABASE --------------------
		if db != nil {
			sqlDB, err := db.DB()
			if err != nil {
				log.Println("Failed to get sql.DB from gorm DB:", err)
			} else {
				if err := sqlDB.Close(); err != nil {
					log.Println("Error closing DB:", err)
				}
			}
		}
	})

}
