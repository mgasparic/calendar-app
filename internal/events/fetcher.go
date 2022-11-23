package events

import (
	"api/internal/commons"
	"database/sql"
)

var dbClient *sql.DB

func SetDbClient(client *sql.DB) {
	dbClient = client
}

func FetchEvents(user commons.User, offset commons.Offset, limit commons.Limit) ([]commons.Event, error) {
	var events []commons.Event
	rows, err := dbClient.Query("SELECT event.id, event.organizer_full_name, event.organizer_email, event.time_start, event.time_end, event.event_summary, event.geo_lat, event.geo_lon "+
		"FROM event INNER JOIN customer ON customer.id=event.owner_id WHERE customer.email=$1 ORDER BY id ASC OFFSET $2 LIMIT $3", user, offset, limit)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var event commons.Event
		err = rows.Scan(&event.Uid, &event.FullName, &event.Email, &event.Start, &event.End, &event.Summary, &event.GeoLat, &event.GeoLon)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	_ = rows.Close()

	return events, nil
}
