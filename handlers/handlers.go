package handlers

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/byuoitav/central-event-system/hub/base"
	"github.com/byuoitav/central-event-system/messenger"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/common/v2/events"
	"github.com/byuoitav/room-auth-ms/structs"
	"github.com/byuoitav/wso2services/wso2requests"
	"github.com/go-rpio"
	"github.com/labstack/echo"
)

// Beep activates the beep
func Beep(ctx echo.Context) error {
	fmt.Printf("BEEP\n")
	pin := rpio.Pin(5)
	pin.Output()
	pin.Low()

	return nil
}

// BeepOff turns off the beep
func BeepOff(ctx echo.Context) error {
	fmt.Printf("BEEP OFF\n")
	pin := rpio.Pin(5)
	pin.Output()
	pin.High()

	return nil
}

// Green activates a green light
func Green(ctx echo.Context) error {
	fmt.Printf("GREEN\n")
	pin := rpio.Pin(11)
	pin.Output()
	pin.Low()

	return nil
}

// GreenOff deactivates a green light
func GreenOff(ctx echo.Context) error {
	fmt.Printf("GREEN OFF\n")
	pin := rpio.Pin(11)
	pin.Output()
	pin.High()

	return nil
}

// Red activates a red light
func Red(ctx echo.Context) error {
	fmt.Printf("RED\n")
	pin := rpio.Pin(6)
	pin.Output()
	pin.Low()

	return nil
}

// RedOff deactivates a red light
func RedOff(ctx echo.Context) error {
	fmt.Printf("RED OFF\n")
	pin := rpio.Pin(6)
	pin.Output()
	pin.High()

	return nil
}

// Hold .
func Hold(ctx echo.Context) error {
	fmt.Printf("HOLD\n")
	pin := rpio.Pin(13)
	pin.Output()
	pin.Low()

	return nil
}

// HoldOff .
func HoldOff(ctx echo.Context) error {
	fmt.Printf("HOLD OFF\n")
	pin := rpio.Pin(13)
	pin.Output()
	pin.High()

	return nil
}

// StartListen .
func StartListen() {
	go ReadIn()
	//go EdgeDetect()
}

// ReadIn .
func ReadIn() {
	messenger, er := messenger.BuildMessenger("ITB-1101-CP4:10023", base.Messenger, 5000)
	if er != nil {
		log.L.Fatalf("failed to build messenger: %s", er)
	}
	pin0 := rpio.Pin(23)
	pin0.Input()
	pin1 := rpio.Pin(24)
	pin1.Input()

	last0 := rpio.High
	last1 := rpio.High
	lastRead := time.Now()
	printy := false
	bitcount := 0
	var bytes []int64
	for {
		r0 := pin0.Read()
		r1 := pin1.Read()

		if r0 == rpio.Low && r1 == rpio.High && last0 == rpio.High {
			fmt.Printf("0")
			bytes = append(bytes, 0)
			lastRead = time.Now()
			printy = false
			bitcount++
		}

		if r1 == rpio.Low && r0 == rpio.High && last1 == rpio.High {
			fmt.Printf("1")
			bytes = append(bytes, 1)
			lastRead = time.Now()
			printy = false
			bitcount++
		}

		if time.Now().Sub(lastRead).Seconds() >= 1 {
			if !printy {
				fmt.Printf("\nBits: %v\n", bitcount)
				if bitcount != 48 {
					fmt.Println("Bad Read")
				} else {
					var num int64
					for i := len(bytes) - 1; i > 24; i-- {
						num += int64(math.Exp2(float64((len(bytes)-i)-1))) * bytes[i]
					}
					if num%2 == 1 {
						num--
					}
					num /= 2
					netID, err := GetNetID(fmt.Sprintf("%d", num))
					if err != nil {
						fmt.Printf("Ruh Roh!: %v\n", err.Error())
					}
					fmt.Printf("NetID: %s\n", netID)
					SendEvent(netID, *messenger)
				}
				bytes = bytes[:0]
				printy = true
				bitcount = 0
			}

		}
		last0 = r0
		last1 = r1
	}
}

// GetNetID takes the Card Serial Number and uses the Person API to return their info
func GetNetID(cardNumber string) (string, *nerr.E) {
	//this is where we get the NetID

	var output structs.WSO2CredentialResponse
	log.L.Debugf("%s\n", cardNumber)
	err := wso2requests.MakeWSO2Request("GET", "https://api.byu.edu:443/byuapi/persons/v3/?credentials.credential_type=SEOS_CARD&credentials.credential_id="+cardNumber, "", &output)
	if err != nil {
		log.L.Debugf("Error when making WSO2 request %v", err)
		return "", err
	}
	log.L.Debugf("this is the output %v", output)
	NetID := output.Values[0].Basic.NetID.Value
	return NetID, nil
}

// EdgeDetect .
func EdgeDetect() {
	pin0 := rpio.Pin(24)
	pin0.Input()
	pin1 := rpio.Pin(23)
	pin1.Input()

	pin0.PullUp()
	pin1.PullUp()

	pin0.Detect(rpio.FallEdge)
	pin1.Detect(rpio.FallEdge)

	lastRead := time.Now()
	printy := false
	bitcount := 0
	for {

		if pin0.EdgeDetected() && pin1.Read() == rpio.High {
			fmt.Printf("0")
			lastRead = time.Now()
			printy = false
			bitcount++
		}

		if pin1.EdgeDetected() && pin0.Read() == rpio.High {
			fmt.Printf("1")
			lastRead = time.Now()
			printy = false
			bitcount++
		}

		if !printy {
			if time.Now().Sub(lastRead).Seconds() >= 1 {
				fmt.Printf("\nBits: %v\n", bitcount)
				printy = true
				bitcount = 0

			}
		}

	}
}

// SendEvent sends an event
func SendEvent(netid string, runner messenger.Messenger) {

	room := os.Getenv("SYSTEM_ID")
	a := strings.Split(room, "-")
	roominfo := events.BasicRoomInfo{}
	if len(a) == 3 {
		roominfo = events.BasicRoomInfo{
			BuildingID: a[0],
			RoomID:     a[0] + "-" + a[1],
		}
	} else {
		roominfo = events.BasicRoomInfo{
			BuildingID: room,
			RoomID:     room,
		}
	}

	basicdevice := events.BasicDeviceInfo{
		BasicRoomInfo: roominfo,
		DeviceID:      os.Getenv("SYSTEM_ID"),
	}

	Event := events.Event{
		GeneratingSystem: os.Getenv("SYSTEM_ID"),
		Timestamp:        time.Now(),
		Key:              "Login",
		Value:            "True",
		User:             netid,
		TargetDevice:     basicdevice,
		AffectedRoom:     roominfo,
		EventTags: []string{
			events.Heartbeat,
		},
	}

	runner.SendEvent(Event)

}
