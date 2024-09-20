package main

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"

	"github.com/gin-gonic/gin"
)

// Defining JSON body for deleteEvent(). It requires 1 Query Parameter eventID.
type DeleteEventParameters struct {
	EventID string `form:"eventID" binding:"required"`
}

// Deletes an Event along with its media
func deleteEvent(c *gin.Context) {

	// Variables for DB and Error
	var db *sql.DB
	var err error

	// Creating an instance of the struct, DeleteBookDetailsParameters
	var deleteEventParameters DeleteEventParameters

	// Bind to the struct's members. If any member is invalid, binding does not happen and an error will be returned. Then its rejected with 400
	if c.Bind(&deleteEventParameters) != nil {
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

	// Query the DB and result is held into the variable, result
	queryToCheckExistingEvent := `SELECT ID FROM EVENTPLANNER WHERE ID = $1;`
	result := db.QueryRow(queryToCheckExistingEvent, deleteEventParameters.EventID)

	// Defining a struct to hold all the values from the Query result
	type GetEventDetails struct {
		ID string
	}

	// Creating an instance of the struct, GetEventDetails
	var getEventDetails GetEventDetails

	// Scan the query result into the struct's members
	result.Scan(&getEventDetails.ID)

	// If the length of getEventDetails.ID is greater than 0, means the query returned a result, so there is a Event by that ID
	// We delete that Event along with its media
	// Else, its rejected with a 404 as there is no event by that ID
	if len(getEventDetails.ID) > 0 {

		// Delete the DB entry for the ID
		queryToDeleteExistingEvent := `DELETE FROM EVENTPLANNER WHERE ID=$1;`
		db.QueryRow(queryToDeleteExistingEvent, deleteEventParameters.EventID)

		// Delete the media folder for the ID
		err := os.RemoveAll("./media/" + deleteEventParameters.EventID)
		if err != nil {
			c.JSON(200, gin.H{"status": "Could not delete the media for , " + deleteEventParameters.EventID})
			return
		}

		c.JSON(200, gin.H{"status": "Event with ID, " + deleteEventParameters.EventID + " and all its media are deleted."})

	} else {
		c.JSON(404, gin.H{"status": "No Event by ID, " + deleteEventParameters.EventID + " exists."})
	}
}
