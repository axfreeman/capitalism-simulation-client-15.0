package models

import "gorilla-client/utils"

// Wrappers for lists of Tables

// List of the user's Simulations.
//
//	u: the user
//	returns:
//	 Slice of SimulationsList
//	 If the user has no simulations, an empty slice
func (u User) SimulationsList() *[]Simulation {
	list := u.Simulation.Table.(*[]Simulation)
	if len(*list) == 0 {
		var fakeList []Simulation = *new([]Simulation)
		return &fakeList
	}
	return list
}

func (u User) Commodities() *[]Commodity {
	return (*u.TableSets[u.ViewedTimeStamp])["commodities"].Table.(*[]Commodity)
}

func (u User) CommodityViews() *[]CommodityView {
	v := (*u.TableSets[u.ViewedTimeStamp])["commodities"].Table.(*[]Commodity)
	c := (*u.TableSets[u.ComparatorTimeStamp])["commodities"].Table.(*[]Commodity)
	return NewCommodityViews(v, c)
}

func (u User) Industries() *[]Industry {
	return (*u.TableSets[u.ViewedTimeStamp])["industries"].Table.(*[]Industry)
}

func (u User) IndustryViews() *[]IndustryView {
	v := (*u.TableSets[u.ViewedTimeStamp])["industries"].Table.(*[]Industry)
	c := (*u.TableSets[u.ComparatorTimeStamp])["industries"].Table.(*[]Industry)

	return NewIndustryViews(u.ViewedTimeStamp, u.ComparatorTimeStamp, v, c)
}

func (u User) ClassViews() *[]ClassView {
	v := (*u.TableSets[u.ViewedTimeStamp])["classes"].Table.(*[]Class)
	c := (*u.TableSets[u.ComparatorTimeStamp])["classes"].Table.(*[]Class)

	return NewClassViews(u.ViewedTimeStamp, u.ComparatorTimeStamp, v, c)
}

func (u User) Classes() *[]Class {
	return (*u.TableSets[u.ViewedTimeStamp])["classes"].Table.(*[]Class)
}

// Wrapper for the IndustryStockList
func (u User) IndustryStocks(timeStamp int) *[]Industry_Stock {
	return (*u.TableSets[timeStamp])["industry stocks"].Table.(*[]Industry_Stock)
}

// Wrapper for the ClassStockList
func (u User) ClassStocks(timeStamp int) *[]Class_Stock {
	return (*u.TableSets[timeStamp])["class stocks"].Table.(*[]Class_Stock)
}

// Wrapper for the TraceList
func (u User) Traces(timeStamp int) *[]Trace {
	if len(u.TableSets) == 0 {
		return nil
	}
	table, ok := (*u.TableSets[timeStamp])["trace"]
	if !ok {
		return nil
	}
	return table.Table.(*[]Trace)
}

// supplies outputData to be passed into Templates for display
//
//	u: a user
//
//	returns:
//      if the user has no simulations, just the template list
//      otherwise, the output data the users current simulation
func (u *User) TemplateData(message string) OutputData {
	slist := u.SimulationsList()
	state := u.Get_current_state()
	utils.TraceInfof(utils.BrightYellow, "Entering TemplateData for user %s with simulationID %d", u.UserName, u.CurrentSimulationID)
	if u.CurrentSimulationID == 0 {
		utils.TraceInfo(utils.BrightYellow, "User has no simulations")
		return OutputData{
			Title:           "Hello",
			Simulations:     nil,
			Templates:       &TemplateList,
			Count:           0,
			Username:        u.UserName,
			State:           state,
			CommodityViews:  nil,
			IndustryViews:   nil,
			ClassViews:      nil,
			Industry_Stocks: nil,
			Class_Stocks:    nil,
			Trace:           nil,
			Message:         message,
		}
	}
	return OutputData{
		Title:           "Hello",
		Simulations:     slist,
		Templates:       &TemplateList,
		Count:           len(*slist),
		Username:        u.UserName,
		State:           state,
		CommodityViews:  u.CommodityViews(),
		IndustryViews:   u.IndustryViews(),
		ClassViews:      u.ClassViews(),
		Industry_Stocks: u.IndustryStocks(u.ViewedTimeStamp),
		Class_Stocks:    u.ClassStocks(u.ViewedTimeStamp),
		Trace:           u.Traces(u.ViewedTimeStamp),
		Message:         message,
	}
}
