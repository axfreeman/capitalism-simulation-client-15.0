// PATH: go-auth/controllers/auth.go

package controllers

import (
	"encoding/json"
	"fmt"
	"gorilla-client/api"
	"gorilla-client/config"
	"gorilla-client/db"
	"gorilla-client/models"
	"gorilla-client/utils"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var Store = sessions.NewCookieStore([]byte("super-secret-password"))
var Tpl *template.Template

// Convenience type for passing Messages into templates
type MessageData struct {
	Message string
}

// registerHandler serves form for registering new users
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	utils.TraceInfo(utils.BrightGreen, "Enter RegisterHandler")
	Tpl.ExecuteTemplate(w, "register.html", nil)
}

// registerAuthHandler creates a new user
func RegisterAuthHandler(w http.ResponseWriter, r *http.Request) {
	utils.TraceInfo(utils.BrightGreen, "Enter RegisterAuthHandler")

	// validate the username
	r.ParseForm()

	username := r.FormValue("username")
	if len(username) < 2 {
		Tpl.ExecuteTemplate(w, "register.html", "Username is too short")
		return
	}

	utils.TraceInfo(utils.BrightGreen, "User Name is valid")

	// check if username already exists for availability
	_, err := db.DataBase.FindRegisteredUser(username)
	if err == nil {
		utils.TraceInfo(utils.BrightGreen, "User already exists")
		Tpl.ExecuteTemplate(w, "register.html", MessageData{"User already exists"})
		return
	}
	utils.TraceInfo(utils.BrightGreen, "User Name is new")

	// create hash from password
	password := r.FormValue("password")
	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		utils.TraceError(fmt.Sprint("bcrypt err:", err))
		Tpl.ExecuteTemplate(w, "register.html", MessageData{fmt.Sprintf("Encryption problem. Please report this to the developer\n%v", err)})
		return
	}
	utils.TraceInfo(utils.BrightGreen, "Pasword is valid")

	// Create a skeleton user to save on the server. The server will fill out details
	registeredUser := models.NewRegisteredUser(username, string(hash), "")
	registeredUserServerRequest := models.RegisteredUserServerRequest{UserName: username}

	utils.TraceInfo(utils.BrightGreen, "Prototype new registered user created")

	// send the skeleton details to the server to construct a fullblown user
	body, _ := json.Marshal(registeredUserServerRequest)
	status, err := api.AdminPostRequest(config.Config.ApiSource+"/admin/register", body)

	if status == http.StatusConflict {
		utils.TraceInfo(utils.BrightGreen, "The server already knows about this user. No worries")
	} else {
		if err != nil {
			message := fmt.Sprintf("The server didn't like the request. It said %v", err)
			utils.TraceError(message)
			Tpl.ExecuteTemplate(w, "register.html", MessageData{message})
			return
		}
	}

	utils.TraceInfof(utils.BrightGreen, "Server says the registration can go ahead")

	// Fetch the user's API key
	status, err = api.AdminGetRequest(config.Config.ApiSource+"/admin/user/"+username, &registeredUser)
	utils.TraceInfof(utils.BrightWhite, "Fetching the server record which returned status %v and err %v", status, err)
	if err != nil {
		message := fmt.Sprintf("The server didn't send an api key for the new user. It said %v", err)
		utils.TraceError(message)
		Tpl.ExecuteTemplate(w, "register.html", MessageData{message})
		return
	}
	utils.TraceInfof(utils.BrightWhite, "The server returned user details as follows:%s", registeredUser.Write())

	// Save the user to the database
	db.DataBase.CreateRegisteredUser((registeredUser))
	Tpl.ExecuteTemplate(w, "login.html", nil)
}

// loginHandler serves form for users to login with
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	utils.TraceInfo(utils.BrightGreen, "Enter LoginHandler")
	Tpl.ExecuteTemplate(w, "login.html", nil)
	utils.TraceInfo(utils.BrightGreen, "Exit LoginHandler")
}

// loginAuthHandler authenticates user login
func LoginAuthHandler(w http.ResponseWriter, r *http.Request) {
	var registeredUser *models.RegisteredUser
	var err error

	utils.TraceInfo(utils.BrightGreen, "Enter LoginAuthHandler")

	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	utils.TraceInfo(utils.BrightGreen, fmt.Sprintf("Request to log in from User %s with password %s", username, password))
	if registeredUser, err = db.DataBase.FindRegisteredUser(username); err != nil {
		utils.TraceError(fmt.Sprintf("User %s is not registered", username))
		Tpl.ExecuteTemplate(w, "login.html", "Check the username and the password")
		return
	}
	utils.TraceInfof(utils.BrightGreen, "The api key for the retrieved RegisteredUser is %s", registeredUser.ApiKey)

	err = bcrypt.CompareHashAndPassword([]byte(registeredUser.Password), []byte(password))
	if err != nil {
		utils.TraceError("Incorrect password")
		Tpl.ExecuteTemplate(w, "login.html", nil)
		return
	}

	// Send the RegisteredUser to the server and retrieve a fullblown user from it
	// This ensures we capture any updates since we registered
	// For example if the user logs in from two devices or browsers
	user := models.NewUser(username)
	status, err := api.AdminGetRequest(config.Config.ApiSource+"/admin/user/"+username, &user)
	utils.TraceInfo(utils.BrightGreen, fmt.Sprintf("The server responded with status %d and error %v", status, err))
	if status != http.StatusOK {
		utils.TraceError("The server doesn't know this user, sorry")
		Tpl.ExecuteTemplate(w, "login.html", "Check username and password")
		return
	}

	// save the name in the authentication store
	session, _ := Store.Get(r, "session") // session struct has field make(map[interface{}]interface{})
	session.Values["userID"] = username
	session.Save(r, w) // save before writing to response/return from handler
	utils.TraceInfo(utils.BrightGreen, fmt.Sprintf("User %s has successfully logged in ", registeredUser.UserName))

	// set the user's api key which was fixed when the user registered
	utils.TraceInfof(utils.BrightGreen, "collecting the api key from the RegisteredUser record. The apikey is:%s", registeredUser.ApiKey)
	user.ApiKey = registeredUser.ApiKey
	utils.TraceInfof(utils.BrightGreen, "the user api key is now:%s", user.ApiKey)

	// Add the fullblown user to the client list of logged-in users
	models.ClientLoggedInUsers[username] = user

	//Grab all the templates from the server
	//See note in DOCS folder
	api.FetchRemoteTemplates()

	//Grab this user's data from the server TODO degrade gracefully if this doesn't work
	utils.TraceInfof(utils.BrightGreen, "the user's current simulation is %d", user.CurrentSimulationID)
	if user.CurrentSimulationID != 0 {
		api.FetchUserObjects(user)
	}

	// display the welcome screen
	user.CurrentPage = "welcome.html"
	Tpl.ExecuteTemplate(w, user.CurrentPage, user.TemplateData(""))
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	utils.TraceInfo(utils.BrightGreen, "Entered LogoutHandler")
	session, _ := Store.Get(r, "session")
	delete(session.Values, "userID")
	session.Save(r, w)
	Tpl.ExecuteTemplate(w, "login.html", "Logged Out")
}

// Auth adds authentication code to handler before returning handler
// func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request)
func Auth(HandlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := Store.Get(r, "session")
		content, ok := session.Values["userID"]
		if !ok {
			http.Redirect(w, r, "auth/login", http.StatusFound)
			return
		}
		utils.TraceInfof(utils.BrightGreen, "Auth was called and retrieved %s", content)

		// Check that the cookie refers to a logged in user
		_, ok = models.ClientLoggedInUsers[content.(string)]
		if !ok {
			http.Redirect(w, r, "auth/login", http.StatusFound)
			return
		}

		// ServeHTTP calls f(w, r)
		// func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request)
		HandlerFunc.ServeHTTP(w, r)
	}
}

// check session for logged in done with middleware Auth()
func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	utils.TraceInfo(utils.BrightGreen, "Enter WelcomeHandler")
	user := CurrentUser(r)
	Tpl.ExecuteTemplate(w, "welcome.html", user.TemplateData(""))
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Enter AboutHandler")
	// check session for logged in done with middleware Auth()
	/*
		session, _ := store.Get(r, "session")
		_, ok := session.Values["userID"]
		if !ok {
			http.Redirect(w, r, "/login", 302)
			return
		}
	*/
	Tpl.ExecuteTemplate(w, "about.html", "Logged In")
}

func CurrentUser(r *http.Request) *models.User {
	session, _ := Store.Get(r, "session")
	content := session.Values["userID"]
	return models.ClientLoggedInUsers[content.(string)]
}
