package network

type UdpMessageType uint

const (
	UdpMessageTypeBroadcast UdpMessageType = iota
	UdpMessageTypeMessage
)

type UdpMessage struct {
	MessageType UdpMessageType `json:"message_type"`
	Username    *string        `json:"username"`
	Message     *string        `json:"message"`
}
