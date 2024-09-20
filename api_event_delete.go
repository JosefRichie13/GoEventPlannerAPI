package main

import (
	"database/sql"
	"os"
	"strconv"

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

	// Creating an instance of the struct, DeleteEventParameters
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

// Deletes only an Event's Media
func deleteEventMedia(c *gin.Context) {

	// Variables for DB and Error
	var db *sql.DB
	var err error

	// Creating an instance of the struct, DeleteEventParameters
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
	// We delete that Event's media
	// Else, its rejected with a 404 as there is no event by that ID
	if len(getEventDetails.ID) > 0 {

		// Delete the media folder for the ID
		err := os.RemoveAll("./media/" + deleteEventParameters.EventID)
		if err != nil {
			c.JSON(200, gin.H{"status": "Could not delete the media for , " + deleteEventParameters.EventID})
			return
		}
		// Recreate the media folder for the ID
		err = os.Mkdir("./media/"+deleteEventParameters.EventID, 0755)
		if err != nil {
			c.JSON(500, gin.H{"status": "Error creating media folder"})
			return
		}

		c.JSON(200, gin.H{"status": "All Media for the event with ID, " + deleteEventParameters.EventID + " is deleted."})

	} else {
		c.JSON(404, gin.H{"status": "No Event by ID, " + deleteEventParameters.EventID + " exists."})
	}
}

// Defining JSON body for deleteSpecificEventMedia(). It requires 2 Query Parameters eventID, mediaNo .
// mediaNo is the index of the media, i.e. 2nd media in the folder etc
type DeleteSpecificEventMediaParameters struct {
	EventID string `form:"eventID" binding:"required"`
	MediaNo int    `form:"mediaNo" binding:"required,numeric,gt=0"`
}

// Deletes only a specific media of an event
func deleteSpecificEventMedia(c *gin.Context) {

	// Variables for DB and Error
	var db *sql.DB
	var err error

	// Creating an instance of the struct, DeleteSpecificEventMediaParameters
	var deleteSpecificEventMediaParameters DeleteSpecificEventMediaParameters

	// Bind to the struct's members. If any member is invalid, binding does not happen and an error will be returned. Then its rejected with 400
	if c.Bind(&deleteSpecificEventMediaParameters) != nil {
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
	result := db.QueryRow(queryToCheckExistingEvent, deleteSpecificEventMediaParameters.EventID)

	// Defining a struct to hold all the values from the Query result
	type GetEventDetails struct {
		ID string
	}

	// Creating an instance of the struct, GetEventDetails
	var getEventDetails GetEventDetails

	// Scan the query result into the struct's members
	result.Scan(&getEventDetails.ID)

	// If the length of getEventDetails.ID is greater than 0, means the query returned a result, so there is a Event by that ID
	// We delete the specific media of that event
	// Else, its rejected with a 404 as there is no event by that ID
	if len(getEventDetails.ID) > 0 {

		// Subtract the mediaNo by 1 to bring the mediaNo in line with index, cause index starts at 0
		fileIndex := deleteSpecificEventMediaParameters.MediaNo - 1

		// Check the number of files in the folder
		// If no of files is less that or equal to the provided mediaNo, we reject with a 404
		// As this is a out of index error
		noOfMedia, _ := os.ReadDir("./media/" + deleteSpecificEventMediaParameters.EventID)
		if len((noOfMedia)) <= fileIndex {
			c.JSON(404, gin.H{"status": "Media number " + strconv.Itoa(deleteSpecificEventMediaParameters.MediaNo) + " not available"})
			return
		}

		// Read all the files in the folder
		files, err := os.ReadDir("./media/" + deleteSpecificEventMediaParameters.EventID)
		if err != nil {
			c.JSON(500, gin.H{"status": "Could not read media for , " + deleteSpecificEventMediaParameters.EventID})
			return
		}

		// Remove the specific file by its index
		err = os.Remove("./media/" + deleteSpecificEventMediaParameters.EventID + "/" + files[fileIndex].Name())
		if err != nil {
			c.JSON(500, gin.H{"status": "Could not delete the media at the index for " + deleteSpecificEventMediaParameters.EventID})
			return
		}

		c.JSON(200, gin.H{"status": "Media number " + strconv.Itoa(deleteSpecificEventMediaParameters.MediaNo) + " for the event with ID, " + deleteSpecificEventMediaParameters.EventID + " is deleted."})

	} else {

		c.JSON(404, gin.H{"status": "No Event by ID, " + deleteSpecificEventMediaParameters.EventID + " exists."})

	}
}
