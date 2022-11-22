package events

import (
	"api/internal/commons"
	"fmt"
	"reflect"
	"testing"
)

func TestAddEvent(t *testing.T) {
	type test struct {
		user      commons.User
		event     commons.Event
		expEvents []commons.Event
	}

	userEvents = map[commons.User][]commons.Event{
		"sir.oliver@gmail.com": {{
			Uid:      "3f145806-6a72-11ed-a1eb-0242ac120002",
			FullName: "Sir Oliver",
			Email:    "sir.oliver@gmail.com",
			Start:    "20221231T200000Z",
			End:      "20221231T220000Z",
			Summary:  "Birthday Party",
			GeoLat:   46.419953,
			GeoLon:   15.869688,
		}},
	}

	tests := []test{
		{
			user: "new.user@yahoo.com",
			event: commons.Event{
				Uid:      "3f145806-6a72-11ed-a1eb-0242ac120002",
				FullName: "Sir Oliver",
				Email:    "sir.oliver@gmail.com",
				Start:    "20221231T200000Z",
				End:      "20221231T220000Z",
				Summary:  "Birthday Party",
				GeoLat:   46.419953,
				GeoLon:   15.869688,
			},
			expEvents: []commons.Event{{
				Uid:      "3f145806-6a72-11ed-a1eb-0242ac120002",
				FullName: "Sir Oliver",
				Email:    "sir.oliver@gmail.com",
				Start:    "20221231T200000Z",
				End:      "20221231T220000Z",
				Summary:  "Birthday Party",
				GeoLat:   46.419953,
				GeoLon:   15.869688,
			}},
		},
		{
			user: "sir.oliver@gmail.com",
			event: commons.Event{
				Uid:      "8063f145-6002-01ed-41ea-20000242ac12",
				FullName: "Sir Oliver",
				Email:    "sir.oliver@gmail.com",
				Start:    "20230101T080000Z",
				End:      "20230101T090000Z",
				Summary:  "New Year",
				GeoLat:   46.419953,
				GeoLon:   15.869688,
			},
			expEvents: []commons.Event{{
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
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			addEvent(tt.user, tt.event)
			if !reflect.DeepEqual(tt.expEvents, userEvents[tt.user]) {
				t.Errorf("cached events mismatch: expected '%v', received '%v'", tt.expEvents, userEvents[tt.user])
			}
		})
	}
}

func TestDeleteEvent(t *testing.T) {
	type test struct {
		user      commons.User
		eventUid  commons.Uid
		expEvents []commons.Event
	}

	userEvents = map[commons.User][]commons.Event{
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
		}},
	}

	tests := []test{
		{
			user:      "non.existing@hotmail.com",
			eventUid:  "3f145806-6a72-11ed-a1eb-0242ac120002",
			expEvents: nil,
		},
		{
			user:     "new.user@yahoo.com",
			eventUid: "8063f145-6002-01ed-41ea-20000242ac12",
			expEvents: []commons.Event{{
				Uid:      "3f145806-6a72-11ed-a1eb-0242ac120002",
				FullName: "Sir Oliver",
				Email:    "sir.oliver@gmail.com",
				Start:    "20221231T200000Z",
				End:      "20221231T220000Z",
				Summary:  "Birthday Party",
				GeoLat:   46.419953,
				GeoLon:   15.869688,
			}},
		},
		{
			user:      "new.user@yahoo.com",
			eventUid:  "3f145806-6a72-11ed-a1eb-0242ac120002",
			expEvents: []commons.Event{},
		},
		{
			user:     "sir.oliver@gmail.com",
			eventUid: "8063f145-6002-01ed-41ea-20000242ac12",
			expEvents: []commons.Event{{
				Uid:      "3f145806-6a72-11ed-a1eb-0242ac120002",
				FullName: "Sir Oliver",
				Email:    "sir.oliver@gmail.com",
				Start:    "20221231T200000Z",
				End:      "20221231T220000Z",
				Summary:  "Birthday Party",
				GeoLat:   46.419953,
				GeoLon:   15.869688,
			}},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			deleteEvent(tt.user, tt.eventUid)
			if !reflect.DeepEqual(tt.expEvents, userEvents[tt.user]) {
				t.Errorf("cached events mismatch: expected '%v', received '%v'", tt.expEvents, userEvents[tt.user])
			}
		})
	}
}

func TestUpdateEvent(t *testing.T) {
	type test struct {
		user      commons.User
		event     commons.Event
		expEvents []commons.Event
	}

	userEvents = map[commons.User][]commons.Event{
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
		}},
	}

	tests := []test{
		{
			user: "non.existing@hotmail.com",
			event: commons.Event{
				Uid:      "3f145806-6a72-11ed-a1eb-0242ac120002",
				FullName: "Sir Oliver",
				Email:    "sir.oliver@gmail.com",
				Start:    "20221230T200000Z",
				End:      "20221230T220000Z",
				Summary:  "Birthday Party",
				GeoLat:   45.419953,
				GeoLon:   16.869688,
			},
			expEvents: nil,
		},
		{
			user: "new.user@yahoo.com",
			event: commons.Event{
				Uid:      "3f145806-6a72-11ed-a1eb-0242ac120002",
				FullName: "Sir Oliver",
				Email:    "sir.oliver@gmail.com",
				Start:    "20221230T200000Z",
				End:      "20221230T220000Z",
				Summary:  "Birthday Party",
				GeoLat:   45.419953,
				GeoLon:   16.869688,
			},
			expEvents: []commons.Event{{
				Uid:      "3f145806-6a72-11ed-a1eb-0242ac120002",
				FullName: "Sir Oliver",
				Email:    "sir.oliver@gmail.com",
				Start:    "20221230T200000Z",
				End:      "20221230T220000Z",
				Summary:  "Birthday Party",
				GeoLat:   45.419953,
				GeoLon:   16.869688,
			}},
		},
		{
			user: "sir.oliver@gmail.com",
			event: commons.Event{
				Uid:      "8063f145-6002-01ed-41ea-20000242ac12",
				FullName: "Alan Ford",
				Email:    "sir.oliver@gmail.com",
				Start:    "20230101T080000Z",
				End:      "20230101T090000Z",
				Summary:  "New Year",
				GeoLat:   46.419953,
				GeoLon:   15.869688,
			},
			expEvents: []commons.Event{{
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
				FullName: "Alan Ford",
				Email:    "sir.oliver@gmail.com",
				Start:    "20230101T080000Z",
				End:      "20230101T090000Z",
				Summary:  "New Year",
				GeoLat:   46.419953,
				GeoLon:   15.869688,
			}},
		},
		{
			user: "sir.oliver@gmail.com",
			event: commons.Event{
				Uid:      "3f145806-6a72-11ed-a1eb-0242ac120002",
				FullName: "Sir Oliver",
				Email:    "sir.oliver@gmail.com",
				Start:    "20221231T200000Z",
				End:      "20221231T220000Z",
				Summary:  "Birthday Party Cancelled",
				GeoLat:   46.419953,
				GeoLon:   15.869688,
			},
			expEvents: []commons.Event{{
				Uid:      "3f145806-6a72-11ed-a1eb-0242ac120002",
				FullName: "Sir Oliver",
				Email:    "sir.oliver@gmail.com",
				Start:    "20221231T200000Z",
				End:      "20221231T220000Z",
				Summary:  "Birthday Party Cancelled",
				GeoLat:   46.419953,
				GeoLon:   15.869688,
			}, {
				Uid:      "8063f145-6002-01ed-41ea-20000242ac12",
				FullName: "Alan Ford",
				Email:    "sir.oliver@gmail.com",
				Start:    "20230101T080000Z",
				End:      "20230101T090000Z",
				Summary:  "New Year",
				GeoLat:   46.419953,
				GeoLon:   15.869688,
			}},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			updateEvent(tt.user, tt.event)
			if !reflect.DeepEqual(tt.expEvents, userEvents[tt.user]) {
				t.Errorf("cached events mismatch: expected '%v', received '%v'", tt.expEvents, userEvents[tt.user])
			}
		})
	}
}

func TestGetEvents(t *testing.T) {
	type test struct {
		user        commons.User
		offset      commons.Offset
		limit       commons.Limit
		expEvents   []commons.Event
		expRetError bool
	}

	userEvents = map[commons.User][]commons.Event{
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
		}, {
			Uid:      "01000123-0000-0000-41ea-20000242ac12",
			FullName: "User 1",
			Email:    "sir.oliver@gmail.com",
			Start:    "20230101T080000Z",
			End:      "20230101T090000Z",
			Summary:  "New Year",
			GeoLat:   46.419953,
			GeoLon:   15.869688,
		}, {
			Uid:      "02000123-0000-0000-41ea-20000242ac12",
			FullName: "User 2",
			Email:    "sir.oliver@gmail.com",
			Start:    "20230101T080000Z",
			End:      "20230101T090000Z",
			Summary:  "New Year",
			GeoLat:   46.419953,
			GeoLon:   15.869688,
		}, {
			Uid:      "03000123-0000-0000-41ea-20000242ac12",
			FullName: "User 3",
			Email:    "sir.oliver@gmail.com",
			Start:    "20230101T080000Z",
			End:      "20230101T090000Z",
			Summary:  "New Year",
			GeoLat:   46.419953,
			GeoLon:   15.869688,
		}, {
			Uid:      "04000123-0000-0000-41ea-20000242ac12",
			FullName: "User 4",
			Email:    "sir.oliver@gmail.com",
			Start:    "20230101T080000Z",
			End:      "20230101T090000Z",
			Summary:  "New Year",
			GeoLat:   46.419953,
			GeoLon:   15.869688,
		}, {
			Uid:      "05000123-0000-0000-41ea-20000242ac12",
			FullName: "User 5",
			Email:    "sir.oliver@gmail.com",
			Start:    "20230101T080000Z",
			End:      "20230101T090000Z",
			Summary:  "New Year",
			GeoLat:   46.419953,
			GeoLon:   15.869688,
		}, {
			Uid:      "06000123-0000-0000-41ea-20000242ac12",
			FullName: "User 6",
			Email:    "sir.oliver@gmail.com",
			Start:    "20230101T080000Z",
			End:      "20230101T090000Z",
			Summary:  "New Year",
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
		}, {
			Uid:      "08000123-0000-0000-41ea-20000242ac12",
			FullName: "User 8",
			Email:    "sir.oliver@gmail.com",
			Start:    "20230101T080000Z",
			End:      "20230101T090000Z",
			Summary:  "New Year",
			GeoLat:   46.419953,
			GeoLon:   15.869688,
		}, {
			Uid:      "09000123-0000-0000-41ea-20000242ac12",
			FullName: "User 9",
			Email:    "sir.oliver@gmail.com",
			Start:    "20230101T080000Z",
			End:      "20230101T090000Z",
			Summary:  "New Year",
			GeoLat:   46.419953,
			GeoLon:   15.869688,
		}, {
			Uid:      "10000123-0000-0000-41ea-20000242ac12",
			FullName: "User 10",
			Email:    "sir.oliver@gmail.com",
			Start:    "20230101T080000Z",
			End:      "20230101T090000Z",
			Summary:  "New Year",
			GeoLat:   46.419953,
			GeoLon:   15.869688,
		}},
		"new.user@yahoo.com": {{
			Uid:      "11000123-0000-0000-41ea-20000242ac12",
			FullName: "User 11",
			Email:    "sir.oliver@gmail.com",
			Start:    "20230101T080000Z",
			End:      "20230101T090000Z",
			Summary:  "New Year",
			GeoLat:   46.419953,
			GeoLon:   15.869688,
		}},
	}

	tests := []test{
		{
			user:        "non.existing@hotmail.com",
			offset:      0,
			limit:       10,
			expEvents:   nil,
			expRetError: true,
		},
		{
			user:   "new.user@yahoo.com",
			offset: 0,
			limit:  10,
			expEvents: []commons.Event{{
				Uid:      "11000123-0000-0000-41ea-20000242ac12",
				FullName: "User 11",
				Email:    "sir.oliver@gmail.com",
				Start:    "20230101T080000Z",
				End:      "20230101T090000Z",
				Summary:  "New Year",
				GeoLat:   46.419953,
				GeoLon:   15.869688,
			}},
			expRetError: false,
		},
		{
			user:   "sir.oliver@gmail.com",
			offset: 0,
			limit:  10,
			expEvents: []commons.Event{
				{
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
				}, {
					Uid:      "01000123-0000-0000-41ea-20000242ac12",
					FullName: "User 1",
					Email:    "sir.oliver@gmail.com",
					Start:    "20230101T080000Z",
					End:      "20230101T090000Z",
					Summary:  "New Year",
					GeoLat:   46.419953,
					GeoLon:   15.869688,
				}, {
					Uid:      "02000123-0000-0000-41ea-20000242ac12",
					FullName: "User 2",
					Email:    "sir.oliver@gmail.com",
					Start:    "20230101T080000Z",
					End:      "20230101T090000Z",
					Summary:  "New Year",
					GeoLat:   46.419953,
					GeoLon:   15.869688,
				}, {
					Uid:      "03000123-0000-0000-41ea-20000242ac12",
					FullName: "User 3",
					Email:    "sir.oliver@gmail.com",
					Start:    "20230101T080000Z",
					End:      "20230101T090000Z",
					Summary:  "New Year",
					GeoLat:   46.419953,
					GeoLon:   15.869688,
				}, {
					Uid:      "04000123-0000-0000-41ea-20000242ac12",
					FullName: "User 4",
					Email:    "sir.oliver@gmail.com",
					Start:    "20230101T080000Z",
					End:      "20230101T090000Z",
					Summary:  "New Year",
					GeoLat:   46.419953,
					GeoLon:   15.869688,
				}, {
					Uid:      "05000123-0000-0000-41ea-20000242ac12",
					FullName: "User 5",
					Email:    "sir.oliver@gmail.com",
					Start:    "20230101T080000Z",
					End:      "20230101T090000Z",
					Summary:  "New Year",
					GeoLat:   46.419953,
					GeoLon:   15.869688,
				}, {
					Uid:      "06000123-0000-0000-41ea-20000242ac12",
					FullName: "User 6",
					Email:    "sir.oliver@gmail.com",
					Start:    "20230101T080000Z",
					End:      "20230101T090000Z",
					Summary:  "New Year",
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
				}, {
					Uid:      "08000123-0000-0000-41ea-20000242ac12",
					FullName: "User 8",
					Email:    "sir.oliver@gmail.com",
					Start:    "20230101T080000Z",
					End:      "20230101T090000Z",
					Summary:  "New Year",
					GeoLat:   46.419953,
					GeoLon:   15.869688,
				},
			},
			expRetError: false,
		},
		{
			user:   "sir.oliver@gmail.com",
			offset: 10,
			limit:  10,
			expEvents: []commons.Event{
				{
					Uid:      "09000123-0000-0000-41ea-20000242ac12",
					FullName: "User 9",
					Email:    "sir.oliver@gmail.com",
					Start:    "20230101T080000Z",
					End:      "20230101T090000Z",
					Summary:  "New Year",
					GeoLat:   46.419953,
					GeoLon:   15.869688,
				}, {
					Uid:      "10000123-0000-0000-41ea-20000242ac12",
					FullName: "User 10",
					Email:    "sir.oliver@gmail.com",
					Start:    "20230101T080000Z",
					End:      "20230101T090000Z",
					Summary:  "New Year",
					GeoLat:   46.419953,
					GeoLon:   15.869688,
				},
			},
			expRetError: false,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			retEvents, retErr := GetEvents(tt.user, tt.offset, tt.limit)
			if !reflect.DeepEqual(tt.expEvents, retEvents) {
				t.Errorf("returned events mismatch: expected '%v', received '%v'", tt.expEvents, retEvents)
			}
			if tt.expRetError != (retErr != nil) {
				t.Errorf("return error mismatch: expected is error '%v', received error '%v'", tt.expRetError, retErr)
			}
		})
	}
}
