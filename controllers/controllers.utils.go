// controllers.display.go

package controllers

import (
	"encoding/json"
	"gorilla-client/models"
	"gorilla-client/utils"

	"net/http"
)

// Display the data that is available for the user who made this call
// Fetch the data from the client local store, not from the server
func GetClientData(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	utils.TraceInfof(utils.Green, "GetClientData for user %s", user.UserName)
	data, _ := json.MarshalIndent(user, " ", " ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// Fetch and display a list of all users
// Do not display their data, This is handled by Fetch()
//
//	TODO should this be Premium ie reserved for admin user?
func DisplayUsers(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	user.Set_current_state("INVEST")
	Tpl.ExecuteTemplate(w, "users.html", user.TemplateData(""))
}

// helper function to obtain the state of the current simulation
// to be replaced by inline call
func Get_current_state(username string) string {
	return models.ClientLoggedInUsers[username].Get_current_state()
}

// Report an error by redisplaying the current template with an error message
//
//	user.CurrentPage must be set with the template name
//
//	user: the current user
//	w: the ResponseWriter to which the message should be sent
//	message: the error message
func ReportError(user *models.User, w http.ResponseWriter, message string) {
	t := user.TemplateData(message)
	utils.TraceError(t.Message)

	// use standard error page if no Current Page is set
	if len(user.CurrentPage) < 1 {
		user.CurrentPage = "errors.html"
	}
	Tpl.ExecuteTemplate(w, user.CurrentPage, t)
}

// The state which follows each action.
var nextStates = map[string]string{
	`demand`:  `SUPPLY`,
	`supply`:  `TRADE`,
	`trade`:   `PRODUCE`,
	`produce`: `CONSUME`,
	`consume`: `INVEST`,
	`invest`:  `DEMAND`,
}

// pages for which redirection is OK.
func useLastVisited(last string) bool {
	if last == "" {
		return false
	}
	switch last {
	case
		`commodities.html`,
		`industries.html`,
		`classes.html`,
		`industry_stocks.html`,
		`class_stocks.html`,
		`index.html`,
		`/`:
		return true
	}
	return false
}
