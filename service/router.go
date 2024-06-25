package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rddl-network/shamir-shareholder-service/config"
	"github.com/rddl-network/shamir-shareholder-service/types"
)

func (ss *ShamirService) configureRouter() {
	ss.router.Use(gin.Logger())
	ss.router.Use(gin.Recovery())
}

func (ss *ShamirService) registerRoutes() {
	ss.router.GET("/mnemonic", ss.getMnemonic)
	ss.router.POST("/mnemonic", ss.postMnemonic)
}

func (ss *ShamirService) getMnemonic(c *gin.Context) {
	cipheredMnemonic, err := ss.GetMnemonic()
	if err != nil {
		ss.logger.Error("msg", "error while fetching mnemonic")
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "error while fetching mnemonic"})
		return
	}

	var resBody types.MnemonicBody
	resBody.Mnemonic, err = ss.DecryptMnemonic(cipheredMnemonic, config.GetConfig().KeyPhrase)
	if err != nil {
		ss.logger.Error("msg", "error while decrypting mnemonic")
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "error while decrypting mnemonic"})
		return
	}

	c.JSON(http.StatusOK, resBody)
}

func (ss *ShamirService) postMnemonic(c *gin.Context) {
	var reqBody types.MnemonicBody
	if err := c.BindJSON(&reqBody); err != nil {
		ss.logger.Error("msg", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	cipheredMnemonic, err := ss.EncryptMnemonic(reqBody.Mnemonic, config.GetConfig().KeyPhrase)
	if err != nil {
		ss.logger.Error("msg", "error while encrypting mnemonic")
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "error while encrypting mnemonic"})
		return
	}

	if err = ss.PutMnemonic(cipheredMnemonic); err != nil {
		ss.logger.Error("msg", "error while storing mnemonic")
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "error while storing mnemonic"})
	}
}
