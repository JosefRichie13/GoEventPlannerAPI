package main

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	request := gin.Default()
	request.MaxMultipartMemory = 8 << 20

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true // Allow all origins for development (consider origin restriction for production)
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	request.Use(cors.New(corsConfig))

	request.Static("/media", "./media")
	request.GET("/", landingPage)
	request.POST("/createAnEvent", createAnEvent)
	request.PUT("/startEvent", startEvent)
	request.PUT("/finishEvent", finishEvent)
	request.PUT("/resetEventDates", resetEventDates)
	request.PUT("/updateEvent", updateEvent)
	request.PUT("/uploadMedia", uploadMedia)
	request.GET("/getAllEvents", getAllEvents)
	request.GET("/getEventDetails", getEventDetails)
	request.GET("/getEventMedia", getEventMedia)
	request.Run(":8083")

}

// Landing page route
func landingPage(c *gin.Context) {
	c.JSON(200, gin.H{"status": "Welcome to Event Planner API, currently being built !!!"})
}

// Landing page route
func getEventMedia(c *gin.Context) {
	c.File("./media/Event.html")
}

// Defining JSON body for createAnEvent(). It requires 2 JSON key's eventName, eventDescription.
type CreateAnEventParameters struct {
	EventName        string `json:"eventName" binding:"required"`
	EventDescription string `json:"eventDescription" binding:"required"`
}

// Create an Event
func createAnEvent(c *gin.Context) {

	// Variables for DB and Error
	var db *sql.DB
	var err error

	// Creating an instance of the struct, CreateAnEventParameters
	var createAnEventParameters CreateAnEventParameters

	// Bind to the struct's members. If any member is invalid, binding does not happen and an error will be returned. Then its rejected with 400
	if c.Bind(&createAnEventParameters) != nil {
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

	// Generate Unique ID for Event and use the same for the event's media folder
	generatedID := uniqueIDGenerator()
	folderName := "./media/" + generatedID

	// Create a folder with the generated folder name. If there is any creating the folder, throw a 500 error and return
	err = os.Mkdir(folderName, 0755)
	if err != nil {
		c.JSON(500, gin.H{"status": "Error creating media folder"})
		return
	}

	// Add the event to the DB
	queryToAddAnEvent := `INSERT INTO EVENTPLANNER (ID, NAME, DESCRIPTION, ATENDEES, START_DATE, END_DATE, MEDIA) Values ($1, $2, $3, $4, $5, $6, $7)`
	db.QueryRow(queryToAddAnEvent, generatedID, sanitizeString(createAnEventParameters.EventName), sanitizeString(createAnEventParameters.EventDescription), 0, 0, 0, folderName)
	c.JSON(200, gin.H{"status": "Event Created", "eventID": generatedID})

}
