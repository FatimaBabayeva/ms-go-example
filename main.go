package main

import (
	"github.com/gorilla/mux"
	"github.com/jessevdk/go-flags"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"ms-go-example/handler"
	"ms-go-example/properties"
	"ms-go-example/repo"
	"net/http"
	"strconv"
)

var opts struct {
	Profile string `short:"p" long:"profile" default:"default" description:"Application run profile"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}

	initLogger()
	initEnvVars()
	properties.LoadConfig()
	applyLoggerLevel()

	log.Info("Application is starting with profile: ", opts.Profile)

	err = repo.MigrateDb()
	if err != nil {
		panic(err)
	}
	repo.InitDb()

	router := mux.NewRouter()
	handler.NewMessageHandler(router)
	handler.HandleHealthRequest(router)

	port := strconv.Itoa(properties.Props.Port)
	log.Info("Starting server at port: ", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func initEnvVars() {
	if godotenv.Load("profiles/default.env") != nil {
		log.Fatal("Error in loading environment variables from: profiles/default.env")
	} else {
		log.Info("Environment variables loaded from: profiles/default.env")
	}

	if opts.Profile != "default" {
		profileFileName := "profiles/" + opts.Profile + ".env"
		if godotenv.Overload(profileFileName) != nil {
			log.Fatal("Error in loading environment variables from: ", profileFileName)
		} else {
			log.Info("Environment variables overloaded from: ", profileFileName)
		}
	}
}

func initLogger() {
	log.SetLevel(log.InfoLevel)
	if opts.Profile == "default" {
		log.SetFormatter(&log.JSONFormatter{})
	}
}

func applyLoggerLevel() {
	loglevel, err := log.ParseLevel(properties.Props.LogLevel)
	if err != nil {
		loglevel = log.InfoLevel
	}
	log.SetLevel(loglevel)
}
