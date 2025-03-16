package state

import (
	"fmt"
	"net"
	"time"
)

type User struct {
	Id             string
	Address        *net.UDPAddr
	Name           string
	Messages       []Message
	DetectedAt     time.Time
	LastActivityAt *time.Time
}

func NewUser(name string, address *net.UDPAddr) User {
	return User{
		Id:             generateId(name),
		Name:           name,
		DetectedAt:     time.Now(),
		LastActivityAt: nil,
		Address:        address,
		Messages:       []Message{},
	}
}

func generateId(name string) string {
	return fmt.Sprintf("%s_%d", name, time.Now().Unix())
}

func (u *User) AppendMessage(appState *AppState, content string, source MessageSource) {
	message := NewMessage(content, source)
	u.Messages = append(u.Messages, message)
	now := time.Now()
	u.LastActivityAt = &now
	appState.updateUser(u)
}
