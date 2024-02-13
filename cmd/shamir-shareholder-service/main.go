package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rddl-network/shamir-shareholder-service/service"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	db, err := service.InitDB("./data")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	service := service.NewShamirService(router, db)
	err = service.Run()
	if err != nil {
		log.Fatalf("fatal error spinning up service: %s", err)
	}
}
