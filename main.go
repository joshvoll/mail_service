package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/disturb16/kafkaclient"
	"sanservices.git.beanstalkapp.com/goproposal.git/internal/handlers"
	"sanservices.git.beanstalkapp.com/goproposal.git/internal/services"
	"sanservices.git.beanstalkapp.com/goproposal.git/settings"
)

func main() {
	//Get project configurations
	config, err := settings.GetConfiguration()

	if err != nil {
		log.Fatal(err)
	}

	if config.Port == "" {
		log.Fatal("No port was given")
	}

	// initialize service
	service, err := services.New(config.SMTPServer)
	if err != nil {
		log.Fatalf("Could not start smtp service... %v", err)
	}

	//Create main handler
	handler := handlers.New(service)

	// listen to mail_service topic
	client := kafkaclient.New(config.Kafka.Server, config.Kafka.Group, handler.OnMessageReceived)
	go client.ListenToTopics([]string{"mail_service"})

	//Start server
	fmt.Println("server running on port " + config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, handler.Router()))
}
