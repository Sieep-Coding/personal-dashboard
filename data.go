package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Entry struct {
	Date    time.Time `json:"date"`
	Study   float64   `json:"study"`
	Workout float64   `json:"workout"`
	Code    float64   `json:"code"`
	Game    float64   `json:"game"`
	Leisure float64   `json:"leisure"`
	Commits int       `json:"commits"`
	Jobs    int       `json:"jobs"`
	Workday bool      `json:"workday"`
}

func parseDate(dateStr string) (time.Time, error) {
	layout := "1/2/2006"
	return time.Parse(layout, dateStr)
}

func getData() ([]Entry, error) {
	file, err := os.Open("./data.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	var entries []Entry
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}
		date, err := parseDate(record[0])
		if err != nil {
			fmt.Println("Error parsing date:", err)
			continue
		}
		study, _ := strconv.ParseFloat(record[1], 64)
		workout, _ := strconv.ParseFloat(record[2], 64)
		code, _ := strconv.ParseFloat(record[3], 64)
		game, _ := strconv.ParseFloat(record[4], 64)
		leisure, _ := strconv.ParseFloat(record[5], 64)
		commits, _ := strconv.Atoi(record[6])
		jobs, _ := strconv.Atoi(record[7])

		entries = append(entries, Entry{
			Date:    date,
			Study:   study,
			Workout: workout,
			Code:    code,
			Game:    game,
			Leisure: leisure,
			Commits: commits,
			Jobs:    jobs,
		})
	}
	return entries, nil
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		return
	}
	entries, err := getData()
	if err != nil {
		http.Error(w, "Error fetching data:", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

func main() {
	http.HandleFunc("/data", dataHandler)
	fmt.Println("localhost:8080")
	http.ListenAndServe(":8080", nil)
}
