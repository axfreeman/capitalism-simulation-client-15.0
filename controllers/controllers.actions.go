// display.actions.go
// This module processes the actions that take the simulation through
// a circuit - Demand, Supply, Trade, Produce, Consume, Invest

package controllers

import (
	"gorilla-client/api"
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
//	user.CurrentPageDetail.Url will be used to display errors if set
//	otherwise, a standard error page will be displayed
func ActionHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var action string
	var ok bool

	user := CurrentUser(r)
	utils.TraceInfof(utils.Green, "Processing action for user %s", user.UserName)

	if action, ok = mux.Vars(r)["action"]; !ok {
		ReportError(user, w, "Poorly specified action in the URL")
		return
	}
	utils.TraceInfof(utils.Green, "User requested action %s", action)

	if _, err = api.UserGetRequest(user.ApiKey, `/action/`+action); err != nil {
		ReportError(user, w, "The server could not complete the action")
		return
	}

	// The action was taken. Advance the TimeStamp and the ViewedTimeStamp.
	// Create a new TableSet and Append it to Datasets.
	// Set the TimeStamps
	*user.GetComparatorTimeStamp() = *user.GetTimeStamp() // Temporary transitional
	*user.GetTimeStamp() += 1                             // Temporary transitional
	*user.GetViewedTimeStamp() = *user.GetTimeStamp()     // Temporary transitional

	// Now refresh the data from the server
	if err = api.FetchTables(user); err != nil {
		ReportError(user, w, "The server completed the action but did not send back any data.")
		return
	}

	utils.TraceInfof(utils.Green, "Fetched the tables")

	// Temporary diagnostics
	tableSetDiagnosticString := user.Write()
	utils.TraceLogf(utils.BrightGreen, "The user data is now \n%s", tableSetDiagnosticString)

	// Set the state so that the simulation can proceed to the next action.
	user.SetCurrentState(nextStates[action])
	utils.TraceInfof(utils.Green, "The last page this user visited was %v ", user.CurrentPage.Url)

	if useLastVisited(user.CurrentPage.Url) {
		Tpl.ExecuteTemplate(w, user.CurrentPage.Url, user.TemplateData(""))
	} else {
		Tpl.ExecuteTemplate(w, "user-dashboard.html", user.TemplateData(""))
	}
}

// Display the previous state of the simulation
// Do nothing if we are already at the earliest stage
func Back(w http.ResponseWriter, r *http.Request) {
	utils.TraceInfo(utils.Green, "Back was requested")
	u := CurrentUser(r)
	if *u.GetViewedTimeStamp() > 0 {
		*u.GetViewedTimeStamp()--
	}
	if *u.GetComparatorTimeStamp() > 0 {
		*u.GetComparatorTimeStamp()--
	}
	utils.TraceInfof(utils.Green, "Viewing %d with comparator %d", *u.GetViewedTimeStamp(), *u.GetComparatorTimeStamp())
	if useLastVisited(u.CurrentPage.Url) {
		Tpl.ExecuteTemplate(w, u.CurrentPage.Url, u.TemplateData(""))
	} else {
		Tpl.ExecuteTemplate(w, "index.html", u.TemplateData(""))
	}
}

// Display the next state of the simulation
// Do nothing if we are already viewing the most recent state
// Ensure the comparator stamp is one step behind the view stamp
func Forward(w http.ResponseWriter, r *http.Request) {
	utils.TraceInfo(utils.Green, "Forward was requested")
	u := CurrentUser(r)

	if *u.GetViewedTimeStamp() < *u.GetTimeStamp() {
		*u.GetViewedTimeStamp()++
	}
	if *u.GetComparatorTimeStamp() != 0 {
		*u.GetComparatorTimeStamp()++
	}

	utils.TraceInfof(utils.Green, "Viewing %d with comparator %d", *u.GetViewedTimeStamp(), *u.GetComparatorTimeStamp())
	if useLastVisited(u.CurrentPage.Url) {
		Tpl.ExecuteTemplate(w, u.CurrentPage.Url, u.TemplateData(""))
	} else {
		Tpl.ExecuteTemplate(w, "index.html", u.TemplateData(""))
	}
}

// TODO not working yet
func SwitchSimulation(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	Tpl.ExecuteTemplate(w, user.CurrentPage.Url, user.TemplateData("Sorry, Switching Simulations is not ready yet"))
}

// TODO not working yet
func DeleteSimulation(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	Tpl.ExecuteTemplate(w, user.CurrentPage.Url, user.TemplateData("Sorry, Deleting a Simulation is not ready yet"))

}

// TODO not working yet
func RestartSimulation(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	Tpl.ExecuteTemplate(w, user.CurrentPage.Url, user.TemplateData("Sorry, Restarting a Simulation is not ready yet"))
}
