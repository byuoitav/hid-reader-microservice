package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/byuoitav/common"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/hid-reader-microservice/event"
	"github.com/byuoitav/hid-reader-microservice/handlers"
	"github.com/byuoitav/hid-reader-microservice/hid"
	"github.com/byuoitav/hid-reader-microservice/wiegand"

	// "github.com/stianeikeland/go-rpio"
	"github.com/stianeikeland/go-rpio/v4"
)

func main() {
	port := ":10023"

	// Setup messenger
	sender, err := event.NewSender(os.Getenv("HUB_ADDRESS"), os.Getenv("SYSTEM_ID"))
	if err != nil {
		log.L.Fatalf("Failed to start the messenger: %s", err)
	}

	// Setup GPIO for other operations
	err = rpio.Open()
	if err != nil {
		fmt.Printf("ERROR OPENING GPIO: %s\n", err)
		return
	}

	defer rpio.Close()

	// Setup card reader
	cardChan := make(chan string, 1)

	reader := wiegand.Reader{
		Data0Pin:   14,
		Data1Pin:   15,
		BufferSize: 52, // Buffering the 48 bits we expect in case of garbage
		Notifier:   cardChan,
	}

	err = reader.Setup()
	if err != nil {
		log.L.Fatalf("Failed to start listening to card reader: %s", err)
	}

	// Listen for card read events and convert them into Card ID's
	go func() {
		for {
			cardBinary := <-cardChan

			numBits := len(cardBinary)

			// Send off bad wiring events for reads of less than 24 bits
			// This is done intentionally to try to separate out bad reads
			// caused by software vs bad reads caused by wiring
			// Note: 24 was picked arbitrarily from experience. Most bad reads
			// due to the non-realtime nature of the pi are typically only a few
			// (1-4) bits off.
			if numBits < 24 {
				log.L.Debugf("Wiring issue, only read %d bits", numBits)
				sender.SendWiringErrorEvent(numBits)
				continue
			}

			log.L.Debugf("Card binary: %s, bits: %d", cardBinary, numBits)
			cardID, err := hid.GetCardID(cardBinary)
			if err != nil {
				log.L.Debugf("Card Read Error: %s", err)
				sender.SendCardReadErrorEvent(numBits)
				continue
			}

			log.L.Debugf("Read Card ID: %s", cardID)
			sender.SendCardReadEvent(cardID)
		}
	}()

	router := common.NewRouter()

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
	router.StartServer(&server)
}
