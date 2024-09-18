package main

import (
	"database/sql"
	"os"
	"strconv"

	_ "modernc.org/sqlite"

	"github.com/gin-gonic/gin"
)

// Returns all the Event Details
func getAllEvents(c *gin.Context) {

	// Variables for DB and Error
	var db *sql.DB
	var err error

	// Connect to the DB. If there is any issue connecting to the DB, throw a 500 error and return
	db, err = sql.Open("sqlite", "./EVENTPLANNER.db")
	if err != nil {
		c.JSON(500, gin.H{"status": "Could not connect to DB"})
		return
	}
	defer db.Close()

	// Query the DB and result is held into the variable, result
	queryToCheckAllEvent := `SELECT * FROM EVENTPLANNER;`
	result, error := db.Query(queryToCheckAllEvent)

	// If there's any error when querying, return it
	if error != nil {
		c.JSON(500, gin.H{"status": "Could not execute Query"})
		return
	}
	defer result.Close()

	// Defining a struct to hold the data queried by the query and scanning into it
	type GetEventDetails struct {
		ID            string   `json:"eventId"`
		Name          string   `json:"eventName"`
		Description   string   `json:"eventDescription"`
		Atendees      int      `json:"eventAttendees"`
		StartDate     string   `json:"eventStartDate"`
		FinishDate    string   `json:"eventFinishDate"`
		Media         string   `json:"eventMedia"`
		MediaLocation []string `json:"eventMediaLocation"`
	}

	// Creating a slice from the struct
	getEventDetails := []GetEventDetails{}

	// Creating a slice to hold the image URL's
	getEventMediaLocation := []string{}

	// Iterating over the results
	for result.Next() {

		//Creating a new struct
		GetEventDetails := GetEventDetails{}

		// Scan the results into the struct
		result.Scan(&GetEventDetails.ID, &GetEventDetails.Name, &GetEventDetails.Description, &GetEventDetails.Atendees, &GetEventDetails.StartDate, &GetEventDetails.FinishDate, &GetEventDetails.Media)

		//Converting Date which is in String to Integer and into DD-MMM-YYYY format
		dateStartConversion, errs := strconv.Atoi(GetEventDetails.StartDate)
		if errs != nil {
			c.JSON(500, gin.H{"status": "Error Processing Start Date"})
			return
		}

		//Converting Date which is in String to Integer and into DD-MMM-YYYY format
		dateFinishConversion, errs := strconv.Atoi(GetEventDetails.FinishDate)
		if errs != nil {
			c.JSON(500, gin.H{"status": "Error Processing Finish Date"})
			return
		}

		//Adding the converted dates
		GetEventDetails.StartDate = convertEpochToDate(dateStartConversion)
		GetEventDetails.FinishDate = convertEpochToDate(dateFinishConversion)

		// Open the Media Folder and get names of the files in the folder
		dir, err := os.Open(GetEventDetails.Media)
		if err != nil {
			c.JSON(500, gin.H{"status": "Unable to open media folder"})
			return
		}
		defer dir.Close()

		files, err := dir.Readdir(-1)
		if err != nil {
			c.JSON(500, gin.H{"status": "Unable to read media folder"})
			return
		}

		// Iterate over the files and append the file URL to the slice
		// Sample URL that is generated http://192.168.1.109:8083/media/0fb966f09d3a4a18b7ae65b57eb01e96/1726650904_QA_PFE's.PNG
		for _, file := range files {
			getEventMediaLocation = append(getEventMediaLocation, getLocalIP()+":8083"+GetEventDetails.Media[1:]+"/"+file.Name())
		}

		// Add the slice to the struct
		GetEventDetails.MediaLocation = getEventMediaLocation

		// Append to the slice
		getEventDetails = append(getEventDetails, GetEventDetails)

		// Clear the slice for next set of URL's
		getEventMediaLocation = nil
	}

	c.JSON(200, gin.H{"allEventDetails": getEventDetails})

}
