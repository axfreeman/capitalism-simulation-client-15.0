package models

import (
	"gorilla-client/utils"
)

// Commonly-used Views to pass into templates
type DisplayData struct {
	Title                 string
	Simulations           *[]Simulation
	Templates             *[]Simulation
	CommodityViews        *[]Viewer
	IndustryViews         *[]Viewer
	NewClassViews         *[]View
	NewIndustryStockViews *[]View
	NewClassStockViews    *[]View
	ClassViews            *[]OldClassViewer      //Deprecated Phase out
	ClassStocks           *[]OldClassStockViewer //Deprecated Phase out
	Trace                 *[]Trace
	Count                 int
	Username              string
	State                 string
	Message               string
}

// Supplies data to pass into Templates for display
//
//		u: a user
//
//		returns:
//	     if the user has no simulations, just the template list
//	     otherwise, the output data the users current simulation
func (u *User) CreateDisplayData(message string) DisplayData {
	slist := u.SimulationsList()
	state := u.GetCurrentState()

	if u.CurrentSimulationID == 0 {
		utils.TraceInfo(utils.BrightYellow, "User has no simulations")
		return DisplayData{
			Title:                 "No simulations",
			Simulations:           nil,
			Templates:             &TemplateList,
			Count:                 0,
			Username:              u.UserName,
			State:                 state,
			CommodityViews:        nil,
			IndustryViews:         nil,
			ClassViews:            nil,
			NewIndustryStockViews: nil,
			ClassStocks:           nil,
			NewClassStockViews:    nil,
			Trace:                 nil,
			Message:               message,
		}
	}
	utils.TraceInfof(utils.BrightYellow, "TemplateData is retrieving data for user %s with simulationID %d", u.UserName, u.CurrentSimulationID)

	// retrieve comparator and viewed records for all data objects
	// to prepare for entry into Views in the DisplayData object
	cv := (*u.TableSets[*u.GetViewedTimeStamp()])["commodities"].Table.(*[]Commodity)
	cc := (*u.TableSets[*u.GetComparatorTimeStamp()])["commodities"].Table.(*[]Commodity)
	iv := (*u.TableSets[*u.GetViewedTimeStamp()])["industries"].Table.(*[]Industry)
	ic := (*u.TableSets[*u.GetComparatorTimeStamp()])["industries"].Table.(*[]Industry)
	clv := (*u.TableSets[*u.GetViewedTimeStamp()])["classes"].Table.(*[]Class)
	clc := (*u.TableSets[*u.GetComparatorTimeStamp()])["classes"].Table.(*[]Class)
	isv := (*u.TableSets[*u.GetViewedTimeStamp()])["industry stocks"].Table.(*[]IndustryStock)
	isc := (*u.TableSets[*u.GetComparatorTimeStamp()])["industry stocks"].Table.(*[]IndustryStock)
	csv := (*u.TableSets[*u.GetViewedTimeStamp()])["class stocks"].Table.(*[]ClassStock)
	csc := (*u.TableSets[*u.GetComparatorTimeStamp()])["class stocks"].Table.(*[]ClassStock)

	// Create the DisplayData object
	return DisplayData{
		Title:                 "Hello",
		Simulations:           slist,
		Templates:             &TemplateList,
		Count:                 len(*slist),
		Username:              u.UserName,
		State:                 state,
		CommodityViews:        CommodityViews(cv, cc),
		IndustryViews:         IndustryViews(iv, ic),
		NewClassViews:         NewClassViews(clv, clc),
		NewIndustryStockViews: IndustryStockViews(isv, isc),
		NewClassStockViews:    NewClassStockViews(csv, csc),
		ClassViews:            OldClassViews(clv, clc),      // Depracated phase out
		ClassStocks:           OldClassStockViews(csv, csc), // Depracated phase out
		Trace:                 u.Traces(*u.GetViewedTimeStamp()),
		Message:               message,
	}
}

// Get a CommodityData to display a single commodity in the commodity.html template
//
//	u: the user
//	message: any message
//	id: the id of the commodity to display
//
//	returns: CommodityData which references this commodity, and embeds an OutputData
func (u User) CommodityDisplayData(message string, id int) CommodityData {
	return CommodityData{
		u.CreateDisplayData(message),
		*u.Commodity(id),
	}
}

// Get a ClassData to display a single social class in the class.html template
//
//	u: the user
//	message: any message
//	id: the id of the social class to display
//
//	returns: classData which references this class, and embeds an OutputData
func (u User) ClassDisplayData(message string, id int) ClassData {
	return ClassData{
		u.CreateDisplayData(message),
		*u.Class(id),
	}
}

// Get an IndustryData to display a single industry in the industry.html template
//
//	u: the user
//	message: any message
//	id: the id of the industry item to display
//
//	returns: industryData which references this industry, and embeds an OutputData
func (u User) IndustryDisplayData(message string, id int) IndustryData {
	return IndustryData{
		u.CreateDisplayData(message),
		*u.Industry(id),
	}
}
