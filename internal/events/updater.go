package events

import (
	"api/internal/commons"
	"github.com/google/uuid"
	"log"
	"math/rand"
	"sync"
	"time"
)

const icsTimeFormat = "20060102T150405Z"

var (
	usersConfig               []commons.User
	eventsCountConfig         commons.EventsCount
	possibleFullNamesConfig   []commons.FullName
	possibleEmailsConfig      []commons.Email
	possibleSummariesConfig   []commons.Summary
	timeIntervalSizeInHours   = 24 * 30 * 12 * 5
	timeIntervalStartInHours  = -24 * 30 * 12 * 3
	maxEventDurationInMinutes = 1440
)

func GenerateInitialEvents(users []commons.User, eventsCount commons.EventsCount, possibleFullNames []commons.FullName, possibleEmails []commons.Email, possibleSummaries []commons.Summary) {
	usersConfig = users
	eventsCountConfig = eventsCount
	possibleFullNamesConfig = possibleFullNames
	possibleEmailsConfig = possibleEmails
	possibleSummariesConfig = possibleSummaries

	var wg sync.WaitGroup

	for _, user := range users {
		wg.Add(1)
		go func(u commons.User) {
			defer wg.Done()
			for i := 0; i < int(eventsCount); i++ {
				addEvent(u, getRandomEvent())
			}
		}(user)
	}

	wg.Wait()
}

func ContinuouslyUpdateEvents(pauseLength time.Duration, replacementEventsCount, updateEventsCount commons.EventsCount) {
	for {
		time.Sleep(pauseLength)
		var wg sync.WaitGroup
		for _, user := range usersConfig {
			wg.Add(1)
			go func(u commons.User) {
				defer wg.Done()
				for i := 0; i < int(replacementEventsCount); i++ {
					events, err := GetEvents(u, commons.Offset(rand.Intn(int(eventsCountConfig))), 1)
					if err != nil {
						log.Fatal(err)
					}
					deleteEvent(u, events[0].Uid)
					addEvent(u, getRandomEvent())
				}
				for i := 0; i < int(updateEventsCount); i++ {
					events, err := GetEvents(u, commons.Offset(rand.Intn(int(eventsCountConfig))), 1)
					if err != nil {
						log.Fatal(err)
					}
					randomEvent := getRandomEvent()
					randomEvent.Uid = events[0].Uid
					updateEvent(u, randomEvent)
				}
			}(user)
		}

		wg.Wait()
	}
}

func getRandomEvent() commons.Event {
	start := time.Now().Add(time.Duration(rand.Intn(timeIntervalSizeInHours)) * time.Hour).Add(time.Duration(timeIntervalStartInHours) * time.Hour)
	end := start.Add(time.Duration(maxEventDurationInMinutes) * time.Minute)
	return commons.Event{
		Uid:      commons.Uid(uuid.New().String()),
		FullName: possibleFullNamesConfig[rand.Intn(len(possibleFullNamesConfig))],
		Email:    possibleEmailsConfig[rand.Intn(len(possibleEmailsConfig))],
		Start:    commons.Start(start.Format(icsTimeFormat)),
		End:      commons.End(end.Format(icsTimeFormat)),
		Summary:  possibleSummariesConfig[rand.Intn(len(possibleSummariesConfig))],
		GeoLat:   commons.GeoLat(float32(rand.Intn(3600)) / 20),
		GeoLon:   commons.GeoLon(float32(rand.Intn(3600)) / 20),
	}
}
