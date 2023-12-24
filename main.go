package main

import (
	"bytes"
	_ "embed"
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
	//go:embed about.html
	about string

	problemPath string
	lastUpdated time.Time
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(about))
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
		log.Println("Updating problem path...")
		problemPath = getProblemPath()
		lastUpdated = time.Now()
		log.Println("Problem path updated to " + problemPath)
	}

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, max-age=0")
	http.Redirect(w, r, baseUrl+problemPath, http.StatusTemporaryRedirect)
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
