package controllers

import (
	"encoding/json"
	"fmt"
	"gorilla-client/api"
	"gorilla-client/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Structure to hold result of the server's response to a clone request
type CloneResult struct {
	Message       string `json:"message"`
	StatusCode    int    `json:"statusCode"`
	Simulation_id int    `json:"simulation_id"`
}

// Creates a new simulation for the user, from the template specified by the 'id' parameter.
// This can be scaled up when and if login is introduced.
func CreateSimulation(w http.ResponseWriter, r *http.Request) {
	var s string
	var ok bool
	var err error
	var body []byte

	user := CurrentUser(r)
	utils.TraceInfof(utils.Green, "Clone Simulation was called by user %s", user.UserName)
	user.CurrentPage = "user-dashboard.html"

	if s, ok = mux.Vars(r)["id"]; !ok {
		ReportError(user, w, "Unrecognisable URL. Please report this to the developer")
		return
	}

	requestedSimulation, _ := strconv.Atoi(s)
	utils.TraceInfof(utils.Green, "Request to clone simulation %d", requestedSimulation)

	// Ask the server to create the clone and tell us the simulation id
	// Do not load it yet
	if body, err = api.UserGetRequest(user.ApiKey, `/clone/`+s); err != nil {
		ReportError(user, w, fmt.Sprintf("There was a problem. Please report this to the developer%v", err))
		return
	}

	utils.TraceInfo(utils.Green, ("Server responded to clone request with the following message:"))
	utils.TraceInfo(utils.Green, ` `+string(body))

	// read the simulation id
	var result CloneResult
	if err = json.Unmarshal(body, &result); err != nil {
		ReportError(user, w, fmt.Sprintf("There was a problem with the server's response. Please report this to the developer%v", err))
		return
	}

	// Set the current simulation
	utils.TraceInfof(utils.Green, "Setting current simulation to be %d", result.Simulation_id)
	user.CurrentSimulationID = result.Simulation_id
	utils.TraceInfo(utils.Green, ("Setting current state to DEMAND"))
	user.Set_current_state("DEMAND")

	// Fetch the whole (new) TableSet from the server
	// (until now we only told the server to create it - now we want it)
	if err = api.FetchUserObjects(user); err != nil {
		ReportError(user, w, fmt.Sprintf("There was a problem. Please report this to the developer\n%v", err))
		return
	}

	// Initialise the timeStamp so that we are viewing the first TableSet.
	// As the user moves through the circuit, this timestamp will move forwards.
	// Each time we move forward, a new TableSet will be created.
	// This allows the user to view and compare with previous stages of the simulation.
	user.ViewedTimeStamp = 0
	Tpl.ExecuteTemplate(w, user.CurrentPage, user.TemplateData(""))
}
