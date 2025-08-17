package api

import (
	"github.com/Komilov31/calendar-service/internal/handler"
	"github.com/Komilov31/calendar-service/internal/middleware"
	"github.com/Komilov31/calendar-service/internal/repository"
	"github.com/Komilov31/calendar-service/internal/service"
	"github.com/gin-gonic/gin"
)

type APIServer struct {
	addr string
}

func NewServer(addr string) *APIServer {
	return &APIServer{addr: addr}
}

func (s *APIServer) Run() error {
	router := gin.Default()
	router.Use(middleware.LoggingMiddleware()) // навесили всем хэндлерам middleware для логирования

	repository := repository.New()
	service := service.New(repository)
	handler := handler.New(service)

	router.POST("/create_event", handler.CreateEvent)
	router.POST("/update_event", handler.UpdateEvent)
	router.POST("/delete_event", handler.DeleteEvent)
	router.GET("/events_for_day", handler.GetEventsForDay)
	router.GET("/events_for_week", handler.GetEventsForWeek)
	router.GET("/events_for_month", handler.GetEventsForMonth)

	return router.Run(s.addr)
}
