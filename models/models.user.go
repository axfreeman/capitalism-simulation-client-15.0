// PATH: go-auth/models/User.go

package models

import (
	"encoding/json"
)

// A User is a 'fullblown' player with everything involved in a simulation
type User struct {
	UserName            string         `json:"username"` // Repeats the key in the map,for ease of use
	Email               string         `json:"email"`
	ApiKey              string         `json:"api_key"` // The api key allocated to this user
	Password            string         `json:"password"`
	Role                string         `json:"role"`
	CurrentSimulationID int            `json:"current_simulation_id"` // the id of the simulation that this user is currently using
	CurrentPage         string         // Remember what the user was looking at (used when an action is requested)
	TableSets           []*TableSet    // Repository for the data objects generated during the simulation
	TimeStamp           int            // Indexes Datasets. Selects the stage that the simulation has reached
	ViewedTimeStamp     int            // Indexes Datasets. Selects what the user is viewing
	ComparatorTimeStamp int            // Indexes Datasets. Selects what Viewed items are compared with.
	Simulation          TableObject    // Details of the current simulation
	Simulations         []*TableObject // List of current simulations
	IsLocked            bool           `json:"is_locked"` // TODO REDUNDANT
}

// Constructor for a standard initial User.
func NewUser(username string) *User {
	new_user := User{
		UserName:            username,
		Password:            "",
		ApiKey:              "",
		CurrentSimulationID: 0,
		CurrentPage:         "",
		TimeStamp:           0,
		ViewedTimeStamp:     0,
		ComparatorTimeStamp: 0,
		TableSets:           []*TableSet{},
		IsLocked:            false,
		Simulation: TableObject{
			ApiUrl: `/simulations/current`,
			Table:  new([]Simulation),
		},
		Simulations: []*TableObject{},
	}
	new_dataset := NewTableSet()
	new_simulation := TableObject{
		ApiUrl: `/simulations/current`,
		Table:  new([]Simulation),
	}
	new_user.TableSets = append(new_user.TableSets, &new_dataset)
	new_user.Simulations = append(new_user.Simulations, &new_simulation)
	return &new_user
}

// A RegisteredUser is used for local authentication
// A User is a logged-in RegisteredUser
type RegisteredUser struct {
	UserName string
	ApiKey   string `json:"api_key"` // The api key will be retrieved from the server
	Password string // hashed
	Cookie   string // supplied by the register process TODO NOT USED DEPRECATE
}

// A RegisteredUserServerRequest is used to send a RegisteredUser to the server
type RegisteredUserServerRequest struct {
	UserName string `json:"username"` // Only this field is sent to the server, for security reasons
}

func (u RegisteredUser) Write() string {
	result, _ := json.MarshalIndent(u, " ", " ")
	return string(result)
}

// Convenience type with commonly-used objects, to pass into templates
type OutputData struct {
	Title           string
	Simulations     *[]Simulation
	Templates       *[]Simulation
	CommodityViews  *[]CommodityView
	IndustryViews   *[]IndustryView
	ClassViews      *[]ClassView
	Industry_Stocks *[]Industry_Stock
	Class_Stocks    *[]Class_Stock
	Trace           *[]Trace
	Count           int
	Username        string
	State           string
	Message         string
}

func (u *User) AsString() string {
	if u == nil {
		return "no such user"
	}
	s, _ := json.MarshalIndent(*u, " ", " ")
	return string(s)
}

// Defines a data object to be synchronised with the server
// ApiUrl is the endpoint on the server which fetches the Table
type TableObject struct {
	ApiUrl string
	Table  interface{} //All the data for one Table (eg Commodity, Industry, etc)
}

var ClientLoggedInUsers = make(map[string]*User) // Every user's simulation data

// Contains all the tables in one stage of one simulation
// Indexed by the name of the table (commodity, industry, etc)
type TableSet map[string]TableObject

// Constructor for a dataset object
// objects are "commodities", "industries", etc
// TODO apiKey seems redundant. Why doesn't the compiler object?
func NewTableSet() TableSet {
	return map[string]TableObject{
		"commodities": {
			ApiUrl: `/commodity`,
			Table:  new([]Commodity),
		},
		"industries": {
			ApiUrl: `/industry`,
			Table:  new([]Industry),
		},
		"classes": {
			ApiUrl: `/classes`,
			Table:  new([]Class),
		},
		"industry stocks": {
			ApiUrl: `/stocks/industry`,
			Table:  new([]Industry_Stock),
		},
		"class stocks": {
			ApiUrl: `/stocks/class`,
			Table:  new([]Class_Stock),
		},
		"trace": {
			ApiUrl: `/trace`,
			Table:  new([]Trace),
		},
	}
}

func NewRegisteredUser(username string, password string, apikey string) *RegisteredUser {
	new_RegisteredUser := RegisteredUser{
		UserName: username,
		Password: password,
		ApiKey:   apikey,
		Cookie:   "",
	}
	return &new_RegisteredUser
}

// Simple diagnostic printout of a user
func (u User) Print() string {
	result, _ := json.MarshalIndent(u, " ", " ")
	return string(result)
}
