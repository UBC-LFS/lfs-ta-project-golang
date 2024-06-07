# LFS teaching assignment w/ Go

The script compiles a CSV file containing course details, instructor names, and their email addresses. When a specific term or session is chosen, the script retrieves the courses and the corresponding instructors with their email addresses for that selected term or session.

# Installation
1. Install Go: [Download](https://go.dev/doc/install)

2. Create a .env file
```
academicURL = "https://stg.api.ubc.ca/academic/v4/"
ClientID = ""
ClientSecret = ""

academicEXPURL = "https://stg.api.ubc.ca/academic-exp/v2/"
expClientID = ""
expClientSecret = ""
```

# How to use

1. Run the program using `go run main.go` OR build an executable
```go
 // to run the code
go run main.go

// to build an executable
go build <`go.mod` module name> // e.g. go build teaching-assignment
```
2. The program will ask you to enter a **year**

Example:
```
Enter a year [1996, 2032]: 2023
```

3. Open **courses.csv** to view your data