package main

// use 'go get github.com/go-sql-driver/mysql' to get the library

// example of a post request from client to microservice:

// [POST] http://<hostname>:8081/api/v1/log/event

// {
// 	"serviceName" : "testService",
// 	"eventName" : "testEvent",
// 	"eventDateTime" : "2018-09-22T12:42:31Z",
// 	"payload" : "testPayload"
// }

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// struct to represent the json request from the client

type ClientRequest struct {
	ServiceName   string    `json:"serviceName"`
	EventName     string    `json:"eventName"`
	EventDateTime time.Time `json:"eventDateTime"`
	Payload       string    `json:"payload"`
}

func main() {

	log.Println("Starting Audit Service...")

	var err error
	db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/db_integral_audit")

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/api/v1/log/event", handleRequests)
	http.ListenAndServe(":8081", nil)

	defer db.Close()

	log.Println("Shutting down Audit Service...")

}

func handleRequests(writer http.ResponseWriter, reader *http.Request) {

	var clientRequest ClientRequest

	err := json.NewDecoder(reader.Body).Decode(&clientRequest)

	if err != nil {
		log.Println(err)
	}

	insert, err := db.Query("INSERT INTO tbl_audit (service_name, event_name, event_timestamp, payload) values ('" + clientRequest.ServiceName + "', '" + clientRequest.EventName + "', now() , '" + clientRequest.Payload + "')")

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Added event to event log")
	insert.Close()
}
