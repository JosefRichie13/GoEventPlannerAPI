# Backend for an Event Planner App using Gin Gonic, SQLite and Go

This repo has the code for an Event Planner App Backend. <br><br>

The basic flow of the App is 

* User creates an event, receives an ID in return
* User can start the event
* User can finish only a started event
* User can upload media files to the event
* User can get event details
* User can delete event, event media or specific media <br><br>

The below REST API endpoints are exposed

* POST /createAnEvent -- Creates an Event

* PUT /startEvent -- Starts an Event
  
* PUT /finishEvent -- Finishes an Event
  
* PUT /resetEventDates -- Resets an Event's Start and Finish Dates
  
* PUT /updateEvent -- Updates an Event's details
  
* PUT /uploadMedia -- Uploads media to an Event
  
* GET /getAllEvents -- Returns all the available Event's details
  
* GET /getEventDetails -- Returns a specific Event's details
  
* GET /getEventDetailsPage -- Returns a specific Event's details in an HTML page
  
* DELETE /deleteEvent -- Deletes an Event
 
* DELETE /deleteEventMedia -- Deletes all the media of an Event

* DELETE /deleteSpecificEventMedia -- Deletes a specific media of an Event<br><br>

The entire suite of endpoints with payloads are available in this HAR, [GoEventPlannerAPI.har](GoEventPlannerAPI.har)
