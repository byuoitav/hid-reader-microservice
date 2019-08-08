package main

import (
	"fmt"
	"net/http"

	"github.com/byuoitav/common"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/hid-reader-microservice/handlers"
	"github.com/stianeikeland/go-rpio"
)

func main() {
	port := ":10023"

	err := rpio.Open()
	if err != nil {
		fmt.Printf("ERROR OPENING GPIO: %s\n", err)
		return
	}

	defer rpio.Close()

	router := common.NewRouter()
	log.SetLevel("debug")

	router.POST("/beep", handlers.Beep)
	router.POST("/beepoff", handlers.BeepOff)
	router.POST("/green", handlers.Green)
	router.POST("/red", handlers.Red)
	router.POST("/greenoff", handlers.GreenOff)
	router.POST("/redoff", handlers.RedOff)
	router.POST("/hold", handlers.Hold)
	router.POST("/holdoff", handlers.HoldOff)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}
	handlers.StartListen()
	router.StartServer(&server)
}
