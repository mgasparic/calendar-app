package events

import (
	"api/internal/commons"
	"bufio"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

func ContinuouslySynchronizeEvents(generatorServiceUrlFormat commons.GeneratorServiceUrlFormat, dbClient *sql.DB, pauseLength time.Duration, users []commons.User, eventsCount commons.EventsCount) {
	for {
		var wg sync.WaitGroup

		eventIds := make(map[commons.Uid]struct{})
		var mux sync.RWMutex
		rows, err := dbClient.Query("SELECT id FROM event")
		if err != nil {
			log.Print(err)
			continue
		}
		for rows.Next() {
			var id uuid.UUID
			err := rows.Scan(&id)
			if err != nil {
				log.Print(err)
				continue
			}
			eventIds[commons.Uid(id.String())] = struct{}{}
		}
		_ = rows.Close()

		for _, user := range users {
			var userId uuid.UUID
			err = dbClient.QueryRow("SELECT id FROM customer WHERE email=$1", user).Scan(&userId)
			if err != nil {
				if err == sql.ErrNoRows {
					userId = uuid.New()
					_, err := dbClient.Exec("INSERT INTO customer (id, email) VALUES ($1, $2)", userId.String(), user)
					if err != nil {
						log.Print(err)
						continue
					}
				} else {
					log.Print(err)
					continue
				}
			}

			for i := 0; i < int(math.Ceil(float64(eventsCount)/1000)); i++ {
				wg.Add(1)
				go func(u commons.User, o commons.Offset) {
					defer wg.Done()
					response, err := http.Get(fmt.Sprintf(string(generatorServiceUrlFormat), u, o, 1000))
					if err != nil {
						log.Print(err)
						return
					}
					if response.StatusCode != http.StatusOK {
						log.Printf("received %d status code", response.StatusCode)
						return
					}
					var (
						events      []commons.Event
						activeEvent commons.Event
					)
					scanner := bufio.NewScanner(response.Body)
					for scanner.Scan() {
						line := scanner.Text()
						if line == "END:VCALENDAR" {
							events = append(events, activeEvent)
							activeEvent = commons.Event{}
						} else if strings.HasPrefix(line, "UID:") {
							activeEvent.Uid = commons.Uid(line[4:])
						} else if strings.HasPrefix(line, "ORGANIZER;") {
							fields := strings.Split(line[13:], ":")
							if len(fields) != 3 {
								log.Printf("organizer has a wrong format: %s", line)
								continue
							}
							activeEvent.FullName = commons.FullName(fields[0])
							activeEvent.Email = commons.Email(fields[2])
						} else if strings.HasPrefix(line, "DTSTART:") {
							activeEvent.Start = commons.Start(line[8:])
						} else if strings.HasPrefix(line, "DTEND:") {
							activeEvent.End = commons.End(line[6:])
						} else if strings.HasPrefix(line, "SUMMARY:") {
							activeEvent.Summary = commons.Summary(line[8:])
						} else if strings.HasPrefix(line, "GEO:") {
							coordinates := strings.Split(line[4:], ";")
							if len(coordinates) != 2 {
								log.Printf("coordinates have a wrong format: %s", line)
								continue
							}
							geoLat, err := strconv.ParseFloat(coordinates[0], 32)
							if err != nil {
								log.Print(err)
								continue
							}
							activeEvent.GeoLat = commons.GeoLat(geoLat)
							geoLon, err := strconv.ParseFloat(coordinates[1], 32)
							if err != nil {
								log.Print(err)
								continue
							}
							activeEvent.GeoLon = commons.GeoLon(geoLon)
						}
					}
					for _, event := range events {
						storedEvent := commons.Event{Uid: event.Uid}
						err = dbClient.QueryRow("SELECT organizer_full_name, organizer_email, time_start, time_end, event_summary, geo_lat, geo_lon "+
							"FROM event WHERE id=$1", event.Uid).Scan(&storedEvent.FullName, &storedEvent.Email, &storedEvent.Start, &storedEvent.End, &storedEvent.Summary, &storedEvent.GeoLat, &storedEvent.GeoLon)
						if err != nil {
							if err == sql.ErrNoRows {
								_, err = dbClient.Exec("INSERT INTO event (id, owner_id, organizer_full_name, organizer_email, time_start, time_end, event_summary, geo_lat, geo_lon) "+
									"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", event.Uid, userId.String(), event.FullName, event.Email, event.Start, event.End, event.Summary, event.GeoLat, event.GeoLon)
								if err != nil {
									log.Print(err)
								}
							} else {
								log.Print(err)
							}
							continue
						}
						mux.Lock()
						delete(eventIds, event.Uid)
						mux.Unlock()
						if storedEvent != event {
							_, err = dbClient.Exec("UPDATE event SET organizer_full_name=$1, organizer_email=$2, time_start=$3, time_end=$4, event_summary=$5, geo_lat=$6, geo_lon=$7 "+
								"WHERE id=$8", event.FullName, event.Email, event.Start, event.End, event.Summary, event.GeoLat, event.GeoLon, event.Uid)
							if err != nil {
								log.Print(err)
							}
						}
					}
				}(user, commons.Offset(i*1000))
			}
		}
		wg.Wait()

		for eventId := range eventIds {
			_, err = dbClient.Exec("DELETE FROM event WHERE id=$1", eventId)
			if err != nil {
				log.Print(err)
			}
		}

		log.Print("synchronization complete, waiting for the next round")
		time.Sleep(pauseLength)
	}
}
