package broadcaster

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/pjmessi/udp_chat/src/network"
)

const BROADCAST_ADDRESS = "255.255.255.255:34354"

const BROADCAST_INTERVAL_IN_SEC time.Duration = 2 * time.Second

func Broadcast(ctx context.Context, username string) error {
	conn, err := net.Dial("udp", BROADCAST_ADDRESS)
	if err != nil {
		return fmt.Errorf("error dialing broadcast address")
	}
	defer conn.Close()

	request := network.UdpMessage{
		MessageType: network.UdpMessageTypeBroadcast,
		Username:    &username,
	}
	requestBuf, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("error while marshalling udp message: %s", err)
	}

	for {
		if _, err := conn.Write(requestBuf); err != nil {
			return fmt.Errorf("error broadcasting username: %s", err)
		}

		time.Sleep(BROADCAST_INTERVAL_IN_SEC)
	}
}
