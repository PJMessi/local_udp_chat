package state

import (
	"fmt"
	"slices"
	"sort"
	"sync"
	"time"
)

const useSort bool = false

type AppState struct {
	users        []User
	UsersMutex   sync.RWMutex
	SelectedUser *User
}

// Creates new AppState with dummy users and messages.
// TODO: Remove dummy variables.
func NewAppState() *AppState {
	// Generate dummy users and messages.
	users := []User{
		{Name: "Alice", Messages: []Message{}},
		{Name: "Bob", Messages: []Message{}},
		{Name: "Charlie", Messages: []Message{}},
		{Name: "Dave", Messages: []Message{}},
	}

	for i := range users {
		users[i].Messages = append(users[i].Messages, Message{
			Source:    MessageSourceSelf,
			Content:   fmt.Sprintf("Hi %s, how are you?", users[i].Name),
			Timestamp: time.Now().Add(-time.Hour),
		})
		users[i].Messages = append(users[i].Messages, Message{
			Source:    MessageSourceUser,
			Content:   "I'm doing well, thanks for asking!",
			Timestamp: time.Now().Add(-time.Minute * 55),
		})
	}

	appState := &AppState{
		users:        users,
		UsersMutex:   sync.RWMutex{},
		SelectedUser: nil,
	}

	appState.SelectUser(0)

	return appState
}

// Returns a copy of the list of users from the app state.
func (a *AppState) GetUsers() []User {
	a.UsersMutex.RLock()
	defer a.UsersMutex.RUnlock()

	if useSort {
		// Sort by LastActivityAt / DetectedAt
		sort.Slice(a.users, func(i, j int) bool {
			dateI := a.users[i].DetectedAt
			if a.users[i].LastActivityAt != nil {
				dateI = *a.users[i].LastActivityAt
			}

			dateJ := a.users[j].DetectedAt
			if a.users[j].LastActivityAt != nil {
				dateJ = *a.users[j].LastActivityAt
			}

			return dateI.After(dateJ)
		})
	}

	usersCopy := make([]User, len(a.users))
	copy(usersCopy, a.users)
	return usersCopy
}

// Sets the user with the given index as the selected user, and returns the copy of the user.
func (a *AppState) SelectUser(i uint) *User {
	a.UsersMutex.Lock()
	defer a.UsersMutex.Unlock()

	if len(a.users) == 0 {
		return nil
	}

	a.SelectedUser = &a.users[i]
	selectedUserCopy := *a.SelectedUser
	return &selectedUserCopy
}

// Adds a user in the app state.
func (a *AppState) AddUser(name string) {
	a.UsersMutex.Lock()
	defer a.UsersMutex.Unlock()

	newUser := NewUser(name)
	a.users = append(a.users, newUser)
}

// Removes the user with the given name from the app state.
func (a *AppState) RemoveUser(name string) {
	a.UsersMutex.Lock()
	defer a.UsersMutex.Unlock()

	for i, user := range a.users {
		if user.Name == name {
			a.users = slices.Delete(a.users, i, i+1)
			return
		}
	}
}

// Returns a copy of the user with the given name from the app state.
func (a *AppState) GetUser(name string) *User {
	a.UsersMutex.RLock()
	defer a.UsersMutex.Unlock()

	var userCopy User
	for _, user := range a.users {
		if user.Name == name {
			userCopy = user
			return &userCopy
		}
	}

	return nil
}

// Returns a copy of the currently selected user from the app state.
func (a *AppState) GetSelectedUser() *User {
	a.UsersMutex.RLock()
	defer a.UsersMutex.RUnlock()

	if a.SelectedUser == nil {
		return nil
	}

	selectedUserCopy := *a.SelectedUser
	return &selectedUserCopy
}

func (a *AppState) updateUser(updatedUser *User) {
	a.UsersMutex.Lock()
	defer a.UsersMutex.Unlock()

	for index := range a.users {
		if a.users[index].Name == updatedUser.Name {
			a.users[index].Id = updatedUser.Id
			a.users[index].Name = updatedUser.Name
			a.users[index].Messages = updatedUser.Messages
			a.users[index].DetectedAt = updatedUser.DetectedAt
			a.users[index].LastActivityAt = updatedUser.LastActivityAt
			return
		}
	}
}
