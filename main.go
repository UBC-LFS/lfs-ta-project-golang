package main

import (
	"encoding/json"
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

type course struct {
	courseName  string
	instructors []string
}

func pullTerms() []string {
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

	// body = result from api get request
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	// Converts data to an interface
	var academicTermData map[string]interface{}
	err = json.Unmarshal([]byte(string(body)), &academicTermData)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	pageItems := academicTermData["pageItems"].([]interface{})

	terms := make([]string, 0)
	// fmt.Println(result["pageItems"])
	for i := 0; i < len(pageItems); i++ {
		item := pageItems[i].(map[string]interface{})
		termName := item["academicPeriod"].(map[string]interface{})["academicPeriodName"]
		terms = append(terms, termName.(string))
	}

	return terms
}

func pullCourseSectionData(reference string) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	ClientID := os.Getenv("expClientID")
	ClientSecret := os.Getenv("expClientSecret")

	client := &http.Client{}
	URL := "https://sat.api.ubc.ca/academic-exp/v2/course-section-details?courseSectionId=" + reference
	// filter out for courses in a specific term

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

	// body = result from api get request
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	// Converts data to an interface
	var courseSectionData map[string]interface{}
	err = json.Unmarshal([]byte(string(body)), &courseSectionData)
	if err != nil {
		fmt.Println(err)
		return
	}

	section := courseSectionData["pageItems"].([]interface{})[0]
	teachingAssignments := section.(map[string]interface{})["teachingAssignments"].([]interface{})
	// for each person teaching
	for i := 0; i < len(teachingAssignments); i++ {
		// if person is instructor
		if teachingAssignments[i].(map[string]interface{})["assignableRole"].(map[string]interface{})["code"] == "Instructor Teaching" {
			identifiers := teachingAssignments[i].(map[string]interface{})["identifiers"].([]interface{})
			// for each instructor
			for j := 0; j < len(identifiers); j++ {
				worker := identifiers[i].(map[string]interface{})["worker"].(map[string]interface{})["personNames"].([]interface{})
				email := identifiers[i].(map[string]interface{})["worker"].(map[string]interface{})["communicationChannel"].(map[string]interface{})["emails"].([]interface{})

				// get the worker's preferred name
				workerUndefined := true
				firstName := ""
				lastName := ""

				// fmt.Println(worker)
				// fmt.Println(email)
				for k := 0; k < len(worker); k++ {
					if worker[k].(map[string]interface{})["nameType"] == "Preferred Name" {
						firstName = worker[k].(map[string]interface{})["givenName"].(string)
						lastName = worker[k].(map[string]interface{})["familyName"].(string)
						workerUndefined = false
						break
					}
				}

				if workerUndefined {
					firstName = worker[0].(map[string]interface{})["givenName"].(string)
					lastName = worker[0].(map[string]interface{})["familyName"].(string)
				}

				fmt.Println(firstName)
				fmt.Println(lastName)
				fmt.Println(email)
			}

		}
	}
}

func getDeptCourses(dept string, selectedTerm string, year string) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	ClientID := os.Getenv("expClientID")
	ClientSecret := os.Getenv("expClientSecret")

	client := &http.Client{}
	URL := "https://sat.api.ubc.ca/academic-exp/v2/course-section-details?academicYear=" + year + "&courseSubject=" + dept + "_V&page=1&pageSize=500"

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

	// body = result from api get request
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	// Converts data to an interface
	var academicRecordData map[string]interface{}
	err = json.Unmarshal([]byte(string(body)), &academicRecordData)
	if err != nil {
		fmt.Println(err)
		return
	}

	courseItems := academicRecordData["pageItems"].([]interface{})

	for i := 0; i < len(courseItems); i++ {
		item := courseItems[i].(map[string]interface{})
		termDetails := item["academicPeriod"].(map[string]interface{})
		term := termDetails["academicPeriodName"]

		// filter out for courses in a specific term
		if term == selectedTerm {
			fmt.Println(dept + " " + item["course"].(map[string]interface{})["courseNumber"].(string) + " " + item["sectionNumber"].(string))
			pullCourseSectionData(item["courseSectionId"].(string))

			// courseInfo := course{
			// 	courseName: dept + " " + item["course"].(map[string]interface{})["courseNumber"].(string) + " " + item["sectionNumber"].(string),
			// 	// instructors: ,
			// }
		}
	}
}

func pullCourses(selectedTerm string) {
	LFSDepts := [10]string{"APBI", "FNH", "FOOD", "FRE", "GRS", "HUNU", "LFS", "LWS", "PLNT", "SOIL"}
	// for loop of getDeptCourses
	year := string(selectedTerm[0:4])
	for i := 0; i < len(LFSDepts); i++ {
		getDeptCourses(LFSDepts[i], selectedTerm, year)
	}
}

func main() {
	// make this return the terms
	terms := pullTerms()
	// fmt.Println(terms)
	selectedTerm := terms[7]

	// pull course data
	// selectedTerm := "2024-25 Winter Term 1 (UBC-V)" // Testing
	pullCourses(selectedTerm)

	// run function to convert course data to .csv
}
