// PATH: go-auth/models/User.go

package models

import (
	"encoding/json"
)

// A User is a 'fullblown' player with everything involved in a simulation
type User struct {
	UserName            string             `json:"username"` // Repeats the key in the map,for ease of use
	Email               string             `json:"email"`
	ApiKey              string             `json:"api_key"` // The api key allocated to this user
	Password            string             `json:"password"`
	Role                string             `json:"role"`
	CurrentSimulationID int                `json:"current_simulation_id"` // the id of the simulation that this user is currently using
	CurrentPage         string             // Remember what the user was looking at (used when an action is requested)
	TableSets           []*TableSet        // Repository for the data objects generated during the simulation
	TableRepositories   []*TableRepository // List of this user's repositories, indexed by the Simulation ID
	TimeStamp           int                // Indexes Datasets. Selects the stage that the simulation has reached
	ViewedTimeStamp     int                // Indexes Datasets. Selects what the user is viewing
	ComparatorTimeStamp int                // Indexes Datasets. Selects what Viewed items are compared with.
	Simulation          TableStruct        // Details of all simulations
	IsLocked            bool               `json:"is_locked"` // TODO REDUNDANT
}

type TableRepository struct {
	SimulationID int         // ID of the simulation that this repository represents
	TableSets    []*TableSet // Repository for the data objects generated during one simulation
}

// Constructor for a standard initial User.
func NewUser(username string) *User {
	newUser := User{
		UserName:            username,
		Password:            "",
		ApiKey:              "",
		CurrentSimulationID: 0,
		CurrentPage:         "",
		TimeStamp:           0,
		ViewedTimeStamp:     0,
		ComparatorTimeStamp: 0,
		TableSets:           []*TableSet{},
		TableRepositories:   []*TableRepository{},
		IsLocked:            false,
		Simulation: TableStruct{
			ApiUrl: `/simulations`,
			Table:  new([]Simulation),
			Name:   "Simulations",
		},
	}
	// newTableSet := NewTableSet()
	// newUser.TableSets = append(newUser.TableSets, &newTableSet)
	// newTableRepository := TableRepository{}
	// newUser.TableRepositories = append(newUser.TableRepositories, &newTableRepository)
	return &newUser
}

// List of LoggedInUsers
var ClientLoggedInUsers = make(map[string]*User) // Every user's simulation data

// A RegisteredUser is used for local authentication
// A User is a logged-in RegisteredUser
type RegisteredUser struct {
	UserName string
	ApiKey   string `json:"api_key"` // The api key will be retrieved from the server
	Password string // hashed
	Cookie   string // TODO NOT USED DEPRECATE
}

// A RegisteredUserServerRequest is used to send a RegisteredUser to the server
type RegisteredUserServerRequest struct {
	UserName string `json:"username"` // Only this field is sent to the server, for security reasons
}

func (u RegisteredUser) Write() string {
	result, _ := json.MarshalIndent(u, " ", " ")
	return string(result)
}

// Commonly-used Views and Tables, to pass into templates
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

// Defines Table to be synchronised with the server
//
//	ApiUrl:the endpoint on the server which fetches the Table
//	Table: one of Commodity, Industry, Class, etc etc
//	Name: convenience field for diagnostics
type TableStruct struct {
	ApiUrl string      //Url to use when requesting data from the server
	Table  interface{} //All the data for one Table (eg Commodity, Industry, etc)
	Name   string      //The name of the table (for convenience, may be redundant)
}

// Contains all the tables in one stage of one simulation
// Indexed by the name of the table (commodity, industry, etc)
type TableSet map[string]TableStruct

// Constructor for a TableSet. This contains all the Tables in one stage
// required for one stage of one simulation. Tables are "commodities",
// "industries", etc
func NewTableSet() TableSet {
	return map[string]TableStruct{
		"commodities": {
			ApiUrl: `/commodity`,
			Table:  new([]Commodity),
			Name:   `Commodity`,
		},
		"industries": {
			ApiUrl: `/industry`,
			Table:  new([]Industry),
			Name:   `Industry`,
		},
		"classes": {
			ApiUrl: `/classes`,
			Table:  new([]Class),
			Name:   `Class`,
		},
		"industry stocks": {
			ApiUrl: `/stocks/industry`,
			Table:  new([]Industry_Stock),
			Name:   `Industry_Stock`,
		},
		"class stocks": {
			ApiUrl: `/stocks/class`,
			Table:  new([]Class_Stock),
			Name:   `Class_Stock`,
		},
		// TODO this is very verbose. Restore it later
		// "trace": {
		// 	ApiUrl: `/trace`,
		// 	Table:  new([]Trace),
		// 	Name:   `Trace`,
		// },
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
