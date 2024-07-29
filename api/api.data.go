// api.data.go
// DataObject is the intermediary between the client and the server.

package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorilla-client/models"
	"gorilla-client/utils"
)

// Retrieves the data for a single simulation object from the server.
// Unmarshals the server response into the DataList of the receiver
//
//	d: target of the data
//	Return: nil if it worked
//	Return: error string if there was an error
func Fetch(u *models.User, d *models.TableObject) error {
	utils.TraceInfo(utils.BrightCyan, fmt.Sprintf("Fetching data from server with api key %s and path %s", u.ApiKey, d.ApiUrl))

	response, err := UserGetRequest(u.ApiKey, d.ApiUrl)

	if err != nil {
		errorReport := fmt.Sprintf("ServerRequest produced the error %v", err)
		utils.TraceInfo(utils.Red, errorReport)
		return errors.New(errorReport)
	}

	if len(string(response)) == 0 {
		utils.TraceInfo(utils.BrightCyan, "INFORMATION: a server response to a fetch request was empty")
		return nil // no response is not an error, but don't process the result
	}

	// Populate the data object
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

	// DEPRECATED
	// Fetch the simulation object
	err = Fetch(user, &user.Simulation)
	if err != nil {
		utils.TraceError("Simulations could not be fetched from the server")
		return errors.New("simulations could not be fetched from the server")
	}

	// Fetch all the user data in the dataset for the current stage of the current simulation
	// Reminder: a dataset is a repository for all objects at one stage of the simulation.
	dataSet := *user.TableSets[user.TimeStamp]

	for key, value := range dataSet {
		err = Fetch(user, &value)
		if err != nil {
			utils.TraceErrorf("Could not retrieve server data for the new dataset with key %s", key)
		}
	}
	utils.TraceInfo(utils.BrightCyan, "Refresh complete")
	return nil
}
