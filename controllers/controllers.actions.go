// display.actions.go
// This module processes the actions that take the simulation through
// a circuit - Demand, Supply, Trade, Produce, Consume, Invest

package controllers

import (
	"gorilla-client/api"
	"gorilla-client/models"
	"gorilla-client/utils"
	"net/http"

	"github.com/gorilla/mux"
)

// Handles requests for the server to take an action comprising a stage
// of the circuit (demand,supply, trade, produce, invest), corresponding
// to a button press. This is specified by the URL parameter 'act'.
//
// Having requested the action from ths server, sets 'state' to the next
// stage of the circuit and redisplays whatever the user was looking at.
//
//	user.CurrentPage will be used to display errors if set
//	otherwise, a standard error page will be displayed
func ActionHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	user := CurrentUser(r)
	utils.TraceInfof(utils.Green, "Processing action for user %s", user.UserName)

	action, ok := mux.Vars(r)["action"]
	if !ok {
		ReportError(user, w, "Unrecognised URL")
		return
	}
	utils.TraceInfof(utils.Green, "User requested action %s", action)

	_, err = api.UserGetRequest(user.ApiKey, `/action/`+action)
	if err != nil {
		ReportError(user, w, "The server could not complete the action")
		return
	}

	// The action was taken. Advance the TimeStamp and the ViewedTimeStamp.
	// Create a new TableSet. Place the next fetched TableSet in the new
	// record, preserving the previous record.
	new_dataset := models.NewTableSet()

	// Append it to Datasets.
	// NOTE we are assuming it is appended as element user.TimeStamp+1
	// but as yet I haven't found documentation confirming this.
	user.TableSets = append(user.TableSets, &new_dataset)

	// Set the Comparator TimeStamp to compare with the effect of the previous action
	user.ComparatorTimeStamp = user.TimeStamp

	// Advance the TimeStamp to refer to the effect of this action.
	user.TimeStamp += 1

	// Reset viewed time stamp to point to the results of this action.
	user.ViewedTimeStamp = user.TimeStamp

	// Now refresh the data from the server
	err = api.FetchUserObjects(user)
	if err != nil {
		ReportError(user, w, "The server completed the action but did not send back any data.")
		return
	}

	// Set the state so that the simulation can proceed to the next action.
	user.Set_current_state(nextStates[action])
	utils.TraceInfof(utils.Green, "The last page this user visited was %v ", user.CurrentPage)
	if useLastVisited(user.CurrentPage) {
		Tpl.ExecuteTemplate(w, user.CurrentPage, user.TemplateData(""))
	} else {
		Tpl.ExecuteTemplate(w, "user-dashboard.html", user.TemplateData("")) //TODO this should be the index page, but it isn't ready yet
	}
}

// Display the previous state of the simulation
// Do nothing if we are already at the earliest stage
func Back(w http.ResponseWriter, r *http.Request) {
	utils.TraceInfo(utils.Green, "Back was requested")
	user := CurrentUser(r)
	if user.ViewedTimeStamp > 0 {
		user.ViewedTimeStamp--
	}
	if user.ComparatorTimeStamp > 0 {
		user.ComparatorTimeStamp--
	}
	utils.TraceInfof(utils.Green, "Viewing %d with comparator %d", user.ViewedTimeStamp, user.ComparatorTimeStamp)
	if useLastVisited(user.CurrentPage) {
		Tpl.ExecuteTemplate(w, user.CurrentPage, user.TemplateData(""))
	} else {
		Tpl.ExecuteTemplate(w, "index.html", user.TemplateData(""))
	}
}

// Display the next state of the simulation
// Do nothing if we are already viewing the most recent state
// Ensure the comparator stamp is one step behind the view stamp
func Forward(w http.ResponseWriter, r *http.Request) {
	utils.TraceInfo(utils.Green, "Forward was requested")
	user := CurrentUser(r)

	if user.ViewedTimeStamp < user.TimeStamp {
		user.ViewedTimeStamp++
	}
	if user.ComparatorTimeStamp != 0 {
		user.ComparatorTimeStamp++
	}

	utils.TraceInfof(utils.Green, "Viewing %d with comparator %d", user.ViewedTimeStamp, user.ComparatorTimeStamp)
	if useLastVisited(user.CurrentPage) {
		Tpl.ExecuteTemplate(w, user.CurrentPage, user.TemplateData(""))
	} else {
		Tpl.ExecuteTemplate(w, "index.html", user.TemplateData(""))
	}
}
