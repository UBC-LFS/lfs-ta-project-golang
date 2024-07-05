package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Course struct {
	code          string
	firstName     string
	lastName      string
	workEmail     string
	personalEmail string
}

func checkError(err error, msg string) {
	if err != nil {
		fmt.Println(err)
		log.Fatal(msg)
		os.Exit(1)
	}
}

func getData(url string, path string, params string) []map[string]interface{} {
	items := make([]map[string]interface{}, 0)
	hasNextPage := true
	page := 1

	for hasNextPage {
		pg := strconv.Itoa(page)
		req, err := http.NewRequest("GET", os.Getenv(url)+path+"?pageSize=500&page="+pg+params, nil)
		checkError(err, "Error: NewRequest")

		req.Header.Add("x-client-id", os.Getenv("ClientID"))
		req.Header.Add("x-client-secret", os.Getenv("ClientSecret"))

		client := &http.Client{}
		res, err := client.Do(req)
		checkError(err, "Error retrieving terms")

		nextPage, err := strconv.ParseBool(res.Header["X-Next-Page"][0])
		checkError(err, "Error retrieving terms")

		body, err := io.ReadAll(res.Body)
		checkError(err, "Unable to read result from API")

		var data map[string]interface{}
		err = json.Unmarshal([]byte(string(body)), &data)
		checkError(err, "Cannot convert data to JSON")

		pageItems := data["pageItems"].([]interface{})
		for i := 0; i < len(pageItems); i++ {
			item := pageItems[i].(map[string]interface{})
			items = append(items, item)
		}
		hasNextPage = nextPage
		page++
	}

	return items
}

func getTerms() (map[string][]string, int, int) {
	pageItems := getData("academicURL", "academic-periods", "")
	minYear := 9999
	maxYear := 0
	items := make(map[string][]string)
	for i := 0; i < len(pageItems); i++ {
		academicPeriod := pageItems[i]["academicPeriod"]
		id := academicPeriod.(map[string]interface{})["academicPeriodId"].(string)
		name := academicPeriod.(map[string]interface{})["academicPeriodName"].(string)

		termYear, err := strconv.Atoi(name[0:4])
		checkError(err, "Failed to parse int")

		if termYear < minYear {
			minYear = termYear
		}
		if termYear > maxYear {
			maxYear = termYear
		}

		if strings.Contains(name, "UBC-V") {
			items[name[0:4]] = append(items[name[0:4]], id+"|"+name)
		}
	}
	return items, minYear, maxYear
}

func getCourses(year string) map[string][]map[string]interface{} {
	fmt.Println("\nStarted fetching all the courses through APIs.")

	subjects := [13]string{"AGEC", "AANB", "APBI", "AQUA", "FNH", "FOOD", "FRE", "GRS", "HUNU", "LFS", "LWS", "PLNT", "SOIL"}

	items := make(map[string][]map[string]interface{})
	for _, subject := range subjects {
		pageItems := getData("academicEXPURL", "course-section-details", "&academicYear="+year+"&courseSubject="+subject+"_V")
		fmt.Printf("Read %s =====> %d pageItems\n", subject, len(pageItems))

		for _, item := range pageItems {
			id := item["academicPeriod"].(map[string]interface{})["academicPeriodId"].(string)
			items[id] = append(items[id], item)
		}
	}

	return items
}

func getCourseInfo(items []map[string]interface{}) ([]Course, []Course) {
	VALID_TYPES := []string{"Lecture", "Research"}
	EXCEPTIONS := []string{"FNH_V 326", "FNH_V 426"}

	courses := make([]Course, 0)
	term2Courses := make([]Course, 0)

	for _, item := range items {
		courseSubject := item["course"].(map[string]interface{})["courseSubject"].(map[string]interface{})["code"].(string)
		courseNumber := item["course"].(map[string]interface{})["courseNumber"].(string)
		sectionNumber := item["sectionNumber"].(string)
		tempCourse := courseSubject + " " + courseNumber
		startDate := item["startDate"].(string)
		endDate := item["endDate"].(string)

		if slices.Contains(VALID_TYPES, item["courseComponent"].(map[string]interface{})["instructionalFormat"].(map[string]interface{})["code"].(string)) || slices.Contains(EXCEPTIONS, tempCourse) {
			tas := item["teachingAssignments"].([]interface{})
			if len(tas) > 0 {
				for _, ta := range tas {
					if ta.(map[string]interface{})["assignableRole"].(map[string]interface{})["code"] == "Instructor Teaching" {
						identifiers := ta.(map[string]interface{})["identifiers"].([]interface{})

						firstName := ""
						lastName := ""
						workEmail := ""
						personalEmail := ""
						for _, identifier := range identifiers {
							_, ok := identifier.(map[string]interface{})["worker"]
							if ok {
								personNames := identifier.(map[string]interface{})["worker"].(map[string]interface{})["personNames"].([]interface{})
								emails := identifier.(map[string]interface{})["worker"].(map[string]interface{})["communicationChannel"].(map[string]interface{})["emails"].([]interface{})

								if len(personNames) > 0 {
									firstName = personNames[0].(map[string]interface{})["givenName"].(string)
									lastName = personNames[0].(map[string]interface{})["familyName"].(string)
								}

								for _, email := range emails {
									if email.(map[string]interface{})["channelType"].(map[string]interface{})["code"].(string) == "Work" {
										workEmail = email.(map[string]interface{})["emailAddress"].(string)
									}
									if email.(map[string]interface{})["channelType"].(map[string]interface{})["code"].(string) == "Personal" {
										personalEmail = email.(map[string]interface{})["emailAddress"].(string)
									}
								}

								c := Course{
									// code:          strings.Replace(courseSubject, "_V", "", -1) + " " + courseNumber + " " + sectionNumber,
									code:          courseSubject + " " + courseNumber + " " + sectionNumber,
									firstName:     firstName,
									lastName:      lastName,
									workEmail:     workEmail,
									personalEmail: personalEmail,
								}

								courses = append(courses, c)

								if startDate[0:4] != endDate[0:4] {
									term2Courses = append(term2Courses, c)
								}
							}
						}
					}
				}
			}
		}
	}

	return courses, term2Courses
}

func saveCSV(dir string, term string, courses []Course) {
	sort.Slice(courses, func(i, j int) bool {
		return courses[i].code < courses[j].code
	})

	fileName := "/output/" + term + " - Courses.csv"
	file, err := os.Create(dir + fileName)
	checkError(err, "\nFailed creating file")

	writer := csv.NewWriter(file)

	data := [][]string{
		{"Course", "Instructor First Name", "Instructor Last Name", "Work Email", "Personal Email"},
	}

	for _, course := range courses {
		temp := []string{course.code, course.firstName, course.lastName, course.workEmail, course.personalEmail}
		data = append(data, temp)
	}

	for _, row := range data {
		_ = writer.Write(row)
	}

	writer.Flush()
	file.Close()

	fmt.Println("Data on courses saved as CSV!")
}

func main() {
	_, filename, _, working := runtime.Caller(0)
	if !working {
		fmt.Println("Could not find the path of this file")
		os.Exit(1)
	}

	dir := filepath.Dir(filename)
	err := godotenv.Load(dir + "/.env")
	checkError(err, "Error: loading .env file")

	allTerms, minYear, maxYear := getTerms()

	var year string
	fmt.Printf("Enter a year [%d, %d]: ", minYear, maxYear)
	fmt.Scan(&year)

	if len(allTerms) > 0 {
		terms, ok := allTerms[year]
		if ok {
			fmt.Printf("\nFound %d term(s). Please see the list below.\n", len(terms))
			slices.Sort(terms)
			for _, term := range terms {
				term_sp := strings.Split(term, "|")
				fmt.Println("-", term_sp[1])
			}

			courses := getCourses(year)
			for _, term := range terms {
				term_sp := strings.Split(term, "|")
				items, ok := courses[term_sp[0]]
				if ok {
					fmt.Printf("\nGet started writing for %s.\n", term_sp[1])
					courseInfo, term2CourseInfo := getCourseInfo(items)
					if len(courseInfo) > 0 {
						saveCSV(dir, term_sp[1], courseInfo)
					} else {
						fmt.Printf("\nThere are no courses in this term - %s.\n", term_sp[1])
					}
					if len(term2CourseInfo) > 0 {
						fmt.Printf("\nThere are %d Term 1+2 courses in this term - %s.\n", len(term2CourseInfo), term_sp[1])
						temp := strings.Split(term_sp[1], " ")
						saveCSV(dir, temp[0]+" "+temp[1]+" Term 1+2 (UBC-V)", term2CourseInfo)
					}
				} else {
					fmt.Printf("\nThis term - %s - does not exist in the Academic-exp courses.\n", term_sp[1])
				}
			}

		} else {
			fmt.Printf("%s term does not exist. Please try again.", year)
		}
	} else {
		fmt.Println("No academic terms found.")
	}
}
