package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Data on a course
type course struct {
	courseName       string
	instructorsArray []instructor
}

// Contact info of an instructor
type instructor struct {
	instructorFirstName      string
	instructorLastName       string
	instructorWorkEmail      string
	instructorSecondaryEmail string
}

// Returns a list of all school terms
func pullTerms() []string {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	ClientID := os.Getenv("ClientID")
	ClientSecret := os.Getenv("ClientSecret")

	client := &http.Client{}
	URL := os.Getenv("academicURL") + "academic-periods"
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

	// Converts data to an interface for extracting
	var academicTermData map[string]interface{}
	err = json.Unmarshal([]byte(string(body)), &academicTermData)
	if err != nil {
		fmt.Println(err)
	}
	pageItems := academicTermData["pageItems"].([]interface{})

	// Initalizes a list of terms
	terms := make([]string, 0)

	// Adds terms to the list of terms
	for i := 0; i < len(pageItems); i++ {
		item := pageItems[i].(map[string]interface{})
		termName := item["academicPeriod"].(map[string]interface{})["academicPeriodName"]
		terms = append(terms, termName.(string))
	}

	return terms
}

// Pull instructor data on courses and returns the array
func pullCourseSectionData(reference string) []instructor {
	// Initalizes an array of instructors
	instructorArray := make([]instructor, 0)

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	ClientID := os.Getenv("expClientID")
	ClientSecret := os.Getenv("expClientSecret")

	client := &http.Client{}
	URL := os.Getenv("academicEXPURL") + "course-section-details?courseSectionId=" + reference

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
				worker := identifiers[j].(map[string]interface{})["worker"].(map[string]interface{})["personNames"].([]interface{})
				email := identifiers[j].(map[string]interface{})["worker"].(map[string]interface{})["communicationChannel"].(map[string]interface{})["emails"].([]interface{})

				// Initalizes variables
				workerUndefined := true
				firstName := ""
				lastName := ""
				workEmail := ""
				secondaryEmail := ""

				// Search for the instructor's preferred name
				for k := 0; k < len(worker); k++ {
					if worker[k].(map[string]interface{})["nameType"] == "Preferred Name" {
						firstName = worker[k].(map[string]interface{})["givenName"].(string)
						lastName = worker[k].(map[string]interface{})["familyName"].(string)
						workerUndefined = false
						break
					}
				}
				// If no preferred name is set, use their given name
				if workerUndefined {
					firstName = worker[0].(map[string]interface{})["givenName"].(string)
					lastName = worker[0].(map[string]interface{})["familyName"].(string)
				}

				// Iterates through the instructor's emails for a work and secondary email address
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

// retrieves a list of courses within the department specified
func getDeptCourses(dept string, selectedTerm string, year string) []course {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	ClientID := os.Getenv("expClientID")
	ClientSecret := os.Getenv("expClientSecret")

	client := &http.Client{}

	URL := os.Getenv("academicEXPURL") + "course-section-details?academicYear=" + year + "&courseSubject=" + dept + "_V&page=1&pageSize=500"

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

	// Converts data to an interface so we can extract it
	var academicRecordData map[string]interface{}
	err = json.Unmarshal([]byte(string(body)), &academicRecordData)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	courseItems := academicRecordData["pageItems"].([]interface{})

	// Initalizes a list of courses in the department
	var deptCourses = make([]course, 0)

	for i := 0; i < len(courseItems); i++ {
		item := courseItems[i].(map[string]interface{})
		termDetails := item["academicPeriod"].(map[string]interface{})
		term := termDetails["academicPeriodName"]

		// filter out for courses in the selected term
		if term == selectedTerm {
			// Retrieves a list of instructors from the course
			instructorArray := pullCourseSectionData(item["courseSectionId"].(string))

			courseInfo := course{
				courseName:       dept + " " + item["course"].(map[string]interface{})["courseNumber"].(string) + " " + item["sectionNumber"].(string),
				instructorsArray: instructorArray,
			}

			deptCourses = append(deptCourses, courseInfo)
		}
	}

	return deptCourses
}

// Pull courses from each department
func pullCourses(selectedTerm string) []course {
	// Initalizes an array to store all the courses
	var allCourses = make([]course, 0)

	LFSDepts := [10]string{"APBI", "FNH", "FOOD", "FRE", "GRS", "HUNU", "LFS", "LWS", "PLNT", "SOIL"}

	// For each department in the LFS, get their courses
	year := string(selectedTerm[0:4])
	for i := 0; i < len(LFSDepts); i++ {
		deptCourses := getDeptCourses(LFSDepts[i], selectedTerm, year)

		// for each course retrieved, add it to allCourses array
		for k := 0; k < len(deptCourses); k++ {
			allCourses = append(allCourses, deptCourses[k])
		}
	}

	return allCourses
}

func main() {
	var sessionIndex int

	// Pull all terms from the API
	terms := pullTerms()

	fmt.Println("Please select the session number you would like to retrieve course information from:")
	// Display all options for terms
	for option, term := range terms {
		fmt.Println("[" + fmt.Sprint(option) + "] - " + term)
	}

	for true {
		fmt.Println("Session number: ")
		fmt.Scan(&sessionIndex) // retrieves user's input
		if 0 <= sessionIndex && sessionIndex <= len(terms)-1 {
			fmt.Println("Fetching for data on courses, please wait...")
			break
		} else {
			fmt.Println("This session does not exist, please try again.")
		}
	}

	selectedTerm := terms[sessionIndex]

	csvFile, err := os.Create("courses.csv")

	if err != nil {
		fmt.Printf("\nFailed creating file: %s \n", err)
	}

	csvwriter := csv.NewWriter(csvFile)

	// Initalizes the CSV file headers
	csvData := [][]string{
		{"Course Code", "Instructor First Name", "Instructor Last Name", "Work Email", "Secondary Email"},
	}

	// If selected term has no term specified, run function twice with both terms
	selectedTerms := make([]string, 0)
	if !(strings.Contains(selectedTerm, "Term")) {
		selectedTerms = append(selectedTerms, strings.Replace(selectedTerm, "Session", "Term 1", 1))
		selectedTerms = append(selectedTerms, strings.Replace(selectedTerm, "Session", "Term 2", 1))
	} else {
		selectedTerms = append(selectedTerms, selectedTerm)
	}

	// For each selected term, pull courses data
	var largestInstructorCount int = 0
	for _, termSelected := range selectedTerms {
		allCoursesData := pullCourses(termSelected)

		for _, row := range allCoursesData {
			courseCSV := []string{
				row.courseName,
			}

			// for each instructor, add it to the courseCSV array
			for instructorCount, instructorData := range row.instructorsArray {
				courseCSV = append(courseCSV, instructorData.instructorFirstName)
				courseCSV = append(courseCSV, instructorData.instructorLastName)
				courseCSV = append(courseCSV, instructorData.instructorWorkEmail)
				courseCSV = append(courseCSV, instructorData.instructorSecondaryEmail)
				// updates the largestInstructorCount
				largestInstructorCount = max(largestInstructorCount, instructorCount)
			}

			csvData = append(csvData, courseCSV)
		}
	}

	for instructorCount := 0; instructorCount < largestInstructorCount; instructorCount++ {
		csvData[0] = append(csvData[0], "Instructor "+fmt.Sprint(instructorCount+2)+" First Name")
		csvData[0] = append(csvData[0], "Instructor "+fmt.Sprint(instructorCount+2)+" Last Name")
		csvData[0] = append(csvData[0], "Instructor "+fmt.Sprint(instructorCount+2)+" Work Email")
		csvData[0] = append(csvData[0], "Instructor "+fmt.Sprint(instructorCount+2)+" Secondary Email")
	}

	// Converts the rows generated to a CSV datasheet
	for _, row := range csvData {
		_ = csvwriter.Write(row)
	}

	csvwriter.Flush()
	csvFile.Close()

	fmt.Println("Data on courses saved!")
}
