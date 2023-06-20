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

func send(client *messagebird.Client, sender string, recipient string, text string) {
	msg, err := sms.Create(
		client,
		sender,
		[]string{recipient},
		text,
		nil,
	)
	if err != nil {
		mbErr, _ := err.(messagebird.ErrorResponse)
	
		fmt.Println("Code:", mbErr.Errors[0].Code)
		fmt.Println("Description:", mbErr.Errors[0].Description)
		fmt.Println("Parameter:", mbErr.Errors[0].Parameter)
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

func template(message string, context gonja.Context) string {
	tpl, err := gonja.FromString(message)
	if err != nil {
		panic(err)
	}

	text, err := tpl.Execute(context)
	if err != nil {
		panic(err)
	}

	return text
}

func processRecords(fields []string, columns []string, phoneColumn string, message string) (string, string) {
	zipped_record := gonja.Context{}

	for index, field := range fields {
		zipped_record[columns[index]] = field
	}

	phone := zipped_record[phoneColumn]

	text := template(message, zipped_record)

	return phone.(string), text

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

	for {
		fields, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error reading CSV: %v", err)
		}
		phone, text := processRecords(fields, columns, phoneColumn, message)
		send(client, sender, phone, text)

	}
	file.Close()
}
