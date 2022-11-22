package main

import (
	"api/internal/commons"
	"api/internal/events"
	"api/internal/handlers"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"html/template"
	"log"
	"os"
	"time"
)

type EnvVars struct {
	ServicePort          int    `envconfig:"SERVICE_PORT" default:"9000"`
	ServiceEnvironment   string `envconfig:"SERVICE_ENVIRONMENT" default:"local"`
	CalendarTemplatePath string `envconfig:"CALENDAR_TEMPLATE_PATH"`
	GeneratorConfigPath  string `envconfig:"GENERATOR_CONFIG_PATH"`
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var envVars EnvVars
	err := envconfig.Process("", &envVars)
	if err != nil {
		log.Fatal(err)
	}

	temp, err := template.ParseFiles(envVars.CalendarTemplatePath)
	if err != nil {
		log.Fatal(err)
	}
	ge := handlers.GeneratorEnvironment{IcsTemplate: temp}

	generatorConfigRaw, err := os.ReadFile(envVars.GeneratorConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	var generatorConfig commons.GeneratorConfig
	err = json.Unmarshal(generatorConfigRaw, &generatorConfig)
	if err != nil {
		log.Fatal(err)
	}

	events.GenerateInitialEvents(generatorConfig.Users, generatorConfig.EventsCount, generatorConfig.Options.FullNames, generatorConfig.Options.Emails, generatorConfig.Options.Summaries)
	go events.ContinuouslyUpdateEvents(time.Duration(generatorConfig.ModificationRateInSeconds)*time.Second, generatorConfig.ReplacementRate, generatorConfig.ReplacementRate)

	if envVars.ServiceEnvironment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.GET("/events/:user", ge.Generator)

	log.Fatal(router.Run(fmt.Sprintf(":%d", envVars.ServicePort)))
}
