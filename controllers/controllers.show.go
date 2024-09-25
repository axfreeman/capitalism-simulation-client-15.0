// display.objects.go
// handlers to display the objects of the simulation on the user's browser

package controllers

import (
	"fmt"
	"gorilla-client/models"
	"gorilla-client/utils"
	"net/http"
)

// display all commodities in the current simulation
func ShowCommodities(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	user.CurrentPage = models.CurrentPageType{Url: "commodities.html", Id: 0}

	utils.TraceInfof(utils.BrightYellow, "Fetching commodities for user %s", user.UserName)
	Tpl.ExecuteTemplate(w, user.CurrentPage.Url, user.CreateTemplateData(""))
}

// display all industries in the current simulation
func ShowIndustries(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	user.CurrentPage = models.CurrentPageType{Url: "industries.html", Id: 0}

	utils.TraceInfof(utils.BrightYellow, "Fetching industries for user %s", user.UserName)
	Tpl.ExecuteTemplate(w, user.CurrentPage.Url, user.CreateTemplateData(""))
}

// display all classes in the current simulation
func ShowClasses(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	user.CurrentPage = models.CurrentPageType{Url: "classes.html", Id: 0}

	utils.TraceInfo(utils.BrightYellow, fmt.Sprintf("Fetching classes for user %s", user.UserName))
	Tpl.ExecuteTemplate(w, user.CurrentPage.Url, user.CreateTemplateData(""))
}

// display all industry stocks in the current simulation
func ShowIndustryStocks(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	user.CurrentPage = models.CurrentPageType{Url: "industry_stocks.html", Id: 0}

	utils.TraceInfof(utils.BrightYellow, "Fetching industry stocks for user %s", user.UserName)
	Tpl.ExecuteTemplate(w, user.CurrentPage.Url, user.CreateTemplateData(""))
}

// display all the class stocks in the current simulation
func ShowClassStocks(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	user.CurrentPage = models.CurrentPageType{Url: "class_stocks.html", Id: 0}

	utils.TraceInfof(utils.BrightYellow, "Fetching class stocks for user %s", user.UserName)
	Tpl.ExecuteTemplate(w, user.CurrentPage.Url, user.CreateTemplateData(""))
}

// display all Trace records in the current simulation
func ShowTrace(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	user.CurrentPage = models.CurrentPageType{Url: "trace.html", Id: 0}
	utils.TraceInfof(utils.BrightYellow, "Fetching classes for user %s", user.UserName)
	t := models.Traces(user)
	utils.UNUSED(t)
	Tpl.ExecuteTemplate(w, user.CurrentPage.Url, models.TemplateData{Trace: t}) // nearly empty TemplateData object with just Trace in it
}

// Display one specific commodity
func ShowCommodity(w http.ResponseWriter, r *http.Request) {
	var err error
	var id int
	user := CurrentUser(r)
	if id, err = FetchIDfromURL(r); err != nil {
		ReportError(user, w, err.Error())
	}
	user.CurrentPage = models.CurrentPageType{Url: "commodity.html", Id: id}

	utils.TraceInfof(utils.BrightYellow, "Fetching commodity %d for user %s", id, user.UserName)
	Tpl.ExecuteTemplate(w,
		user.CurrentPage.Url,
		user.CommodityDisplayData("", id))
}

// Display one specific industry
func ShowIndustry(w http.ResponseWriter, r *http.Request) {
	var err error
	var id int
	user := CurrentUser(r)
	if id, err = FetchIDfromURL(r); err != nil {
		ReportError(user, w, err.Error())
	}
	user.CurrentPage = models.CurrentPageType{Url: "industry.html", Id: id}

	utils.TraceInfof(utils.BrightYellow, "Fetching industry %d for user %s", id, user.UserName)
	Tpl.ExecuteTemplate(w, user.CurrentPage.Url, user.IndustryDisplayData("", id))
}

// Display one specific class
func ShowClass(w http.ResponseWriter, r *http.Request) {
	var err error
	var id int
	user := CurrentUser(r)
	if id, err = FetchIDfromURL(r); err != nil {
		ReportError(user, w, err.Error())
	}
	user.CurrentPage = models.CurrentPageType{Url: "class.html", Id: id}

	utils.TraceInfof(utils.BrightYellow, "Fetching class %d for user %s", id, user.UserName)
	Tpl.ExecuteTemplate(w, user.CurrentPage.Url, user.ClassDisplayData("", id))
}

// Displays a snapshot of the economy
func ShowIndexPage(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	user.CurrentPage = models.CurrentPageType{Url: "index.html", Id: 0}

	utils.TraceInfo(utils.BrightYellow, fmt.Sprintf("Showing Index Page for user %s", user.UserName))
	Tpl.ExecuteTemplate(w, user.CurrentPage.Url, user.CreateTemplateData(""))
}

func UserDashboard(w http.ResponseWriter, r *http.Request) {
	user := CurrentUser(r)
	user.CurrentPage = models.CurrentPageType{Url: "user-dashboard.html", Id: 0}

	Tpl.ExecuteTemplate(w, user.CurrentPage.Url, user.CreateTemplateData(""))
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	Tpl.ExecuteTemplate(w, "404.html", "")
}

// check session for logged in done with middleware Auth()
func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	utils.TraceInfo(utils.BrightGreen, "Enter WelcomeHandler")
	user := CurrentUser(r)
	Tpl.ExecuteTemplate(w, "welcome.html", user.CreateTemplateData(""))
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Enter AboutHandler")
	Tpl.ExecuteTemplate(w, "about.html", "Logged In")
}
