package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// When run without any options
// 		Returns csv list of courses, instructor names, instructor emails
// Script should have options to,
// 		list available terms from API
// 		supply a term to get assignment for that term
// 		Supply a list of unique instructor names and email in csv format

func pullTerms() {
	err := godotenv.Load(".env")
	if err != nil {
		// Handle error, e.g., log it or exit the program
		fmt.Println("Error loading .env file")
	}

	ClientID := os.Getenv("ClientID")
	ClientSecret := os.Getenv("ClientSecret")

	client := &http.Client{}
	URL := "https://stg.api.ubc.ca/academic/v4/academic-periods"
	req, err := http.NewRequest("GET", URL, nil)

	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("x-client-id", ClientID)
	req.Header.Add("x-client-secret", ClientSecret)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(body))
}

func main() {
	pullTerms()
}
