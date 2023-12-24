package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const (
	baseUrl    = "https://leetcode.com"
	graphQlUrl = baseUrl + "/graphql"
)

var (
	problemPath string
	lastUpdated time.Time
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	http.HandleFunc("/", handler)
	log.Println("Server is listening on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("Error:", err)
	}

}

func handler(w http.ResponseWriter, r *http.Request) {
	y1, m1, d1 := time.Now().UTC().Date()
	y2, m2, d2 := lastUpdated.UTC().Date()
	if y1 != y2 || m1 != m2 || d1 != d2 {
		problemPath = getProblemPath()
		lastUpdated = time.Now()
	}
	http.Redirect(w, r, baseUrl+problemPath, http.StatusMovedPermanently)
}

type LeetCodeRes struct {
	Data struct {
		ActiveDailyCodingChallengeQuestion struct {
			Link string `json:"link"`
		} `json:"activeDailyCodingChallengeQuestion"`
	} `json:"data"`
}

func getProblemPath() string {
	payload := map[string]string{
		"query": `
		    query questionOfToday { activeDailyCodingChallengeQuestion { link }}
		`,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("There was an error marshaling the JSON instance %v", err)
	}
	req, err := http.NewRequest("POST", graphQlUrl, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: time.Second * 5}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("There was an error executing the request%v", err)
	}
	defer res.Body.Close()
	var leetCodeRes LeetCodeRes
	json.NewDecoder(res.Body).Decode(&leetCodeRes)
	problemPath = leetCodeRes.Data.ActiveDailyCodingChallengeQuestion.Link
	return problemPath
}
