package main

import (
   "fmt"
   "os"
   "encoding/csv"
   "io"
   "log"

   "github.com/messagebird/go-rest-api"
   "github.com/messagebird/go-rest-api/sms"
   "github.com/ilyakaznacheev/cleanenv"


)

func send(client *messagebird.Client, sender string, recipient string, first string, message string) {
  text := fmt.Sprintf(message, first)
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
  Api       string `env:"MESSAGEBIRD_API"`
  Message   string `yaml: "message" env:"MESSAGEBIRD_MESSAGE"`
  Sender    string `yaml: "sender" env:"MESSAGEBIRD_SENDER"`
}

func getConfig() Configuration {
  var configuration Configuration

  err := cleanenv.ReadConfig("config.yml", &configuration)

  if err != nil {
    fmt.Println("error:", err)
  }
  return configuration
}

func main(){
  configuration := getConfig()
  api := configuration.Api
  message := configuration.Message
  sender := configuration.Sender

  client := messagebird.New(api)

  
  csvFileName := "names.csv"

	// Open the CSV file
	file, err := os.Open(csvFileName)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

  for {
		record, err := reader.Read()
		if err == io.EOF {
			break // End of file, break the loop
		}
		if err != nil {
			log.Fatalf("Error reading CSV: %v", err)
		}

    phone := record[2]
    name := record[0]

		// Process the row (record) here
    send(client, sender, phone, name, message)
	}
}

