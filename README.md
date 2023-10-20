# LFS teaching assignment w/ Go

The script compiles a CSV file containing course details, instructor names, and their email addresses. When a specific term or session is chosen, the script retrieves the courses and the corresponding instructors with their email addresses for that selected term or session.

# Installation
1. Install Go: [Download](https://go.dev/doc/install)

2. Create a .env file
```
academicURL = "https://stg.api.ubc.ca/academic/v4/"
ClientID = ""
ClientSecret = ""

academicEXPURL = "https://sat.api.ubc.ca/academic-exp/v2/"
expClientID = ""
expClientSecret = ""
```

# How to use

1. Run the program using `go run main.go`
2. The program will ask you to select the session number you want to retrieve course information from

Example:
```
Please select the session number you would like to retrieve course information from:
[0] - 2025-26 Winter Session (UBC-V)
[1] - 2024-25 Winter Session (UBC-V)
[2] - 2029-30 Winter Session (UBC-O)
[3] - 2024-25 Winter Term 2 (UBC-V)
[4] - 2026-27  Winter Term 2 (UBC-O)
[5] - 2025-26 Winter Term 2 (UBC-V)
[6] - 2027-28  Winter Term 2 (UBC-O)
[7] - 2024-25 Winter Term 1 (UBC-V)
[8] - 2031 Summer Session (UBC-O)
[9] - 2025-26 Winter Term 2 (UBC-O)
Session number:
```
3. If you want to select `2024-25 Winter Term 1 (UBC-V)`, you would enter `7`

Note: Selecting `2024-25 Winter Session (UBC-V)` will get you data for both term 1 and term 2

4. Wait for the script to fetch data for you
5. Open `courses.csv` to view your data!

# Developer Notes

The API being used is incomplete. Only the 1st instructor is being added to the CSV file since no data retrieved from the API includes 2 or more instructors, so this could not be implemented & tested.