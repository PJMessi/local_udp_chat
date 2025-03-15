package state

import "time"

type MessageSource uint

const (
	MessageSourceUser MessageSource = iota
	MessageSourceSelf
)

type Message struct {
	Content   string
	Source    MessageSource
	Timestamp time.Time
}

func NewMessage(content string, source MessageSource) Message {
	return Message{
		Content:   content,
		Source:    source,
		Timestamp: time.Now(),
	}
}
