package pkg

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/cors"
	"time"
)

func GetCors() app.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS", "PUT", "PATCH"}, // Allowed request methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},          // Allowed request headers
		ExposeHeaders:    []string{"Content-Length"},                                   // Request headers allowed in the upload_file
		AllowCredentials: true,                                                         // Whether cookies are attached
		MaxAge:           12 * time.Hour,                                               // Maximum length of upload_file-side cache preflash requests (seconds)
	})
}
