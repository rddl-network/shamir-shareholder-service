package service

import "github.com/gin-gonic/gin"

type ShamirService struct {
	router *gin.Engine
}

func NewShamirService(router *gin.Engine) *ShamirService {
	service := &ShamirService{router: router}
	service.configureRouter()
	service.registerRoutes()
	return service
}

func (ss *ShamirService) Run() (err error) {
	return ss.router.Run()
}
