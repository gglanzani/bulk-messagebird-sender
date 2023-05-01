package main

import (
   "fmt"
   "os"
   "encoding/csv"
   "encoding/json"
   "io"
   "log"

   "github.com/messagebird/go-rest-api"
   "github.com/messagebird/go-rest-api/sms"

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
  Api       string `json: "api"`
  Message   string `json: "message"`
  Sender    string `json: "sender"`
}

func getConfig() Configuration {


  file, _ := os.Open("config.json")
  defer file.Close()
  decoder := json.NewDecoder(file)
  configuration := Configuration{}
  err := decoder.Decode(&configuration)
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

