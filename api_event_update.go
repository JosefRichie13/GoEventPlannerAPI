package main

import (
	"database/sql"

	_ "modernc.org/sqlite"

	"github.com/gin-gonic/gin"
)

// Defining JSON body for startEvent(). It requires 2 JSON key's, eventID, startDate.
type StartEventParameters struct {
	EventID   string `json:"eventID" binding:"required"`
	StartDate string `json:"startDate" binding:"required"`
}

// Starts an unfinished Event
func startEvent(c *gin.Context) {

	// Variables for DB and Error
	var db *sql.DB
	var err error

	// Creating an instance of the struct, StartEventParameters
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
	queryToCheckExistingBook := `SELECT ID, START_DATE FROM EVENTPLANNER WHERE ID=$1;`
	result := db.QueryRow(queryToCheckExistingBook, startEventParameters.EventID)

	// Defining a struct to hold the data queried by the query and scanning into it
	type CheckEvent struct {
		checkID        string
		checkStartDate int
	}

	var checkEvent CheckEvent
	result.Scan(&checkEvent.checkID, &checkEvent.checkStartDate)

	// If the length of checkResult is greater than 0, means the query returned a result, so there is an event by that ID
	// So, update its start date
	// Else, its rejected with a 404 as there is no event by that ID
	if len(checkEvent.checkID) > 0 {

		// Check if the event's start date is not 0, if yes, means event is already started, reject with 403
		if checkEvent.checkStartDate != 0 {
			c.JSON(403, gin.H{"status": "Event already started"})
			return
		}

		// Update event start date
		queryToUpdateAnEvent := `UPDATE EVENTPLANNER SET START_DATE = $1 WHERE ID = $2;`
		db.QueryRow(queryToUpdateAnEvent, convertDateToEpoch(startEventParameters.StartDate), startEventParameters.EventID)
		c.JSON(200, gin.H{"status": "Event, " + startEventParameters.EventID + " started."})

	} else {
		c.JSON(404, gin.H{"status": "No Event with ID, " + startEventParameters.EventID + " exists"})
	}

}

// Defining JSON body for finishEvent(). It requires 2 JSON key's, eventID, finishDate.
type FinishEventParameters struct {
	EventID    string `json:"eventID" binding:"required"`
	FinishDate string `json:"finishDate" binding:"required"`
}

// Finishes a started Event
func finishEvent(c *gin.Context) {

	// Variables for DB and Error
	var db *sql.DB
	var err error

	// Creating an instance of the struct, FinishEventParameters
	var finishEventParameters FinishEventParameters

	// Bind to the struct's members. If any member is invalid, binding does not happen and an error will be returned. Then its rejected with 400
	if c.BindJSON(&finishEventParameters) != nil {
		c.JSON(400, gin.H{"status": "Incorrect parameters, please provide all required parameters"})
		return
	}

	// Checks if the supplied date is in DD-MMM-YYYY format
	if !checkDateFormat(finishEventParameters.FinishDate) {
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

	// Check if the EventID exists in the DB by querying for the ID and START_DATE
	// Result is scanned into the variable, checkResult
	queryToCheckExistingBook := `SELECT ID, START_DATE, END_DATE FROM EVENTPLANNER WHERE ID=$1;`
	result := db.QueryRow(queryToCheckExistingBook, finishEventParameters.EventID)

	// Defining a struct to hold the data queried by the query and scanning into it
	type CheckEvent struct {
		checkID         string
		checkStartDate  int
		checkFinishDate int
	}

	var checkEvent CheckEvent
	result.Scan(&checkEvent.checkID, &checkEvent.checkStartDate, &checkEvent.checkFinishDate)

	// If the length of checkResult is greater than 0, means the query returned a result, so there is an event by that ID
	// So, update its finish date
	// Else, its rejected with a 404 as there is no event by that ID
	if len(checkEvent.checkID) > 0 {

		// Check if the event is started (by checking if the start date is greater than 0), if no, reject with 403
		if checkEvent.checkStartDate <= 0 {
			c.JSON(403, gin.H{"status": "Event not started, cannot finish it"})
			return
		}

		// Check if the start date is greater than the provided finish date, if yes, reject with 403 as finish date cannot be greater than start date
		if checkEvent.checkStartDate > convertDateToEpoch(finishEventParameters.FinishDate) {
			c.JSON(403, gin.H{"status": "Finish date cannot be before start date"})
			return
		}

		// Check if the event's end date is not 0, if yes, means event is already finished, reject with 403
		if checkEvent.checkFinishDate != 0 {
			c.JSON(403, gin.H{"status": "Event already finished"})
			return
		}

		// Update event finish date
		queryToUpdateAnEvent := `UPDATE EVENTPLANNER SET END_DATE = $1 WHERE ID = $2;`
		db.QueryRow(queryToUpdateAnEvent, convertDateToEpoch(finishEventParameters.FinishDate), finishEventParameters.EventID)
		c.JSON(200, gin.H{"status": "Event, " + finishEventParameters.EventID + " finished."})

	} else {
		c.JSON(404, gin.H{"status": "No Event with ID, " + finishEventParameters.EventID + " exists"})
	}

}

// Defining JSON body for resetEventDates(). It requires 1 JSON key, eventID.
type ResetEventParameters struct {
	EventID string `json:"eventID" binding:"required"`
}

// Starts an unfinished Event
func resetEventDates(c *gin.Context) {

	// Variables for DB and Error
	var db *sql.DB
	var err error

	// Creating an instance of the struct, ResetEventParameters
	var resetEventParameters ResetEventParameters

	// Bind to the struct's members. If any member is invalid, binding does not happen and an error will be returned. Then its rejected with 400
	if c.BindJSON(&resetEventParameters) != nil {
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
	queryToCheckExistingBook := `SELECT ID FROM EVENTPLANNER WHERE ID=$1;`
	result := db.QueryRow(queryToCheckExistingBook, resetEventParameters.EventID)
	var checkResult string
	result.Scan(&checkResult)

	// If the length of checkResult is greater than 0, means the query returned a result, so there is an event by that ID
	// So, reset its start and end date
	// Else, its rejected with a 404 as there is no event by that ID
	if len(checkResult) > 0 {

		// Reset event start and end date
		queryToUpdateAnEvent := `UPDATE EVENTPLANNER SET START_DATE = $1, END_DATE = $2 WHERE ID = $3;`
		db.QueryRow(queryToUpdateAnEvent, 0, 0, resetEventParameters.EventID)
		c.JSON(200, gin.H{"status": "Event, " + resetEventParameters.EventID + " reset."})

	} else {
		c.JSON(404, gin.H{"status": "No Event with ID, " + resetEventParameters.EventID + " exists"})
	}

}
