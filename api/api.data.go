// api.data.go
// DataObject is the intermediary between the client and the server.

package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorilla-client/models"
	"gorilla-client/utils"
	"strconv"
)

// Retrieves the data for a single table from the server.
// Unmarshals the server response into the DataList of the receiver
//
//	apiKey: sent to the server to identify and authorize the user
//	d: target of the data
//
//	Return: nil if it worked
//	Return: error string if there was an error
func Fetch(apiKey string, d *models.TableStruct) error {
	utils.TraceInfo(utils.BrightCyan, fmt.Sprintf("Fetching a table from server with api key %s and path %s", apiKey, d.ApiUrl))

	response, err := UserGetRequest(apiKey, d.ApiUrl)
	if err != nil {
		errorReport := fmt.Sprintf("ServerRequest produced the error %v", err)
		utils.TraceInfo(utils.Red, errorReport)
		return errors.New(errorReport)
	}

	if len(string(response)) == 0 {
		utils.TraceInfo(utils.BrightCyan, "INFORMATION: a server response to a fetch request was empty")
		return nil // no response is not an error, but don't process the result
	}

	// Populate the table
	jsonErr := json.Unmarshal(response, &d.Table)
	if jsonErr != nil {
		utils.TraceInfof(utils.Red, "Server response could not be unmarshalled because: %v", jsonErr)
		utils.TraceInfof(utils.Red, "The server response was %s\n", response)
		return errors.New("server response could not be unmarshalled")
	}
	return nil
}

// Iterates through ApiList to refresh all user objects for one user
//
//	Returns: false if any table fails.
//	Returns: true if all tables succeed.
func FetchUserObjects(user *models.User) error {
	var err error
	utils.TraceInfof(utils.BrightCyan, "Fetching details from server for user %s", user.UserName)

	// Fetch the simulation table
	// NOTE the api server identifies the user from the apiKey
	// This key is supplied by the api server
	err = Fetch(user.ApiKey, &user.Simulation)
	if err != nil {
		utils.TraceError("Simulations could not be fetched from the server")
		return errors.New("simulations could not be fetched from the server")
	}

	// Fetch all the user data in the tableset for the current stage of the current simulation
	// Reminder - a tableset contains all tables at one stage of the simulation.
	tableset := *user.TableSets[user.TimeStamp]
	for key, value := range tableset {
		err = Fetch(user.ApiKey, &value)
		if err != nil {
			utils.TraceErrorf("Could not retrieve server data for the new tableset with key %s", key)
		}
	}
	utils.TraceInfo(utils.BrightCyan, "Refresh complete")
	return nil
}

// Fetches a simulation and all associated tables from the server.
// Creates new objects that hold the results. The parameters uniquely
// identify the required data for the api server to process.
//
//	apikey: uniquely identifies the user
//	simulationID: uniquely identifies the simulation belonging to the user
//
//	returns:
//	 Pointer to a simulation of type TableStruct
//	 Pointer to a TableSet containing the Tables in this simulation
//	 err if anything goes wrong
func FetchSimulationAndTables(apiKey string, simulationID int) (*models.Simulation, *models.TableSet, error) {
	newSimulation := new(models.Simulation)
	newTableStruct := models.TableStruct{
		ApiUrl: `/simulations/by_id/` + strconv.Itoa(simulationID),
		Table:  newSimulation,
		Name:   `Simulation`,
	}

	err := Fetch(apiKey, &newTableStruct)
	if err != nil {
		utils.TraceError("Simulations could not be fetched from the server")
		return nil, nil, errors.New("simulations could not be fetched from the server")
	}

	//TODO NewTableSet should return a pointer
	newTableSet := models.NewTableSet()

	for key, value := range newTableSet {
		err = Fetch(apiKey, &value)
		if err != nil {
			utils.TraceErrorf("Could not retrieve server data for a new tableset with key %s", key)
		}
	}

	return newSimulation, &newTableSet, nil
}
