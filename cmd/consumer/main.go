package main

import (
	"api/internal/commons"
	"api/internal/events"
	"api/internal/graph"
	"api/internal/graph/generated"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

type EnvVars struct {
	ServicePort               int    `envconfig:"SERVICE_PORT" default:"9000"`
	ServiceEnvironment        string `envconfig:"SERVICE_ENVIRONMENT" default:"local"`
	GeneratorServiceUrlFormat string `envconfig:"GENERATOR_SERVICE_URL_FORMAT"`
	Users                     string `envconfig:"USERS"`
	Interval                  int    `envconfig:"INTERVAL"`
	EventsCount               int    `envconfig:"EVENTS_COUNT"`
	DbHost                    string `envconfig:"DB_HOST"`
	DbPort                    int    `envconfig:"DB_PORT"`
	DbUser                    string `envconfig:"DB_USER"`
	DbPassword                string `envconfig:"DB_PASSWORD"`
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var envVars EnvVars
	err := envconfig.Process("", &envVars)
	if err != nil {
		log.Fatal(err)
	}

	var users []commons.User
	err = json.Unmarshal([]byte(envVars.Users), &users)
	if err != nil {
		log.Fatal(err)
	}

	dbClient, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", envVars.DbHost, envVars.DbPort, envVars.DbUser, envVars.DbPassword))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = dbClient.Close()
	}()

	for i := 0; i < 5; i++ {
		err := dbClient.Ping()
		if err != nil {
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
	err = dbClient.Ping()
	if err != nil {
		log.Fatal("db not reachable")
	}

	go events.ContinuouslySynchronizeEvents(commons.GeneratorServiceUrlFormat(envVars.GeneratorServiceUrlFormat), dbClient, time.Duration(envVars.Interval)*time.Minute, users, commons.EventsCount(envVars.EventsCount))

	events.SetDbClient(dbClient)
	server := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", server)
	log.Printf("connect to http://localhost:%d/ for GraphQL playground", envVars.ServicePort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", envVars.ServicePort), nil))
}
