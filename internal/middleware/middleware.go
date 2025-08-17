package middleware

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware() gin.HandlerFunc {
	logger := logrus.New()

	file, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("could not open file for logs: " + err.Error())
	}

	logger.SetOutput(file)

	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	return func(c *gin.Context) {
		requestTime := time.Now()

		c.Next()

		logger.WithFields(logrus.Fields{
			"method":       c.Request.Method,
			"url":          c.Request.URL.Path,
			"request_time": requestTime.Format(time.DateOnly),
		}).Info("Request received")
	}
}
