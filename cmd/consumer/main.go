package main

import (
	"api/internal/commons"
	"api/internal/events"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"log"
	"time"
)

type EnvVars struct {
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

	events.ContinuouslySynchronizeEvents(commons.GeneratorServiceUrlFormat(envVars.GeneratorServiceUrlFormat), dbClient, time.Duration(envVars.Interval)*time.Minute, users, commons.EventsCount(envVars.EventsCount))
}
