package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

func updateStatus(fileName string) {
	for {
		status := Status{
			Water: rand.Intn(100) + 1,
			Wind:  rand.Intn(100) + 1,
		}

		jsonData, err := json.MarshalIndent(status, "", "  ")
		if err != nil {
			fmt.Println("Error encoding JSON:", err)
			continue
		}

		err = os.WriteFile(fileName, jsonData, 0644)
		if err != nil {
			fmt.Println("Error writing to file:", err)
		}

		fmt.Println("Status updated:", status)

		time.Sleep(15 * time.Second)
	}
}

func getStatusHandlerJSON(w http.ResponseWriter, r *http.Request) {
    fileName := "status.json"
    data, err := os.ReadFile(fileName)
    if err != nil {
        http.Error(w, "Error reading status file", http.StatusInternalServerError)
        return
    }

    var status Status
    err = json.Unmarshal(data, &status)
    if err != nil {
        http.Error(w, "Error decoding JSON", http.StatusInternalServerError)
        return
    }

    statusMessage := getStatusMessage(status)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(struct {
        Status Status `json:"status"`
        Message string `json:"message"`
    }{status, statusMessage})
}


func getStatusHandlerHTML(w http.ResponseWriter, r *http.Request) {
    htmlFile, err := os.ReadFile("index.html")
    if err != nil {
        http.Error(w, "Error reading HTML file", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "text/html")
    w.Write(htmlFile)
}




func getStatusMessage(status Status) string {
	var waterStatus, windStatus string

	if status.Water < 5 {
		waterStatus = "Aman"
	} else if status.Water >= 6 && status.Water <= 8 {
		waterStatus = "Siaga"
	} else {
		waterStatus = "Bahaya"
	}

	if status.Wind < 6 {
		windStatus = "Aman"
	} else if status.Wind >= 7 && status.Wind <= 15 {
		windStatus = "Siaga"
	} else {
		windStatus = "Bahaya"
	}

	return fmt.Sprintf("Status Air: %s, Status Angin: %s", waterStatus, windStatus)
}

func main() {
	go updateStatus("status.json")

	var port = ":8080"
	

	http.HandleFunc("/", getStatusHandlerHTML)
	http.HandleFunc("/status/json", getStatusHandlerJSON)
	fmt.Println("Server started on http://localhost" + port)
	http.ListenAndServe(port, nil)
}
