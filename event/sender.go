package event

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/byuoitav/central-event-system/hub/base"
	"github.com/byuoitav/central-event-system/messenger"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/v2/events"
)

// Sender represents a messenger that possesses all the information necessary
// to send events to the central event hub
type Sender struct {
	m        *messenger.Messenger
	roomInfo events.BasicRoomInfo
	device   events.BasicDeviceInfo
	sysID    string
}

// NewSender returns a new sender that will send messages to the given messenger
// address as the given system ID
func NewSender(addr, sysID string) (*Sender, error) {

	sender := Sender{}

	m, err := messenger.BuildMessenger(addr, base.Messenger, 5000)
	if err != nil {
		return nil, fmt.Errorf("Error while trying to build messenger: %w", err)
	}

	sender.m = m

	// Setup all the room/device info structs
	sender.sysID = sysID
	a := strings.Split(sysID, "-")
	if len(a) == 3 {
		sender.roomInfo = events.BasicRoomInfo{
			BuildingID: a[0],
			RoomID:     a[0] + "-" + a[1],
		}
	} else {
		sender.roomInfo = events.BasicRoomInfo{
			BuildingID: sysID,
			RoomID:     sysID,
		}
	}

	sender.device = events.BasicDeviceInfo{
		BasicRoomInfo: sender.roomInfo,
		DeviceID:      sysID,
	}

	return &sender, nil

}

// SendCardReadEvent sends a card read event for the given card id
func (s *Sender) SendCardReadEvent(cardID string) {

	e := events.Event{
		GeneratingSystem: s.sysID,
		Timestamp:        time.Now(),
		Key:              "card-read",
		Value:            cardID,
		TargetDevice:     s.device,
		AffectedRoom:     s.roomInfo,
		EventTags: []string{
			events.Heartbeat,
		},
	}

	log.L.Debugf("Sending event: %v+", e)

	s.m.SendEvent(e)

}

// SendCardReadErrorEvent sends a card-read-error event
func (s *Sender) SendCardReadErrorEvent(bits int) {

	e := events.Event{
		GeneratingSystem: s.sysID,
		Timestamp:        time.Now(),
		Key:              "card-read-error",
		Value:            strconv.Itoa(bits),
		TargetDevice:     s.device,
		AffectedRoom:     s.roomInfo,
		EventTags: []string{
			events.Heartbeat,
		},
	}

	log.L.Debugf("Sending event: %v+", e)

	s.m.SendEvent(e)
}

// SendWiringErrorEvent sends a card-reader-wiring-error event
func (s *Sender) SendWiringErrorEvent(bits int) {

	e := events.Event{
		GeneratingSystem: s.sysID,
		Timestamp:        time.Now(),
		Key:              "card-reader-wiring-error",
		Value:            strconv.Itoa(bits),
		TargetDevice:     s.device,
		AffectedRoom:     s.roomInfo,
		EventTags: []string{
			events.Heartbeat,
		},
	}

	log.L.Debugf("Sending event: %v+", e)

	s.m.SendEvent(e)
}
