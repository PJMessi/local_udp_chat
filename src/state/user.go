package state

import (
	"net"
	"time"
)

type User struct {
	Address        *net.UDPAddr
	Name           string
	Messages       []Message
	DetectedAt     time.Time
	LastActivityAt *time.Time
}

func NewUser(name string, address *net.UDPAddr) User {
	return User{
		Name:           name,
		DetectedAt:     time.Now(),
		LastActivityAt: nil,
		Address:        address,
		Messages:       []Message{},
	}
}

func (u *User) AppendMessage(appState *AppState, content string, source MessageSource) {
	message := NewMessage(content, source)
	u.Messages = append(u.Messages, message)
	now := time.Now()
	u.LastActivityAt = &now
	appState.updateUser(u)
}
