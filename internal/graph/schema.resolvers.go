package graph

import (
	"api/internal/commons"
	"api/internal/events"
	"api/internal/graph/generated"
	"api/internal/graph/model"
	"context"
	"log"
)

func (r *queryResolver) Events(_ context.Context, input model.Page) ([]*model.Event, error) {
	var parsedEvents []*model.Event
	storedEvents, err := events.FetchEvents(commons.User(input.UserEmail), commons.Offset(input.Offset), commons.Limit(input.Limit))
	if err != nil {
		log.Print(err)
		return nil, err
	}
	for _, storedEvent := range storedEvents {
		event := model.Event{
			UID:      string(storedEvent.Uid),
			FullName: string(storedEvent.FullName),
			Email:    string(storedEvent.Email),
			Start:    string(storedEvent.Start),
			End:      string(storedEvent.End),
			Summary:  string(storedEvent.Summary),
			GeoLat:   float64(storedEvent.GeoLat),
			GeoLon:   float64(storedEvent.GeoLon),
		}
		parsedEvents = append(parsedEvents, &event)
	}
	return parsedEvents, nil
}

func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
