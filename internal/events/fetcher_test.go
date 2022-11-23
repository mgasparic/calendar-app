package events

import (
	"api/internal/commons"
	"database/sql"
	"fmt"
	"log"
	"os/exec"
	"reflect"
	"testing"
	"time"
)

func TestFetchEvents(t *testing.T) {
	type test struct {
		user        commons.User
		offset      commons.Offset
		limit       commons.Limit
		expEvents   []commons.Event
		expRetError bool
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

	_, _ = dbClient.Exec("INSERT INTO customer (id, email) VALUES ('0eabee48-0756-4e09-8b88-4f86be2e8786', 'jimmy@accelerator-apps.com')")
	_, _ = dbClient.Exec("INSERT INTO customer (id, email) VALUES ('47a0f0e0-3e3f-42a3-b960-866a5c277192', 'johnny@accelerator-apps.com')")

	_, _ = dbClient.Exec("INSERT INTO event (id, owner_id, organizer_full_name, organizer_email, time_start, time_end, event_summary, geo_lat, geo_lon) VALUES ('a0d1485e-1b01-45b4-868c-bf0ea4d6a9d6', '0eabee48-0756-4e09-8b88-4f86be2e8786', 'Kanye West', 'Kanye.West@gmail.com', '20201213T020436Z', '20201214T020436Z', 'Fancy Party', 136.4, 133.85)")
	_, _ = dbClient.Exec("INSERT INTO event (id, owner_id, organizer_full_name, organizer_email, time_start, time_end, event_summary, geo_lat, geo_lon) VALUES ('09293e4a-66b2-4a77-b373-ea2baceb0228', '47a0f0e0-3e3f-42a3-b960-866a5c277192', 'Minnie Mouse', 'John.Oliver@gmail.com', '20210712T040436Z', '20210713T040436Z', 'Happy Day', 54.35, 64.4)")
	_, _ = dbClient.Exec("INSERT INTO event (id, owner_id, organizer_full_name, organizer_email, time_start, time_end, event_summary, geo_lat, geo_lon) VALUES ('0f65663e-95ec-4b83-9e36-7519d6ee0980', '47a0f0e0-3e3f-42a3-b960-866a5c277192', 'Kanye West', 'Donald.Duck@gmail.com', '20210207T040436Z', '20210208T040436Z', 'Sad Day', 147.7, 160.25)")
	_, _ = dbClient.Exec("INSERT INTO event (id, owner_id, organizer_full_name, organizer_email, time_start, time_end, event_summary, geo_lat, geo_lon) VALUES ('27cc98be-04c7-481a-9788-b41adc44c0a7', '47a0f0e0-3e3f-42a3-b960-866a5c277192', 'Sheldon Cooper', 'Kanye.West@gmail.com', '20240902T110436Z', '20240903T110436Z', 'Sun Bathing', 154.45, 16.4)")
	_, _ = dbClient.Exec("INSERT INTO event (id, owner_id, organizer_full_name, organizer_email, time_start, time_end, event_summary, geo_lat, geo_lon) VALUES ('ce664e65-4e31-4a83-b226-6d3602b87bf1', '47a0f0e0-3e3f-42a3-b960-866a5c277192', 'Kanye West', 'maria@accelerator-apps.com', '20200214T020436Z', '20200215T020436Z', 'D Day', 89.35, 141.55)")
	_, _ = dbClient.Exec("INSERT INTO event (id, owner_id, organizer_full_name, organizer_email, time_start, time_end, event_summary, geo_lat, geo_lon) VALUES ('e426e402-7718-41dc-9dee-b9ee86a5d41a', '47a0f0e0-3e3f-42a3-b960-866a5c277192', 'Sheldon Cooper', 'Mickey.Mouse@gmail.com', '20230204T220436Z', '20230205T220436Z', 'Fancy Party', 175.3, 124.75)")

	SetDbClient(dbClient)

	tests := []test{
		{
			user:        "not@existing.com",
			offset:      0,
			limit:       3,
			expEvents:   nil,
			expRetError: false,
		},
		{
			user:   "jimmy@accelerator-apps.com",
			offset: 0,
			limit:  3,
			expEvents: []commons.Event{
				{
					Uid:      "a0d1485e-1b01-45b4-868c-bf0ea4d6a9d6",
					FullName: "Kanye West",
					Email:    "Kanye.West@gmail.com",
					Start:    "20201213T020436Z",
					End:      "20201214T020436Z",
					Summary:  "Fancy Party",
					GeoLat:   136.4,
					GeoLon:   133.85,
				},
			},
			expRetError: false,
		},
		{
			user:   "johnny@accelerator-apps.com",
			offset: 3,
			limit:  2,
			expEvents: []commons.Event{
				{
					Uid:      "ce664e65-4e31-4a83-b226-6d3602b87bf1",
					FullName: "Kanye West",
					Email:    "maria@accelerator-apps.com",
					Start:    "20200214T020436Z",
					End:      "20200215T020436Z",
					Summary:  "D Day",
					GeoLat:   89.35,
					GeoLon:   141.55,
				},
				{
					Uid:      "e426e402-7718-41dc-9dee-b9ee86a5d41a",
					FullName: "Sheldon Cooper",
					Email:    "Mickey.Mouse@gmail.com",
					Start:    "20230204T220436Z",
					End:      "20230205T220436Z",
					Summary:  "Fancy Party",
					GeoLat:   175.3,
					GeoLon:   124.75,
				},
			},
			expRetError: false,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			retEvents, retErr := FetchEvents(tt.user, tt.offset, tt.limit)
			if !reflect.DeepEqual(tt.expEvents, retEvents) {
				t.Errorf("returned events mismatch: expected '%v', received '%v'", tt.expEvents, retEvents)
			}
			if tt.expRetError != (retErr != nil) {
				t.Errorf("return error mismatch: expected is error '%v', received error '%v'", tt.expRetError, retErr)
			}
		})
	}
}
