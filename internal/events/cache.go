package events

import (
	"api/internal/commons"
	"fmt"
	"sync"
)

func init() {
	userEvents = map[commons.User][]commons.Event{}
}

var (
	userEvents map[commons.User][]commons.Event
	mux        sync.RWMutex
)

func addEvent(user commons.User, event commons.Event) {
	mux.Lock()
	defer mux.Unlock()

	userEvents[user] = append(userEvents[user], event)
}

func deleteEvent(user commons.User, eventUid commons.Uid) {
	mux.Lock()
	defer mux.Unlock()

	events := userEvents[user]
	for i := range events {
		if events[i].Uid == eventUid {
			userEvents[user] = append(events[0:i], events[i+1:]...)
			return
		}
	}
}

func updateEvent(user commons.User, event commons.Event) {
	mux.Lock()
	defer mux.Unlock()

	events := userEvents[user]
	for i := range events {
		if events[i].Uid == event.Uid {
			followingEvents := append([]commons.Event{event}, events[i+1:]...)
			userEvents[user] = append(events[0:i], followingEvents...)
			return
		}
	}
}

func GetEvents(user commons.User, offset commons.Offset, limit commons.Limit) ([]commons.Event, error) {
	mux.RLock()
	defer mux.RUnlock()

	events, ok := userEvents[user]
	if !ok {
		return nil, fmt.Errorf("user does not exist")
	}
	if len(events) <= int(offset) {
		return []commons.Event{}, nil
	}
	lastIndex := int(offset) + int(limit)
	if len(events) < lastIndex {
		lastIndex = len(events)
	}

	page := make([]commons.Event, lastIndex-int(offset))
	copy(page, events[offset:lastIndex])
	return page, nil
}
