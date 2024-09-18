package main

import (
	"database/sql"
	"os"
	"strconv"
	"time"

	_ "modernc.org/sqlite"

	"github.com/gin-gonic/gin"
)

// Defining JSON body for uploadMedia(). It requires 1 JSON key, eventID.
type UploadMediaParameters struct {
	EventID string `form:"eventID" binding:"required"`
}

// Uploads media to an event
func uploadMedia(c *gin.Context) {

	// Variables for DB and Error
	var db *sql.DB
	var err error

	// Creating an instance of the struct, UploadMediaParameters
	var uploadMediaParameters UploadMediaParameters

	// Bind to the struct's members. If any member is invalid, binding does not happen and an error will be returned. Then its rejected with 400
	if c.Bind(&uploadMediaParameters) != nil {
		c.JSON(400, gin.H{"status": "Incorrect parameters, please provide all required parameters"})
		return
	}

	// Connect to the DB. If there is any issue connecting to the DB, throw a 500 error and return
	db, err = sql.Open("sqlite", "./EVENTPLANNER.db")
	if err != nil {
		c.JSON(500, gin.H{"status": "Could not connect to DB"})
		return
	}
	defer db.Close()

	// Check if the EventID exists in the DB by querying for the ID
	// Result is scanned into the variable, checkEvent
	queryToCheckExistingEvent := `SELECT ID, MEDIA FROM EVENTPLANNER WHERE ID=$1;`
	result := db.QueryRow(queryToCheckExistingEvent, uploadMediaParameters.EventID)

	// Defining a struct to hold the data queried by the query and scanning into it
	type CheckEvent struct {
		checkID       string
		mediaLocation string
	}

	var checkEvent CheckEvent
	result.Scan(&checkEvent.checkID, &checkEvent.mediaLocation)

	// If the length of checkEvent.checkID is greater than 0, means the query returned a result, so there is an event by that ID
	// So, upload media
	// Else, its rejected with a 404 as there is no event by that ID
	if len(checkEvent.checkID) > 0 {

		// Accept a multipart form upload for files
		form, _ := c.MultipartForm()
		files := form.File["upload[]"]

		// Loop through the uploaded files
		for _, file := range files {

			// Check the number of files in the folder and if its 5, stop upload. Currently it supports only 5 files
			noOfMedia, _ := os.ReadDir(checkEvent.mediaLocation)
			if len((noOfMedia)) > 4 {
				c.JSON(403, gin.H{"status": "Unable to upload new media, limit of 5 media reached"})
				return
			}

			// Save the uploaded file with current timestamp and file name in its name and sleep for 10 secs
			c.SaveUploadedFile(file, checkEvent.mediaLocation+"/"+strconv.Itoa(currentTimestamp())+"_"+file.Filename)
			time.Sleep(5 * time.Second)
		}
		c.JSON(200, gin.H{"status": "Media uploaded for Event ID, " + uploadMediaParameters.EventID})

	} else {
		c.JSON(404, gin.H{"status": "No Event with ID, " + uploadMediaParameters.EventID + " exists"})
	}

}
