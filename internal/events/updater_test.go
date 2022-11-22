package events

import (
	"api/internal/commons"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestGenerateInitialEvents(t *testing.T) {
	type test struct {
		users             []commons.User
		eventsCount       commons.EventsCount
		possibleFullNames []commons.FullName
		possibleEmails    []commons.Email
		possibleSummaries []commons.Summary
		expEventsCount    map[commons.User]int
	}

	tests := []test{{
		users:             []commons.User{"new.user@yahoo.com", "sir.oliver@gmail.com"},
		eventsCount:       10,
		possibleFullNames: []commons.FullName{"Alan Ford", "Sir Oliver", "Bob Rock", "Broj Jedan"},
		possibleEmails:    []commons.Email{"kladionica@grupa.tnt", "new.user@yahoo.com", "sir.oliver@gmail.com"},
		possibleSummaries: []commons.Summary{"Christmas", "Easter", "New Year", "Birthday"},
		expEventsCount:    map[commons.User]int{"new.user@yahoo.com": 10, "sir.oliver@gmail.com": 10},
	}}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			GenerateInitialEvents(tt.users, tt.eventsCount, tt.possibleFullNames, tt.possibleEmails, tt.possibleSummaries)
			for user, expCount := range tt.expEventsCount {
				retCount := len(userEvents[user])
				if expCount != retCount {
					t.Errorf("cached events mismatch: expected '%d', received '%d'", expCount, retCount)
				}
			}
		})
	}
}

func TestContinuouslyUpdateEvents(t *testing.T) {
	type test struct {
		pauseLength            time.Duration
		replacementEventsCount commons.EventsCount
		updateEventsCount      commons.EventsCount
	}

	usersConfig = []commons.User{"new.user@yahoo.com", "sir.oliver@gmail.com"}
	eventsCountConfig = 2
	possibleFullNamesConfig = []commons.FullName{"Alan Ford", "Sir Oliver", "Bob Rock", "Broj Jedan"}
	possibleEmailsConfig = []commons.Email{"kladionica@grupa.tnt", "new.user@yahoo.com", "sir.oliver@gmail.com"}
	possibleSummariesConfig = []commons.Summary{"Christmas", "Easter", "New Year", "Birthday"}
	initialUserEvents := map[commons.User][]commons.Event{
		"sir.oliver@gmail.com": {{
			Uid:      "3f145806-6a72-11ed-a1eb-0242ac120002",
			FullName: "Sir Oliver",
			Email:    "sir.oliver@gmail.com",
			Start:    "20221231T200000Z",
			End:      "20221231T220000Z",
			Summary:  "Birthday Party",
			GeoLat:   46.419953,
			GeoLon:   15.869688,
		}, {
			Uid:      "8063f145-6002-01ed-41ea-20000242ac12",
			FullName: "Sir Oliver",
			Email:    "sir.oliver@gmail.com",
			Start:    "20230101T080000Z",
			End:      "20230101T090000Z",
			Summary:  "New Year",
			GeoLat:   46.419953,
			GeoLon:   15.869688,
		}},
		"new.user@yahoo.com": {{
			Uid:      "3f145806-6a72-11ed-a1eb-0242ac120002",
			FullName: "Sir Oliver",
			Email:    "sir.oliver@gmail.com",
			Start:    "20221231T200000Z",
			End:      "20221231T220000Z",
			Summary:  "Birthday Party",
			GeoLat:   46.419953,
			GeoLon:   15.869688,
		}, {
			Uid:      "07000123-0000-0000-41ea-20000242ac12",
			FullName: "User 7",
			Email:    "sir.oliver@gmail.com",
			Start:    "20230101T080000Z",
			End:      "20230101T090000Z",
			Summary:  "New Year",
			GeoLat:   46.419953,
			GeoLon:   15.869688,
		}},
	}
	for user, events := range initialUserEvents {
		userEvents[user] = append([]commons.Event{}, events...)
	}

	tests := []test{{
		pauseLength:            5 * time.Second,
		replacementEventsCount: 2,
		updateEventsCount:      2,
	}}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			go ContinuouslyUpdateEvents(tt.pauseLength, tt.replacementEventsCount, tt.updateEventsCount)
			time.Sleep(7 * time.Second)
			for user, events := range initialUserEvents {
				if reflect.DeepEqual(events, userEvents[user]) {
					t.Errorf("cached events were not updated for user %s", user)
				}
			}
		})
	}
}
