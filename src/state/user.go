package state

import (
	"fmt"
	"time"
)

type User struct {
	Id             string
	Name           string
	Messages       []Message
	DetectedAt     time.Time
	LastActivityAt *time.Time
}

func NewUser(name string) User {
	return User{
		Id:             generateId(name),
		Name:           name,
		DetectedAt:     time.Now(),
		LastActivityAt: nil,
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
