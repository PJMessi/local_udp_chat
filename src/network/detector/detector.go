package detector

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/pjmessi/udp_chat/src/network"
	"github.com/pjmessi/udp_chat/src/state"
)

type SearchCompleteEvent struct {
	Username string
	Address  *net.UDPAddr
}

type NewMessageEvent struct {
	User           state.User
	MessageContent string
}

const BROADCAST_ADDRESS = "0.0.0.0:34354"
const SEARCH_DURATION = 2 * time.Second

func Listen(
	ctx context.Context,
	appState *state.AppState,
	searchCompleteCh chan<- SearchCompleteEvent,
	newMsgCh chan<- NewMessageEvent,
) error {
	addr, err := net.ResolveUDPAddr("udp", BROADCAST_ADDRESS)
	if err != nil {
		return fmt.Errorf("error resolving search addr: %s", err)
	}

	listener, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("error creating broadcast listener: %s", err)
	}

	for {
		var msgBuf []byte = make([]byte, 1024)
		n, clientAddr, err := listener.ReadFromUDP(msgBuf)
		if err != nil {
			return fmt.Errorf("error reading broadcast msg: %s", err)
		}

		msgBuf = msgBuf[:n]

		response := network.UdpMessage{}
		err = json.Unmarshal(msgBuf, &response)
		if err != nil {
			return fmt.Errorf("error while unmarshalling udp message: %s", err)
		}

		if response.MessageType == network.UdpMessageTypeBroadcast {
			if response.Username != nil && *response.Username != appState.SelfUsername {
				user := appState.GetUser(*response.Username)
				if user == nil {
					appState.AddUser(*response.Username, clientAddr)
					searchCompleteCh <- SearchCompleteEvent{
						Username: *response.Username,
						Address:  clientAddr,
					}
				}
			}
		} else if response.MessageType == network.UdpMessageTypeMessage {
			if response.Username != nil && response.Message != nil {
				user := appState.GetUser(*response.Username)
				if user != nil {
					user.AppendMessage(appState, *response.Message, state.MessageSourceUser)
					newMsgCh <- NewMessageEvent{
						User:           *user,
						MessageContent: *response.Message,
					}
				}
			}
		}
	}
}
