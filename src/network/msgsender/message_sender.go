package msgsender

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	"github.com/pjmessi/udp_chat/src/network"
	"github.com/pjmessi/udp_chat/src/state"
)

func SendMessage(ctx context.Context, selfUsername string, user *state.User, messageContent string) error {
	addr := fmt.Sprintf("%s:34354", user.Address.IP)
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return fmt.Errorf("error dialing broadcast address")
	}
	defer conn.Close()

	request := network.UdpMessage{
		Username:    &selfUsername,
		MessageType: network.UdpMessageTypeMessage,
		Message:     &messageContent,
	}
	requestBuf, err := json.Marshal(request)
	if err != nil {
		return err
	}

	if _, err := conn.Write(requestBuf); err != nil {
		return fmt.Errorf("error broadcasting username: %s", err)
	}

	return nil
}
