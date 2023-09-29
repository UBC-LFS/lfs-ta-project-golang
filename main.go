package main

import (
	"fmt"
	"log"
	"net/http"
)

// When run without any options
// 		Returns csv list of courses, instructor names, instructor emails
// Script should have options to,
// 		list available terms from API
// 		supply a term to get assignment for that term
// 		Supply a list of unique instructor names and email in csv format

func pullTerms() {
	URL := "https://sat.api.ubc.ca/academic-exp/v2"
	resp, err := http.Get(URL)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(resp)
}

func main() {
	fmt.Println("Hello, ")
	pullTerms()
}
