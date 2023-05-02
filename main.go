package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	messagebird "github.com/messagebird/go-rest-api"
	"github.com/messagebird/go-rest-api/sms"
	"github.com/noirbizarre/gonja"
)

func send(client *messagebird.Client, sender string, recipient string, context gonja.Context, message string) {
	tpl, err := gonja.FromString(message)
	if err != nil {
		panic(err)
	}

	text, err := tpl.Execute(context)
	if err != nil {
		panic(err)
	}

	msg, err := sms.Create(
		client,
		sender,
		[]string{recipient},
		text,
		nil,
	)
	if err != nil {
		log.Println(err)
	}
	// You can log the msg variable for development, or discard it by assigning it to `_`
	log.Println(msg)
}

type Configuration struct {
	Api         string   `env:"MESSAGEBIRD_API"`
	Message     string   `yaml:"message" env:"MESSAGEBIRD_MESSAGE"`
	Sender      string   `yaml:"sender" env:"MESSAGEBIRD_SENDER"`
	Columns     []string `yaml:"columns"`
	FileName    string   `yaml:"filename"`
	PhoneColumn string   `yaml:"phoneColumn"`
}

func getConfig() Configuration {
	var configuration Configuration

	err := cleanenv.ReadConfig("config.yml", &configuration)

	if err != nil {
		fmt.Println("error:", err)
	}
	return configuration
}

func processRecords(reader *csv.Reader, columns []string, phoneColumn string, message string, sender string, client *messagebird.Client) {
	for {
		fields, err := reader.Read()
		if err == io.EOF {
			break // End of file, break the loop
		}
		if err != nil {
			log.Fatalf("Error reading CSV: %v", err)
		}

		zipped_record := gonja.Context{}

		for index, field := range fields {
			zipped_record[columns[index]] = field
		}

		phone := zipped_record[phoneColumn]

		// Process the row (record) here
		send(client, sender, phone.(string), zipped_record, message)
	}
}

func main() {
	configuration := getConfig()
	api := configuration.Api
	message := configuration.Message
	sender := configuration.Sender
	phoneColumn := configuration.PhoneColumn
	columns := configuration.Columns

	client := messagebird.New(api)

	// Open the CSV file
	file, err := os.Open(configuration.FileName)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	// Create a CSV reader
	reader := csv.NewReader(file)

	processRecords(reader, columns, phoneColumn, message, sender, client)

	file.Close()
}
