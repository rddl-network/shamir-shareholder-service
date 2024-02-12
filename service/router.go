package service

import "github.com/gin-gonic/gin"

func (ss *ShamirService) configureRouter() {
	ss.router.Use(gin.Logger())
	ss.router.Use(gin.Recovery())
}

func (ss *ShamirService) registerRoutes() {
	ss.router.GET("/mnemonic", ss.getMnemonic)
	ss.router.POST("/mnemonic", ss.postMnemonic)
}

func (ss *ShamirService) getMnemonic(_ *gin.Context) {

}

func (ss *ShamirService) postMnemonic(_ *gin.Context) {

}
