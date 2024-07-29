// display.objects.go
// handlers to display the objects of the simulation on the user's browser

package controllers

import (
	"fmt"
	"gorilla-client/utils"
	"net/http"
)

// display all commodities in the current simulation
func ShowCommodities(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	user.CurrentPage = "commodities.html"
	utils.TraceInfof(utils.BrightYellow, "Fetching commodities for user %s", user.UserName)
	Tpl.ExecuteTemplate(w, user.CurrentPage, user.TemplateData(""))
}

// display all industries in the current simulation
func ShowIndustries(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	user.CurrentPage = "industries.html"
	utils.TraceInfof(utils.BrightYellow, "Fetching industries for user %s", user.UserName)
	Tpl.ExecuteTemplate(w, user.CurrentPage, user.TemplateData(""))
}

// display all classes in the current simulation
func ShowClasses(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	user.CurrentPage = "classes.html"
	utils.TraceInfo(utils.BrightYellow, fmt.Sprintf("Fetching classes for user %s", user.UserName))
	Tpl.ExecuteTemplate(w, user.CurrentPage, user.TemplateData(""))
}

// display all industry stocks in the current simulation
func ShowIndustryStocks(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	user.CurrentPage = "industry_stocks.html"
	utils.TraceInfof(utils.BrightYellow, "Fetching industry stocks for user %s", user.UserName)
	Tpl.ExecuteTemplate(w, user.CurrentPage, user.TemplateData(""))
}

// display all the class stocks in the current simulation
func ShowClassStocks(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	user.CurrentPage = "class_stocks.html"
	utils.TraceInfof(utils.BrightYellow, "Fetching class stocks for user %s", user.UserName)
	Tpl.ExecuteTemplate(w, user.CurrentPage, user.TemplateData(""))
}

// display all Trace records in the current simulation
func ShowTrace(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	user.CurrentPage = "trace.html"
	utils.TraceInfof(utils.BrightYellow, "Fetching classes for user %s", user.UserName)
	Tpl.ExecuteTemplate(w, user.CurrentPage, user.TemplateData(""))
}

// Display one specific commodity
func ShowCommodity(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	user.CurrentPage = "commodity.html"
	utils.TraceInfo(utils.BrightYellow, fmt.Sprintf("Fetching a commodity for user %s", user.UserName))
	Tpl.ExecuteTemplate(w, user.CurrentPage, user.TemplateData(""))
}

// Display one specific industry
func ShowIndustry(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	user.CurrentPage = "industry.html"
	utils.TraceInfo(utils.BrightYellow, fmt.Sprintf("Fetching an industry for user %s", user.UserName))
	Tpl.ExecuteTemplate(w, user.CurrentPage, user.TemplateData(""))
}

// Display one specific class
func ShowClass(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	user.CurrentPage = "class.html"
	utils.TraceInfof(utils.BrightYellow, "Fetching a class for user %s", user.UserName)
	Tpl.ExecuteTemplate(w, user.CurrentPage, user.TemplateData(""))
}

// Displays a snapshot of the economy
func ShowIndexPage(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	user.CurrentPage = "index.html"
	utils.TraceInfo(utils.BrightYellow, fmt.Sprintf("Showing Index Page for user %s", user.UserName))
	Tpl.ExecuteTemplate(w, user.CurrentPage, user.TemplateData(""))
}

// TODO not working yet
func SwitchSimulation(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	utils.UNUSED(user)
}

// TODO not working yet
func DeleteSimulation(w http.ResponseWriter, r *http.Request) {
}

// TODO not working yet
func RestartSimulation(w http.ResponseWriter, r *http.Request) {
}

func UserDashboard(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	user.CurrentPage = "user-dashboard.html"
	Tpl.ExecuteTemplate(w, user.CurrentPage, user.TemplateData(""))
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	Tpl.ExecuteTemplate(w, "404.html", "")
}
