# LFS teaching assignment w/ Go

# Installation
1. Install Go: [Download](https://go.dev/doc/install)
```

```

2. Create a .env file
```
academicURL = "https://stg.api.ubc.ca/academic/v4/"
ClientID = ""
ClientSecret = ""

academicEXPURL = "https://sat.api.ubc.ca/academic-exp/v2/"
expClientID = ""
expClientSecret = ""
```

# Developer Notes

The API being used is incomplete. Only the 1st instructor is being added to the CSV file since no data retrieved from the API includes 2 or more instructors, so this could not be implemented & tested.