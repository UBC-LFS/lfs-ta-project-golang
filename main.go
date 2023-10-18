package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"encoding/csv"

	"github.com/joho/godotenv"
)

// When run without any options
// 		Returns csv list of courses, instructor names, instructor emails
// Script should have options to,
// 		list available terms from API
// 		supply a term to get assignment for that term
// 		Supply a list of unique instructor names and email in csv format

type course struct {
	courseName               string
	instructorFirstName      string
	instructorLastName       string
	instructorWorkEmail      string
	instructorSecondaryEmail string
}

type instructor struct {
	instructorFirstName      string
	instructorLastName       string
	instructorWorkEmail      string
	instructorSecondaryEmail string
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

func pullCourseSectionData(reference string) []instructor {

	instructorArray := make([]instructor, 0)

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
		return nil
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
				workEmail := ""
				secondaryEmail := ""

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

				for k := 0; k < len(email); k++ {
					if email[k].(map[string]interface{})["channelType"].(map[string]interface{})["description"].(string) == "Work" {
						workEmail = email[k].(map[string]interface{})["emailAddress"].(string)
					} else {
						secondaryEmail = email[k].(map[string]interface{})["emailAddress"].(string)
					}
				}

				instructorData := instructor{
					instructorFirstName:      firstName,
					instructorLastName:       lastName,
					instructorWorkEmail:      workEmail,
					instructorSecondaryEmail: secondaryEmail,
				}
				instructorArray = append(instructorArray, instructorData)
			}

		}
	}

	return instructorArray
}

func getDeptCourses(dept string, selectedTerm string, year string) []course {
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
		return nil
	}

	courseItems := academicRecordData["pageItems"].([]interface{})

	var deptCourses = make([]course, 0)

	for i := 0; i < len(courseItems); i++ {
		item := courseItems[i].(map[string]interface{})
		termDetails := item["academicPeriod"].(map[string]interface{})
		term := termDetails["academicPeriodName"]

		// filter out for courses in a specific term
		if term == selectedTerm {
			// fmt.Println(dept + " " + item["course"].(map[string]interface{})["courseNumber"].(string) + " " + item["sectionNumber"].(string))
			instructorArray := pullCourseSectionData(item["courseSectionId"].(string))

			instructorFirstName := ""
			instructorLastName := ""
			instructorWorkEmail := ""
			instructorSecondaryEmail := ""

			// If data on instructors exist
			if len(instructorArray) > 0 {
				instructorData := instructorArray[0]
				instructorFirstName = instructorData.instructorFirstName
				instructorLastName = instructorData.instructorLastName
				instructorWorkEmail = instructorData.instructorWorkEmail
				instructorSecondaryEmail = instructorData.instructorSecondaryEmail
			}

			courseInfo := course{
				courseName:               dept + " " + item["course"].(map[string]interface{})["courseNumber"].(string) + " " + item["sectionNumber"].(string),
				instructorFirstName:      instructorFirstName,
				instructorLastName:       instructorLastName,
				instructorWorkEmail:      instructorWorkEmail,
				instructorSecondaryEmail: instructorSecondaryEmail,
			}

			deptCourses = append(deptCourses, courseInfo)
		}
	}

	return deptCourses
}

func pullCourses(selectedTerm string) []course {
	var allCourses = make([]course, 0)
	LFSDepts := [10]string{"APBI", "FNH", "FOOD", "FRE", "GRS", "HUNU", "LFS", "LWS", "PLNT", "SOIL"}
	// for loop of getDeptCourses
	year := string(selectedTerm[0:4])
	for i := 0; i < len(LFSDepts); i++ {
		deptCourses := getDeptCourses(LFSDepts[i], selectedTerm, year)

		for k := 0; k < len(deptCourses); k++ {
			allCourses = append(allCourses, deptCourses[k])
		}
	}

	return allCourses
}

func main() {
	// make this return the terms
	terms := pullTerms()
	// fmt.Println(terms)
	selectedTerm := terms[7]

	// pull course data
	// selectedTerm := "2024-25 Winter Term 1 (UBC-V)" // Testing
	allCoursesData := pullCourses(selectedTerm)

	fmt.Println(allCoursesData)
	// run function to convert course data to .csv

	csvFile, err := os.Create("courses.csv")

	if err != nil {
		fmt.Println("Failed creating file: %s", err)
	}

	csvwriter := csv.NewWriter(csvFile)

	csvData := [][]string{
		{"Course Code", "Instructor First Name", "Instructor Last Name", "Work Email", "Secondary Email"},
	}

	for _, row := range allCoursesData {
		courseCSV := []string{
			row.courseName,
			row.instructorFirstName,
			row.instructorLastName,
			row.instructorWorkEmail,
			row.instructorSecondaryEmail,
		}
		csvData = append(csvData, courseCSV)
	}

	// Converts the rows we generated to a CSV datasheet
	for _, row := range csvData {
		_ = csvwriter.Write(row)
	}

	csvwriter.Flush()
	csvFile.Close()
}
