package event

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/byuoitav/central-event-system/hub/base"
	"github.com/byuoitav/central-event-system/messenger"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/v2/events"
)

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

	m, err := messenger.BuildMessenger(os.Getenv("HUB_ADDRESS"), base.Messenger, 5000)
	if err != nil {
		return nil, fmt.Errorf("Error while trying to build messenger: %w", err)
	}

	sender.m = m

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
func (s *Sender) SendCardReadErrorEvent() {

	e := events.Event{
		GeneratingSystem: s.sysID,
		Timestamp:        time.Now(),
		Key:              "card-read-error",
		Value:            "true",
		TargetDevice:     s.device,
		AffectedRoom:     s.roomInfo,
		EventTags: []string{
			events.Heartbeat,
		},
	}

	log.L.Debugf("Sending event: %v+", e)

	s.m.SendEvent(e)
}
