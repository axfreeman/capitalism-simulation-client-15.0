package models

import (
	"gorilla-client/utils"
)

// Commonly-used Views to pass into templates
type TemplateData struct {
	Title              string
	Simulations        *[]Simulation
	Templates          *[]Simulation
	CommodityViews     *[]Viewer
	IndustryViews      *[]Viewer
	ClassViews         *[]Viewer
	IndustryStockViews *[]Viewer
	ClassStockViews    *[]Viewer
	Trace              *[]Trace
	Count              int
	Username           string
	State              string
	Message            string
}

// Supplies data to pass into Templates for display
//
//		u: a user
//
//		returns:
//	     if the user has no simulations, just the template list
//	     otherwise, the output data the users current simulation
func (u *User) CreateDisplayData(message string) TemplateData {
	slist := u.SimulationsList()
	state := u.GetCurrentState()

	if u.CurrentSimulationID == 0 {
		utils.TraceInfo(utils.BrightYellow, "User has no simulations")
		return TemplateData{
			Title:              "No simulations",
			Simulations:        nil,
			Templates:          &TemplateList,
			Count:              0,
			Username:           u.UserName,
			State:              state,
			CommodityViews:     nil,
			IndustryViews:      nil,
			IndustryStockViews: nil,
			ClassStockViews:    nil,
			Trace:              nil,
			Message:            message,
		}
	}
	utils.TraceInfof(utils.BrightYellow, "TemplateData is retrieving data for user %s with simulationID %d", u.UserName, u.CurrentSimulationID)

	// retrieve comparator and viewed records for all data objects
	// to prepare for entry into Views in the DisplayData object

	// OLD CODE
	// cv := (*u.TableSets[*u.GetViewedTimeStamp()])["commodities"].Table.(*[]Commodity)
	// cc := (*u.TableSets[*u.GetComparatorTimeStamp()])["commodities"].Table.(*[]Commodity)
	// iv := (*u.TableSets[*u.GetViewedTimeStamp()])["industries"].Table.(*[]Industry)
	// ic := (*u.TableSets[*u.GetComparatorTimeStamp()])["industries"].Table.(*[]Industry)
	// clv := (*u.TableSets[*u.GetViewedTimeStamp()])["classes"].Table.(*[]Class)
	// clc := (*u.TableSets[*u.GetComparatorTimeStamp()])["classes"].Table.(*[]Class)
	// isv := (*u.TableSets[*u.GetViewedTimeStamp()])["industry stocks"].Table.(*[]IndustryStock)
	// isc := (*u.TableSets[*u.GetComparatorTimeStamp()])["industry stocks"].Table.(*[]IndustryStock)
	// csv := (*u.TableSets[*u.GetViewedTimeStamp()])["class stocks"].Table.(*[]ClassStock)
	// csc := (*u.TableSets[*u.GetComparatorTimeStamp()])["class stocks"].Table.(*[]ClassStock)

	cv := ViewedObjects[Commodity](*u, `commodities`)
	cc := ComparedObjects[Commodity](*u, `commodities`)
	iv := ViewedObjects[Industry](*u, `industries`)
	ic := ComparedObjects[Industry](*u, `industries`)
	clv := ViewedObjects[Class](*u, `classes`)
	clc := ComparedObjects[Class](*u, `classes`)
	isv := ViewedObjects[IndustryStock](*u, `industry stocks`)
	isc := ComparedObjects[IndustryStock](*u, `industry stocks`)
	csv := ViewedObjects[ClassStock](*u, `class stocks`)
	csc := ComparedObjects[ClassStock](*u, `class stocks`)

	// Create the DisplayData object
	templateData := TemplateData{
		Title:              "Hello",
		Simulations:        slist,
		Templates:          &TemplateList,
		Count:              len(*slist),
		Username:           u.UserName,
		State:              state,
		CommodityViews:     CommodityViews(cv, cc),
		IndustryViews:      IndustryViews(iv, ic),
		ClassViews:         ClassViews(clv, clc),
		IndustryStockViews: IndustryStockViews(isv, isc),
		ClassStockViews:    ClassStockViews(csv, csc),
		Trace:              u.Traces(*u.GetViewedTimeStamp()),
		Message:            message,
	}

	// classTest := (*displayData.ClassViews)[0]
	// industryTest := (*displayData.IndustryViews)[0]

	// fmt.Println("class test is ", classTest)
	// fmt.Println("industry test is ", industryTest)

	return templateData
}

// Create a CommodityData to display a single commodity in the
// commodity.html template. This is added dynamically to the DisplayData
// template when the Commodity view is requested
//
//	u: the user
//	message: any message
//	id: the id of the commodity to display
//
//	returns: CommodityData which references this commodity, and embeds an OutputData
func (u User) CommodityDisplayData(message string, id int) CommodityData {
	return CommodityData{
		u.CreateDisplayData(message),
		*ViewedObject[Commodity](u, `commodities`, id),
	}
}

// Create a CommodityData to display a single class in the
// class.html template. This is added dynamically to the DisplayData
// template when the Commodity view is requested
//
//	u: the user
//	message: any message
//	id: the id of the social class to display
//
//	returns: classData which references this class, and embeds an OutputData
func (u User) ClassDisplayData(message string, id int) ClassData {
	return ClassData{
		u.CreateDisplayData(message),
		*ViewedObject[Class](u, `classes`, id),
	}
}

// Create an IndustryData to display a single industry in the
// industry.html template. This is added dynamically to the DisplayData
// template when the Commodity view is requested
//
//	u: the user
//	message: any message
//	id: the id of the industry item to display
//
//	returns: industryData which references this industry, and embeds an OutputData
func (u User) IndustryDisplayData(message string, id int) IndustryData {
	return IndustryData{
		u.CreateDisplayData(message),
		*ViewedObject[Industry](u, `industries`, id),
	}
}

// Embedded data for a single commodity, to pass into templates
type CommodityData struct {
	TemplateData
	Commodity Commodity
}

// Embedded data for a single class, to pass into templates
type ClassData struct {
	TemplateData
	Class Class
}

// Embedded data for a single industry, to pass into templates
type IndustryData struct {
	TemplateData
	Industry Industry
}
