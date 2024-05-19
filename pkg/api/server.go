package api

import (
	"crazydocker/pkg/core"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Run() error {
	// set the ce
	executor, err := core.NewCommandExecutor()
	if err != nil {
		log.Fatal(err)
	}
	ce = executor
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"Origin", "Content-Type"},
		ExposeHeaders: []string{"Content-Length"},
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))
	router.GET("/containers", listContainers)
	router.GET("/images", listImages)
	router.GET("/container/action", performActionOnContainer)
	router.GET("/container/stream", streamContainer)
	router.GET("/image/stream", streamImage)
	router.GET("/config", getConfig)
	router.POST("/container/create", createContainer)
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, map[string]string{"status": "ok"})
	})
	router.GET("/machine/exec", execIntoMachine)
	router.GET("/container/exec", execIntoContainer)
	router.GET("/container/log", streamContainerLogs)
	return router.Run(fmt.Sprintf(":%s", os.Getenv("API_SERVER_PORT")))
}
