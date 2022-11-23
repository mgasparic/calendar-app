package events

import (
	"api/internal/commons"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"reflect"
	"testing"
	"time"
)

func TestContinuouslySynchronizeEvents(t *testing.T) {
	type test struct {
		generatorServiceUrlFormat commons.GeneratorServiceUrlFormat
		dbClient                  *sql.DB
		pauseLength               time.Duration
		users                     []commons.User
		eventsCount               commons.EventsCount
		expEventsFirst            []commons.Event
		expEventsSecond           []commons.Event
	}

	_ = exec.Command("docker", "build", "-t", "test-db", "--file", "../../build/package/db/Dockerfile", "../../scripts").Run()
	_ = exec.Command("docker", "container", "run", "--publish", "5432:5432", "--detach", "--env", "POSTGRES_PASSWORD=password123", "--name", "test-db", "test-db").Run()
	dbClient, _ := sql.Open("postgres", "host=localhost port=5432 user=postgres password=password123 sslmode=disable")
	defer func() {
		_ = dbClient.Close()
		_ = exec.Command("docker", "container", "rm", "-f", "test-db").Run()
		_ = exec.Command("docker", "image", "rm", "-f", "test-db").Run()
	}()

	for i := 0; i < 5; i++ {
		err := dbClient.Ping()
		if err != nil {
			log.Print("db not ready yet, will try again")
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}

	start := time.Now()
	eventsServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.RawQuery == "offset=0&limit=100" {
			if request.URL.Path == "/events/test@domain.net" {
				if start.Add(3 * time.Second).After(time.Now()) {
					_, _ = writer.Write([]byte(
						"BEGIN:VCALENDAR" +
							"\nVERSION:2.0" +
							"\nPRODID:-//accelerator/app//NONSGML v1.0//EN" +
							"\nBEGIN:VEVENT" +
							"\nUID:2f54ceca-d963-4ecb-b784-14fd215cca1d" +
							"\nORGANIZER;CN=Donald Duck:MAILTO:jimmy@accelerator-apps.com" +
							"\nDTSTART:20231214T234338Z" +
							"\nDTEND:20231215T234338Z" +
							"\nSUMMARY:Birth Day" +
							"\nGEO:65.012;97.012" +
							"\nEND:VEVENT" +
							"\nEND:VCALENDAR" +
							"\nBEGIN:VCALENDAR" +
							"\nVERSION:2.0" +
							"\nPRODID:-//accelerator/app//NONSGML v1.0//EN" +
							"\nBEGIN:VEVENT" +
							"\nUID:07e38ade-5b57-4dc2-ab5b-653fa5a713b2" +
							"\nORGANIZER;CN=James Bond:MAILTO:johnny@accelerator-apps.com" +
							"\nDTSTART:20230720T184338Z" +
							"\nDTEND:20230721T184338Z" +
							"\nSUMMARY:D Day" +
							"\nGEO:58;25" +
							"\nEND:VEVENT" +
							"\nEND:VCALENDAR" +
							"\nBEGIN:VCALENDAR" +
							"\nVERSION:2.0" +
							"\nPRODID:-//accelerator/app//NONSGML v1.0//EN" +
							"\nBEGIN:VEVENT" +
							"\nUID:84ff5ac7-fb80-4072-bd37-7b8e3724e3b4" +
							"\nORGANIZER;CN=Donald Duck:MAILTO:Bill.Maher@gmail.com" +
							"\nDTSTART:20240301T094338Z" +
							"\nDTEND:20240302T094338Z" +
							"\nSUMMARY:Independence Day" +
							"\nGEO:101;7" +
							"\nEND:VEVENT" +
							"\nEND:VCALENDAR" +
							"\nBEGIN:VCALENDAR" +
							"\nVERSION:2.0" +
							"\nPRODID:-//accelerator/app//NONSGML v1.0//EN" +
							"\nBEGIN:VEVENT" +
							"\nUID:9156a41f-5fc1-400d-85e4-5a9ef7c0a9a4" +
							"\nORGANIZER;CN=Minnie Mouse:MAILTO:Kanye.West@gmail.com" +
							"\nDTSTART:20220901T014338Z" +
							"\nDTEND:20220902T014338Z" +
							"\nSUMMARY:Sad Day" +
							"\nGEO:167;26" +
							"\nEND:VEVENT" +
							"\nEND:VCALENDAR" +
							"\nBEGIN:VCALENDAR" +
							"\nVERSION:2.0" +
							"\nPRODID:-//accelerator/app//NONSGML v1.0//EN" +
							"\nBEGIN:VEVENT" +
							"\nUID:85dee3e4-9078-4170-94db-ef4fe2841e17" +
							"\nORGANIZER;CN=James Bond:MAILTO:Bill.Maher@gmail.com" +
							"\nDTSTART:20200221T094338Z" +
							"\nDTEND:20200222T094338Z" +
							"\nSUMMARY:Sun Bathing" +
							"\nGEO:179;121" +
							"\nEND:VEVENT" +
							"\nEND:VCALENDAR"))
				} else {
					_, _ = writer.Write([]byte(
						"BEGIN:VCALENDAR" +
							"\nVERSION:2.0" +
							"\nPRODID:-//accelerator/app//NONSGML v1.0//EN" +
							"\nBEGIN:VEVENT" +
							"\nUID:2f54ceca-d963-4ecb-b784-14fd215cca1d" +
							"\nORGANIZER;CN=Donald Duck:MAILTO:jimmy@accelerator-apps.com" +
							"\nDTSTART:20231215T234338Z" +
							"\nDTEND:20231216T234338Z" +
							"\nSUMMARY:Birth Day" +
							"\nGEO:65.012;97.012" +
							"\nEND:VEVENT" +
							"\nEND:VCALENDAR" +
							"\nBEGIN:VCALENDAR" +
							"\nVERSION:2.0" +
							"\nPRODID:-//accelerator/app//NONSGML v1.0//EN" +
							"\nBEGIN:VEVENT" +
							"\nUID:07e38ade-5b57-4dc2-ab5b-653fa5a713b2" +
							"\nORGANIZER;CN=James Bond:MAILTO:johnny@accelerator-apps.com" +
							"\nDTSTART:20230720T184338Z" +
							"\nDTEND:20230721T184338Z" +
							"\nSUMMARY:D Day" +
							"\nGEO:58.12;25.34" +
							"\nEND:VEVENT" +
							"\nEND:VCALENDAR" +
							"\nBEGIN:VCALENDAR" +
							"\nVERSION:2.0" +
							"\nPRODID:-//accelerator/app//NONSGML v1.0//EN" +
							"\nBEGIN:VEVENT" +
							"\nUID:84ff5ac7-fb80-4072-bd37-7b8e3724e3b4" +
							"\nORGANIZER;CN=Donald Duck:MAILTO:Bill.Maher@gmail.com" +
							"\nDTSTART:20240301T094338Z" +
							"\nDTEND:20240302T094338Z" +
							"\nSUMMARY:Independence Day" +
							"\nGEO:101;7" +
							"\nEND:VEVENT" +
							"\nEND:VCALENDAR" +
							"\nBEGIN:VCALENDAR" +
							"\nVERSION:2.0" +
							"\nPRODID:-//accelerator/app//NONSGML v1.0//EN" +
							"\nBEGIN:VEVENT" +
							"\nUID:9156a41f-5fc1-400d-85e4-5a9ef7c0a9a4" +
							"\nORGANIZER;CN=Minnie Mouse:MAILTO:Kanye.West@gmail.com" +
							"\nDTSTART:20220901T014338Z" +
							"\nDTEND:20220902T014338Z" +
							"\nSUMMARY:Sad Day" +
							"\nGEO:167;26" +
							"\nEND:VEVENT" +
							"\nEND:VCALENDAR" +
							"\nBEGIN:VCALENDAR" +
							"\nVERSION:2.0" +
							"\nPRODID:-//accelerator/app//NONSGML v1.0//EN" +
							"\nBEGIN:VEVENT" +
							"\nUID:3699ee55-f35f-4c27-85db-dbc5899c6cf7" +
							"\nORGANIZER;CN=Donald Duck:MAILTO:Kanye.West@gmail.com" +
							"\nDTSTART:20210718T204338Z" +
							"\nDTEND:20210719T204338Z" +
							"\nSUMMARY:Happy Day" +
							"\nGEO:132;171" +
							"\nEND:VEVENT" +
							"\nEND:VCALENDAR"))
				}
			} else if request.URL.Path == "/events/another@domain.net" {
				_, _ = writer.Write([]byte(
					"BEGIN:VCALENDAR" +
						"\nVERSION:2.0" +
						"\nPRODID:-//accelerator/app//NONSGML v1.0//EN" +
						"\nBEGIN:VEVENT" +
						"\nUID:79a9e21d-6557-48df-85cf-f3c73d3e64e8" +
						"\nORGANIZER;CN=John Oliver:MAILTO:Mickey.Mouse@gmail.com" +
						"\nDTSTART:20200629T074338Z" +
						"\nDTEND:20200630T074338Z" +
						"\nSUMMARY:Fancy Party" +
						"\nGEO:59;137" +
						"\nEND:VEVENT" +
						"\nEND:VCALENDAR" +
						"\nBEGIN:VCALENDAR" +
						"\nVERSION:2.0" +
						"\nPRODID:-//accelerator/app//NONSGML v1.0//EN" +
						"\nBEGIN:VEVENT" +
						"\nUID:54b8fcc0-4517-430b-9f70-7bd608a69f3c" +
						"\nORGANIZER;CN=Minnie Mouse:MAILTO:Kanye.West@gmail.com" +
						"\nDTSTART:20240305T024338Z" +
						"\nDTEND:20240306T024338Z" +
						"\nSUMMARY:Sun Bathing" +
						"\nGEO:96;117" +
						"\nEND:VEVENT" +
						"\nEND:VCALENDAR" +
						"\nBEGIN:VCALENDAR" +
						"\nVERSION:2.0" +
						"\nPRODID:-//accelerator/app//NONSGML v1.0//EN" +
						"\nBEGIN:VEVENT" +
						"\nUID:74b93f6f-9414-45dd-bce3-ae7a1723d7a8" +
						"\nORGANIZER;CN=Sheldon Cooper:MAILTO:Mickey.Mouse@gmail.com" +
						"\nDTSTART:20210625T114338Z" +
						"\nDTEND:20210626T114338Z" +
						"\nSUMMARY:D Day" +
						"\nGEO:122;67" +
						"\nEND:VEVENT" +
						"\nEND:VCALENDAR" +
						"\nBEGIN:VCALENDAR" +
						"\nVERSION:2.0" +
						"\nPRODID:-//accelerator/app//NONSGML v1.0//EN" +
						"\nBEGIN:VEVENT" +
						"\nUID:1a14045e-7bb8-44a2-9476-305e08207994" +
						"\nORGANIZER;CN=Minnie Mouse:MAILTO:Mickey.Mouse@gmail.com" +
						"\nDTSTART:20220428T054338Z" +
						"\nDTEND:20220429T054338Z" +
						"\nSUMMARY:Skiing" +
						"\nGEO:113;45" +
						"\nEND:VEVENT" +
						"\nEND:VCALENDAR" +
						"\nBEGIN:VCALENDAR" +
						"\nVERSION:2.0" +
						"\nPRODID:-//accelerator/app//NONSGML v1.0//EN" +
						"\nBEGIN:VEVENT" +
						"\nUID:8d68ae32-a3a7-442b-873c-776044455db7" +
						"\nORGANIZER;CN=Sheldon Cooper:MAILTO:anna@accelerator-apps.com" +
						"\nDTSTART:20230131T024338Z" +
						"\nDTEND:20230201T024338Z" +
						"\nSUMMARY:Christmas" +
						"\nGEO:157;46" +
						"\nEND:VEVENT" +
						"\nEND:VCALENDAR"))
			}
		}
	}))

	tests := []test{
		{
			generatorServiceUrlFormat: commons.GeneratorServiceUrlFormat(fmt.Sprintf("%s/%s", eventsServer.URL, "events/%s?offset=%d&limit=%d")),
			dbClient:                  dbClient,
			pauseLength:               5 * time.Second,
			users:                     []commons.User{"test@domain.net", "another@domain.net"},
			eventsCount:               5,
			expEventsFirst: []commons.Event{
				{
					Uid:      "07e38ade-5b57-4dc2-ab5b-653fa5a713b2",
					FullName: "James Bond",
					Email:    "johnny@accelerator-apps.com",
					Start:    "20230720T184338Z",
					End:      "20230721T184338Z",
					Summary:  "D Day",
					GeoLat:   58,
					GeoLon:   25,
				},
				{
					Uid:      "1a14045e-7bb8-44a2-9476-305e08207994",
					FullName: "Minnie Mouse",
					Email:    "Mickey.Mouse@gmail.com",
					Start:    "20220428T054338Z",
					End:      "20220429T054338Z",
					Summary:  "Skiing",
					GeoLat:   113,
					GeoLon:   45,
				},
				{
					Uid:      "2f54ceca-d963-4ecb-b784-14fd215cca1d",
					FullName: "Donald Duck",
					Email:    "jimmy@accelerator-apps.com",
					Start:    "20231214T234338Z",
					End:      "20231215T234338Z",
					Summary:  "Birth Day",
					GeoLat:   65.012,
					GeoLon:   97.012,
				},
				{
					Uid:      "54b8fcc0-4517-430b-9f70-7bd608a69f3c",
					FullName: "Minnie Mouse",
					Email:    "Kanye.West@gmail.com",
					Start:    "20240305T024338Z",
					End:      "20240306T024338Z",
					Summary:  "Sun Bathing",
					GeoLat:   96,
					GeoLon:   117,
				},
				{
					Uid:      "74b93f6f-9414-45dd-bce3-ae7a1723d7a8",
					FullName: "Sheldon Cooper",
					Email:    "Mickey.Mouse@gmail.com",
					Start:    "20210625T114338Z",
					End:      "20210626T114338Z",
					Summary:  "D Day",
					GeoLat:   122,
					GeoLon:   67,
				},
				{
					Uid:      "79a9e21d-6557-48df-85cf-f3c73d3e64e8",
					FullName: "John Oliver",
					Email:    "Mickey.Mouse@gmail.com",
					Start:    "20200629T074338Z",
					End:      "20200630T074338Z",
					Summary:  "Fancy Party",
					GeoLat:   59,
					GeoLon:   137,
				},
				{
					Uid:      "84ff5ac7-fb80-4072-bd37-7b8e3724e3b4",
					FullName: "Donald Duck",
					Email:    "Bill.Maher@gmail.com",
					Start:    "20240301T094338Z",
					End:      "20240302T094338Z",
					Summary:  "Independence Day",
					GeoLat:   101,
					GeoLon:   7,
				},
				{
					Uid:      "85dee3e4-9078-4170-94db-ef4fe2841e17",
					FullName: "James Bond",
					Email:    "Bill.Maher@gmail.com",
					Start:    "20200221T094338Z",
					End:      "20200222T094338Z",
					Summary:  "Sun Bathing",
					GeoLat:   179,
					GeoLon:   121,
				},
				{
					Uid:      "8d68ae32-a3a7-442b-873c-776044455db7",
					FullName: "Sheldon Cooper",
					Email:    "anna@accelerator-apps.com",
					Start:    "20230131T024338Z",
					End:      "20230201T024338Z",
					Summary:  "Christmas",
					GeoLat:   157,
					GeoLon:   46,
				},
				{
					Uid:      "9156a41f-5fc1-400d-85e4-5a9ef7c0a9a4",
					FullName: "Minnie Mouse",
					Email:    "Kanye.West@gmail.com",
					Start:    "20220901T014338Z",
					End:      "20220902T014338Z",
					Summary:  "Sad Day",
					GeoLat:   167,
					GeoLon:   26,
				},
			},
			expEventsSecond: []commons.Event{
				{
					Uid:      "07e38ade-5b57-4dc2-ab5b-653fa5a713b2",
					FullName: "James Bond",
					Email:    "johnny@accelerator-apps.com",
					Start:    "20230720T184338Z",
					End:      "20230721T184338Z",
					Summary:  "D Day",
					GeoLat:   58.12,
					GeoLon:   25.34,
				},
				{
					Uid:      "1a14045e-7bb8-44a2-9476-305e08207994",
					FullName: "Minnie Mouse",
					Email:    "Mickey.Mouse@gmail.com",
					Start:    "20220428T054338Z",
					End:      "20220429T054338Z",
					Summary:  "Skiing",
					GeoLat:   113,
					GeoLon:   45,
				},
				{
					Uid:      "2f54ceca-d963-4ecb-b784-14fd215cca1d",
					FullName: "Donald Duck",
					Email:    "jimmy@accelerator-apps.com",
					Start:    "20231215T234338Z",
					End:      "20231216T234338Z",
					Summary:  "Birth Day",
					GeoLat:   65.012,
					GeoLon:   97.012,
				},
				{
					Uid:      "3699ee55-f35f-4c27-85db-dbc5899c6cf7",
					FullName: "Donald Duck",
					Email:    "Kanye.West@gmail.com",
					Start:    "20210718T204338Z",
					End:      "20210719T204338Z",
					Summary:  "Happy Day",
					GeoLat:   132,
					GeoLon:   171,
				},
				{
					Uid:      "54b8fcc0-4517-430b-9f70-7bd608a69f3c",
					FullName: "Minnie Mouse",
					Email:    "Kanye.West@gmail.com",
					Start:    "20240305T024338Z",
					End:      "20240306T024338Z",
					Summary:  "Sun Bathing",
					GeoLat:   96,
					GeoLon:   117,
				},
				{
					Uid:      "74b93f6f-9414-45dd-bce3-ae7a1723d7a8",
					FullName: "Sheldon Cooper",
					Email:    "Mickey.Mouse@gmail.com",
					Start:    "20210625T114338Z",
					End:      "20210626T114338Z",
					Summary:  "D Day",
					GeoLat:   122,
					GeoLon:   67,
				},
				{
					Uid:      "79a9e21d-6557-48df-85cf-f3c73d3e64e8",
					FullName: "John Oliver",
					Email:    "Mickey.Mouse@gmail.com",
					Start:    "20200629T074338Z",
					End:      "20200630T074338Z",
					Summary:  "Fancy Party",
					GeoLat:   59,
					GeoLon:   137,
				},
				{
					Uid:      "84ff5ac7-fb80-4072-bd37-7b8e3724e3b4",
					FullName: "Donald Duck",
					Email:    "Bill.Maher@gmail.com",
					Start:    "20240301T094338Z",
					End:      "20240302T094338Z",
					Summary:  "Independence Day",
					GeoLat:   101,
					GeoLon:   7,
				},
				{
					Uid:      "8d68ae32-a3a7-442b-873c-776044455db7",
					FullName: "Sheldon Cooper",
					Email:    "anna@accelerator-apps.com",
					Start:    "20230131T024338Z",
					End:      "20230201T024338Z",
					Summary:  "Christmas",
					GeoLat:   157,
					GeoLon:   46,
				},
				{
					Uid:      "9156a41f-5fc1-400d-85e4-5a9ef7c0a9a4",
					FullName: "Minnie Mouse",
					Email:    "Kanye.West@gmail.com",
					Start:    "20220901T014338Z",
					End:      "20220902T014338Z",
					Summary:  "Sad Day",
					GeoLat:   167,
					GeoLon:   26,
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			go ContinuouslySynchronizeEvents(tt.generatorServiceUrlFormat, tt.dbClient, tt.pauseLength, tt.users, tt.eventsCount)
			time.Sleep(3 * time.Second)
			retEventsFirst := getEventsInDb(dbClient)
			if !reflect.DeepEqual(tt.expEventsFirst, retEventsFirst) {
				t.Errorf("stored events mismatch: expected '%v', received '%v'", tt.expEventsFirst, retEventsFirst)
			}
			time.Sleep(5 * time.Second)
			retEventsSecond := getEventsInDb(dbClient)
			if !reflect.DeepEqual(tt.expEventsSecond, retEventsSecond) {
				t.Errorf("stored events mismatch: expected '%v', received '%v'", tt.expEventsSecond, retEventsSecond)
			}
		})
	}
}

func getEventsInDb(dbClient *sql.DB) []commons.Event {
	var events []commons.Event
	rows, _ := dbClient.Query("SELECT id, organizer_full_name, organizer_email, time_start, time_end, event_summary, geo_lat, geo_lon FROM event ORDER BY id ASC")
	for rows.Next() {
		var event commons.Event
		_ = rows.Scan(&event.Uid, &event.FullName, &event.Email, &event.Start, &event.End, &event.Summary, &event.GeoLat, &event.GeoLon)
		events = append(events, event)
	}
	_ = rows.Close()
	return events
}
