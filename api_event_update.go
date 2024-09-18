package main

import (
	"database/sql"

	_ "modernc.org/sqlite"

	"github.com/gin-gonic/gin"
)

// Defining JSON body for startEvent(). It requires 2 JSON key's, eventID. startDate.
type StartEventParameters struct {
	EventID   string `json:"eventID" binding:"required"`
	StartDate string `json:"startDate" binding:"required"`
}

func startEvent(c *gin.Context) {

	// Variables for DB and Error
	var db *sql.DB
	var err error

	// Creating an instance of the struct, UpdateBookDetailsParameters
	var startEventParameters StartEventParameters

	// Bind to the struct's members. If any member is invalid, binding does not happen and an error will be returned. Then its rejected with 400
	if c.BindJSON(&startEventParameters) != nil {
		c.JSON(400, gin.H{"status": "Incorrect parameters, please provide all required parameters"})
		return
	}

	// Checks if the supplied date is in DD-MMM-YYYY format
	if !checkDateFormat(startEventParameters.StartDate) {
		c.JSON(400, gin.H{"status": "Incorrect Date format, Date should be in DD-MMM-YYYY format, e.g., 27-Aug-2024"})
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
	// Result is scanned into the variable, checkResult
	queryToCheckExistingBook := `SELECT ID FROM EVENTPLANNER WHERE ID=$1;`
	result := db.QueryRow(queryToCheckExistingBook, startEventParameters.EventID)
	var checkResult string
	result.Scan(&checkResult)

	// If the length of checkResult is greater than 0, means the query returned a result, so there is an event by that ID
	// So, update its start date
	// Else, its rejected with a 404 as there is no event by that ID
	if len(checkResult) > 0 {

		// Update event start date
		queryToUpdateAnEvent := `UPDATE EVENTPLANNER SET START_DATE = $1 WHERE ID = $2;`
		db.QueryRow(queryToUpdateAnEvent, convertDateToEpoch(startEventParameters.StartDate), startEventParameters.EventID)
		c.JSON(200, gin.H{"status": "Event, " + startEventParameters.EventID + " started."})

	} else {
		c.JSON(404, gin.H{"status": "No Event with ID, " + startEventParameters.EventID + " exists"})
	}

}
